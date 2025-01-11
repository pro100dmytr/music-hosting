package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TODO: create DBConfig struct

// TODO: create ServerConfig struct

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

func LoadConfig(configPath string) (*Config, error) {
	filename, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{}

	if err = yaml.Unmarshal(yamlFile, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}
