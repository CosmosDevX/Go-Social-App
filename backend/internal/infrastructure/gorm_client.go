// Package infrastructure
package infrastructure

import (
	"log"
	"myapp/internal/model"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormClient struct {
	db *gorm.DB
}

func (c *GormClient) Setup() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error during loading .env file")
	}

	dsn := os.Getenv("CONNECTION_STRING")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{}, &model.PostLike{}); err != nil {
		log.Fatal(err)
	}

	c.db = db
}

func (c GormClient) GetDB() *gorm.DB {
	return c.db
}

func (c GormClient) Shutdown() error {
	sqliteDB, err := c.db.DB()
	if err != nil {
		return err
	}

	if err := sqliteDB.Close(); err != nil {
		return err
	}

	return nil
}
