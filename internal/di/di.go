package di

import (
	"github.com/seth16888/wxcommon/hc"
	"github.com/seth16888/wxcommon/logger"
	"github.com/seth16888/wxcommon/redis"
	"github.com/seth16888/wxtoken/internal/biz"
	"github.com/seth16888/wxtoken/internal/cache"
	"github.com/seth16888/wxtoken/internal/config"
	"github.com/seth16888/wxtoken/internal/data"
	"github.com/seth16888/wxtoken/internal/service"
	"go.uber.org/zap"
)

var DI *Container

type Container struct {
	Conf *config.Bootstrap
  DB                 *data.Data   // 数据库连接
	Log  *zap.Logger
	Svc  *service.WXTokenService
  Redis *redis.RedisClient
}

func NewContainer(configFile string) *Container {
  conf:= config.ReadConfigFromFile(configFile)
  log := logger.InitLogger(conf.Log)

  db := data.NewData(conf.Database, log)

  redis.ConnectRedis(conf.Redis.Addr, conf.Redis.Username,
    conf.Redis.Password, conf.Redis.DB, log)

  cache.CacheRepo = cache.NewRedisCache()
  cache.CacheManage = cache.NewGoCache(redis.Redis.Client)

  hc := hc.NewClient(hc.DefaultTimeout, hc.DefaultIdleConnTimeout, hc.CommonCheckRedirect)
  appRepo := data.NewAppRepo(db, log)
  repo := data.NewAccessTokenRepo(db, log)
  cache := cache.NewRedisCache()
  uc := biz.NewTokenUsecase(repo, log, cache, appRepo, hc)

  svc := service.NewWXTokenService(log, uc)

	DI = &Container{
    Conf: conf,
    DB: db,
    Log: log,
    Svc: svc,
    Redis: redis.Redis,
  }
	return DI
}
