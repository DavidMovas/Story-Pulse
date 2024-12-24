package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"story-pulse/internal/shared/config"
	srvConfig "story-pulse/internal/users-service/config"
	"story-pulse/internal/users-service/server"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.NewConfig[srvConfig.Config]()
	if err != nil {
		slog.Warn("failed to parse config", "err", err)
	}

	cfg.SetDefaults()

	srv, err := server.NewServer(ctx, cfg)
	if err != nil {
		slog.Warn("failed to create server", "err", err)
		os.Exit(1)
	}

	go func() {
		signalCh := make(chan os.Signal, 1)

		signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

		<-signalCh
		slog.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.GracefulShutdownTimeout)*time.Second)
		defer cancel()

		if err = srv.Stop(ctx); err != nil {
			slog.Warn("Error stopping server: %v", err)
		} else {
			slog.Info("Server gracefully stopped")
		}
	}()

	if err = srv.Run(); err != nil {
		slog.Error("Error starting server: %v", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")
}
