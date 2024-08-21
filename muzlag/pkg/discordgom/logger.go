package discordgom

type Logger interface {
	Warn(msg string, args ...any)
}

type noopLogger struct{}

func (noopLogger) Warn(string, ...any) {}
