package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Regions         []string  `yaml:"regions"`
	Services        []Service `yaml:"services"`
	IgnoredServices []string  `yaml:"ignored_services"`
}

type Service struct {
	Name                 string         `yaml:"name"`
	ResourceTypes        []ResourceType `yaml:"resource_types"`
	IgnoredResourceTypes []string       `yaml:"ignored_resource_types"`
}

type ResourceType struct {
	Name    string   `yaml:"name"`
	Regions []string `yaml:"regions"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read '%s' file: %w", filename, err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to parse yaml: %w", err)
	}

	return &cfg, nil
}
