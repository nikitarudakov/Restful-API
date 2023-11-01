package controller

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/registry"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/usecase"
	"github.com/labstack/echo/v4"
)

type AdminEndpointHandler interface {
	InitRoutes(e *echo.Echo, adminCfg *config.Admin)
	Auth(username, password string, adminCfg *config.Admin) (bool, error)
	GetUserProfiles(ctx echo.Context) error
}

type AuthEndpointHandler interface {
	InitRoutes(e *echo.Echo)
	Register(ctx echo.Context) error
}

type UserEndpointsHandler interface {
	InitRoutes(e *echo.Echo)
	Auth(username string, password string) (bool, error)
	UpdatePassword(ctx echo.Context) error
	CreateProfile(ctx echo.Context) error
	UpdateProfile(ctx echo.Context) error
}

type AppController struct {
	User  UserEndpointsHandler
	Admin AdminEndpointHandler
	Auth  AuthEndpointHandler
}

func NewAppController(r *registry.Registry) AppController {
	uu := usecase.NewUserUsecase(r.Ur, r.Pr)

	return AppController{
		User:  NewUserController(uu),
		Admin: NewAdminController(r.Ar),
		Auth:  NewAuthController(uu),
	}
}
