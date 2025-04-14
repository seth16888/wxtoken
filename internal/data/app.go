package data

import (
	"context"
	"fmt"
	"time"

	"github.com/seth16888/wxtoken/internal/biz"
	"github.com/seth16888/wxtoken/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AppRepo struct {
	col  *mongo.Collection
	data *Data
	log  *zap.Logger
	db   *mongo.Database
}

func NewAppRepo(data *Data, logger *zap.Logger) biz.AppRepo {
	appDB := data.client.Database("wxbusiness")
	collection := appDB.Collection("platform_apps")
	return &AppRepo{
		data: data,
		col:  collection,
		log:  logger,
		db:   appDB,
	}
}

func (p *AppRepo) GetSecret(ctx context.Context, id string) (*entities.PlatformApp, error) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

  objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": objectID}
  // 记录查询信息
  p.log.Info("MongoDB query",
    zap.String("database", p.db.Name()),
    zap.String("collection", p.col.Name()),
    zap.Any("filter", filter),
  )

	app := &entities.PlatformApp{}
	if err := p.col.FindOne(c, filter).Decode(&app); err != nil {
		return nil, err
	}

	return app, nil
}
