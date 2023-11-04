package usecase

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	repository2 "git.foxminded.ua/foxstudent106092/user-management/internal/business/usecase/repository"
)

type UserUsecase struct {
	ur repository2.UserRepository
	pr repository2.ProfileRepository
}

// NewUserUsecase implicitly links UserManager to userUsecase struct
func NewUserUsecase(ur repository2.UserRepository, pr repository2.ProfileRepository) *UserUsecase {
	return &UserUsecase{ur: ur, pr: pr}
}

// CreateUser creates new user and stores it in DB with repository.UserRepository
func (uu *UserUsecase) CreateUser(u *model.User) error {
	_, err := uu.ur.Create(u)
	if err != nil {
		return err
	}

	return nil
}

// CreateProfile creates profile for model.User and stores it in DB with repository.ProfileRepository
func (uu *UserUsecase) CreateProfile(p *model.Profile) (interface{}, error) {
	insertedID, err := uu.pr.Create(p)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

// Find looks up for user (u *model.User) and returns it if found
func (uu *UserUsecase) Find(u *model.User) (*model.User, error) {
	userFromDB, err := uu.ur.Find(u)
	if err != nil {
		return nil, err
	}

	return userFromDB, nil
}

func (uu *UserUsecase) UpdateUsername(newUsername, oldUsername string) error {
	var u model.User

	u.Username = newUsername
	if err := uu.ur.UpdateUsername(&u, oldUsername); err != nil {
		return err
	}

	return nil
}

func (uu *UserUsecase) UpdatePassword(u *model.User) error {
	if err := uu.ur.UpdatePassword(u); err != nil {
		return err
	}

	return nil
}

// UpdateProfile updates profile of model.User in DB with repository.ProfileRepository
func (uu *UserUsecase) UpdateProfile(p *model.Update, authUsername string) error {
	err := uu.pr.Update(p, authUsername)
	if err != nil {
		return err
	}

	return nil
}
