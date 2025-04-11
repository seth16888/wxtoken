package biz

import (
	"github.com/google/wire"
	"github.com/seth16888/wxcommon/hc"
	"github.com/seth16888/wxcommon/redis"
	"github.com/seth16888/wxtoken/internal/cache"
	"github.com/seth16888/wxtoken/internal/config"
	"go.uber.org/zap"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewTokenUsecase, NewHttpClient, NewRedisClient, NewCacheRepo)

func NewHttpClient() *hc.Client {
	return hc.NewClient(hc.DefaultTimeout, hc.DefaultIdleConnTimeout, hc.CommonCheckRedirect)
}

func NewRedisClient(conf *config.Bootstrap, logger *zap.Logger) *redis.RedisClient {
	redis.ConnectRedis(conf.Redis.Addr, conf.Redis.Username, conf.Redis.Password, int(conf.Redis.DB), logger)
	return redis.Redis
}

func NewCacheRepo(redis *redis.RedisClient, logger *zap.Logger) cache.Cache {
	cache.CacheRepo = cache.NewRedisCache()
	cache.CacheManage = cache.NewGoCache(redis.Client)
	return cache.CacheRepo
}
