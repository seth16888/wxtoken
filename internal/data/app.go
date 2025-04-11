package data

import (
	"context"
	"time"

	"github.com/seth16888/wxtoken/internal/biz"
	"github.com/seth16888/wxtoken/internal/entities"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppRepo struct {
	data *gorm.DB
	log  *zap.Logger
}

// NewAppRepo .
func NewAppRepo(data *gorm.DB, logger *zap.Logger) biz.AppRepo {
	return &AppRepo{
		data: data,
		log:  logger,
	}
}

func (p *AppRepo) GetSecret(ctx context.Context, id uint64) (*entities.PlatformAppEntity, error) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	app := &entities.PlatformAppEntity{}
	if err := p.data.WithContext(c).Where("id = ?", id).First(app).Error; err != nil {
		return nil, err
	}

	return app, nil
}
