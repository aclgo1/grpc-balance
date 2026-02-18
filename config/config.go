package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	Path         string
	ApiPort      string `mapstructure:"API_PORT"`
	DbUrl        string `mapstructure:"DB_URL"`
	DbName       string `mapstructure:"DB_NAME"`
	DbCollection string `mapstructure:"DB_COLLECTION"`
}

func NewConfig(path string) *config {
	return &config{
		Path: path,
	}
}

func (c *config) Load() error {
	v := viper.New()

	v.AddConfigPath(c.Path)
	v.SetConfigName("grpc-balance")
	v.SetConfigFile(".env")
	v.SetConfigType("env")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("v.ReadInConfig: %w", err)
	}

	if err := v.Unmarshal(c); err != nil {
		return fmt.Errorf("v.Unmarshal: %w", err)
	}

	return nil
}
