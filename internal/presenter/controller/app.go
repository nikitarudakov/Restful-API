package controller

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"github.com/labstack/echo/v4"
)

type AuthEndpointHandler interface {
	InitAuthRoutes(e *echo.Echo)
	UpdatePassword(ctx echo.Context) error
	Login(ctx echo.Context) error
	Register(ctx echo.Context) error
}

type UserEndpointsHandler interface {
	InitUserRoutes(e *echo.Echo, cfg *config.Config)
	UpdateUserAndProfile(ctx echo.Context) error
	ListProfiles(ctx echo.Context) error
	DeleteUserAndProfile(ctx echo.Context) error
}

type VoteEndpointsHandler interface {
	InitVoteRoutes(e *echo.Echo, cfg *config.Config)
	Vote(ctx echo.Context) error
	RetractVote(ctx echo.Context) error
	GetRating(ctx echo.Context) error
}
