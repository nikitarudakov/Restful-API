package presenter

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type userPresenter struct {
	ctx echo.Context
}

func NewUserPresenter(ctx echo.Context) usecase.UserOutput {
	return &userPresenter{ctx}
}

func (up *userPresenter) Render(i interface{}) error {
	return up.ctx.JSON(http.StatusOK, i)
}
