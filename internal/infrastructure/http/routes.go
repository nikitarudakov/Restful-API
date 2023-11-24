package http

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/controller"
	validator2 "git.foxminded.ua/foxstudent106092/user-management/tools/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutesWithControllers(r *registry.Registry, cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = &validator2.CustomValidator{
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	authController := controller.NewAuthController(r.UserUseCase, r.ProfileUseCase, cfg)
	authController.InitAuthRoutes(e, cfg)

	voteController := controller.NewVoteController(r.VoteUseCase, r.CacheDB)
	voteController.InitVoteRoutes(e, cfg)

	userController := controller.NewUserController(r.UserUseCase, r.ProfileUseCase, r.CacheDB)
	userController.InitUserRoutes(e, cfg)

	return e
}
