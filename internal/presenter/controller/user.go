package controller

import (
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/repository"
	"git.foxminded.ua/foxstudent106092/user-management/tools/hashing"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

type UserController struct {
	userUsecase UserManager
	AuthEndpointHandler
}

// UserManager contains methods for performing operations on User/Profile datatype
type UserManager interface {
	CreateUser(u *model.User) (*repository.InsertResult, error)
	CreateProfile(p *model.Profile) (*repository.InsertResult, error)
	Find(u *model.User) (*model.User, error)
	StoreVote(v *model.Vote) (*repository.VoteInsertResult, error)
	RetractVote(u *model.Update, sender string) error
	GetRating(target string) (*model.Rating, error)
	UpdateUsername(newUsername, oldUsername string) error
	UpdatePassword(u *model.User) error
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

	userRouter.PUT("/profiles/:target/vote", func(ctx echo.Context) error {
		return uc.Vote(ctx)
	})

	userRouter.DELETE("/profiles/:target/vote/retract", func(ctx echo.Context) error {
		return uc.RetractVote(ctx)
	})

	userRouter.GET("/profiles/:target/rating", func(ctx echo.Context) error {
		return uc.GetRating(ctx)
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

	return ctx.JSON(http.StatusOK, nil)
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

func (uc *UserController) Vote(ctx echo.Context) error {
	vote, err := uc.parseUserVote(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	result, err := uc.userUsecase.StoreVote(vote)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, result)
}

func (uc *UserController) RetractVote(ctx echo.Context) error {
	sender := fmt.Sprintf("%v", ctx.Get("username"))
	target := ctx.Param("target")

	var update = &model.Update{Nickname: target}

	err := uc.userUsecase.RetractVote(update, sender)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (uc *UserController) GetRating(ctx echo.Context) error {
	target := ctx.Param("target")

	rating, err := uc.userUsecase.GetRating(target)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, rating)
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

func (uc *UserController) parseUserVote(ctx echo.Context) (*model.Vote, error) {
	sender := fmt.Sprintf("%v", ctx.Get("username"))
	target := ctx.Param("target")

	if sender == target {
		return nil, errors.New("self-voting is forbidden (bad request)")
	}

	voteStr := ctx.FormValue("vote")
	vote, err := strconv.Atoi(voteStr)
	if err != nil || (vote != -1 && vote != 1) {
		return nil, ctx.String(http.StatusBadRequest, err.Error())
	}

	now := time.Now().Unix()

	voteObj := model.Vote{
		Sender:  sender,
		Target:  target,
		Vote:    vote,
		VotedAt: &now,
	}

	return &voteObj, nil
}
