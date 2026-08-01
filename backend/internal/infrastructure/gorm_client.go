// Package infrastructure
package infrastructure

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormClient struct {
	db *gorm.DB
}

func (c *GormClient) Setup(connStr string) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		/*Logger: logger.Default.LogMode(logger.Silent),*/
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(); err != nil {
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
