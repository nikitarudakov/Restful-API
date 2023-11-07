package http

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/usecase"
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

	r := registry.NewRegistry(db, &cfg.Database)

	uu := usecase.NewUserUsecase(r.Ur, r.Pr)

	authController := controller.NewAuthController(uu, cfg)
	authController.InitRoutes(e)

	userController := controller.NewUserController(uu, authController)
	userController.InitRoutes(e)

	adminController := controller.NewAdminController(r.Ur, r.Pr, &cfg.Admin, userController, authController)
	adminController.InitRoutes(e)

	return e
}
