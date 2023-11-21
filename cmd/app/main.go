package main

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/http"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	"git.foxminded.ua/foxstudent106092/user-management/logger"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.InitConfig(".config")
	if err != nil {
		panic(err)
	}

	logger.InitLogger(&cfg.Logger)

	conn, err := grpc.Dial(cfg.Dao.Server+":"+cfg.Dao.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Warn().Err(err).Send()
	}

	r := registry.NewRegistry(conn, cfg)

	e := http.InitRoutesWithControllers(r, cfg)

	log.Info().Msg("Server is running at http://localhost" + ":" + cfg.Server.Port)
	if err = e.Start(":" + cfg.Server.Port); err != nil {
		log.Fatal().Err(err).Send()
	}
}
