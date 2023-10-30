package usecase

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/repository"
	"github.com/rs/zerolog/log"
)

type userUsecase struct {
	ur repository.UserRepository
	pr repository.ProfileRepository
}

// UserManager contains methods for performing operations on User/Profile datatype
type UserManager interface {
	Create(u *model.User) error
	Find(u *model.User) (*model.User, error)
	UpdatePassword(u *model.User) error
	CreateProfile(p *model.Profile) (interface{}, error)
	UpdateProfile(p *model.Profile) error
}

// NewUserUsecase implicitly links UserManager to userUsecase struct
func NewUserUsecase(ur repository.UserRepository, pr repository.ProfileRepository) UserManager {
	return &userUsecase{ur: ur, pr: pr}
}

// Create creates new user and stores it in DB with repository.UserRepository
func (uu *userUsecase) Create(u *model.User) error {
	_, err := uu.ur.Create(u)
	if err != nil {
		return err
	}

	return nil
}

// Find looks up for user (u *model.User) and returns it if found
func (uu *userUsecase) Find(u *model.User) (*model.User, error) {
	userFromDB, err := uu.ur.Find(u)
	if err != nil {
		return nil, err
	}

	return userFromDB, nil
}

func (uu *userUsecase) UpdatePassword(u *model.User) error {
	if err := uu.ur.UpdatePassword(u); err != nil {
		return err
	}

	return nil
}

// CreateProfile creates profile for model.User and stores it in DB with repository.ProfileRepository
func (uu *userUsecase) CreateProfile(p *model.Profile) (interface{}, error) {
	insertedID, err := uu.pr.Create(p)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

// UpdateProfile updates profile of model.User in DB with repository.ProfileRepository
func (uu *userUsecase) UpdateProfile(p *model.Profile) error {
	err := uu.pr.Update(p)
	if err != nil {
		return err
	}

	if p.Nickname != p.AuthUsername {
		var u model.User

		u.Username = p.Nickname

		if err = uu.ur.UpdateUsername(&u, p.AuthUsername); err != nil {
			return err
		}

		log.Error().Err(err).
			Str("auth_username", p.AuthUsername).
			Str("new_username", p.Nickname).
			Msg("error updating username")
	}

	return nil
}
