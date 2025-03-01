package config

type Config struct {
	HttpAddress string `env:"HTTP_ADDRESS" envDefault:"0.0.0.0:8000"`
	PostgresDSN string `env:"POSTGRES_DSN" envDefault:"postgres://postgres:postgres@db:5432/hich_db?sslmode=disable"`
	RedisDSN    string `env:"REDIS_DSN" envDefault:"redis://localhost:6379"`
}

var cfg Config

func Get() *Config {
	return &cfg
}
