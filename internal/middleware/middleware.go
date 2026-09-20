package middleware

import (
	"log/slog"
	"net/http"
)

func Wrapper(log *slog.Logger) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		wrappedHandler, _ := recovery(traceID(requestID(logger(next, log))))

		return wrappedHandler
	}
}
