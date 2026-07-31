// Package config
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SecretKey          string
	DBConnectionString string
}

func (c *Config) Load() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error during loading .env file")
	}

	c.SecretKey = os.Getenv("SECRET_KEY")
	c.DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
}
