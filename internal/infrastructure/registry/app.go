package registry

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/usecase"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/profileDao"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/userDao"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/voteDao"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/repository"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type Registry struct {
	UserUseCase    *usecase.UserUseCase
	ProfileUseCase *usecase.ProfileUseCase
	VoteUseCase    *usecase.VoteUsecase
	CacheDB        *cache.Database
}

type RepoRegistry struct {
	UserRepository    *repository.UserRepository
	ProfileRepository *repository.ProfileRepository
	VoteRepository    *repository.VoteRepository
}

func NewRepoRegistry(db *mongo.Database, cfg *config.Config) *RepoRegistry {
	ur := repository.NewUserRepository(db.Collection(cfg.Database.UserRepo))
	pr := repository.NewProfileRepository(db.Collection(cfg.Database.ProfileRepo))
	vr := repository.NewVoteRepository(db.Collection(cfg.Database.VoteRepo))

	return &RepoRegistry{ur, pr, vr}
}

func NewRegistry(conn *grpc.ClientConn, cfg *config.Config) *Registry {
	profileClient := profileDao.NewProfileRepoClient(conn)
	userClient := userDao.NewUserRepoClient(conn)
	voteClient := voteDao.NewVoteRepoClient(conn)

	uu := usecase.NewUserUsecase(userClient)
	pu := usecase.NewProfileUseCase(profileClient)
	vu := usecase.NewVoteUsecase(voteClient)

	cacheDB, err := cache.NewCacheDatabase(&cfg.Cache)
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error connecting to Redis")
	}

	return &Registry{uu, pu,
		vu, cacheDB}
}
