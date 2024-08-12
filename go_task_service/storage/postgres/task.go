package postgres

import (
	"context"
	"database/sql"
	"fmt"
	task "go_task_service/genproto/task_service"
	"go_task_service/grpc/client"
	"go_task_service/pkg"
	"go_task_service/storage"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type taskRepo struct {
	db       *pgxpool.Pool
	services client.GrpcClientI
}

func NewTaskRepo(db *pgxpool.Pool, srvs client.GrpcClientI) storage.TaskRepoI {
	return &taskRepo{
		db:       db,
		services: srvs,
	}
}

func generateExternalId(db *pgxpool.Pool, ctx context.Context) (string, error) {
	var nextVal int
	err := db.QueryRow(ctx, "SELECT nextval('task_external_id_seq')").Scan(&nextVal)
	if err != nil {
		return "", err
	}
	task_id := "T" + fmt.Sprintf("%05d", nextVal)
	return task_id, nil
}

func (c *taskRepo) Create(ctx context.Context, req *task.CreateTask) (*task.GetTask, error) {

	id := uuid.NewString()
	external_id, err := generateExternalId(c.db, ctx)
	if err != nil {
		log.Println("error while creating external_id", err)
		return nil, err
	}

	comtag, err := c.db.Exec(ctx, `
		INSERT INTO tasks (
			id,
			user_id,
			external_id,
			title,
			task_status,
			task_description,
			deadline
		) VALUES ($1,$2,$3,$4,$5,$6,$7
		)`,
		id,
		req.UserId,
		external_id,
		req.Title,
		req.TaskStatus,
		req.TaskDescription,
		req.Deadline)
	if err != nil {
		log.Println("error while creating tasks", comtag)
		return nil, err
	}

	tasks, err := c.GetById(ctx, &task.TaskPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting tasks by id")
		return nil, err
	}
	return tasks, nil
}

func (c *taskRepo) Update(ctx context.Context, req *task.UpdateTask) (*task.GetTask, error) {

	_, err := c.db.Exec(ctx, `
		UPDATE tasks SET
		user_id = $1,
		title = $2,
		task_description = $3,
		deadline = $4,
		updated_at = NOW()
		WHERE id = $5
		`,
		req.UserId,
		req.Title,
		req.TaskDescription,
		req.Deadline,
		req.Id)
	if err != nil {
		log.Println("error while updating tasks")
		return nil, err
	}

	tasks, err := c.GetById(ctx, &task.TaskPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting tasks by id")
		return nil, err
	}
	return tasks, nil
}

func (c *taskRepo) GetById(ctx context.Context, id *task.TaskPrimaryKey) (*task.GetTask, error) {
	var (
		deadline   sql.NullString
		task       task.GetTask
		created_at sql.NullString
		updated_at sql.NullString
	)

	query := `SELECT
				id,
				user_id,
				external_id,
				title,
				task_status,
				task_description,
				deadline,
				created_at,
				updated_at
			FROM tasks
			WHERE id = $1 AND deleted_at IS NULL`

	rows := c.db.QueryRow(ctx, query, id.Id)

	if err := rows.Scan(
		&task.Id,
		&task.UserId,
		&task.ExternalId,
		&task.Title,
		&task.TaskStatus,
		&task.TaskDescription,
		&deadline,
		&created_at,
		&updated_at); err != nil {
		return &task, err
	}
	task.Deadline = pkg.NullStringToString(deadline)
	task.CreatedAt = pkg.NullStringToString(created_at)
	task.UpdatedAt = pkg.NullStringToString(updated_at)

	return &task, nil
}

func (c *taskRepo) Delete(ctx context.Context, id *task.TaskPrimaryKey) (emptypb.Empty, error) {

	_, err := c.db.Exec(ctx, `
		UPDATE tasks SET
		deleted_at = NOW()
		WHERE id = $1
		`,
		id.Id)

	if err != nil {
		return emptypb.Empty{}, err
	}
	return emptypb.Empty{}, nil
}

func (c *taskRepo) GetList(ctx context.Context, req *task.GetListTaskRequest) (*task.GetListTaskResponse, error) {
	tasks := task.GetListTaskResponse{}
	var (
		deadline   sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	filter_by_date := ""
	offset := (req.Offset - 1) * req.Limit
	if req.FromDate != "" {
		filter_by_date = fmt.Sprintln(` AND deadline BETWEEN '`, req.FromDate, `' AND '`, req.ToDate, `'`)
	}
	query := `SELECT
				id,
				user_id,
				external_id,
				title,
				task_status,
				task_description,
				deadline,
				created_at,
				updated_at
			FROM tasks 
			WHERE user_id = $1 AND deleted_at is null ` + filter_by_date + ` ORDER BY deadline OFFSET $2 LIMIT $3`

	rows, err := c.db.Query(ctx, query, req.OwnerId, offset, req.Limit)

	if err != nil {
		log.Println("error while getting all tasks")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			task task.GetListTask
		)

		if err = rows.Scan(
			&task.Id,
			&task.UserId,
			&task.ExternalId,
			&task.Title,
			&task.TaskStatus,
			&task.TaskDescription,
			&deadline,
			&created_at,
			&updated_at,
		); err != nil {
			return &tasks, err
		}
		task.Deadline = pkg.NullStringToString(deadline)
		tasks.Tasks = append(tasks.Tasks, &task)
	}

	err = c.db.QueryRow(ctx, `  SELECT count(*) FROM tasks  
								WHERE TRUE AND deleted_at is null  `+filter_by_date+``).Scan(&tasks.Count)
	if err != nil {
		return &tasks, err
	}

	return &tasks, nil
}

func (c *taskRepo) GetByExternalId(ctx context.Context, id *task.TaskPrimaryKey) (*task.GetTask, error) {
	var (
		deadline   sql.NullString
		task       task.GetTask
		created_at sql.NullString
		updated_at sql.NullString
	)

	query := `SELECT
				id,
				user_id,
				external_id,
				title,
				task_status,
				task_description,
				deadline,
				created_at,
				updated_at
			FROM tasks
			WHERE external_id = $1 AND deleted_at IS NULL`

	rows := c.db.QueryRow(ctx, query, id.Id)

	if err := rows.Scan(
		&task.Id,
		&task.UserId,
		&task.ExternalId,
		&task.Title,
		&task.TaskStatus,
		&task.TaskDescription,
		&deadline,
		&created_at,
		&updated_at); err != nil {
		return &task, err
	}
	task.Deadline = pkg.NullStringToString(deadline)
	task.CreatedAt = pkg.NullStringToString(created_at)
	task.UpdatedAt = pkg.NullStringToString(updated_at)

	return &task, nil
}

func (c *taskRepo) ChangeStatus(ctx context.Context, req *task.TaskChangeStatus) (*task.TaskChangeStatusResp, error) {
	var resp task.TaskChangeStatusResp

	_, err := c.db.Exec(ctx, `
		UPDATE tasks SET
		task_status = $1,
		updated_at = NOW()
		WHERE external_id = $2 AND deleted_at IS NULL
		`,
		req.NewStatus,
		req.TaskId)
	if err != nil {
		log.Println("error while updating task status")
		return nil, err
	}
	resp.Comment = "status changed successfully"

	return &resp, nil
}
