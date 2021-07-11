package logger

import (
	"github.com/rs/zerolog"
	"os"
)

func NewConsoleLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
}
