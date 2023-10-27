package repository

import "git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"

type ProfileRepository interface {
	Create(p *model.Profile) (interface{}, error)
	Find(p *model.Profile) (*model.Profile, error)
	Update(p *model.Profile) error
}
