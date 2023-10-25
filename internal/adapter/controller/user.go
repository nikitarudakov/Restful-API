package controller

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/hashing"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/repository"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type userController struct {
	userRepository repository.UserRepository
}

type User interface {
	RegisterUser(ctx echo.Context) error
}

func NewUserController(ur repository.UserRepository) User {
	return &userController{userRepository: ur}
}

func (uc *userController) RegisterUser(ctx echo.Context) error {
	var params model.User

	params.Username = ctx.FormValue("username")

	hashedPassword, err := hashing.HashPassword(ctx.FormValue("password"))
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	params.HashedPassword = hashedPassword

	/*	if err := ctx.Bind(&params); err != nil {
		return err
	}*/

	log.Trace().Str("user", fmt.Sprintf("%+v\n", params)).Send()

	u, err := uc.userRepository.Create(&params)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, u)
}
