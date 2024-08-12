package service

import (
	"context"
	"go_task_service/config"
	"go_task_service/genproto/task_service"
	"go_task_service/genproto/user_service"

	"go_task_service/grpc/client"
	"go_task_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.GrpcClientI
}

func NewTaskService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.GrpcClientI) *TaskService {
	return &TaskService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (f *TaskService) Create(ctx context.Context, req *task_service.CreateTask) (*task_service.GetTask, error) {

	f.log.Info("---CreateTask--->>>", logger.Any("req", req))
	_, err := f.services.User().Check(ctx, &user_service.UserPrimaryKey{Id: req.UserId})
	if err != nil {
		f.log.Error("error while check group")
		return nil, err
	}
	resp, err := f.strg.Task().Create(ctx, req)
	if err != nil {
		f.log.Error("---CreateTask--->>>", logger.Error(err))
		return &task_service.GetTask{}, err
	}

	return resp, nil
}
func (f *TaskService) Update(ctx context.Context, req *task_service.UpdateTask) (*task_service.GetTask, error) {

	f.log.Info("---UpdateTask--->>>", logger.Any("req", req))

	resp, err := f.strg.Task().Update(ctx, req)
	if err != nil {
		f.log.Error("---UpdateTask--->>>", logger.Error(err))
		return &task_service.GetTask{}, err
	}

	return resp, nil
}

func (f *TaskService) GetList(ctx context.Context, req *task_service.GetListTaskRequest) (*task_service.GetListTaskResponse, error) {
	f.log.Info("---GetListTask--->>>", logger.Any("req", req))

	resp, err := f.strg.Task().GetList(ctx, req)
	if err != nil {
		f.log.Error("---GetListTask--->>>", logger.Error(err))
		return &task_service.GetListTaskResponse{}, err
	}

	return resp, nil
}

func (f *TaskService) GetByID(ctx context.Context, id *task_service.TaskPrimaryKey) (*task_service.GetTask, error) {
	f.log.Info("---GetTask--->>>", logger.Any("req", id))

	resp, err := f.strg.Task().GetById(ctx, id)
	if err != nil {
		f.log.Error("---GetTask--->>>", logger.Error(err))
		return &task_service.GetTask{}, err
	}

	return resp, nil
}

func (f *TaskService) GetByExternalId(ctx context.Context, id *task_service.TaskPrimaryKey) (*task_service.GetTask, error) {
	f.log.Info("---GetTask--->>>", logger.Any("req", id))

	resp, err := f.strg.Task().GetByExternalId(ctx, id)
	if err != nil {
		f.log.Error("---GetTask--->>>", logger.Error(err))
		return &task_service.GetTask{}, err
	}

	return resp, nil
}

func (f *TaskService) Delete(ctx context.Context, req *task_service.TaskPrimaryKey) (*emptypb.Empty, error) {

	f.log.Info("---DeleteTask--->>>", logger.Any("req", req))

	_, err := f.strg.Task().Delete(ctx, req)
	if err != nil {
		f.log.Error("---DeleteTask--->>>", logger.Error(err))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (f *TaskService) ChangeStatus(ctx context.Context, req *task_service.TaskChangeStatus) (*task_service.TaskChangeStatusResp, error) {

	f.log.Info("---UpdateTask--->>>", logger.Any("req", req))

	resp, err := f.strg.Task().ChangeStatus(ctx, req)
	if err != nil {
		f.log.Error("---UpdateTask--->>>", logger.Error(err))
		return &task_service.TaskChangeStatusResp{}, err
	}

	return resp, nil
}
