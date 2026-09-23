package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver/response"
)

func Recovery(log *slog.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("recovered from panic",
						slog.Any("error", rec),
						slog.String("stack", string(debug.Stack())),
					)

					response.RespondStatus(
						w,
						http.StatusInternalServerError,
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
