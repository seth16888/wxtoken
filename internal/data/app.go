package data

import (
	"context"
	"time"

	"github.com/seth16888/wxtoken/internal/biz"
	"github.com/seth16888/wxtoken/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AppRepo struct {
	col  *mongo.Collection
	data *Data
	log  *zap.Logger
}

func NewAppRepo(data *Data, logger *zap.Logger) biz.AppRepo {
	collection := data.db.Collection("platform_app")
	return &AppRepo{
		data: data,
		col:  collection,
		log:  logger,
	}
}

func (p *AppRepo) GetSecret(ctx context.Context, id string) (*entities.PlatformApp, error) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}

	app := &entities.PlatformApp{}
	if err := p.col.FindOne(c, filter).Decode(&app); err != nil {
		return nil, err
	}

	return app, nil
}
