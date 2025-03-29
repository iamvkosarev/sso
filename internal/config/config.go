package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type CorsOptions struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type Server struct {
	RestPrefix  string      `yaml:"rest_prefix"`
	RESTPort    string      `yaml:"rest_port"`
	GRPCPort    string      `yaml:"grpc_port"`
	CorsOptions CorsOptions `yaml:"cors"`
}

type App struct {
	Secret   string        `yaml:"secret"`
	TokenTTL time.Duration `yaml:"token_ttl"`
}

type Config struct {
	Env    string `yaml:"env" env-default:"development"`
	Server `yaml:"server"`
	App    `yaml:"app"`
}

func MustLoad() *Config {
	path := os.Getenv("CONFIG_PATH")

	if path == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(path); err != nil {
		log.Fatalf("CONFIG_PATH does not exist at: %s\n", path)
	}
	var config Config
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	return &config
}
