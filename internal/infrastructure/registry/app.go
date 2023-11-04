package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/usecase/repository"
	repoAdapter "git.foxminded.ua/foxstudent106092/user-management/internal/presenter/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Registry struct {
	Ur repository.UserRepository
	Pr repository.ProfileRepository
}

func NewRegistry(db *mongo.Database, dbCfg *config.Database) *Registry {
	return &Registry{
		Ur: repoAdapter.NewUserRepository(db.Collection(dbCfg.UserRepo)),
		Pr: repoAdapter.NewProfileRepository(db.Collection(dbCfg.ProfileRepo)),
	}
}
