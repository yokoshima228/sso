package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	Grpc
}

type Grpc struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func InitConfig() *Config {
	configPath := flag.String("cfg", "./config/local.yaml", "path to config")
	flag.Parse()

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("config is not exists %s", *configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		log.Fatalf("Error reading config %s", err.Error())
	}

	return &cfg
}
