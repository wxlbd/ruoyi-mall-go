package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemOAuth2Client OAuth2 客户端
type SystemOAuth2Client struct {
	ID                          int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ClientID                    string `gorm:"column:client_id;not null;comment:客户端编号" json:"clientId"`
	ClientSecret                string `gorm:"column:client_secret;not null;comment:客户端密钥" json:"clientSecret"`
	Name                        string `gorm:"column:name;not null;comment:应用名" json:"name"`
	Logo                        string `gorm:"column:logo;comment:应用图标" json:"logo"`
	Description                 string `gorm:"column:description;comment:应用描述" json:"description"`
	Status                      int    `gorm:"column:status;not null;default:0;comment:状态" json:"status"`
	AccessTokenValiditySeconds  int    `gorm:"column:access_token_validity_seconds;not null;comment:访问令牌的有效期" json:"accessTokenValiditySeconds"`
	RefreshTokenValiditySeconds int    `gorm:"column:refresh_token_validity_seconds;not null;comment:刷新令牌的有效期" json:"refreshTokenValiditySeconds"`
	RedirectUris                string `gorm:"column:redirect_uris;comment:可重定向的 URI 地址" json:"redirectUris"`                   // JSON array
	AuthorizedGrantTypes        string `gorm:"column:authorized_grant_types;not null;comment:授权类型" json:"authorizedGrantTypes"` // JSON array
	Scopes                      string `gorm:"column:scopes;not null;comment:授权范围" json:"scopes"`                               // JSON array
	AutoApproveScopes           string `gorm:"column:auto_approve_scopes;comment:自动授权范围" json:"autoApproveScopes"`              // JSON array
	Authorities                 string `gorm:"column:authorities;comment:权限" json:"authorities"`                                // JSON array
	ResourceIDs                 string `gorm:"column:resource_ids;comment:资源" json:"resourceIds"`                               // JSON array
	AdditionalInformation       string `gorm:"column:additional_information;comment:附加信息" json:"additionalInformation"`         // JSON string

	Creator   string         `gorm:"column:creator;size:64;default:'';comment:创建者" json:"creator"`
	Updater   string         `gorm:"column:updater;size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_time;index;comment:删除时间" json:"-"`
	Deleted   BitBool        `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"-"`
}

func (SystemOAuth2Client) TableName() string {
	return "system_oauth2_client"
}
