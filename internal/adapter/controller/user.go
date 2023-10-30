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
	userUsecase usecase.UserManager
}

// UserEndpointsHandler stores methods (handlers) for processing requests from the server
type UserEndpointsHandler interface {
	Register(ctx echo.Context) error
	Auth(username string, password string, ctx echo.Context) (bool, error)
	UpdatePassword(ctx echo.Context) error
	CreateProfile(ctx echo.Context) error
	UpdateProfile(ctx echo.Context) error
}

// NewUserController implicitly links UserEndpointsHandler to userController
// Here to instantiate userController we provide usecase.UserManager
func NewUserController(um usecase.UserManager) UserEndpointsHandler {
	return &userController{userUsecase: um}
}

// Register creates and stores new model.User in DB
func (uc *userController) Register(ctx echo.Context) error {
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

	err = uc.userUsecase.Create(&u)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, u)
}

// Auth is authentication handler for BasicAuth middleware
// It hashes credentials and compares them using subtle.ConstantTimeCompare
// to prevent time attacks. If matches returns true which means
// user was successfully authenticated and BasicAuth header was added
func (uc *userController) Auth(username string, password string, ctx echo.Context) (bool, error) {
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

func (uc *userController) UpdatePassword(ctx echo.Context) error {
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

	return ctx.JSON(http.StatusOK, nil)
}

// CreateProfile checks authentication, parses request data (params) to model.Profile
// and creates model.Profile in DB
func (uc *userController) CreateProfile(ctx echo.Context) error {
	cred, err := uc.CheckUserAuth(ctx)
	if err != nil {
		return ctx.String(http.StatusForbidden, err.Error())
	}

	p, err := uc.ParseUserProfileFromServerRequest(ctx, cred)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	insertedID, err := uc.userUsecase.CreateProfile(p)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, struct {
		Id interface{} `json:"_id"`
	}{insertedID})
}

// UpdateProfile checks authentication, parses request data (params) to model.Profile
// and updates model.Profile in DB
func (uc *userController) UpdateProfile(ctx echo.Context) error {
	cred, err := uc.CheckUserAuth(ctx)
	if err != nil {
		return ctx.String(http.StatusForbidden, err.Error())
	}

	p, err := uc.ParseUserProfileFromServerRequest(ctx, cred)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err = uc.userUsecase.UpdateProfile(p)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

// CheckUserAuth checks authentication of model.User
func (uc *userController) CheckUserAuth(ctx echo.Context) (*[]string, error) {
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
func (uc *userController) ParseUserProfileFromServerRequest(
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
		p.AuthUsername = username
	} else {
		p.AuthUsername = username
	}

	if err := ctx.Validate(p); err != nil {
		return nil, err
	}

	return &p, nil
}
