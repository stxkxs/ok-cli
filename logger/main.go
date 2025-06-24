package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func Setup(debug bool) error {
	var level zerolog.Level
	if debug {
		level = zerolog.DebugLevel
	} else {
		level = zerolog.InfoLevel
	}

	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)

	return nil
}
