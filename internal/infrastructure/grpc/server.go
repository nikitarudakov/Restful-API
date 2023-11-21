package grpc

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/profileDao"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/userDao"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/voteDao"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
)

func newDAOServer(pr profileDao.ProfileRepoController,
	ur userDao.UserRepoController, vr voteDao.VoteRepoController) (*grpc.Server, error) {

	profileDAO := profileDao.Server{ProfileRepo: pr}
	userDAO := userDao.Server{UserRepo: ur}
	voteDAO := voteDao.Server{VoteRepo: vr}

	grpcServer := grpc.NewServer()

	profileDao.RegisterProfileRepoServer(grpcServer, &profileDAO)
	userDao.RegisterUserRepoServer(grpcServer, &userDAO)
	voteDao.RegisterVoteRepoServer(grpcServer, &voteDAO)

	return grpcServer, nil
}

func StartDAOServer(r *registry.RepoRegistry, cfg *config.Config) {
	lis, err := net.Listen("tcp", ":"+cfg.Dao.Port)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	daoServer, err := newDAOServer(r.ProfileRepository, r.UserRepository, r.VoteRepository)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().Msg("Server is running at http://localhost" + ":" + cfg.Dao.Port)
	if err = daoServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Send()
	}
}
