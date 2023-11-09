package usecase

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/repository"
)

type UserUsecase struct {
	ur repository.UserRepoManager
	pr repository.ProfileRepoManager
}

// NewUserUsecase implicitly links UserManager to userUsecase struct
func NewUserUsecase(
	ur repository.UserRepoManager,
	pr repository.ProfileRepoManager,
) *UserUsecase {
	return &UserUsecase{ur: ur, pr: pr}
}

// CreateUser creates new user and stores it in DB with repository.UserRepository
func (uu *UserUsecase) CreateUser(u *model.User) (*repository.InsertResult, error) {
	result, err := uu.ur.Create(u)
	return result, err
}

// CreateProfile creates profile for model.User and stores it in DB with repository.ProfileRepository
func (uu *UserUsecase) CreateProfile(p *model.Profile) (*repository.InsertResult, error) {
	result, err := uu.pr.Create(p)
	return result, err
}

// Find looks up for user (u *model.User) and returns it if found
func (uu *UserUsecase) Find(u *model.User) (*model.User, error) {
	userFromDB, err := uu.ur.Find(u)
	return userFromDB, err
}

func (uu *UserUsecase) UpdateUsername(newUsername, oldUsername string) error {
	var u model.User
	u.Username = newUsername

	return uu.ur.UpdateUsername(&u, oldUsername)
}

func (uu *UserUsecase) UpdatePassword(u *model.User) error {
	return uu.ur.UpdatePassword(u)
}

// UpdateProfile updates profile of model.User in DB with repository.ProfileRepository
func (uu *UserUsecase) UpdateProfile(p *model.Update, authUsername string) error {
	return uu.pr.Update(p, authUsername)
}
