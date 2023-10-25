package router

import "github.com/labstack/echo/v4"

func basicAuth(username, password string, c echo.Context) (bool, error) {
	return true, nil
}
