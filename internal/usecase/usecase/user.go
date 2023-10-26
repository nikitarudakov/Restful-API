package usecase

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/repository"
)

type userUsecase struct {
	ur repository.UserRepository
}

type UserInput interface {
	Create(u *model.User) error
}

type UserOutput interface {
	Render(i interface{}) error
}

func NewUserUsecase(r repository.UserRepository) UserInput {
	return &userUsecase{ur: r}
}

func (uu *userUsecase) Create(u *model.User) error {
	// Create user at database
	_, err := uu.ur.Create(u)
	if err != nil {
		return err
	}

	return nil
}
