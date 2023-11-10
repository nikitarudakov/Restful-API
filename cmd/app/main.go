package main

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/http"
	"git.foxminded.ua/foxstudent106092/user-management/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.InitConfig(".config")
	if err != nil {
		panic(err)
	}

	logger.InitLogger(&cfg.Logger)

	e := http.InitRoutesWithControllers(cfg)

	log.Info().Msg("Server is running at http://localhost" + ":" + cfg.Server.Port)
	if err = e.Start(":" + cfg.Server.Port); err != nil {
		log.Fatal().Err(err).Send()
	}
}
