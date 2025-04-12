package entities

import "gorm.io/gorm"

// PlatformAppEntity 公众号平台账号
type PlatformAppEntity struct {
	Id             string         `json:"id" gorm:"primaryKey"`
	UserId         uint           `json:"userId" gorm:"not null,comment:用户id"`
	Name           string         `json:"name" gorm:"not null,comment:公众号名称"`
	PicUrl         string         `json:"picUrl" gorm:"not null,comment:公众号图标"`
	Type           int            `json:"type" gorm:"not null,comment:公众号类型"`
	Introduction   string         `json:"introduction" gorm:"comment:介绍"`
	Token          string         `json:"token" gorm:"not null,comment:接口验证token"`
	EncodingAesKey string         `json:"encodingAESKey" gorm:"not null,comment:消息加密密钥EncodingAESKey"`
	EncodingType   int            `json:"encodingType" gorm:"not null,comment:消息加密方式:1-明文,2-兼容,3-安全"`
	AppId          string         `json:"appId" gorm:"not null,comment:公众号appid"`
	AppSecret      string         `json:"appSecret" gorm:"not null,comment:公众号appsecret"`
	Status         int            `json:"status" gorm:"not null,default:0,comment:状态:0-未验证,1-验证失败,2-验证成功,3-接入成功"`
	CreatedAt      int64          `json:"createdAt" gorm:"autoCreateTime:milli"`
	UpdatedAt      int64          `json:"updatedAt" gorm:"autoUpdateTime:milli"`
	DeletedAt      gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Version        int64          `json:"-"`
}

func (*PlatformAppEntity) TableName() string {
	return "platform_app"
}
