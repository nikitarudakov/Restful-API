package main

import (
	"context"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/grpc/auth"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/grpc/vote"
	"git.foxminded.ua/foxstudent106092/user-management/logger"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
)

func main() {
	cfg, err := config.InitConfig(".config")
	if err != nil {
		panic(err)
	}

	logger.InitLogger(&cfg.Logger)

	service, err := grpc.NewUserManagementGRPCService(cfg.Grpc.Server + ":" + cfg.Grpc.Port)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}

	login, err := service.AuthClient.Login(context.Background(),
		&auth.Credentials{Username: "user2", Password: "5jVZZ6VN%"})
	if err != nil {
		log.Fatal().Err(err).Msg("login")
	}

	ctx := context.Background()
	md := metadata.Pairs("authorization", login.Token)
	ctx = metadata.NewOutgoingContext(ctx, md)
	ctx = metadata.AppendToOutgoingContext(ctx, "target", "user2")
	ctx = metadata.AppendToOutgoingContext(ctx, "cache", "GET")

	rating, err := service.VoteClient.GetRating(ctx, &vote.Target{Target: "user2"})
	if err != nil {
		log.Warn().Err(err).Send()
	}

	fmt.Printf("%+v\n", rating)
}
