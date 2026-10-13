package shortcode_test

import (
	"strings"
	"testing"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/service/link/shortcode"
)

func TestGenerateShortCode_Length(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size int
	}{
		{name: "zero", size: 0},
		{name: "one", size: 1},
		{name: "typical", size: 8},
		{name: "large", size: 64},
		{name: "very large", size: 1024},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, err := shortcode.GenerateShortCode(tt.size)
			if err != nil {
				t.Fatalf("GenerateShortCode(%d) returned unexpected error: %v", tt.size, err)
			}

			if got := len(code); got != tt.size {
				t.Errorf("GenerateShortCode(%d) length = %d, want %d", tt.size, got, tt.size)
			}
		})
	}
}

func TestGenerateShortCode_Alphabet(t *testing.T) {
	t.Parallel()

	const (
		size    = 256
		allowed = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	)

	code, err := shortcode.GenerateShortCode(size)
	if err != nil {
		t.Fatalf("GenerateShortCode() returned unexpected error: %v", err)
	}

	for i, r := range code {
		if !strings.ContainsRune(allowed, r) {
			t.Errorf("GenerateShortCode() char at position %d = %q, not in allowed alphabet", i, r)
		}
	}
}

func TestGenerateShortCode_Uniqueness(t *testing.T) {
	t.Parallel()

	const (
		iterations = 1000
		size       = 16
	)

	seen := make(map[string]struct{}, iterations)
	for i := 0; i < iterations; i++ {
		code, err := shortcode.GenerateShortCode(size)
		if err != nil {
			t.Fatalf("GenerateShortCode() iteration %d returned unexpected error: %v", i, err)
		}

		if _, exists := seen[code]; exists {
			t.Fatalf("GenerateShortCode() produced duplicate code %q at iteration %d", code, i)
		}
		seen[code] = struct{}{}
	}
}

func TestGenerateShortCode_NegativeSize(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GenerateShortCode(-1) did not panic, want panic")
		}
	}()

	_, _ = shortcode.GenerateShortCode(-1)
}
