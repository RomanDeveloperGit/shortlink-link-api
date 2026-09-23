package middleware

import (
	"net/http"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver/middleware/observability"
)

func RequestID() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, id := observability.ContextWithRequestID(r.Context())

			w.Header().Set(observability.RequestIDKey, id)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
