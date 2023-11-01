package main

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/controller"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/registry"
	"git.foxminded.ua/foxstudent106092/user-management/logger"
	customValidator "git.foxminded.ua/foxstudent106092/user-management/tools/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	logger.InitLogger(&cfg.Logger)

	db, err := datastore.NewDB(&cfg.Database)
	if err != nil {
		panic(err)
	}

	r := registry.NewRegistry(db, &cfg.Database)
	appController := controller.NewAppController(r)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = &customValidator.CustomValidator{
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	appController.User.InitRoutes(e)
	appController.Admin.InitRoutes(e, &cfg.Admin)
	appController.Auth.InitRoutes(e)

	fmt.Println("Server listen at http://localhost" + ":" + cfg.Server.Port)
	if err = e.Start(":" + cfg.Server.Port); err != nil {
		log.Fatalln(err)
	}
}
