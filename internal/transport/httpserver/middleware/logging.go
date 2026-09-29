package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver/middleware/observability"
)

const maxBodyLog = 1024

type responseWriter struct {
	http.ResponseWriter

	status   int
	response string
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.response = string(b)

	return rw.ResponseWriter.Write(b)
}

func getHeaders(header http.Header) string {
	headers := strings.Builder{}

	for k, v := range header {
		headers.WriteString(k)
		headers.WriteString(":")
		headers.WriteString(strings.Join(v, ","))
		headers.WriteString(",")
	}

	return strings.TrimSuffix(headers.String(), ",")
}

func getBody(r *http.Request) string {
	ct := r.Header.Get("Content-Type")

	if r.Body == nil {
		return ""
	}

	if strings.HasPrefix(ct, "application/json") ||
		strings.HasPrefix(ct, "application/x-www-form-urlencoded") ||
		strings.HasPrefix(ct, "text/") {
		bodyBytes, _ := io.ReadAll(io.LimitReader(r.Body, maxBodyLog))
		r.Body.Close()

		r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(bodyBytes), r.Body))

		return string(bodyBytes)
	}

	return ""
}

func Logging(log *slog.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log = observability.NewRequestLogger(r.Context(), log)

			log.Info("http request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.String("headers", getHeaders(r.Header)),
				slog.String("body", getBody(r)),
				slog.String("remote_addr", r.RemoteAddr),
			)

			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)

			log.Info("http response",
				slog.Int("status", rw.status),
				slog.Duration("duration", time.Since(start)),
				slog.String("headers", getHeaders(rw.Header())),
				slog.String("response", rw.response),
			)
		})
	}
}
