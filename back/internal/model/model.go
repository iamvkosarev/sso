package model

import (
	"time"
)

type User struct {
	ID       int
	Email    string
	PassHash []byte
}

type App struct {
	Secret   string        `yaml:"secret"`
	TokenTTL time.Duration `yaml:"token_ttl"`
}
