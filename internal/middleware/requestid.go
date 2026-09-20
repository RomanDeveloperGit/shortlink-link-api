package middleware

import (
	"log/slog"
	"net/http"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/requestid"
)

func requestID(next http.HandlerFunc, logger *slog.Logger) (http.HandlerFunc, *slog.Logger) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, id := requestid.ContextWithRequestID(r.Context())

		w.Header().Set(requestid.RequestIDKey, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	}), logger
}
