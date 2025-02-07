package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env"`
	StoragePath string        `yaml:"storage_path" env-requred:"true"`
	TokenTtl    time.Duration `yaml:"token_ttl" env-requred:"true"`
	Grpc        Grpc
}

type Grpc struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	var configPath string

	flag.StringVar(&configPath, "cfg", "", "path to config")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CfgPath")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config is not exists %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config %s", err.Error())
	}

	return &cfg
}
