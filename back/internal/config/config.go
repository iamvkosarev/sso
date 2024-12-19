package config

import (
	"github.com/iamvkosarev/sso/back/internal/model"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type HTTPServerPort struct {
	Address string `yaml:"address"`
}

type Auth struct {
	SecretKey string `yaml:"secret_key"`
	Algorithm string `yaml:"algorithm"`
}

type Config struct {
	Env            string `yaml:"env" env-default:"development"`
	StoragePath    string `yaml:"storage_path"`
	HTTPServerPort `yaml:"http_server"`
	Auth           `yaml:"auth"`
	App            model.App
}

func MustLoad() *Config {
	path := os.Getenv("CONFIG_PATH")

	if path == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(path); err != nil {
		log.Fatal("CONFIG_PATH does not exist")
	}
	var config Config
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	return &config
}
