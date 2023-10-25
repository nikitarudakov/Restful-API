package logger

import (
	"github.com/rs/zerolog"
)

const loggingLevel = -1

// InitLogger initializes logger with configurations from variable of type *config.Config
func InitLogger() {
	zerolog.SetGlobalLevel(zerolog.Level(loggingLevel))
}
