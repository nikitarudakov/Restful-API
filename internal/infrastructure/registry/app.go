package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Registry struct {
	Ur repository.UserRepoManager
	Pr repository.ProfileRepoManager
	Vr repository.VoteRepoManager
}

func NewRegistry(db *mongo.Database, dbCfg *config.Database, cache *cache.Database) *Registry {
	return &Registry{
		Ur: repository.NewUserRepository(db.Collection(dbCfg.UserRepo)),
		Pr: repository.NewProfileRepository(db.Collection(dbCfg.ProfileRepo), cache),
		Vr: repository.NewVoteRepository(db.Collection(dbCfg.VoteRepo), cache),
	}
}
