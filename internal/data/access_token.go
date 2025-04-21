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
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type AccessTokenRepo struct {
	col  *mongo.Collection
	data *Data
	log  *zap.Logger
}

func NewAccessTokenRepo(data *Data, logger *zap.Logger) biz.AccessTokenRepo {
	collection := data.db.Collection("mp_access_token")
	return &AccessTokenRepo{
		col:  collection,
		data: data,
		log:  logger,
	}
}

// Get 根据mpId获取访问令牌
func (p *AccessTokenRepo) Get(ctx context.Context, mpId string) (*entities.AccessToken, error) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"mp_id": mpId}
	token := &entities.AccessToken{}

	err := p.col.FindOne(c, filter).Decode(token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("data not found")
		}
		return nil, err
	}

	return token, nil
}

// Save 保存或更新访问令牌
func (p *AccessTokenRepo) Save(ctx context.Context, entity *entities.AccessToken) (string, error) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"mp_id": entity.MpId}
	update := bson.M{"$set": entity}
	opts := options.Update().SetUpsert(true)

	result, err := p.col.UpdateOne(c, filter, update, opts)
	if err != nil {
    p.log.Error("Save", zap.Error(err))
		return "", err
	}

	if result.UpsertedID != nil {
		if oid, ok := result.UpsertedID.(primitive.ObjectID); ok {
			return oid.Hex(), nil
		}
	}

	return entity.Id.Hex(), nil
}

// GetNeedRefresh 获取需要刷新的访问令牌列表
func (p *AccessTokenRepo) GetNeedRefresh(ctx context.Context, deadline int64, limit int) ([]*entities.AccessToken, error) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	opts := options.Find().
		SetSort(bson.M{"deadline": 1}).
		SetLimit(int64(limit))

	filter := bson.M{"deadline": bson.M{"$lte": deadline}}

	cursor, err := p.col.Find(c, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var tokens []*entities.AccessToken
	if err = cursor.All(c, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}
