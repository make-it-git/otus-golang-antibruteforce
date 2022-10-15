//go:build integration
// +build integration

package storage

import (
	"net"
	"testing"

	"github.com/make-it-git/otus-antibruteforce/internal/contract"

	"github.com/go-redis/redis"
	"github.com/make-it-git/otus-antibruteforce/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestBucket(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisDSN,
		Password: "",
		DB:       0,
	})
	storage := NewRedisStorage(redisClient)

	t.Run("redis_connection_oj", func(t *testing.T) {
		pong, err := redisClient.Ping().Result()
		assert.NoError(t, err)
		assert.Equal(t, "PONG", pong)
	})

	t.Run("whitelist", func(t *testing.T) {
		storage.ClearLists()
		defer storage.ClearLists()

		checkIp := net.ParseIP("1.2.3.10")
		checkIpNet := "1.2.3.0/24"

		status, err := storage.GetStatus(checkIp)
		assert.NoError(t, err)
		assert.Equal(t, contract.Unknown, status)

		err = storage.WhiteListAdd(checkIpNet)
		assert.NoError(t, err)

		status, err = storage.GetStatus(checkIp)
		assert.NoError(t, err)
		assert.Equal(t, contract.Whitelisted, status)

		err = storage.WhiteListRemove(checkIpNet)
		assert.NoError(t, err)

		status, err = storage.GetStatus(checkIp)
		assert.NoError(t, err)
		assert.Equal(t, contract.Unknown, status)
	})

	t.Run("blacklist", func(t *testing.T) {
		storage.ClearLists()
		defer storage.ClearLists()

		checkIp := net.ParseIP("1.2.4.20")
		checkIpNet := "1.2.4.0/24"

		status, err := storage.GetStatus(checkIp)
		assert.NoError(t, err)
		assert.Equal(t, contract.Unknown, status)

		err = storage.BlackListAdd(checkIpNet)
		assert.NoError(t, err)

		status, err = storage.GetStatus(checkIp)
		assert.NoError(t, err)
		assert.Equal(t, contract.Blacklisted, status)

		err = storage.BlackListRemove(checkIpNet)
		assert.NoError(t, err)

		status, err = storage.GetStatus(checkIp)
		assert.NoError(t, err)
		assert.Equal(t, contract.Unknown, status)
	})
}
