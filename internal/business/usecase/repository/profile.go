package repository

import "git.foxminded.ua/foxstudent106092/user-management/internal/business/model"

type ProfileRepository interface {
	Create(p *model.Profile) (interface{}, error)
	Find(p *model.Profile) (*model.Profile, error)
	Update(p *model.Update, authUsername string) error
	Delete(authUsername string) error
	ListUserProfiles(page int64) ([]model.Profile, error)
}
