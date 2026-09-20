package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/httpjson"
)

func recovery(next http.HandlerFunc, logger *slog.Logger) (http.HandlerFunc, *slog.Logger) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("recovered from panic",
					slog.Any("error", rec),
					slog.String("stack", string(debug.Stack())),
				)

				httpjson.WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	}), logger
}
