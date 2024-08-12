package grpc

import (
	"go_task_service/config"
	"go_task_service/genproto/task_service"

	"go_task_service/grpc/client"
	"go_task_service/grpc/service"
	"go_task_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvc client.GrpcClientI) (grpcServer *grpc.Server) {

	grpcServer = grpc.NewServer()

	task_service.RegisterTaskServiceServer(grpcServer, service.NewTaskService(cfg, log, strg, srvc))

	reflection.Register(grpcServer)
	return
}
