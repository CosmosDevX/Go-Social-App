// Package config
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SecretKey                  string
	DBConnectionString         string
	DBConnectionStringForTests string
	RedisAddres                string
	RedisPassword              string
}

func (c *Config) Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("error during loading .env file")
	}

	c.SecretKey = os.Getenv("SECRET_KEY")
	c.DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
	c.DBConnectionStringForTests = os.Getenv("DB_CONNECTION_STRING_FOR_TESTS")
	c.RedisAddres = os.Getenv("REDIS_ADDRES")
	c.RedisPassword = os.Getenv("REDIS_PASSWORD")
}
