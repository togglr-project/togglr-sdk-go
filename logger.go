package togglr

// Logger interface for logging
type Logger interface {
	Debug(msg string, kv ...any)
	Info(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, kv ...any)
}

// NoOpLogger is a no-op implementation of Logger
type NoOpLogger struct{}

func (NoOpLogger) Debug(msg string, kv ...any) {}
func (NoOpLogger) Info(msg string, kv ...any)  {}
func (NoOpLogger) Warn(msg string, kv ...any)  {}
func (NoOpLogger) Error(msg string, kv ...any) {}
