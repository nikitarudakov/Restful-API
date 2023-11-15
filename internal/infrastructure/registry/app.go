package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Registry struct {
	Ur repository.UserRepoController
	Pr repository.ProfileRepoController
	Vr repository.VoteRepoController
}

func NewRegistry(db *mongo.Database, dbCfg *config.Database) *Registry {
	return &Registry{
		Ur: repository.NewUserRepository(db.Collection(dbCfg.UserRepo)),
		Pr: repository.NewProfileRepository(db.Collection(dbCfg.ProfileRepo)),
		Vr: repository.NewVoteRepository(db.Collection(dbCfg.VoteRepo)),
	}
}
