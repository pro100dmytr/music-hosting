package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TODO: create DBConfig struct
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// TODO: create ServerConfig struct
type ServerConfig struct {
	Port string `yaml:"port"`
}

type Logger struct {
	LogLevel string `yaml:"log_level"`
}

func LoadConfig(configPath string) (*ServerConfig, *DBConfig, *Logger, error) {
	filename, err := filepath.Abs(configPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid config path: %w", err)
	}

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error reading config file: %w", err)
	}

	var rawConfig struct {
		Server   ServerConfig `yaml:"server"`
		Database DBConfig     `yaml:"database"`
		Logger   Logger       `yaml:"logger"`
	}

	if err := yaml.Unmarshal(yamlFile, &rawConfig); err != nil {
		return nil, nil, nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &rawConfig.Server, &rawConfig.Database, &rawConfig.Logger, nil
}
