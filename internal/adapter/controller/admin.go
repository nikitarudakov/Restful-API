package controller

import (
	"crypto/subtle"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strconv"
)

type AdminController struct {
	adminRepository repository.AdminRepository
}

func NewAdminController(ar repository.AdminRepository) *AdminController {
	return &AdminController{adminRepository: ar}
}

func (ac *AdminController) InitRoutes(e *echo.Echo, adminCfg *config.Admin) {
	admin := e.Group("/admin")

	admin.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		return ac.Auth(username, password, adminCfg)
	}))

	admin.GET("/users/profiles", func(ctx echo.Context) error {
		return ac.GetUserProfiles(ctx)
	})
}

func (ac *AdminController) Auth(username, password string, adminCfg *config.Admin) (bool, error) {
	if subtle.ConstantTimeCompare([]byte(username), []byte(adminCfg.Username)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(adminCfg.Password)) == 1 {

		return true, nil
	}

	return false, nil
}

func (ac *AdminController) GetUserProfiles(ctx echo.Context) error {
	profile := model.Profile{}
	profileCollName := profile.TableName()

	var page int64 = 1

	pageStr := ctx.QueryParam("page")
	if pageStr != "" {
		parsedPage, err := strconv.ParseInt(ctx.QueryParam("page"), 10, 64)
		if err != nil {
			return ctx.String(http.StatusBadRequest, err.Error())
		}

		page = parsedPage
	}

	result, err := ac.adminRepository.FindUserProfiles(profileCollName, page)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, result)
}
