package http

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/controller"
	validator2 "git.foxminded.ua/foxstudent106092/user-management/tools/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutesWithControllers(cfg *config.Config) *echo.Echo {
	db, err := datastore.NewDB(&cfg.Database)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = &validator2.CustomValidator{
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	r := registry.NewRegistry(db, cfg)

	authController := controller.NewAuthController(r, cfg)
	authController.InitAuthRoutes(e)

	voteController := controller.NewVoteController(r)
	voteController.InitVoteRoutes(e, cfg)

	userController := controller.NewUserController(r)
	userController.InitUserRoutes(e, cfg)

	return e
}
