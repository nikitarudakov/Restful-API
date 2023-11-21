package usecase

import (
	"context"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/userDao"
)

type UserUseCase struct {
	grpcClient userDao.UserRepoClient
}

// NewUserUsecase implicitly links UserManager to userUsecase struct
func NewUserUsecase(
	userRepo userDao.UserRepoClient,
) *UserUseCase {
	return &UserUseCase{grpcClient: userRepo}
}

// CreateUser creates new user and stores it in DB with repository.UserRepository
func (uu *UserUseCase) CreateUser(u *model.User) (*userDao.InsertResult, error) {
	convertedUser := userDao.MarshalTogRPCUserObj(u)

	insertResult, err := uu.grpcClient.InsertUserToStorage(context.Background(), convertedUser)

	return insertResult, err
}

// FindUser looks up for user (u *model.User) and returns it if found
func (uu *UserUseCase) FindUser(u *model.User) (*model.User, error) {
	gRPCUser, err := uu.grpcClient.FindUserInStorage(context.Background(),
		&userDao.Username{Val: u.Username})

	user := userDao.UnmarshalToUser(gRPCUser)

	return user, err
}

func (uu *UserUseCase) DeleteUser(username string) error {
	_, err := uu.grpcClient.DeleteUserFromStorage(context.Background(),
		&userDao.Username{Val: username})

	return err
}

func (uu *UserUseCase) UpdateUsername(newUsername, oldUsername string) error {
	var u model.User
	u.Username = newUsername

	gRPCUserReplace := userDao.MarshalTogRPCUserReplace(&u, oldUsername)

	_, err := uu.grpcClient.UpdateUsernameInStorage(context.Background(), gRPCUserReplace)

	return err
}

func (uu *UserUseCase) UpdatePassword(u *model.User) error {
	gRPCUserObj := userDao.MarshalTogRPCUserObj(u)

	_, err := uu.grpcClient.UpdatePasswordInStorage(context.Background(), gRPCUserObj)

	return err
}
