package observability_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/observability"
)

func TestContextWithRequestID_ReturnsValidUUID(t *testing.T) {
	t.Parallel()

	_, id := observability.ContextWithRequestID(context.Background())

	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("ContextWithRequestID returned non-UUID id %q: %v", id, err)
	}
}

func TestRequestIDFromContext_ReturnsValidID(t *testing.T) {
	id := "ef9aa731-9674-4877-bcd9-91a80ee9c133"
	ctx := context.WithValue(context.Background(), observability.RequestIDKey, id)

	if got := observability.RequestIDFromContext(ctx); got != id {
		t.Errorf("RequestIDFromContext returned %q, want %q", got, id)
	}
}

func TestRequestIDFromContext_EmptyContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	if got := observability.RequestIDFromContext(ctx); got != "" {
		t.Errorf("RequestIDFromContext(empty) = %q, want %q", got, "")
	}
}
