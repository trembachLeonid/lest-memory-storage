package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
)

var file *os.File

func InitLogger() *slog.Logger {
	file, _ = os.OpenFile("./logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	multiWriter := io.MultiWriter(os.Stdout, file)
	logger := slog.New(slog.NewJSONHandler(multiWriter, nil))
	slog.SetDefault(logger)

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
