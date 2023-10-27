package main

import (
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/router"
	"git.foxminded.ua/foxstudent106092/user-management/internal/registry"
	"git.foxminded.ua/foxstudent106092/user-management/logger"
	"log"
)

func main() {
	logger.InitLogger()

	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	db, err := datastore.NewDB()
	if err != nil {
		panic(err)
	}

	r := registry.NewRegistry(db)

	mainRouter := router.SetUserRouter(r)
	router.SetAdminGroupRouter(r, mainRouter)

	fmt.Println("Server listen at http://localhost" + ":8080")
	if err = mainRouter.Start(":8080"); err != nil {
		log.Fatalln(err)
	}
}
