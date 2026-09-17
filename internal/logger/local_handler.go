package logger

import (
	"context"
	"log/slog"
	"os"
)

type color string

const (
	red    color = "\033[31m"
	yellow color = "\033[33m"
	grey   color = "\033[90m"
	reset  color = "\033[0m"
)

var levelColor = map[slog.Level]color{
	slog.LevelDebug: grey,
	slog.LevelWarn:  yellow,
	slog.LevelError: red,
}

type localHandler struct {
	baseHandler slog.Handler
}

func newLocalHandler(baseHandler slog.Handler) *localHandler {
	return &localHandler{baseHandler}
}

func (lh *localHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return lh.baseHandler.Enabled(ctx, level)
}

func (lh *localHandler) Handle(ctx context.Context, r slog.Record) error {
	if color, ok := levelColor[r.Level]; ok {
		os.Stdout.WriteString(string(color))

		defer os.Stdout.WriteString(string(reset))
	}

	return lh.baseHandler.Handle(ctx, r)
}

func (lh *localHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &localHandler{baseHandler: lh.baseHandler.WithAttrs(attrs)}
}

func (lh *localHandler) WithGroup(name string) slog.Handler {
	return &localHandler{baseHandler: lh.baseHandler.WithGroup(name)}
}
