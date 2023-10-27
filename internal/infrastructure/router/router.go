package router

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/registry"
	customValidator "git.foxminded.ua/foxstudent106092/user-management/tools/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetUserRouter(r registry.Registry) *echo.Echo {
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	router.Validator = &customValidator.CustomValidator{
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	regRouter := router.Group("/auth")

	regRouter.POST("/users/register", func(ctx echo.Context) error {
		return r.NewAppController().User.Register(ctx)
	})

	userRouter := router.Group("/users")

	userRouter.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		return r.NewAppController().User.Auth(username, password, ctx)
	}))

	userRouter.POST("/profiles/create", func(ctx echo.Context) error {
		return r.NewAppController().User.CreateProfile(ctx)
	})

	userRouter.PUT("/profiles/update", func(ctx echo.Context) error {
		return r.NewAppController().User.UpdateProfile(ctx)
	})

	return router
}

func SetAdminGroupRouter(r registry.Registry, parentRouter *echo.Echo) {
	admin := parentRouter.Group("/admin")

	admin.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		return r.NewAppController().Admin.Auth(username, password, ctx)
	}))

	admin.GET("/users/profiles", func(ctx echo.Context) error {
		return r.NewAppController().Admin.GetUserProfiles(ctx)
	})
}
