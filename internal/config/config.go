package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Addr       string `yaml:"address"`
	StorageURL string `yaml:"database_url"`
}

func MustLoad() *Config {
	var configPath string
	flag.StringVar(&configPath, "config", "", "configPath to config file")
	flag.Parse()
	if configPath == "" {
		panic("config file configPath is required")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not found")
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}
	return &cfg
}
