package repository

import "git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"

type UserRepository interface {
	Find(u *model.User) (*model.User, error)
	Create(u *model.User) (interface{}, error)
	UpdateUsername(newUser *model.User, oldVal string) error
	UpdatePassword(newUser *model.User) error
}
