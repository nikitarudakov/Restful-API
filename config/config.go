package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// DatabaseConfig stores values for db connection
type DatabaseConfig struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type AdminCred struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Config Create private data struct to hold config options.
type config struct {
	Database DatabaseConfig `mapstructure:"db"`
	AdminAPI AdminCred      `mapstructure:"admin_api"`
}

var C config

// InitConfig parses .json file to Config struct
func InitConfig() error {
	Config := &C

	viper.SetConfigName(".config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	// read config
	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Send()
		return fmt.Errorf("error reading Config with viper %w", err)
	}

	// parse config to struct
	if err := viper.Unmarshal(&Config); err != nil {
		log.Error().Err(err).Send()
		return fmt.Errorf("error unmarshaling to Config struct %w", err)
	}

	return nil
}
