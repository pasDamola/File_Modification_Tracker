package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Directory  string `mapstructure:"directory"`
	Frequency  int    `mapstructure:"frequency"`
	SocketPath string `mapstructure:"socket_path"`
}

func LoadConfig() (Config, error) {
	var config Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	err := viper.Unmarshal(&config)
	return config, err
}
