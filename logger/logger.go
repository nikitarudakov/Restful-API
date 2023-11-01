package logger

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"github.com/rs/zerolog"
)

// InitLogger initializes logger with configurations from variable of type *config.Config
func InitLogger(loggerCfg *config.Logger) {
	zerolog.SetGlobalLevel(zerolog.Level(loggerCfg.LoggingLevel))
}
