package postgresql

import (
	"fmt"
	"os"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Username, c.Password, c.Host, c.Port, c.Database)
}

func NewConfig() (*Config, error) {
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DB")

	if username == "" || password == "" || host == "" || port == "" || database == "" {
		return nil, fmt.Errorf("отсутствует необходимая конфигурация базы данных")
	}

	return &Config{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}, nil
}
