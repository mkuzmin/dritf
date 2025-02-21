package aws

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Regions         []string        `yaml:"regions"`
	Services        []ServiceConfig `yaml:"services"`
	IgnoredServices []string        `yaml:"ignored_services"`
}

type ServiceConfig struct {
	Name          string               `yaml:"name"`
	ResourceTypes []ResourceTypeConfig `yaml:"resource_types"`
	IgnoredTypes  []string             `yaml:"ignored_types"`
}

type ResourceTypeConfig struct {
	Name string `yaml:"name"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read '%s' file: %w", filename, err)
	}

	var cfg Config
	err = yaml.UnmarshalStrict(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to parse yaml: %w", err)
	}

	//cfg.Regions = []string{"eu-central-1"} //TODO filter regions

	return &cfg, nil
}
