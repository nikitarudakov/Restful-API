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

	r := registry.NewRegistry(db, &cfg.Database, cacheDB)

	uu := usecase.NewUserUsecase(r.Ur, r.Pr)
	vu := usecase.NewVoteUsecase(r.Pr, r.Vr)

	authController := controller.NewAuthController(uu, cfg)
	authController.InitRoutes(e)

	voteController := controller.NewVoteController(vu, authController)
	voteController.InitRoutes(e)

	userController := controller.NewUserController(uu, authController)
	userController.InitRoutes(e)

	adminController := controller.NewAdminController(r.Ur, r.Pr, &cfg.Admin, userController, authController)
	adminController.InitRoutes(e)

	return e
}
