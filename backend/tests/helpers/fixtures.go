package helpers

import (
	"context"
	"testing"
	"time"

	"myapp/internal/domain"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func CreateTestUser(t *testing.T, db *sqlx.DB, username, password string) domain.User {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	require.NoError(t, err)

	var id int
	err = db.QueryRowContext(context.Background(),
		`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`,
		username, string(hashed),
	).Scan(&id)
	require.NoError(t, err)

	return domain.User{
		ID:       id,
		Username: username,
		Password: string(hashed),
	}
}

func CreateTestPost(t *testing.T, db *sqlx.DB, name, description string, creatorID int, imageName string) int {
	t.Helper()
	var id int
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO posts (name, description, creator_id, image_name, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		name, description, creatorID, imageName, time.Now().UTC(),
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func CreateTestComment(t *testing.T, db *sqlx.DB, text string, postID, creatorID int) int {
	t.Helper()
	var id int
	err := db.QueryRowContext(context.Background(),
		`INSERT INTO comments (text, post_id, creator_id) VALUES ($1, $2, $3) RETURNING id`,
		text, postID, creatorID,
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func CreateTestLike(t *testing.T, db *sqlx.DB, userID, postID int) {
	t.Helper()
	_, err := db.ExecContext(context.Background(),
		`INSERT INTO post_likes (liked_user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, postID,
	)
	require.NoError(t, err)
}
