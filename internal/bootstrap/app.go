package bootstrap

import (
	"github.com/seth16888/wxtoken/internal/di"
	"github.com/seth16888/wxtoken/internal/server"
)

func StartApp() error {
  return server.Start(di.DI)
}
