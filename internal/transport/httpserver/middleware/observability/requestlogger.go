package observability

import (
	"context"
	"log/slog"
)

func NewRequestLogger(ctx context.Context, log *slog.Logger) *slog.Logger {
	return log.With(
		slog.String("trace_id", TraceIDFromContext(ctx)),
		slog.String("request_id", RequestIDFromContext(ctx)),
	)
}
