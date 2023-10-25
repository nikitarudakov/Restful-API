package main

import (
	"context"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/adapter/gateway"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/router"
	"git.foxminded.ua/foxstudent106092/user-management/internal/registry"
	"git.foxminded.ua/foxstudent106092/user-management/logger"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	logger.InitLogger()

	err := config.InitConfig()
	if err != nil {
		log.Fatal().Err(err).Send()
		panic(err)
	}

	mongoDB, err := datastore.NewDB()
	if err != nil {
		panic(err)
	}

	r := registry.NewRegistry(gateway.NewDatabase(mongoDB))

	// Router for registering new Users
	mainRouter := echo.New()
	mainRouter = router.NewRouter(mainRouter, r.NewAppConroller())

	auth := mainRouter.Group("/auth")
	auth = router.RegisterRouter(auth, r.NewAppConroller())

	defer mongoDB.Disconnect(context.TODO())

	fmt.Println("Server listen at http://localhost" + ":8080")
	if err := mainRouter.Start(":8080"); err != nil {
		log.Fatal().Err(err).Send()
	}

	return
}
