package main

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	"git.foxminded.ua/foxstudent106092/user-management/logger"
)

func main() {
	cfg, err := config.InitConfig(".config")
	if err != nil {
		panic(err)
	}

	logger.InitLogger(&cfg.Logger)

	db, err := datastore.NewDB(&cfg.Database)
	if err != nil {
		panic(err)
	}

	r := registry.NewRegistry(db, cfg)

	grpc.ServeGRPCServer(r, cfg)
}
