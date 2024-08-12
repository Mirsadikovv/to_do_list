package storage

import (
	"context"
	"go_task_service/genproto/task_service"

	"google.golang.org/protobuf/types/known/emptypb"
)

type StorageI interface {
	CloseDB()
	Task() TaskRepoI
}

type TaskRepoI interface {
	Create(context.Context, *task_service.CreateTask) (*task_service.GetTask, error)
	Update(context.Context, *task_service.UpdateTask) (*task_service.GetTask, error)
	GetById(context.Context, *task_service.TaskPrimaryKey) (*task_service.GetTask, error)
	GetByExternalId(context.Context, *task_service.TaskPrimaryKey) (*task_service.GetTask, error)
	Delete(context.Context, *task_service.TaskPrimaryKey) (emptypb.Empty, error)
	GetList(context.Context, *task_service.GetListTaskRequest) (*task_service.GetListTaskResponse, error)
	ChangeStatus(context.Context, *task_service.TaskChangeStatus) (*task_service.TaskChangeStatusResp, error)
}
