package middleware

import (
	"net/http"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/observability"
)

func TraceID() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, id := observability.ContextWithTraceID(
				r.Context(),
				r.Header.Get(observability.TraceIDKey),
			)

			w.Header().Set(observability.TraceIDKey, id)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
