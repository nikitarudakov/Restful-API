package controller

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"github.com/labstack/echo/v4"
)

type AdminEndpointHandler interface {
	InitRoutes(e *echo.Echo)
	Auth(username, password string, adminCfg *config.Admin) (bool, error)
	GetUserProfiles(ctx echo.Context) error
	ModifyUserProfile(ctx echo.Context) error
	DeleteUserProfile(ctx echo.Context) error
}

type AuthEndpointHandler interface {
	InitRoutes(e *echo.Echo)
	InitAuthMiddleware(g *echo.Group, accessibleRoles []string)
	Login(ctx echo.Context) error
	Auth(username string, password string) (bool, error)
	Register(ctx echo.Context) error
}

type UserEndpointsHandler interface {
	InitRoutes(e *echo.Echo)
	Auth(username string, password string) (bool, error)
	UpdatePassword(ctx echo.Context) error
	UpdateUserProfile(ctx echo.Context) error
}
