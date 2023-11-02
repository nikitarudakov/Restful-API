package usecase

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	repository2 "git.foxminded.ua/foxstudent106092/user-management/internal/business/usecase/repository"
	"github.com/rs/zerolog/log"
)

type UserUsecase struct {
	ur repository2.UserRepository
	pr repository2.ProfileRepository
}

// NewUserUsecase implicitly links UserManager to userUsecase struct
func NewUserUsecase(ur repository2.UserRepository, pr repository2.ProfileRepository) *UserUsecase {
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

func (uu *UserUsecase) UpdateUsername(p *model.Profile, authUsername string) error {
	if p.Nickname != authUsername {
		var u model.User

		u.Username = p.Nickname
		if err := uu.ur.UpdateUsername(&u, authUsername); err != nil {
			return err
		}
	}

	return nil
}

func (uu *UserUsecase) UpdatePassword(u *model.User) error {
	if err := uu.ur.UpdatePassword(u); err != nil {
		return err
	}

	return nil
}

// CreateProfile creates profile for model.User and stores it in DB with repository.ProfileRepository
func (uu *UserUsecase) CreateProfile(p *model.Profile, authUsername string) (interface{}, error) {
	if err := uu.UpdateUsername(p, authUsername); err != nil {
		log.Error().Err(err).
			Str("auth_username", authUsername).
			Str("new_username", p.Nickname).
			Msg("error updating username")
	}

	insertedID, err := uu.pr.Create(p)
	if err != nil {
		return nil, err
	}

	return insertedID, nil
}

// UpdateProfile updates profile of model.User in DB with repository.ProfileRepository
func (uu *UserUsecase) UpdateProfile(p *model.Profile, authUsername string) error {
	if err := uu.UpdateUsername(p, authUsername); err != nil {
		log.Error().Err(err).
			Str("auth_username", authUsername).
			Str("new_username", p.Nickname).
			Msg("error updating username")
	}

	err := uu.pr.Update(p, authUsername)
	if err != nil {
		return err
	}

	return nil
}
