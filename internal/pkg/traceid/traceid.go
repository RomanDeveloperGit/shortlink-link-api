package traceid

import (
	"context"

	"github.com/google/uuid"
)

const (
	TraceIDKey = "X-Trace-ID"
)

func ContextWithTraceID(ctx context.Context, defaultValue string) (context.Context, string) {
	id := defaultValue

	if defaultValue == "" {
		id = uuid.NewString()
	}

	ctx = context.WithValue(ctx, TraceIDKey, id)

	return ctx, id
}

func TraceIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(TraceIDKey).(string)

	if !ok {
		return ""
	}

	return id
}
