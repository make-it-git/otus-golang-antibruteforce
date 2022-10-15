package leakybucket

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ttl = 10 * time.Second

func TestBucket(t *testing.T) {
	t.Run("too_many_login", func(t *testing.T) {
		ctx := context.Background()
		lb := NewLeakyBucket(ctx, 10, 20, 30, ttl)
		for i := int64(0); i < lb.maxLoadLogin; i++ {
			err := lb.Try("login", "password", "127.0.0.1")
			assert.NoError(t, err)
		}
		err := lb.Try("login", "password", "127.0.0.1")
		assert.Equal(t, ErrBlocked, err)
	})
	t.Run("too_many_password", func(t *testing.T) {
		ctx := context.Background()
		lb := NewLeakyBucket(ctx, 20, 10, 30, ttl)
		for i := int64(0); i < lb.maxLoadPassword; i++ {
			err := lb.Try("login", "password", "127.0.0.1")
			assert.NoError(t, err)
		}
		err := lb.Try("login", "password", "127.0.0.1")
		assert.Equal(t, ErrBlocked, err)
	})
	t.Run("too_many_ip", func(t *testing.T) {
		ctx := context.Background()
		lb := NewLeakyBucket(ctx, 30, 20, 10, ttl)
		for i := int64(0); i < lb.maxLoadIP; i++ {
			err := lb.Try("login", "password", "127.0.0.1")
			assert.NoError(t, err)
		}
		err := lb.Try("login", "password", "127.0.0.1")
		assert.Equal(t, ErrBlocked, err)
	})
	t.Run("without_blocking", func(t *testing.T) {
		ctx := context.Background()
		lb := NewLeakyBucket(ctx, 10, 10, 10, ttl)
		for i := int64(0); i < 10; i++ {
			err := lb.Try("login", "password", "127.0.0.1")
			assert.NoError(t, err)
		}
	})
}
