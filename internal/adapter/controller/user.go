package controller

import (
	"crypto/subtle"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/domain/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/auth"
	"git.foxminded.ua/foxstudent106092/user-management/internal/usecase/usecase"
	"git.foxminded.ua/foxstudent106092/user-management/tools/hashing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type UserController struct {
	userUsecase usecase.UserManager
}

// NewUserController implicitly links  *UserController to userController
// Here to instantiate userController we provide usecase.UserManager
func NewUserController(um usecase.UserManager) *UserController {
	return &UserController{userUsecase: um}
}

func (uc *UserController) InitRoutes(e *echo.Echo) {
	userRouter := e.Group("/users")

	userRouter.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
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

// Auth is authentication handler for BasicAuth middleware
// It hashes credentials and compares them using subtle.ConstantTimeCompare
// to prevent time attacks. If matches returns true which means
// user was successfully authenticated and BasicAuth header was added
func (uc *UserController) Auth(username string, password string) (bool, error) {
	var u model.User

	u.Username = username

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

func (uc *UserController) UpdatePassword(ctx echo.Context) error {
	cred, err := uc.CheckUserAuth(ctx)
	if err != nil {
		return ctx.String(http.StatusForbidden, err.Error())
	}

	var u model.User

	u.Username = (*cred)[0]
	u.Password, err = hashing.HashPassword(ctx.FormValue("password"))
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
	cred, err := uc.CheckUserAuth(ctx)
	if err != nil {
		return ctx.String(http.StatusForbidden, err.Error())
	}

	p, err := uc.ParseUserProfileFromServerRequest(ctx, cred)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	insertedID, err := uc.userUsecase.CreateProfile(p, (*cred)[0])
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
	cred, err := uc.CheckUserAuth(ctx)
	if err != nil {
		return ctx.String(http.StatusForbidden, err.Error())
	}

	p, err := uc.ParseUserProfileFromServerRequest(ctx, cred)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err = uc.userUsecase.UpdateProfile(p, (*cred)[0])
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

// CheckUserAuth checks authentication of model.User
func (uc *UserController) CheckUserAuth(ctx echo.Context) (*[]string, error) {
	a, err := auth.ReadAuthHeader(ctx.Request().Header)
	if err != nil {
		return nil, err
	}

	cred, err := auth.DecodeBasicAuthCred(a)
	if err != nil {
		return nil, err
	}

	return cred, nil
}

// ParseUserProfileFromServerRequest parses server request data to model.Profile
func (uc *UserController) ParseUserProfileFromServerRequest(
	ctx echo.Context,
	cred *[]string) (*model.Profile, error) {

	var p model.Profile
	if err := ctx.Bind(&p); err != nil {
		return nil, err
	}

	username := (*cred)[0]

	// check if new Nickname was passed
	if p.Nickname == "" {
		p.Nickname = username // assign User username to Profile nickname
	}

	if err := ctx.Validate(p); err != nil {
		return nil, err
	}

	return &p, nil
}
