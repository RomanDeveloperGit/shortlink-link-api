package observability

import (
	"context"

	"github.com/google/uuid"
)

const (
	TraceIDKey = "X-Trace-ID"
)

func ContextWithTraceID(ctx context.Context, currentId string) (context.Context, string) {
	id := currentId

	if currentId == "" {
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
