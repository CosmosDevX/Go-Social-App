package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myapp/internal/constants"
	"myapp/internal/repository"
	"myapp/tests/helpers"
)

func TestRefreshTokenRepository(t *testing.T) {
	client := helpers.NewTestRedis(t)
	helpers.FlushRedis(t, client)
	repo := repository.NewRefreshTokenRepository(client)
	ctx := context.Background()

	t.Run("Set Get Delete cycle", func(t *testing.T) {
		helpers.FlushRedis(t, client)

		err := repo.Set("42", "token-value", constants.TokenWhiteListPrefix, time.Minute, ctx)
		require.Nil(t, err)

		val, err := repo.Get("42", constants.TokenWhiteListPrefix, ctx)
		require.Nil(t, err)
		assert.Equal(t, "token-value", val)

		err = repo.Delete("42", constants.TokenWhiteListPrefix, ctx)
		require.Nil(t, err)

		val, err = repo.Get("42", constants.TokenWhiteListPrefix, ctx)
		require.Nil(t, err)
		assert.Empty(t, val)
	})

	t.Run("Get missing key returns empty without error", func(t *testing.T) {
		helpers.FlushRedis(t, client)

		val, err := repo.Get("999", constants.TokenWhiteListPrefix, ctx)
		require.Nil(t, err)
		assert.Empty(t, val)
	})

	t.Run("Blacklist and whitelist are independent keys", func(t *testing.T) {
		helpers.FlushRedis(t, client)

		require.Nil(t, repo.Set("7", "white", constants.TokenWhiteListPrefix, time.Minute, ctx))

		w, err := repo.Get("7", constants.TokenWhiteListPrefix, ctx)
		require.Nil(t, err)
		assert.Equal(t, "white", w)
	})
}
