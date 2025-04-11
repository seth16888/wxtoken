package biz

import (
	"context"

	"github.com/seth16888/wxtoken/internal/entities"
)

type AppRepo interface {
	GetSecret(context.Context, uint64) (*entities.PlatformAppEntity, error)
}
