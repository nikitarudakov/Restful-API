package main

import (
	daoGRPC "git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc"
	"github.com/rs/zerolog/log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	daoServer, err := daoGRPC.NewProfileDAOServer()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err := daoServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Send()
	}
}
