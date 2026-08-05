package helpers

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func NewTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	cfg := LoadTestConfig()
	connStr := strings.Replace(cfg.DBConnectionString, "host=localhost", "host=127.0.0.1", 1)
	db, err := sqlx.Connect("postgres", connStr)
	require.NoError(t, err, "failed to connect to test database")

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute)

	err = db.Ping()
	require.NoError(t, err, "failed to ping test database")

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func NewTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	require.NoError(t, err, "failed to ping Redis")

	t.Cleanup(func() {
		_ = client.Close()
	})

	return client
}

func TruncateAllTables(t *testing.T, db *sqlx.DB) {
	t.Helper()
	ctx := context.Background()

	tables := []string{"post_likes", "comments", "posts", "users"}
	for _, table := range tables {
		_, err := db.ExecContext(ctx, "TRUNCATE TABLE "+table+" RESTART IDENTITY CASCADE")
		if err != nil {
			log.Printf("truncate %s: %v", table, err)
			require.NoError(t, err)
		}
	}
}

func FlushRedis(t *testing.T, client *redis.Client) {
	t.Helper()
	ctx := context.Background()
	err := client.FlushDB(ctx).Err()
	require.NoError(t, err, "failed to flush Redis")
}
