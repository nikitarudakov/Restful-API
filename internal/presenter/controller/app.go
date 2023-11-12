package controller

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"github.com/labstack/echo/v4"
)

type AdminEndpointHandler interface {
	InitAdminRoutes(e *echo.Echo, cacheDB *cache.Database, config *config.Config)
	GetUserProfiles(ctx echo.Context) error
	ModifyUserProfile(ctx echo.Context) error
	DeleteUserProfile(ctx echo.Context) error
}

type AuthEndpointHandler interface {
	InitAuthRoutes(e *echo.Echo)
	InitAuthMiddleware(g *echo.Group, accessibleRoles []string)
	UpdatePassword(ctx echo.Context) error
	Login(ctx echo.Context) error
	Register(ctx echo.Context) error
}

type UserEndpointsHandler interface {
	InitUserRoutes(e *echo.Echo)
	UpdateUserAndProfile(ctx echo.Context) error
	ListProfiles(ctx echo.Context) error
	DeleteUserAndProfile(ctx echo.Context) error
}

type VoteEndpointsHandler interface {
	InitVoteRoutes(e *echo.Echo, cacheDB *cache.Database, cfg *config.Config)
	Vote(ctx echo.Context) error
	RetractVote(ctx echo.Context) error
	GetRating(ctx echo.Context) error
}
