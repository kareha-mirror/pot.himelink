package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Route struct {
	Path     string `yaml:"path"`
	Protocol string `yaml:"protocol"`
	API      string `yaml:"api"`
	Raw      string `yaml:"raw,omitempty"`
}

type Config struct {
	App struct {
		Addr string `yaml:"addr"`
	} `yaml:"app"`

	Site struct {
		Name   string `yaml:"name"`
		Static string `yaml:"static"`
	} `yaml:"site"`

	Routes []Route `yaml:"routes"`

	Client struct {
		UserAgent   string `yaml:"user-agent"`
		GitHubToken string `yaml:"github-token"`
	}
}

func Load(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	return &cfg
}
