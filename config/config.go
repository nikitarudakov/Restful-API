package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Database stores values for db connection
type Database struct {
	Name        string `json:"name"`
	UserRepo    string `mapstructure:"user_repo_name"`
	ProfileRepo string `mapstructure:"profile_repo_name"`
	VoteRepo    string `mapstructure:"ratings_repo_name"`
	URL         string `json:"url"`
}

type Cache struct {
	Server         string `json:"server"`
	ExpirationQuan string `mapstructure:"expiration_quantifier"`
	ExpirationUnit string `mapstructure:"expiration_unit"`
	Port           string `json:"PORT"`
}

type Admin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Server struct {
	Port string `json:"PORT"`
}

type Logger struct {
	LoggingLevel int8 `json:"logging_level"`
}

type Auth struct {
	SecretKey string `mapstructure:"secret_key"`
}

type Grpc struct {
	Server string `json:"server"`
	Port   string `json:"PORT"`
}

// Config Create private data struct to hold config options.
type Config struct {
	Database Database `mapstructure:"db"`
	Admin    Admin    `mapstructure:"admin_api"`
	Auth     Auth     `mapstructure:"auth"`
	Cache    Cache    `mapstructure:"cache"`
	Grpc     Grpc     `mapstructure:"grpc"`
	Server   Server   `mapstructure:"server"`
	Logger   Logger   `mapstructure:"logger"`
}

// InitConfig parses .json file to Config struct
func InitConfig(configName string) (*Config, error) {
	config := &Config{}

	viper.SetConfigName(configName)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	// read config
	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Send()
		return nil, fmt.Errorf("error reading Config with viper %w", err)
	}

	// parse config to struct
	if err := viper.Unmarshal(config); err != nil {
		log.Error().Err(err).Send()
		return nil, fmt.Errorf("error unmarshaling to Config struct %w", err)
	}

	return config, nil
}
