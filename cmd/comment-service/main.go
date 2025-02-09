package main

import (
	srvConfig "brain-wave/internal/comment-service/config"
	"brain-wave/internal/comment-service/server"
	"brain-wave/internal/shared/config"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.NewConfig[srvConfig.Config]()
	if err != nil {
		slog.Warn("failed to parse config", "err", err)
	}

	srv, err := server.NewServer(cfg)

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
