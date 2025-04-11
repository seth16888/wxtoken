package cache

import (
	"context"

	gocache "github.com/eko/gocache/lib/v4/cache"

	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
)

var (
	CacheManage *gocache.Cache[interface{}]
	CacheRepo    Cache
)

type RedisCache struct {
}

func NewGoCache(redisCli *redis.Client) *gocache.Cache[interface{}] {
	store := redis_store.NewRedis(redisCli)

	return gocache.New[interface{}](store)
}

func NewRedisCache() Cache {
	return &RedisCache{}
}

func (m *RedisCache) Get(key string) (interface{}, error) {
	ctx := context.Background()
	return CacheManage.Get(ctx, key)
}

func (m *RedisCache) Set(key string, value interface{}) error {
	ctx := context.Background()
	return CacheManage.Set(ctx, key, value)
}

func (m *RedisCache) Delete(key string) error {
	return CacheManage.Delete(context.Background(), key)
}
