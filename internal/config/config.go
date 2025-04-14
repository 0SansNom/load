package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config represents the application settings loaded from a YAML file.
type Config struct {
	ListenAddr          string        `yaml:"listen_addr"`
	Backends            []string      `yaml:"backends"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
	HealthCheckTimeout  time.Duration `yaml:"health_check_timeout"`
}

// Load parses a YAML configuration file into a Config struct.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.UnmarshalStrict(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate checks that the configuration fields have acceptable values.
func (c *Config) validate() error {
	if c.ListenAddr == "" {
		return errors.New("listen_addr is required")
	}
	if len(c.Backends) == 0 {
		return errors.New("at least one backend is required")
	}
	for _, b := range c.Backends {
		if b == "" {
			return errors.New("backend URLs must not be empty")
		}
	}
	if c.HealthCheckInterval <= 0 {
		return errors.New("health_check_interval must be greater than 0")
	}
	if c.HealthCheckTimeout <= 0 {
		return errors.New("health_check_timeout must be greater than 0")
	}
	return nil
}
