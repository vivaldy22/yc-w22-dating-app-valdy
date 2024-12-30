package server

import (
	"golang.org/x/sync/errgroup"

	"yc-w22-dating-app-valdy/di"
	"yc-w22-dating-app-valdy/internal/server/http"
)

func StartServer(di *di.DI) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return http.StartHttpServer(di)
	})

	return eg.Wait()
}
