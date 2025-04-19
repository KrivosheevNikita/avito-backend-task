package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTP HTTPConfig
	GRPC GRPCConfig
	DB   DBConfig
}

type HTTPConfig struct {
	Host string
	Port int
}

type GRPCConfig struct {
	Host string
	Port int
}

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(".env файл не найден")
	}

	httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("parse HTTP_PORT: %w", err)
	}

	grpcPort, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		return nil, fmt.Errorf("parse GRPC_PORT: %w", err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("parse DB_PORT: %w", err)
	}

	cfg := &Config{
		HTTP: HTTPConfig{
			Host: os.Getenv("HTTP_HOST"),
			Port: httpPort,
		},
		GRPC: GRPCConfig{
			Host: os.Getenv("GRPC_HOST"),
			Port: grpcPort,
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}

	return cfg, nil
}

func (c *HTTPConfig) ListenAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *GRPCConfig) ListenAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *DBConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s pool_max_conns=1500",
		c.Host, c.Port, c.Name, c.User, c.Password,
	)
}
