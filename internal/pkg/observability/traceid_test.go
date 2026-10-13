package observability_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/observability"
)

func TestContextWithTraceID_GeneratesNewIDWhenEmpty(t *testing.T) {
	t.Parallel()

	ctx, id := observability.ContextWithTraceID(context.Background(), "")

	if id == "" {
		t.Fatal("ContextWithTraceID(\"\") returned empty id, want generated UUID")
	}

	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("generated id %q is not a valid UUID: %v", id, err)
	}

	if got := observability.TraceIDFromContext(ctx); got != id {
		t.Errorf("TraceIDFromContext() = %q, want %q", got, id)
	}
}

func TestContextWithTraceID_KeepsProvidedID(t *testing.T) {
	t.Parallel()

	const provided = "trace-from-upstream"

	ctx, id := observability.ContextWithTraceID(context.Background(), provided)

	if id != provided {
		t.Errorf(
			"ContextWithTraceID() id = %q, want %q (must not overwrite incoming id)",
			id,
			provided,
		)
	}
	if got := observability.TraceIDFromContext(ctx); got != provided {
		t.Errorf("TraceIDFromContext() = %q, want %q", got, provided)
	}
}

func TestTraceIDFromContext_MissingAndWrongType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "empty context",
			ctx:  context.Background(),
		},
		{
			name: "wrong value type",
			ctx:  context.WithValue(context.Background(), observability.TraceIDKey, 12345),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := observability.TraceIDFromContext(tt.ctx); got != "" {
				t.Errorf("TraceIDFromContext() = %q, want %q", got, "")
			}
		})
	}
}

func BenchmarkContextWithTraceID_Generate(b *testing.B) {
	parent := context.Background()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = observability.ContextWithTraceID(parent, "")
	}
}

func BenchmarkContextWithTraceID_Reuse(b *testing.B) {
	parent := context.Background()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = observability.ContextWithTraceID(parent, "existing-trace")
	}
}
