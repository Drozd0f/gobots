package log

import (
	"log/slog"
	"os"
)

func NewLogger(level slog.Level) *slog.Logger {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(l)

	return l
}

func SlogError(err error) slog.Attr {
	return slog.Any("error", err)
}
