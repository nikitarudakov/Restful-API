package main

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/registry"
	"git.foxminded.ua/foxstudent106092/user-management/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	logger.InitLogger()

	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	client, err := datastore.NewDB()
	if err != nil {
		panic(err)
	}

	r := registry.NewRegistry(client)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/auth", func(ctx echo.Context) error {
		return r.NewAppController().User.Register(ctx)
	})

	fmt.Println("Server listen at http://localhost" + ":8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalln(err)
	}
}
