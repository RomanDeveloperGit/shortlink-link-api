package httpserver

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

type HealthHandler interface {
	Check(w http.ResponseWriter, r *http.Request)
}

type LinkHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByShortCode(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Visit(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	LinkHandler   LinkHandler
	HealthHandler HealthHandler
}

type Options struct {
	Host                 string
	Port                 int
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	GracefulShutdownTime time.Duration
	Logger               *slog.Logger
	Middleware           func(http.HandlerFunc) http.HandlerFunc
	Handlers             Handlers
}

func NewHTTPServer(opts *Options) *httpServer {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", opts.Middleware(opts.Handlers.HealthHandler.Check))

	mux.HandleFunc("POST /api/v1/links", opts.Middleware(opts.Handlers.LinkHandler.Create))
	mux.HandleFunc("GET /api/v1/links", opts.Middleware(opts.Handlers.LinkHandler.GetByShortCode)) // with "short_code" query!
	mux.HandleFunc("GET /api/v1/links/{id}", opts.Middleware(opts.Handlers.LinkHandler.GetByID))

	mux.HandleFunc("GET /{short_code}", opts.Middleware(opts.Handlers.LinkHandler.Visit))

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Handler:      mux,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
	}

	return &httpServer{
		server:                  server,
		logger:                  opts.Logger,
		gracefulShutdownTimeout: opts.GracefulShutdownTime,
	}
}

func (hs *httpServer) BackgroundRun() {
	go func() {
		hs.logger.Debug("HTTP server started",
			slog.String("addr", hs.server.Addr),
		)

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
