package repository

import "git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"

type UserRepository interface {
	Create(u *model.User) (interface{}, error)
}
