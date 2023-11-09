package controller

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/usecase"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	validator2 "git.foxminded.ua/foxstudent106092/user-management/tools/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetupServer() echo.Context {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/users/profiles/:target/rating", nil)
	w := httptest.NewRecorder()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = &validator2.CustomValidator{
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	return e.NewContext(req, w)
}

func GetVoteController() *VoteController {
	cfg, err := config.InitConfig(".controllerConfig.json")
	if err != nil {
		panic(err)
	}

	db, err := datastore.NewDB(&cfg.Database)
	if err != nil {
		panic(err)
	}

	r := registry.NewRegistry(db, &cfg.Database)

	uu := usecase.NewUserUsecase(r.Ur, r.Pr)
	vu := usecase.NewVoteUsecase(r.Pr, r.Vr)

	ac := NewAuthController(uu, cfg)

	voteController := NewVoteController(vu, ac)

	return voteController
}

func TestVoteController_GetRating(t *testing.T) {
	vv := GetVoteController()

	t.Run("get rating of existing user", func(t *testing.T) {
		ctx := SetupServer()
		ctx.SetParamNames("target")
		ctx.SetParamValues("user2")

		if err := vv.GetRating(ctx); err != nil {
			t.Error(err)
		}

		status := ctx.Response().Status
		assert.Equal(t, 200, status)
	})

	t.Run("get rating of non-existent user", func(t *testing.T) {
		ctx := SetupServer()
		ctx.SetParamNames("target")
		ctx.SetParamValues("user-that-does-not-exist")

		if err := vv.GetRating(ctx); err != nil {
			t.Log(err)
		}

		status := ctx.Response().Status
		assert.Equal(t, 400, status)
	})
}
