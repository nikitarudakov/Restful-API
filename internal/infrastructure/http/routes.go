package http

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/usecase"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/controller"
	validator2 "git.foxminded.ua/foxstudent106092/user-management/tools/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func InitRoutesWithControllers(cfg *config.Config) *echo.Echo {
	db, err := datastore.NewDB(&cfg.Database)
	if err != nil {
		panic(err)
	}

	cacheDB, err := cache.NewCacheDatabase(&cfg.Cache)
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error connecting to Redis")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = &validator2.CustomValidator{
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	r := registry.NewRegistry(db, &cfg.Database)

	uu := usecase.NewUserUsecase(r.Ur)
	pu := usecase.NewProfileUseCase(r.Pr)
	vu := usecase.NewVoteUsecase(r.Pr, r.Vr)

	authController := controller.NewAuthController(uu, pu, cfg)
	authController.InitAuthRoutes(e)

	voteController := controller.NewVoteController(vu)
	voteController.InitVoteRoutes(e, cacheDB, cfg)

	userController := controller.NewUserController(uu, pu)
	userController.InitUserRoutes(e, cacheDB, cfg)

	return e
}
