package controller

import (
	"crypto/subtle"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/auth"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/usecase"
	"git.foxminded.ua/foxstudent106092/user-management/tools/hashing"
	"github.com/labstack/echo/v4"
	"net/http"
)

type userController struct {
	userUsecase usecase.UserInput
}

type User interface {
	Register(ctx echo.Context) error
	Auth(username string, password string, ctx echo.Context) (bool, error)
	CreateProfile(ctx echo.Context) error
	UpdateProfile(ctx echo.Context) error
}

func NewUserController(uu usecase.UserInput) User {
	return &userController{userUsecase: uu}
}

func (uc *userController) Register(ctx echo.Context) error {
	var u model.User

	u.Username = hashing.HashUsername(ctx.FormValue("username"))
	u.Password = ctx.FormValue("password")

	if err := ctx.Validate(u); err != nil {
		return ctx.String(http.StatusForbidden, err.Error())
	}

	hashedPassword, err := hashing.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	err = uc.userUsecase.Create(&u)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, u)
}

func (uc *userController) Auth(username string, password string, ctx echo.Context) (bool, error) {
	var u model.User

	u.Username = hashing.HashUsername(username)
	u.Password = password

	userFromDB, err := uc.userUsecase.Find(&u)
	if err != nil {
		return false, fmt.Errorf("user was not found: %w", err)
	}

	if subtle.ConstantTimeCompare([]byte(u.Username), []byte(userFromDB.Username)) == 1 {
		err = hashing.CheckPassword(userFromDB.Password, password)
		if err != nil {
			return false, fmt.Errorf("username/password is incorrect: %w", err)
		}

		return true, nil
	}

	return false, fmt.Errorf("username/password is incorrect: %w", err)
}

func (uc *userController) AuthAndValidateProfile(ctx echo.Context) (*model.Profile, error) {
	a, err := auth.ReadAuthHeader(ctx.Request().Header)
	if err != nil {
		return nil, ctx.String(http.StatusInternalServerError, err.Error())
	}

	cred, err := auth.DecodeBasicAuthCred(a)
	if err != nil {
		return nil, ctx.String(http.StatusInternalServerError, err.Error())
	}

	var p model.Profile
	if err = ctx.Bind(&p); err != nil {
		return nil, ctx.String(http.StatusBadRequest, err.Error())
	}

	p.Nickname = hashing.HashUsername((*cred)[0])

	if err = ctx.Validate(p); err != nil {
		return nil, ctx.String(http.StatusBadRequest, err.Error())
	}

	return &p, nil
}

func (uc *userController) CreateProfile(ctx echo.Context) error {
	p, err := uc.AuthAndValidateProfile(ctx)
	if err != nil {
		return err
	}

	insertedID, err := uc.userUsecase.CreateProfile(p)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, struct {
		Id interface{} `json:"_id"`
	}{insertedID})
}

func (uc *userController) UpdateProfile(ctx echo.Context) error {
	p, err := uc.AuthAndValidateProfile(ctx)
	if err != nil {
		return err
	}

	err = uc.userUsecase.UpdateProfile(p)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
