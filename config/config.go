package config

import "flag"

type Config struct {
	HttpAddress string `env:"HTTP_ADDRESS" envDefault:"0.0.0.0:8000"`
	PostgresDSN string `env:"POSTGRES_DSN" envDefault:"postgres://postgres:postgres@db:5432/hich_db?sslmode=disable"`
	RedisDSN    string `env:"REDIS_DSN" envDefault:"redis://localhost:6379"`
	IsTestEnv   bool
}

func isTestEnv() bool {
	return flag.Lookup("test.v") != nil
}

var cfg *Config

func Get() *Config {
	if cfg != nil {
		return cfg
	}
	cfg = &Config{}
	cfg.IsTestEnv = isTestEnv()
	return cfg
}
