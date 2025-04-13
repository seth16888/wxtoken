package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type AccessToken struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AppId       string             `bson:"app_id" json:"app_id"`
	MpId        string             `bson:"mp_id" json:"mp_id"`
	AccessToken string             `bson:"access_token" json:"access_token"`
	Deadline    int64              `bson:"deadline" json:"deadline"`
	ExpiresIn   uint64             `bson:"expires_in" json:"expires_in"`
	CreatedAt   int64              `bson:"created_at" json:"created_at"`
	UpdatedAt   int64              `bson:"updated_at" json:"updated_at"`
	DeletedAt   int64              `bson:"deleted_at" json:"deleted_at"`
	Version     int64              `bson:"version" json:"-"`
}
