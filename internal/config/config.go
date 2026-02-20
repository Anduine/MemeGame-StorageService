package config

import (
	"flag"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	StoragePath string
}

func MustLoadConfig() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("Config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Config file does not exist: " + path)
	}

	var cfg Config
	data, _ := os.ReadFile(path)

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic("Failed to unmarshal yaml")
	}

	cfg.StoragePath = getStoragePath()

	return &cfg
}

func fetchConfigPath() string {
	var result string

	// --config="path/to/config.yaml"
	flag.StringVar(&result, "config", "./configs/config.yaml", "path to config file")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	return result
}

func getStoragePath() string {
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "storage/"
	}

	return storagePath
}
