package grpc

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/grpc/auth"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/grpc/user"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/grpc/vote"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserManagementGRPCService struct {
	AuthClient auth.AuthControllerClient
	UserClient user.UserControllerClient
	VoteClient vote.VoteControllerClient
}

// NewUserManagementGRPCService creates a new gRPC user service connection using the specified connection string.
func NewUserManagementGRPCService(connString string) (*UserManagementGRPCService, error) {
	conn, err := grpc.Dial(connString, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if conn == nil || err != nil {
		return nil, err
	}
	return &UserManagementGRPCService{
		AuthClient: auth.NewAuthControllerClient(conn),
		UserClient: user.NewUserControllerClient(conn),
		VoteClient: vote.NewVoteControllerClient(conn),
	}, nil
}
