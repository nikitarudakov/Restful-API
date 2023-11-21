package controller

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/auth"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/profileDao"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/userDao"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type UserController struct {
	userUseCase    UserManager
	profileUseCase ProfileManager
	cacheDb        *cache.Database
}

// UserManager contains methods for performing operations on User/Profile datatype
type UserManager interface {
	CreateUser(u *model.User) (*userDao.InsertResult, error)
	FindUser(u *model.User) (*model.User, error)
	UpdateUsername(newUsername, oldUsername string) error
	UpdatePassword(u *model.User) error
	DeleteUser(username string) error
}

// ProfileManager contains methods for performing operations on User/Profile datatype
type ProfileManager interface {
	CreateProfile(p *model.Profile) (*profileDao.InsertResult, error)
	UpdateProfile(p *model.Update, profileName string) error
	DeleteProfile(profileName string) error
	ListProfiles(page int64) ([]*model.Profile, error)
}

// NewUserController implicitly links  *UserController to userController
// Here to instantiate userController we provide usecase.UserManager
func NewUserController(r *registry.Registry) *UserController {
	return &UserController{r.UserUseCase, r.ProfileUseCase, r.CacheDB}
}

func (uc *UserController) InitUserRoutes(e *echo.Echo, cfg *config.Config) {
	userRouter := e.Group("/users")
	userRoles := []string{"admin", "user", "moderator"}

	auth.InitAuthMiddleware(userRouter, &cfg.Auth, userRoles)

	userRouter.PUT("/profiles/:username/update", func(ctx echo.Context) error {
		return uc.UpdateUserAndProfile(ctx)
	})

	adminRouter := userRouter.Group("/admin")
	adminRoles := []string{"admin", "moderator"}

	auth.InitAuthMiddleware(adminRouter, &cfg.Auth, adminRoles)

	var profiles []model.Profile
	adminRouter.Use(cache.Middleware(uc.cacheDb, &profiles, &cfg.Cache))

	adminRouter.GET("/profiles/list", func(ctx echo.Context) error {
		return uc.ListProfiles(ctx)
	})

	adminRouter.DELETE("/profiles/:username/delete", func(ctx echo.Context) error {
		return uc.DeleteUserAndProfile(ctx)
	})
}

func (uc *UserController) DeleteUserAndProfile(ctx echo.Context) error {
	username := ctx.Param("username")

	err := uc.profileUseCase.DeleteProfile(username)
	if err != nil {
		return err
	}

	err = uc.userUseCase.DeleteUser(username)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserController) ListProfiles(ctx echo.Context) error {
	var page int64 = 1

	pageStr := ctx.QueryParam("page")
	if pageStr != "" {
		parsedPage, err := strconv.ParseInt(ctx.QueryParam("page"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		page = parsedPage
	}

	profiles, err := uc.profileUseCase.ListProfiles(page)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cacheKey := ctx.Request().Method + ":" + ctx.Request().RequestURI
	ctx.Set(cacheKey, profiles)

	return ctx.JSON(http.StatusOK, profiles)
}

// UpdateUserAndProfile checks authentication, parses request data (params) to model.Profile
// and updates model.Profile in DB
func (uc *UserController) UpdateUserAndProfile(ctx echo.Context) error {
	username := ctx.Param("username")

	update, err := uc.parseUserAndProfileUpdate(ctx, username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = uc.profileUseCase.UpdateProfile(update, username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = uc.userUseCase.UpdateUsername(update.Nickname, username)
	if err != nil {
		log.Error().Err(err).
			Str("old_username", username).
			Str("new_username", update.Nickname).
			Msg("error updating username")
	}

	return ctx.JSON(http.StatusOK, nil)
}

// ParseUserProfileFromServerRequest parses server request data to model.Profile
func (uc *UserController) parseUserAndProfileUpdate(
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
