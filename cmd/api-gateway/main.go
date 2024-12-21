package main

import (
	"log/slog"
	"os"
	"os/signal"
	srvConfig "story-pulse/internal/api-gateway/config"
	"story-pulse/internal/api-gateway/server"
	"story-pulse/internal/shared/config"
	"syscall"
)

func main() {
	cfg, err := config.NewConfig[srvConfig.Config]()
	if err != nil {
		slog.Warn("failed to parse config", "err", err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		slog.Error("failed to create server", "err", err)
		os.Exit(1)
	}

	go func() {
		signalCh := make(chan os.Signal, 1)

		signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

		<-signalCh
		slog.Info("Shutting down server...")

		if err = srv.Stop(); err != nil {
			slog.Warn("failed to stop server", "err", err)
		}
	}()

	if err = srv.Run(); err != nil {
		slog.Error("Error starting server: %v", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")
}
