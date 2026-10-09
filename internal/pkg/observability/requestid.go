package observability

import (
	"context"

	"github.com/google/uuid"
)

const (
	RequestIDKey = "X-Request-ID"
)

func ContextWithRequestID(ctx context.Context) (context.Context, string) {
	id := uuid.NewString()
	ctx = context.WithValue(ctx, RequestIDKey, id)

	return ctx, id
}

func RequestIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(RequestIDKey).(string)

	if !ok {
		return ""
	}

	return id
}
