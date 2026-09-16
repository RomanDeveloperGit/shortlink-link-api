package logger

import (
	"context"
	"log/slog"
	"os"
)

type Color string

const (
	RedColor    Color = "\033[31m"
	YellowColor Color = "\033[33m"
	GreyColor   Color = "\033[90m"
	ResetColor  Color = "\033[0m"
)

var typeColorMap = map[slog.Level]Color{
	slog.LevelDebug: GreyColor,
	slog.LevelWarn:  YellowColor,
	slog.LevelError: RedColor,
}

type LocalHandler struct {
	baseHandler slog.Handler
}

func NewLocalHandler(baseHandler slog.Handler) *LocalHandler {
	return &LocalHandler{baseHandler}
}

func (lh *LocalHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return lh.baseHandler.Enabled(ctx, level)
}

func (lh *LocalHandler) Handle(ctx context.Context, r slog.Record) error {
	if color, ok := typeColorMap[r.Level]; ok {
		os.Stdout.WriteString(string(color))

		defer os.Stdout.WriteString(string(ResetColor))
	}

	return lh.baseHandler.Handle(ctx, r)
}

func (lh *LocalHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LocalHandler{baseHandler: lh.baseHandler.WithAttrs(attrs)}
}

func (lh *LocalHandler) WithGroup(name string) slog.Handler {
	return &LocalHandler{baseHandler: lh.baseHandler.WithGroup(name)}
}
