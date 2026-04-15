package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Environment represents a single deployment environment (e.g. staging, production).
type Environment struct {
	Name      string            `yaml:"name"`
	BaseURL   string            `yaml:"base_url"`
	Variables map[string]string `yaml:"variables"`
}

// PatchworkConfig is the top-level configuration structure loaded from patchwork.yaml.
type PatchworkConfig struct {
	Version      string        `yaml:"version"`
	Project      string        `yaml:"project"`
	Environments []Environment `yaml:"environments"`
}

// Load reads and parses a patchwork config file from the given path.
func Load(path string) (*PatchworkConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg PatchworkConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// GetEnvironment returns the environment with the given name, or an error if not found.
func (c *PatchworkConfig) GetEnvironment(name string) (*Environment, error) {
	for i := range c.Environments {
		if c.Environments[i].Name == name {
			return &c.Environments[i], nil
		}
	}
	return nil, fmt.Errorf("environment %q not found in config", name)
}

func (c *PatchworkConfig) validate() error {
	if c.Project == "" {
		return fmt.Errorf("project name is required")
	}
	if len(c.Environments) == 0 {
		return fmt.Errorf("at least one environment must be defined")
	}
	seen := make(map[string]bool)
	for _, env := range c.Environments {
		if env.Name == "" {
			return fmt.Errorf("all environments must have a name")
		}
		if seen[env.Name] {
			return fmt.Errorf("duplicate environment name: %q", env.Name)
		}
		seen[env.Name] = true
	}
	return nil
}
