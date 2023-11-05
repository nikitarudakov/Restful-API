package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Registry struct {
	Ur repository.UserRepoManager
	Pr repository.ProfileRepoManager
}

func NewRegistry(db *mongo.Database, dbCfg *config.Database) *Registry {
	return &Registry{
		Ur: repository.NewUserRepository(db.Collection(dbCfg.UserRepo)),
		Pr: repository.NewProfileRepository(db.Collection(dbCfg.ProfileRepo)),
	}
}
