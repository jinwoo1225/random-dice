package config

import "github.com/spf13/viper"

type Config struct {
	// MongoDB
	MongoDB struct {
		Host     string `yaml:"Host"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
		Database string `yaml:"Database"`
	} `yaml:"MongoDB"`

	Server struct {
		Host     string `yaml:"Host"`
		GRPCPort int64  `yaml:"GRPCPort"`
	} `yaml:"Server"`
}

func NewConfig() (*Config, error) {
	var config *Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
