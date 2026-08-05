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

func TestUserRepository(t *testing.T) {
	db := helpers.NewTestDB(t)
	helpers.TruncateAllTables(t, db)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	t.Run("CreateUser and GetUserByName success", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)

		id, err := repo.CreateUser(ctx, dto.UserDTO{
			Username: "alice",
			Password: "hashed-pass",
		})
		require.Nil(t, err)
		assert.Greater(t, id, 0)

		user, err := repo.GetUserByName(ctx, "alice")
		require.Nil(t, err)
		require.NotNil(t, user)
		assert.Equal(t, id, user.ID)
		assert.Equal(t, "alice", user.Username)
		assert.Equal(t, "hashed-pass", user.Password)
	})

	t.Run("GetUserByName not found", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)

		user, err := repo.GetUserByName(ctx, "ghost")
		assert.Nil(t, user)
		require.NotNil(t, err)
		assert.Equal(t, constants.NotFound, err.Code)
	})

	t.Run("GetUsernameByID success", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		user := helpers.CreateTestUser(t, db, "bob", "password1234")

		name, err := repo.GetUsernameByID(ctx, user.ID)
		require.Nil(t, err)
		assert.Equal(t, "bob", name)
	})

	t.Run("GetUsernameByID not found", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)

		name, err := repo.GetUsernameByID(ctx, 99999)
		assert.Empty(t, name)
		require.NotNil(t, err)
		assert.Equal(t, constants.NotFound, err.Code)
	})

	t.Run("GetUsernameByIDs success", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)
		u1 := helpers.CreateTestUser(t, db, "u1", "password1234")
		u2 := helpers.CreateTestUser(t, db, "u2", "password1234")

		names, err := repo.GetUsernameByIDs(ctx, []int{u1.ID, u2.ID})
		require.Nil(t, err)
		assert.ElementsMatch(t, []string{"u1", "u2"}, names)
	})

	t.Run("GetUsernameByIDs empty ids", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)

		names, err := repo.GetUsernameByIDs(ctx, []int{})
		assert.Nil(t, names)
		require.NotNil(t, err)
		assert.Equal(t, constants.NotFound, err.Code)
	})

	t.Run("GetUsernameByIDs no matches", func(t *testing.T) {
		helpers.TruncateAllTables(t, db)

		names, err := repo.GetUsernameByIDs(ctx, []int{999, 1000})
		assert.Nil(t, names)
		require.NotNil(t, err)
		assert.Equal(t, constants.NotFound, err.Code)
	})
}
