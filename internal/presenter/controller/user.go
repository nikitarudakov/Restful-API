package controller

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/tools/hashing"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type UserController struct {
	userUsecase UserManager
	AuthEndpointHandler
}

// UserManager contains methods for performing operations on User/Profile datatype
type UserManager interface {
	CreateUser(u *model.User) error
	Find(u *model.User) (*model.User, error)
	UpdateUsername(newUsername, oldUsername string) error
	UpdatePassword(u *model.User) error
	CreateProfile(p *model.Profile) (interface{}, error)
	UpdateProfile(p *model.Update, authUsername string) error
}

// NewUserController implicitly links  *UserController to userController
// Here to instantiate userController we provide usecase.UserManager
func NewUserController(um UserManager, ac AuthEndpointHandler) *UserController {
	return &UserController{um, ac}
}

func (uc *UserController) InitRoutes(e *echo.Echo) {
	userRouter := e.Group("/users")

	roles := []string{"admin", "user", "moderator"}

	uc.InitAuthMiddleware(userRouter, roles)

	userRouter.PUT("/password/update", func(ctx echo.Context) error {
		return uc.UpdatePassword(ctx)
	})

	userRouter.PUT("/profiles/:username/update", func(ctx echo.Context) error {
		return uc.UpdateUserProfile(ctx)
	})
}

func (uc *UserController) UpdatePassword(ctx echo.Context) error {
	var u model.User

	u.Username = fmt.Sprintf("%v", ctx.Get("username"))

	password, err := hashing.HashPassword(ctx.FormValue("password"))
	u.Password = password
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	if err = uc.userUsecase.UpdatePassword(&u); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, fmt.Sprintf("Password was updated!"))
}

// UpdateUserProfile checks authentication, parses request data (params) to model.Profile
// and updates model.Profile in DB
func (uc *UserController) UpdateUserProfile(ctx echo.Context) error {
	username := ctx.Param("username")

	update, err := uc.parseUserProfileUpdate(ctx, username)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err = uc.userUsecase.UpdateProfile(update, username)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	err = uc.userUsecase.UpdateUsername(update.Nickname, username)
	if err != nil {
		log.Error().Err(err).
			Str("old_username", username).
			Str("new_username", update.Nickname).
			Msg("error updating username")
	}

	return ctx.JSON(http.StatusOK, nil)
}

// ParseUserProfileFromServerRequest parses server request data to model.Profile
func (uc *UserController) parseUserProfileUpdate(
	ctx echo.Context,
	username string) (*model.Update, error) {

	var update model.Update
	if err := ctx.Bind(&update); err != nil {
		return nil, err
	}

	// check if new Nickname was passed
	if update.Nickname == "" {
		update.Nickname = username // assign User username to Update nickname
	}

	return &update, nil
}
