package data

import (
	"context"
	"time"

	"github.com/seth16888/wxtoken/internal/biz"
	"github.com/seth16888/wxtoken/internal/entities"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccessTokenRepo struct {
	data *gorm.DB
	log  *zap.Logger
}

func NewAccessTokenRepo(data *gorm.DB, logger *zap.Logger) biz.AccessTokenRepo {
	return &AccessTokenRepo{
		data: data,
		log:  logger,
	}
}

func (p *AccessTokenRepo) Get(ctx context.Context, mpId string) (*entities.AccessToken, error) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	token := &entities.AccessToken{}
	if err := p.data.WithContext(c).Where("mp_id = ?", mpId).First(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func (p *AccessTokenRepo) Save(ctx context.Context, entity *entities.AccessToken) (uint, error) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx := p.data.WithContext(c).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}, {Name: "app_id"}},
			UpdateAll: true,
		},
	).Create(entity)

	if tx.Error != nil {
		return 0, tx.Error
	}

	return entity.Id, nil
}

func (p *AccessTokenRepo) GetNeedRefresh(ctx context.Context, deadline int64, limit int) ([]*entities.AccessToken, error) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var rows []*entities.AccessToken
	tx := p.data.WithContext(c).Where("deadline <= ?", deadline).Limit(limit).Order("deadline ASC").Find(&rows)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return rows, nil
}
