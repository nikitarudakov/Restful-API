package repository

import "git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"

// AdminRepository provide methods for managing entities for an entire api service
type AdminRepository interface {
	FindUserProfiles(colName string, page int64) ([]model.Profile, error)
}
