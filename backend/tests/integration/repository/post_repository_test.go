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

func TestPostRepository(t *testing.T) {
	db := helpers.NewTestDB(t)
	helpers.TruncateAllTables(t, db)
	repo := repository.NewPostRepository(db)
	ctx := context.Background()

	t.Run("Create and GetByID", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "poster", "password1234")

		id, err := repo.Create(ctx, dto.PostDTO{
			PostName:        "My first post",
			PostDescription: "description here",
			CreatorID:       user.ID,
			ImageName:       "img.jpg",
		})
		require.Nil(t, err)
		assert.Greater(t, id, 0)

		post, err := repo.GetByID(ctx, id)
		require.Nil(t, err)
		require.NotNil(t, post)
		assert.Equal(t, "My first post", post.PostName)
		assert.Equal(t, "description here", post.PostDescription)
		assert.Equal(t, user.ID, post.CreatorID)
		assert.Equal(t, "img.jpg", post.ImageName)
		assert.Equal(t, 0, post.Likes)
	})

	t.Run("GetByID not found", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)

		post, err := repo.GetByID(ctx, 99999)
		assert.Nil(t, post)
		require.NotNil(t, err)
		assert.Equal(t, constants.NotFound, err.Code)
	})

	t.Run("GetAllByID", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "owner", "password1234")
		helpers.CreateTestPost(t, db, "Post A", "desc A", user.ID, "")
		helpers.CreateTestPost(t, db, "Post B", "desc B", user.ID, "pic.png")

		posts, err := repo.GetAllByID(ctx, user.ID)
		require.Nil(t, err)
		assert.Len(t, posts, 2)
		assert.Equal(t, "owner", posts[0].CreatorUsername)
	})

	t.Run("GetAllByUsername", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "nameduser", "password1234")
		helpers.CreateTestPost(t, db, "Named post", "desc", user.ID, "")

		posts, err := repo.GetAllByUsername(ctx, "nameduser")
		require.Nil(t, err)
		assert.Len(t, posts, 1)
		assert.Equal(t, "Named post", posts[0].PostName)
	})

	t.Run("GetAllByUsername empty result does not error", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)

		posts, err := repo.GetAllByUsername(ctx, "nobody")
		require.Nil(t, err)
		assert.Empty(t, posts)
	})

	t.Run("IncrementLikes and DecrementLikes", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "liker", "password1234")
		postID := helpers.CreateTestPost(t, db, "Likeable", "desc", user.ID, "")

		likes, err := repo.IncrementLikes(ctx, postID)
		require.Nil(t, err)
		assert.Equal(t, 1, likes)

		likes, err = repo.IncrementLikes(ctx, postID)
		require.Nil(t, err)
		assert.Equal(t, 2, likes)

		likes, err = repo.DecrementLikes(ctx, postID)
		require.Nil(t, err)
		assert.Equal(t, 1, likes)
	})

	t.Run("DeletePost success and not found", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "deleter", "password1234")
		postID := helpers.CreateTestPost(t, db, "To delete", "desc", user.ID, "img.jpg")

		err := repo.DeletePost(ctx, postID, user.ID)
		require.Nil(t, err)

		err = repo.DeletePost(ctx, postID, user.ID)
		require.NotNil(t, err)
		assert.Equal(t, constants.NotFound, err.Code)

		postID2 := helpers.CreateTestPost(t, db, "Another", "desc", user.ID, "")
		err = repo.DeletePost(ctx, postID2, 999)
		require.NotNil(t, err)
		assert.Equal(t, constants.NotFound, err.Code)
	})

	t.Run("GetPostFeed returns posts", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "feeder", "password1234")
		for i := 0; i < 3; i++ {
			helpers.CreateTestPost(t, db, "Feed post", "desc", user.ID, "")
		}

		posts, err := repo.GetPostFeed(ctx)
		require.Nil(t, err)
		assert.GreaterOrEqual(t, len(posts), 3)
	})

	t.Run("GetImageName", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "imager", "password1234")
		postID := helpers.CreateTestPost(t, db, "With image", "desc", user.ID, "photo.png")

		name, err := repo.GetImageName(ctx, postID)
		require.Nil(t, err)
		assert.Equal(t, "photo.png", name)

		name, err = repo.GetImageName(ctx, 99999)
		assert.Empty(t, name)
		require.NotNil(t, err)
		assert.Equal(t, constants.NotFound, err.Code)
	})
}
