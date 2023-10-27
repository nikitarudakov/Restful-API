package router

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/registry"
	customValidator "git.foxminded.ua/foxstudent106092/user-management/tools/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetRouters(r registry.Registry) *echo.Echo {
	regRouter := echo.New()
	regRouter.Use(middleware.Logger())
	regRouter.Use(middleware.Recover())

	regRouter.Validator = &customValidator.CustomValidator{
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	regRouter.POST("/register", func(ctx echo.Context) error {
		return r.NewAppController().User.Register(ctx)
	})

	mainRouter := regRouter.Group("/users")

	mainRouter.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		return r.NewAppController().User.Auth(username, password, ctx)
	}))

	mainRouter.POST("/profiles/create", func(ctx echo.Context) error {
		return r.NewAppController().User.CreateProfile(ctx)
	})

	mainRouter.PUT("/profiles/update", func(ctx echo.Context) error {
		return r.NewAppController().User.UpdateProfile(ctx)
	})

	return regRouter
}
