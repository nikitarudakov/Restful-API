package auth

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitAuthMiddleware(g *echo.Group, authCfg *config.Auth, accessibleRoles []string) {
	tokenConfig := getTokenConfig(authCfg, accessibleRoles)
	g.Use(echojwt.WithConfig(tokenConfig))
}
