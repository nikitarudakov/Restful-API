package controller

import (
	"crypto/subtle"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/usecase/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type AdminController struct {
	adminRepository repository.AdminRepository
	cfgAdmin        *config.Admin
	AuthEndpointHandler
}

func NewAdminController(ar repository.AdminRepository,
	cfgAdmin *config.Admin, ac AuthEndpointHandler) *AdminController {
	return &AdminController{
		ar,
		cfgAdmin,
		ac,
	}
}

func (ac *AdminController) InitRoutes(e *echo.Echo) {
	admin := e.Group("/admin")

	roles := []string{"admin"}

	ac.InitAuthMiddleware(admin, roles)

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
