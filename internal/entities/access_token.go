package entities

import "gorm.io/gorm"

type AccessToken struct {
	Id          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	AppId       string         `gorm:"column:app_id"`
	MpId        string         `gorm:"column:mp_id"`
	AccessToken string         `gorm:"column:access_token;size:512"`
	Deadline    int64          `gorm:"column:deadline"`
	ExpiresIn   uint64         `gorm:"column:expires_in"`
	CreatedAt   int64          `json:"createdAt" gorm:"autoCreateTime:milli"`
	UpdatedAt   int64          `json:"updatedAt" gorm:"autoUpdateTime:milli"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Version     int64          `json:"-"`
}

func (AccessToken) TableName() string {
	return "mp_access_token"
}
