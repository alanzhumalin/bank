package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Str("app", "minibank").Timestamp().Logger()
	return logger
}
