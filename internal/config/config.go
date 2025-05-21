package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB     DBConfig     `yaml:"db"`
	Server ServerConfig `yaml:"server"`
	Logger Logger       `yaml:"logger"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type Logger struct {
	LogLevel string `yaml:"log_level"`
}

func (l *Logger) GetLogLevel() slog.Level {
	switch l.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
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

	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}
