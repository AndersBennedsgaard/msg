package config

import "github.com/spf13/viper"

type Config struct {
	BasePath string `mapstructure:"basePath"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
