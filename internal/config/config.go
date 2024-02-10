package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
)

// Config holds the application configuration.
type Config struct {
	Server Server
	DB     DB
}

// Server holds the server configuration.
type Server struct {
	Port int    `env:"PORT"`
	Env  string `env:"ENV"`
}

// DB holds the database configuration.
type DB struct {
	Port     int    `env:"DB_PORT"`
	Host     string `env:"DB_HOST"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
	SslMode  string `env:"DB_SSLMODE"`
}

// New returns a new Config.
func New() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	return &cfg, nil
}

// Addr returns the server port.
func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}

// DSN returns the database connection string.
func (c *Config) DSN() string {
	fmt.Println(c.DB.Port)
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
		c.DB.SslMode,
	)
}
