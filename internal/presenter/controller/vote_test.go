package controller

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	validator2 "git.foxminded.ua/foxstudent106092/user-management/tools/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func TestVoteController_GetRating(t *testing.T) {
	cfg, err := config.InitConfig(".controllerConfig.json")
	if err != nil {
		t.Error(err)
	}

	db, err := datastore.NewDB(&cfg.Database)
	if err != nil {
		t.Error(err)
	}

	r := registry.NewRegistry(db, cfg)

	vv := NewVoteController(r)

	t.Run("get rating of existing user", func(t *testing.T) {
		ctx := SetupServer()
		ctx.SetParamNames("target")
		ctx.SetParamValues("user2")

		if err := vv.GetRating(ctx); err != nil {
			t.Error(err)
		}
	})

	t.Run("get rating of non-existent user", func(t *testing.T) {
		ctx := SetupServer()
		ctx.SetParamNames("target")
		ctx.SetParamValues("user-that-does-not-exist")

		if err := vv.GetRating(ctx); err != nil {
			t.Log(err)
		}
	})
}
