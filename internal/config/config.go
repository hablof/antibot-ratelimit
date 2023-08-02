package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WindowSize  time.Duration `yaml:"windowsize"`
	WindowLimit int           `yaml:"windowlimit"`
	BanDuration time.Duration `yaml:"banduration"`

	PrefixSize int `yaml:"prefixsize"`
}

func ReadConfig(path string) (Config, error) {

	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
