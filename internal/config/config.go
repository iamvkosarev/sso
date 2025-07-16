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
	HTTPAddress string      `yaml:"http_address"`
	GRPCAddress string      `yaml:"grpc_address"`
	CorsOptions CorsOptions `yaml:"cors"`
}

type App struct {
	Secret              string        `yaml:"secret"`
	TokenTTL            time.Duration `yaml:"token_ttl"`
	ShuttingDownTimeout time.Duration `yaml:"shutting_down"`
}

type OTelTracing struct {
	AlwaysSample              bool   `yaml:"always_sample"`
	ServiceName               string `yaml:"service_name"`
	TraceGRPCExporterEndpoint string `yaml:"trace_grpc_exporter_endpoint"`
}

type OTel struct {
	Tracing OTelTracing `yaml:"tracing"`
}

type Config struct {
	Env    string `yaml:"env" env-default:"development"`
	OTel   `yaml:"open_telemetry"`
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
