package usecase

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/repository"
)

type UserUseCase struct {
	userRepo repository.UserRepoController
}

// NewUserUsecase implicitly links UserManager to userUsecase struct
func NewUserUsecase(
	userRepo repository.UserRepoController,
) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

// CreateUser creates new user and stores it in DB with repository.UserRepository
func (uu *UserUseCase) CreateUser(u *model.User) (*repository.InsertResult, error) {
	result, err := uu.userRepo.InsertUserToStorage(u)
	return result, err
}

// FindUser looks up for user (u *model.User) and returns it if found
func (uu *UserUseCase) FindUser(u *model.User) (*model.User, error) {
	userFromDB, err := uu.userRepo.FindUserInStorage(u.Username)
	return userFromDB, err
}

func (uu *UserUseCase) DeleteUser(username string) error {
	return uu.userRepo.DeleteUserFromStorage(username)
}

func (uu *UserUseCase) UpdateUsername(newUsername, oldUsername string) error {
	var u model.User
	u.Username = newUsername

	return uu.userRepo.UpdateUsernameInStorage(&u, oldUsername)
}

func (uu *UserUseCase) UpdatePassword(u *model.User) error {
	return uu.userRepo.UpdatePasswordInStorage(u)
}
