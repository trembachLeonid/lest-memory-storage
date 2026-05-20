package logging

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

var file *os.File
var environment = "production"

func InitLogger() *slog.Logger {
	var logger *slog.Logger
	if environment == "development" {
		file, _ = os.OpenFile("./logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		multiWriter := io.MultiWriter(os.Stdout, file)
		logger = slog.New(slog.NewJSONHandler(multiWriter, nil))
		slog.SetDefault(logger)
	} else if environment == "production" {
		programLevel := &slog.LevelVar{}
		programLevel.Set(slog.LevelError)

		logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel}))
		slog.SetDefault(logger)
	} else {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	return logger
}

func CloseLogger() {
	if file != nil {
		file.Close()
	}
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value("logger").(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}

func AddAttribute(logger *slog.Logger, name string, value string) *slog.Logger {
	return logger.With(name, value)
}

func AddCorrelationId(logger *slog.Logger) *slog.Logger {
	if environment == "development" {
		return AddAttribute(logger, "correlation_id", uuid.New().String())
	}
	return logger
}

func AddConnectionId(logger *slog.Logger) *slog.Logger {
	if environment == "development" {
		return AddAttribute(logger, "connection_id", uuid.New().String())
	}
	return logger
}
