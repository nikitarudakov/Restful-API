package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	repoAdapter "git.foxminded.ua/foxstudent106092/user-management/internal/adapter/repository"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Registry struct {
	Ur repository.UserRepository
	Pr repository.ProfileRepository
	Ar repository.AdminRepository
}

func NewRegistry(db *mongo.Database, dbCfg *config.Database) *Registry {
	return &Registry{
		Ur: repoAdapter.NewUserRepository(db.Collection(dbCfg.UserRepo)),
		Pr: repoAdapter.NewProfileRepository(db.Collection(dbCfg.ProfileRepo)),
		Ar: repoAdapter.NewAdminRepository(db),
	}
}
