package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myapp/internal/repository"
	"myapp/tests/helpers"
)

func TestPostLikeRepository(t *testing.T) {
	db := helpers.NewTestDB(t)
	helpers.TruncateAllTables(t, db)
	repo := repository.NewPostLikeRepository(db)
	ctx := context.Background()

	t.Run("CreateLike inserts and returns rows affected", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "liker1", "password1234")
		postID := helpers.CreateTestPost(t, db, "Post", "desc", user.ID, "")

		rows, err := repo.CreateLike(ctx, user.ID, postID)
		require.Nil(t, err)
		assert.Equal(t, 1, rows)

		rows, err = repo.CreateLike(ctx, user.ID, postID)
		require.Nil(t, err)
		assert.Equal(t, 0, rows)
	})

	t.Run("DeleteLike", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "liker2", "password1234")
		postID := helpers.CreateTestPost(t, db, "Post", "desc", user.ID, "")
		helpers.CreateTestLike(t, db, user.ID, postID)

		rows, err := repo.DeleteLike(ctx, user.ID, postID)
		require.Nil(t, err)
		assert.Equal(t, 1, rows)

		rows, err = repo.DeleteLike(ctx, user.ID, postID)
		require.Nil(t, err)
		assert.Equal(t, 0, rows)
	})

	t.Run("GetLikedPostsID", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "liker3", "password1234")
		p1 := helpers.CreateTestPost(t, db, "P1", "d", user.ID, "")
		p2 := helpers.CreateTestPost(t, db, "P2", "d", user.ID, "")
		helpers.CreateTestLike(t, db, user.ID, p1)
		helpers.CreateTestLike(t, db, user.ID, p2)

		ids, err := repo.GetLikedPostsID(ctx, user.ID)
		require.Nil(t, err)
		assert.ElementsMatch(t, []int{p1, p2}, ids)
	})

	t.Run("GetLikedPostsID empty", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "nolikes", "password1234")

		ids, err := repo.GetLikedPostsID(ctx, user.ID)
		require.Nil(t, err)
		assert.Empty(t, ids)
	})
}
