package controller

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/auth"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/repository"
	"git.foxminded.ua/foxstudent106092/user-management/tools/hashing"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type AuthController struct {
	userUsecase UserManager
	cfg         *config.Config
}

type AuthResult struct {
	User    *repository.InsertResult
	Profile *repository.InsertResult
}

func NewAuthController(uu UserManager, cfg *config.Config) *AuthController {
	return &AuthController{userUsecase: uu, cfg: cfg}
}

func (ac *AuthController) InitRoutes(e *echo.Echo) {
	regRouter := e.Group("/auth")

	regRouter.POST("/register", func(ctx echo.Context) error {
		return ac.Register(ctx)
	})

	regRouter.POST("/login", func(ctx echo.Context) error {
		return ac.Login(ctx)
	})
}

func (ac *AuthController) InitAuthMiddleware(g *echo.Group, accessibleRoles []string) {
	tokenConfig := auth.GetTokenConfig(&ac.cfg.Auth, accessibleRoles)
	g.Use(echojwt.WithConfig(tokenConfig))
}

func (ac *AuthController) UpdatePassword(ctx echo.Context) error {
	var u model.User

	u.Username = fmt.Sprintf("%v", ctx.Get("username"))

	password, err := hashing.HashPassword(ctx.FormValue("password"))
	u.Password = password
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err = ac.userUsecase.UpdatePassword(&u); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (ac *AuthController) Login(ctx echo.Context) error {
	var u model.User

	u.Username = ctx.FormValue("username")
	password := ctx.FormValue("password")

	userFromDB, err := ac.userUsecase.Find(&u)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden,
			fmt.Sprintf("user was not found: %s", err.Error()))
	}

	u.Role = userFromDB.Role

	if err = hashing.CheckPassword(userFromDB.Password, password); err != nil ||
		subtle.ConstantTimeCompare([]byte(u.Username), []byte(userFromDB.Username)) != 1 {

		return echo.NewHTTPError(http.StatusForbidden,
			fmt.Sprintf("username/password is incorrect: %s", err.Error()))
	}

	authCfg := ac.cfg.Auth
	token, err := auth.GenerateJWTToken(&u, []byte(authCfg.SecretKey))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx.Set("username", u.Username)

	return ctx.JSON(http.StatusOK, echo.Map{
		"token":    token,
		"username": u.Username,
		"message":  "Successfully logged in!",
	})
}

// Register creates and stores new model.User in DB
func (ac *AuthController) Register(ctx echo.Context) error {
	var u model.User

	userResult, err := ac.registerUser(ctx, u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	profileResult, err := ac.registerProfile(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, AuthResult{User: userResult, Profile: profileResult})
}

func (ac *AuthController) registerUser(ctx echo.Context, u model.User) (*repository.InsertResult, error) {
	u.Username = ctx.FormValue("username")
	u.Password = ctx.FormValue("password")
	u.Role = ctx.FormValue("role")

	log.Trace().Str("username", u.Username).
		Str("role", u.Role).
		Msg("register request")

	if err := ctx.Validate(u); err != nil {
		return nil, err
	}

	if u.Role == "admin" {
		if subtle.ConstantTimeCompare([]byte(u.Username), []byte(ac.cfg.Admin.Username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(u.Password), []byte(ac.cfg.Admin.Password)) != 1 {

			return nil, errors.New("error: unable to register ADMIN user")
		}
	}

	hashedPassword, err := hashing.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = hashedPassword

	result, err := ac.userUsecase.CreateUser(&u)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (ac *AuthController) registerProfile(ctx echo.Context) (*repository.InsertResult, error) {
	username := ctx.FormValue("username")

	p, err := ac.parseValidateUserProfileCreate(ctx, username)
	if err != nil {
		return nil, err
	}

	result, err := ac.userUsecase.CreateProfile(p)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ParseUserProfileFromServerRequest parses server request data to model.Profile
func (ac *AuthController) parseValidateUserProfileCreate(
	ctx echo.Context,
	username string) (*model.Profile, error) {

	var p model.Profile
	if err := ctx.Bind(&p); err != nil {
		return nil, err
	}

	// check if new Nickname was passed
	if p.Nickname == "" {
		p.Nickname = username // assign User username to Update nickname
	}

	if err := ctx.Validate(p); err != nil {
		return nil, err
	}

	return &p, nil
}
