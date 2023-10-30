package controller

import (
	"crypto/subtle"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/repository"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type adminController struct {
	adminRepository repository.AdminRepository
}

type AdminEndpointHandler interface {
	Auth(username, password string, ctx echo.Context) (bool, error)
	GetUserProfiles(ctx echo.Context) error
}

func NewAdminController(ar repository.AdminRepository) AdminEndpointHandler {
	return &adminController{adminRepository: ar}
}

func (a *adminController) Auth(username, password string, ctx echo.Context) (bool, error) {
	if subtle.ConstantTimeCompare([]byte(username), []byte(config.C.AdminAPI.Username)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(config.C.AdminAPI.Password)) == 1 {

		return true, nil
	}

	return false, nil
}

func (a *adminController) GetUserProfiles(ctx echo.Context) error {
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

	result, err := a.adminRepository.FindUserProfiles(profileCollName, page)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, result)
}
