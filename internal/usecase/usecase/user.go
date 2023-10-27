package usecase

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/repository"
)

type userUsecase struct {
	ur repository.UserRepository
	pr repository.ProfileRepository
}

type UserInput interface {
	Create(u *model.User) error
	Find(u *model.User) (*model.User, error)
	CreateProfile(p *model.Profile) (interface{}, error)
	UpdateProfile(p *model.Profile) error
}

type UserOutput interface {
	Render(i interface{}) error
}

func NewUserUsecase(ur repository.UserRepository, pr repository.ProfileRepository) UserInput {
	return &userUsecase{ur: ur, pr: pr}
}

func (uu *userUsecase) Create(u *model.User) error {
	// Create user at database
	_, err := uu.ur.Create(u)
	if err != nil {
		return err
	}

	return nil
}

func (uu *userUsecase) Find(u *model.User) (*model.User, error) {
	// Find user at database
	userFromDB, err := uu.ur.Find(u)
	if err != nil {
		return nil, err
	}

	return userFromDB, nil
}

func (uu *userUsecase) CreateProfile(p *model.Profile) (interface{}, error) {
	insertedID, err := uu.pr.Create(p)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

func (uu *userUsecase) UpdateProfile(p *model.Profile) error {
	err := uu.pr.Update(p)
	if err != nil {
		return err
	}

	return nil
}
