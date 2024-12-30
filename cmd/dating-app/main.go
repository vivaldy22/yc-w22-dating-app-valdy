package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"yc-w22-dating-app-valdy/di"
	"yc-w22-dating-app-valdy/internal/server"
)

func main() {
	d := di.SetupDependencies()

	go func() {
		if err := server.StartServer(d); err != nil {
			log.Printf("Server failed to start: %s\n", err.Error())
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := d.Echo.Shutdown(ctx); err != nil {
		log.Printf("Server failed to shutdown: %s\n", err.Error())
	}

	d.CleanUp()

	log.Println("Server shut down.")
}
