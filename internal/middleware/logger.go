package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/requestid"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/traceid"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func logger(next http.HandlerFunc, logger *slog.Logger) (http.HandlerFunc, *slog.Logger) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestId := requestid.RequestIDFromContext(r.Context())
		traceId := traceid.TraceIDFromContext(r.Context())
		userAgent := r.UserAgent()

		logger.Debug("http",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("user_agent", userAgent),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("request_id", requestId),
			slog.String("trace_id", traceId),
		)

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		logger.Info("http",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("user_agent", userAgent),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Int("status", rw.status),
			slog.Duration("duration", time.Since(start)),
			slog.String("request_id", requestId),
			slog.String("trace_id", traceId),
		)
	}), logger
}
