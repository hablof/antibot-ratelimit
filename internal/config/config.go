package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BucketSize  int           `yaml:"bucketsize"`
	RPMLimit    int           `yaml:"rpmlimit"`
	PrefixSize  int           `yaml:"prefixsize"`
	BanDuration time.Duration `yaml:"banduration"`
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
