package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	AUTH          string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		ServerAddress: viper.GetString("SERVER_ADDRESS"),
		DatabaseURL:   viper.GetString("DATABASE_URL"),
		AUTH:          viper.GetString("AUTH"),
	}

	return cfg, nil
}
