package controller

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	middleware2 "git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/middleware"
	"git.foxminded.ua/foxstudent106092/user-management/tools/hashing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type UserController struct {
	userUsecase UserManager
	AuthEndpointHandler
}

// UserManager contains methods for performing operations on User/Profile datatype
type UserManager interface {
	Create(u *model.User) error
	Find(u *model.User) (*model.User, error)
	UpdateUsername(p *model.Profile, authUsername string) error
	UpdatePassword(u *model.User) error
	CreateProfile(p *model.Profile, authUsername string) (interface{}, error)
	UpdateProfile(p *model.Profile, authUsername string) error
}

// NewUserController implicitly links  *UserController to userController
// Here to instantiate userController we provide usecase.UserManager
func NewUserController(um UserManager, ac AuthEndpointHandler) *UserController {
	return &UserController{um, ac}
}

func (uc *UserController) InitRoutes(e *echo.Echo) {
	userRouter := e.Group("/users")

	userRouter.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		err := middleware2.GetUserAuth(ctx)
		if err != nil {
			return false, err
		}

		return uc.Auth(username, password)
	}))

	userRouter.PUT("/password/update", func(ctx echo.Context) error {
		return uc.UpdatePassword(ctx)
	})

	userRouter.POST("/profiles/create", func(ctx echo.Context) error {
		return uc.CreateProfile(ctx)
	})

	userRouter.PUT("/profiles/update", func(ctx echo.Context) error {
		return uc.UpdateProfile(ctx)
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

// CreateProfile checks authentication, parses request data (params) to model.Profile
// and creates model.Profile in DB
func (uc *UserController) CreateProfile(ctx echo.Context) error {
	username := fmt.Sprintf("%v", ctx.Get("username"))

	p, err := uc.ParseUserProfileFromServerRequest(ctx, username)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	insertedID, err := uc.userUsecase.CreateProfile(p, username)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, struct {
		Id interface{} `json:"_id"`
	}{insertedID})
}

// UpdateProfile checks authentication, parses request data (params) to model.Profile
// and updates model.Profile in DB
func (uc *UserController) UpdateProfile(ctx echo.Context) error {
	username := fmt.Sprintf("%v", ctx.Get("username"))

	p, err := uc.ParseUserProfileFromServerRequest(ctx, username)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err = uc.userUsecase.UpdateProfile(p, username)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

// ParseUserProfileFromServerRequest parses server request data to model.Profile
func (uc *UserController) ParseUserProfileFromServerRequest(
	ctx echo.Context,
	username string) (*model.Profile, error) {

	var p model.Profile
	if err := ctx.Bind(&p); err != nil {
		return nil, err
	}

	// check if new Nickname was passed
	if p.Nickname == "" {
		p.Nickname = username // assign User username to Profile nickname
	}

	if err := ctx.Validate(p); err != nil {
		return nil, err
	}

	return &p, nil
}
