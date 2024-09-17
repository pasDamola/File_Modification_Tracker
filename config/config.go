package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Directory  string `mapstructure:"directory" validate:"required"`
	Frequency  int    `mapstructure:"frequency" validate:"required,min=1"`
	SocketPath string `mapstructure:"socket_path" validate:"required"`
}

var validate *validator.Validate

func LoadConfig() (Config, error) {
	var config Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	validate = validator.New()
	err = validate.Struct(config)
	if err != nil {
		return config, err
	}

	return config, nil
}
