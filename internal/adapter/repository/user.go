package repository

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/gateway"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/repository"
	"github.com/rs/zerolog/log"
)

type userRepository struct {
	db gateway.Database
}

func NewUserRepository(db gateway.Database) repository.UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) Create(u *model.User) (*model.User, error) {
	_, err := ur.db.InsertUpdateItem(u)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	return u, err
}
