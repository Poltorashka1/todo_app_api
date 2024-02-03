package logging

import (
	"log/slog"
	"os"
	colorLog "web/internal/logging/lib"
)

// SetupLogger setup logger for application.
func SetupLogger() *slog.Logger {
	// Create a new logger with a text handler that writes to os.Stdout
	opts := colorLog.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := colorLog.NewPrettyHandler(os.Stdout, opts)
	log := slog.New(handler)

	return log
}
