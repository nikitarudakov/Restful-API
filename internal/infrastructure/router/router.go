package router

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRouter(e *echo.Group, c controller.AppController) *echo.Group {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/users/register", func(ctx echo.Context) error { return c.User.RegisterUser(ctx) })

	return e
}

func NewRouter(e *echo.Echo, c controller.AppController) *echo.Echo {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.BasicAuth(basicAuth))

	e.GET("/users", func(ctx echo.Context) error { return nil })

	return e
}
