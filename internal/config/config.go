package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Config struct {
	Server         `yaml:"server"`
	DatabaseConfig `yaml:"databaseConfig"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Type   string            `yaml:"type"`
	Config map[string]string `yaml:"config"`
}

// NewConfig read and create Config for project
func NewConfig(configFilePath string, log *slog.Logger) *Config {
	//validate configFilePath
	if err := validateConfigPath(configFilePath); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	//read configFile and return cfg
	cfg, err := ReadConfig(configFilePath)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	return cfg
}

// ReadConfig Decode configFile in Config{}
func ReadConfig(configFilePath string) (*Config, error) {
	const op = "config.ReadConfig"

	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err.Error())
	}

	cfg := &Config{}
	err = yaml.NewDecoder(file).Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err.Error())
	}
	return cfg, nil
}

// validateConfigPath validate configFilePath
func validateConfigPath(configFilePath string) error {
	const op = "config.validateConfigPath"
	if _, err := os.Stat(configFilePath); err != nil {
		return fmt.Errorf("%v: %v", op, err.Error())
	}
	return nil
}
