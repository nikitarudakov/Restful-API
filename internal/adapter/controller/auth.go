package controller

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/usecase"
	"git.foxminded.ua/foxstudent106092/user-management/tools/hashing"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthController struct {
	userUsecase usecase.UserManager
}

func NewAuthController(uu usecase.UserManager) *AuthController {
	return &AuthController{userUsecase: uu}
}

func (ac *AuthController) InitRoutes(e *echo.Echo) {
	regRouter := e.Group("/auth")

	regRouter.POST("/register", func(ctx echo.Context) error {
		return ac.Register(ctx)
	})
}

// Register creates and stores new model.User in DB
func (ac *AuthController) Register(ctx echo.Context) error {
	var u model.User

	u.Username = ctx.FormValue("username")
	u.Password = ctx.FormValue("password")

	if err := ctx.Validate(u); err != nil {
		return ctx.String(http.StatusForbidden, err.Error())
	}

	hashedPassword, err := hashing.HashPassword(u.Password)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	u.Password = hashedPassword

	err = ac.userUsecase.Create(&u)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, u)
}
