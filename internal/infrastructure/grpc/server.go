package grpc

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/interceptor"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/grpc/auth"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/grpc/user"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/grpc/vote"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
)

func withServerUnaryInterceptor(cfg *config.Config, db *cache.Database) grpc.ServerOption {
	authInterceptor := interceptor.GetServerAuthInterceptor(&cfg.Auth)
	cacheInterceptor := interceptor.GetServerCacheInterceptor(&cfg.Cache, db)
	return grpc.ChainUnaryInterceptor(authInterceptor, cacheInterceptor)
}

func ServeGRPCServer(r *registry.Registry, cfg *config.Config) {
	lis, err := net.Listen("tcp", ":"+cfg.Grpc.Port)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	grpcServer := grpc.NewServer(
		withServerUnaryInterceptor(cfg, r.CacheDB),
	)

	authServer := auth.GRPCAuthServer{UserUseCase: r.UserUseCase,
		ProfileUseCase: r.ProfileUseCase, Cfg: cfg}

	userServer := user.GRPCUserService{UserUseCase: r.UserUseCase,
		ProfileUseCase: r.ProfileUseCase, CacheDb: r.CacheDB}

	voteServer := vote.GRPCVoteService{VoteUsecase: r.VoteUseCase,
		CacheDB: r.CacheDB}

	auth.RegisterAuthControllerServer(grpcServer, &authServer)
	user.RegisterUserControllerServer(grpcServer, &userServer)
	vote.RegisterVoteControllerServer(grpcServer, &voteServer)

	log.Info().Msg("Server is running at http://localhost" + ":" + cfg.Grpc.Port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Send()
	}
}
