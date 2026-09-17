package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type httpServer struct {
	server                  *http.Server
	logger                  *slog.Logger
	gracefulShutdownTimeout time.Duration
}

func NewHTTPServer(
	host string,
	port int,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	gracefulShutdownTimeout time.Duration,
	logger *slog.Logger,
) *httpServer {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	return &httpServer{
		server,
		logger,
		gracefulShutdownTimeout,
	}
}

func (hs *httpServer) BackgroundRun() {
	go func() {
		if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			hs.logger.Error("failed to start HTTP server", slog.String("error", err.Error()))
		}
	}()
}

func (hs *httpServer) ShutdownGracefully() {
	ctx, cancel := context.WithTimeout(context.Background(), hs.gracefulShutdownTimeout)
	defer cancel()

	if err := hs.server.Shutdown(ctx); err != nil {
		hs.logger.Error("failed to shutdown HTTP server", slog.String("error", err.Error()))
	}

	hs.logger.Debug("HTTP server closed")
}
