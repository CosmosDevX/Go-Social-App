package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/repository"
	"myapp/tests/helpers"
)

func TestCommentRepository(t *testing.T) {
	db := helpers.NewTestDB(t)
	helpers.TruncateAllTables(t, db)
	repo := repository.NewCommentRepository(db)
	ctx := context.Background()

	t.Run("Create and GetAllByPostID", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "commenter", "password1234")
		postID := helpers.CreateTestPost(t, db, "Post", "desc", user.ID, "")

		id, err := repo.Create(ctx, dto.CommentDTO{
			CommentText: "Nice post!",
			PostID:      postID,
			CreatorID:   user.ID,
		})
		require.Nil(t, err)
		assert.Greater(t, id, 0)

		comments, err := repo.GetAllByPostID(ctx, postID)
		require.Nil(t, err)
		require.Len(t, comments, 1)
		assert.Equal(t, "Nice post!", comments[0].CommentText)
		assert.Equal(t, "commenter", comments[0].CreatorUsername)
		assert.Equal(t, user.ID, comments[0].CreatorID)
	})

	t.Run("GetAllByPostID not found", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)

		comments, err := repo.GetAllByPostID(ctx, 99999)
		assert.Nil(t, comments)
		require.NotNil(t, err)
		assert.Equal(t, constants.NotFound, err.Code)
	})

	t.Run("CountCommentsOnPost", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "counter", "password1234")
		postID := helpers.CreateTestPost(t, db, "Post", "desc", user.ID, "")
		helpers.CreateTestComment(t, db, "c1", postID, user.ID)
		helpers.CreateTestComment(t, db, "c2", postID, user.ID)

		count, err := repo.CountCommentsOnPost(ctx, postID)
		require.Nil(t, err)
		assert.Equal(t, 2, count)

		count, err = repo.CountCommentsOnPost(ctx, 99999)
		require.Nil(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Delete success and not found", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "delcomment", "password1234")
		postID := helpers.CreateTestPost(t, db, "Post", "desc", user.ID, "")
		commentID := helpers.CreateTestComment(t, db, "to delete", postID, user.ID)

		err := repo.Delete(ctx, commentID, user.ID)
		require.Nil(t, err)

		err = repo.Delete(ctx, commentID, user.ID)
		require.NotNil(t, err)
		assert.Equal(t, constants.NotFound, err.Code)

		commentID2 := helpers.CreateTestComment(t, db, "another", postID, user.ID)
		err = repo.Delete(ctx, commentID2, 999)
		require.NotNil(t, err)
		assert.Equal(t, constants.NotFound, err.Code)
	})
}
