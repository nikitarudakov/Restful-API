package usecase

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/repository"
	"github.com/rs/zerolog/log"
)

type UserUsecase struct {
	ur repository.UserRepository
	pr repository.ProfileRepository
}

// UserManager contains methods for performing operations on User/Profile datatype
type UserManager interface {
	Create(u *model.User) error
	Find(u *model.User) (*model.User, error)
	UpdatePassword(u *model.User) error
	CreateProfile(p *model.Profile, authUsername string) (interface{}, error)
	UpdateProfile(p *model.Profile, authUsername string) error
}

// NewUserUsecase implicitly links UserManager to userUsecase struct
func NewUserUsecase(ur repository.UserRepository, pr repository.ProfileRepository) *UserUsecase {
	return &UserUsecase{ur: ur, pr: pr}
}

// Create creates new user and stores it in DB with repository.UserRepository
func (uu *UserUsecase) Create(u *model.User) error {
	_, err := uu.ur.Create(u)
	if err != nil {
		return err
	}

	return nil
}

// Find looks up for user (u *model.User) and returns it if found
func (uu *UserUsecase) Find(u *model.User) (*model.User, error) {
	userFromDB, err := uu.ur.Find(u)
	if err != nil {
		return nil, err
	}

	return userFromDB, nil
}

func (uu *UserUsecase) UpdatePassword(u *model.User) error {
	if err := uu.ur.UpdatePassword(u); err != nil {
		return err
	}

	return nil
}

// CreateProfile creates profile for model.User and stores it in DB with repository.ProfileRepository
func (uu *UserUsecase) CreateProfile(p *model.Profile, authUsername string) (interface{}, error) {
	if p.Nickname != authUsername {
		var u model.User

		u.Username = p.Nickname

		if err := uu.ur.UpdateUsername(&u, authUsername); err != nil {
			log.Error().Err(err).
				Str("auth_username", authUsername).
				Str("new_username", p.Nickname).
				Msg("error updating username")

			return nil, err
		}
	}

	insertedID, err := uu.pr.Create(p)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

// UpdateProfile updates profile of model.User in DB with repository.ProfileRepository
func (uu *UserUsecase) UpdateProfile(p *model.Profile, authUsername string) error {
	if p.Nickname != authUsername {
		var u model.User

		u.Username = p.Nickname

		if err := uu.ur.UpdateUsername(&u, authUsername); err != nil {
			log.Error().Err(err).
				Str("auth_username", authUsername).
				Str("new_username", p.Nickname).
				Msg("error updating username")

			return err
		}
	}

	err := uu.pr.Update(p)
	if err != nil {
		return err
	}

	return nil
}
