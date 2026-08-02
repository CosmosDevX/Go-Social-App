// Package infrastructure
package infrastructure

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SQLxClient struct {
	db *sqlx.DB
}

func NewSQLxClient(connStr string) SQLxClient {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	return SQLxClient{
		db: db,
	}
}

func (c SQLxClient) GetDB() *sqlx.DB {
	return c.db
}

func (c SQLxClient) CreateTables(schema string) error {
	_, err := c.db.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}

func (c SQLxClient) Shutdown() error {
	if err := c.db.Close(); err != nil {
		return err
	}

	return nil
}
