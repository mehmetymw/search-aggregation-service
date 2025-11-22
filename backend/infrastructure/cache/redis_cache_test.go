package cache

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestRedis(t *testing.T) *RedisCache {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR not set, skipping integration test")
	}

	client, err := NewRedisCache(addr)
	if err != nil {
		t.Fatalf("failed to create redis client: %v", err)
	}

	return client.(*RedisCache)
}

func TestRedisCache_SetAndGet(t *testing.T) {
	cache := setupTestRedis(t)
	defer cache.Close()

	ctx := context.Background()
	key := "test_key"
	value := map[string]string{"foo": "bar"}

	err := cache.Set(ctx, key, value, time.Minute)
	assert.NoError(t, err)

	var dest map[string]string
	found, err := cache.Get(ctx, key, &dest)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "bar", dest["foo"])
}

func TestRedisCache_Miss(t *testing.T) {
	cache := setupTestRedis(t)
	defer cache.Close()

	ctx := context.Background()
	key := "non_existent_key"

	var dest string
	found, err := cache.Get(ctx, key, &dest)
	assert.NoError(t, err)
	assert.False(t, found)
}
