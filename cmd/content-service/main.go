package main

import (
	"log"
	srvConfig "story-pulse/internal/content-service/config"
	"story-pulse/internal/content-service/server"
	"story-pulse/internal/shared/config"
)

func main() {
	cfg, err := config.NewConfig[srvConfig.Config]()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	srv, err := server.NewServer(cfg)

	defer func() {
		err = srv.Stop()
	}()

	if err = srv.Run(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
