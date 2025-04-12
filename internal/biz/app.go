package biz

import (
	"context"

	"github.com/seth16888/wxtoken/internal/entities"
)

type AppRepo interface {
	GetSecret(context.Context, string) (*entities.PlatformAppEntity, error)
}
