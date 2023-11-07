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
	AddUpdateVoteTarget(u *model.Update) error
	AddUpdateVoteSender(u *model.Update, vote *model.Vote) error
	RetractVote(uTarget *model.Update, uSender *model.Update) error
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
	voteObj, err := uc.parseUserVote(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if voteObj.Target == voteObj.Sender {
		return ctx.JSON(http.StatusBadRequest, errors.New("self-voting is forbidden"))
	}

	senderUpdate := &model.Update{
		Nickname: voteObj.Sender,
	}

	err = uc.userUsecase.AddUpdateVoteSender(senderUpdate, voteObj)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	targetUpdate := &model.Update{
		Nickname: voteObj.Target,
		Rating:   &voteObj.Vote,
	}

	fmt.Printf("%+v\n", targetUpdate)

	err = uc.userUsecase.AddUpdateVoteTarget(targetUpdate)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	log.Info().
		Str("sender", voteObj.Sender).
		Str("target", voteObj.Target).
		Str("service", "/users/profile/:username/vote").
		Msg("vote has been recorded")

	return ctx.JSON(http.StatusOK, nil)
}

func (uc *UserController) RetractVote(ctx echo.Context) error {
	sender := fmt.Sprintf("%v", ctx.Get("username"))
	target := ctx.Param("target")

	targetUpdate := &model.Update{
		Nickname: target,
	}

	senderUpdate := &model.Update{
		Nickname: sender,
	}

	err := uc.userUsecase.RetractVote(targetUpdate, senderUpdate)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
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

func (uc *UserController) parseUserVote(ctx echo.Context) (*model.Vote, error) {
	sender := fmt.Sprintf("%v", ctx.Get("username"))
	target := ctx.Param("target")

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
