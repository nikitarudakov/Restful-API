package controller

import (
	"github.com/labstack/echo/v4"
)

type AdminEndpointHandler interface {
	InitRoutes(e *echo.Echo)
	GetUserProfiles(ctx echo.Context) error
	ModifyUserProfile(ctx echo.Context) error
	DeleteUserProfile(ctx echo.Context) error
	UserEndpointsHandler
}

type AuthEndpointHandler interface {
	InitRoutes(e *echo.Echo)
	InitAuthMiddleware(g *echo.Group, accessibleRoles []string)
	UpdatePassword(ctx echo.Context) error
	Login(ctx echo.Context) error
	Register(ctx echo.Context) error
}

type UserEndpointsHandler interface {
	InitRoutes(e *echo.Echo)
	UpdateUserProfile(ctx echo.Context) error
}

type VoteEndpointsHandler interface {
	InitRoutes(g *echo.Group)
	Vote(ctx echo.Context) error
	RetractVote(ctx echo.Context) error
	GetRating(ctx echo.Context) error
}
