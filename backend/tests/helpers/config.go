package helpers

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"myapp/internal/config"

	"github.com/joho/godotenv"
)

func LoadTestConfig() config.Config {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "../..")
	envPath := filepath.Join(root, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("failed to load .env for tests: %v", err)
	}

	return config.Config{
		SecretKey:                  os.Getenv("SECRET_KEY"),
		DBConnectionStringForTests: os.Getenv("DB_CONNECTION_STRING_FOR_TESTS"),
	}
}
