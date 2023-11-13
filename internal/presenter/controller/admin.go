package controller

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AdminController struct {
	userController UserEndpointsHandler
	AuthEndpointHandler
	cfgAdmin *config.Admin
}

func NewAdminController(userHandler UserEndpointsHandler,
	authHandler AuthEndpointHandler, cfgAdmin *config.Admin) *AdminController {
	return &AdminController{userHandler, authHandler, cfgAdmin}
}

func (ac *AdminController) InitAdminRoutes(e *echo.Echo, cacheDB *cache.Database, cfg *config.Config) {
	admin := e.Group("/admin")

	roles := []string{"admin", "moderator"}

	var profiles []model.Profile
	admin.Use(cache.Middleware(cacheDB, &profiles, &cfg.Cache))

	ac.InitAuthMiddleware(admin, roles)

	admin.GET("/users/profiles", func(ctx echo.Context) error {
		return ac.GetUserProfiles(ctx)
	})

	admin.PUT("/users/profiles/:username/modify", func(ctx echo.Context) error {
		return ac.ModifyUserProfile(ctx)
	})

	admin.DELETE("/users/profiles/:username/delete", func(ctx echo.Context) error {
		return ac.DeleteUserProfile(ctx)
	})
}

func (ac *AdminController) GetUserProfiles(ctx echo.Context) error {
	return ac.userController.ListProfiles(ctx)
}

func (ac *AdminController) ModifyUserProfile(ctx echo.Context) error {
	err := ac.userController.UpdateUserAndProfile(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (ac *AdminController) DeleteUserProfile(ctx echo.Context) error {
	return ac.userController.DeleteUserAndProfile(ctx)
}
