package controller

import (
	"git.foxminded.ua/foxstudent106092/user-management/hashing"
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/presenter"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/usecase"
	"github.com/labstack/echo/v4"
)

type userController struct {
	userUsecase usecase.UserInput
}

type User interface {
	Register(ctx echo.Context) error
}

func NewUserController(uu usecase.UserInput) User {
	return &userController{userUsecase: uu}
}

func (uc *userController) Register(ctx echo.Context) error {
	var params model.User

	authCredentials := ctx.FormValue("username") + ":" + ctx.FormValue("password")

	hashedAuthCredentials, err := hashing.HashAuthCredentials(authCredentials)
	if err != nil {
		return err
	}

	params.HashID = hashedAuthCredentials

	err = uc.userUsecase.Create(&params)
	if err != nil {
		return err
	}

	outputPort := presenter.NewUserPresenter(ctx)

	return outputPort.Render(params)
}
