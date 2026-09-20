package middleware

import (
	"log/slog"
	"net/http"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/traceid"
)

func traceID(next http.HandlerFunc, logger *slog.Logger) (http.HandlerFunc, *slog.Logger) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, id := traceid.ContextWithTraceID(r.Context(), traceid.TraceIDFromContext(r.Context()))

		w.Header().Set(traceid.TraceIDKey, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	}), logger
}
