package grpc

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/profile"
	"google.golang.org/grpc"
)

func NewProfileDAOServer() (*grpc.Server, error) {
	profileDAO := profile.ProfileDAOServer{}

	grpcServer := grpc.NewServer()

	profile.RegisterProfileRepoServer(grpcServer, &profileDAO)

	return grpcServer, nil
}
