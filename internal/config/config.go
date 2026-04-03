// Package config handles loading and validation of the application configuration.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RedisConfig holds the Redis connection settings.
type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	DB   int    `yaml:"db"`
}

// ListsConfig maps each supported CSV type to its target Redis list key.
type ListsConfig struct {
	Starling  string `yaml:"starling"`
	Amex      string `yaml:"amex"`
	Monzo     string `yaml:"monzo"`
	MonzoFlex string `yaml:"monzo_flex"`
}

// Config is the top-level configuration structure.
type Config struct {
	Redis RedisConfig `yaml:"redis"`
	Lists ListsConfig `yaml:"lists"`
}

// Load reads a YAML configuration file from the given path and returns a Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if cfg.Redis.Host == "" {
		cfg.Redis.Host = "localhost"
	}
	if cfg.Redis.Port == 0 {
		cfg.Redis.Port = 6379
	}

	return &cfg, nil
}

// ListKey returns the Redis list key for the given CSV type.
func (c *Config) ListKey(csvType string) (string, error) {
	switch csvType {
	case "starling":
		return c.Lists.Starling, nil
	case "amex":
		return c.Lists.Amex, nil
	case "monzo":
		return c.Lists.Monzo, nil
	case "monzo-flex":
		return c.Lists.MonzoFlex, nil
	default:
		return "", fmt.Errorf("unsupported CSV type %q", csvType)
	}
}
