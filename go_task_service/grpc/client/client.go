package client

import (
	"go_task_service/config"

	"log"

	admin_service "go_task_service/genproto/admin_service"
	user_service "go_task_service/genproto/user_service"

	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientI interface {
	User() user_service.UserServiceClient
	Admin() admin_service.AdminServiceClient
}

type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

func New(cfg config.Config) (*GrpcClient, error) {

	connUser, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("user service dial host: %v port :%v err:%v",
			cfg.UserServiceHost, cfg.UserServicePort, err)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"user_service":  user_service.NewUserServiceClient(connUser),
			"admin_service": admin_service.NewAdminServiceClient(connUser),
		},
	}, nil
}

func (g *GrpcClient) User() user_service.UserServiceClient {
	client, ok := g.connections["user_service"].(user_service.UserServiceClient)
	if !ok {
		log.Println("failed to assert type for user_service")
		return nil
	}
	return client
}

func (g *GrpcClient) Admin() admin_service.AdminServiceClient {
	client, ok := g.connections["admin_service"].(admin_service.AdminServiceClient)
	if !ok {
		log.Println("failed to assert type for admin_service")
		return nil
	}
	return client
}
