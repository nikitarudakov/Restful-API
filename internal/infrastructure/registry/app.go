package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/usecase"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/repository"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/controller"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type Registry struct {
	UserUseCase    controller.UserManager
	UserRepo       repository.UserRepoController
	ProfileUseCase controller.ProfileManager
	ProfileRepo    repository.ProfileRepoController
	VoteUseCase    controller.VoteManager
	VoteRepo       repository.VoteRepoController
	CacheDB        *cache.Database
}

func NewRegistry(db *mongo.Database, cfg *config.Config) *Registry {
	cacheDB, err := cache.NewCacheDatabase(&cfg.Cache)
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error connecting to Redis")
	}

	ur := repository.NewUserRepository(db.Collection(cfg.Database.UserRepo))
	pr := repository.NewProfileRepository(db.Collection(cfg.Database.ProfileRepo))
	vr := repository.NewVoteRepository(db.Collection(cfg.Database.VoteRepo))

	uu := usecase.NewUserUsecase(ur)
	pu := usecase.NewProfileUseCase(pr)
	vu := usecase.NewVoteUsecase(pr, vr)

	return &Registry{
		uu, ur,
		pu, pr,
		vu, vr, cacheDB,
	}
}
