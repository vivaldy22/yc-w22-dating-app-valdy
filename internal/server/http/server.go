package http

import (
	"log"

	"yc-w22-dating-app-valdy/di"
)

type server struct {
	di *di.DI
	h  *handler
}

func StartHttpServer(di *di.DI) error {
	log.Println("Starting Http Server...")

	s := &server{
		di: di,
	}

	s.SetupHandlers()
	//s.SetupMiddlewares()
	s.SetupRoutes()

	err := s.di.Echo.Start(di.Configuration.GetHttpPort())
	if err != nil {
		return err
	}

	return nil
}
