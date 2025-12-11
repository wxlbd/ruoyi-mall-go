package req

import (
	"backend-go/internal/pkg/core"
)

type OAuth2ClientCreateReq struct {
	ClientID                    string   `json:"clientId" binding:"required"`
	ClientSecret                string   `json:"clientSecret" binding:"required"`
	Name                        string   `json:"name" binding:"required"`
	Logo                        string   `json:"logo"`
	Description                 string   `json:"description"`
	Status                      int      `json:"status" binding:"required"`
	AccessTokenValiditySeconds  int      `json:"accessTokenValiditySeconds" binding:"required"`
	RefreshTokenValiditySeconds int      `json:"refreshTokenValiditySeconds" binding:"required"`
	RedirectUris                []string `json:"redirectUris"`
	AuthorizedGrantTypes        []string `json:"authorizedGrantTypes" binding:"required"`
	Scopes                      []string `json:"scopes" binding:"required"`
	AutoApproveScopes           []string `json:"autoApproveScopes"`
	Authorities                 []string `json:"authorities"`
	ResourceIDs                 []string `json:"resourceIds"`
	AdditionalInformation       string   `json:"additionalInformation"`
}

type OAuth2ClientUpdateReq struct {
	ID                          int64    `json:"id" binding:"required"`
	ClientID                    string   `json:"clientId" binding:"required"`
	ClientSecret                string   `json:"clientSecret" binding:"required"`
	Name                        string   `json:"name" binding:"required"`
	Logo                        string   `json:"logo"`
	Description                 string   `json:"description"`
	Status                      int      `json:"status" binding:"required"`
	AccessTokenValiditySeconds  int      `json:"accessTokenValiditySeconds" binding:"required"`
	RefreshTokenValiditySeconds int      `json:"refreshTokenValiditySeconds" binding:"required"`
	RedirectUris                []string `json:"redirectUris"`
	AuthorizedGrantTypes        []string `json:"authorizedGrantTypes" binding:"required"`
	Scopes                      []string `json:"scopes" binding:"required"`
	AutoApproveScopes           []string `json:"autoApproveScopes"`
	Authorities                 []string `json:"authorities"`
	ResourceIDs                 []string `json:"resourceIds"`
	AdditionalInformation       string   `json:"additionalInformation"`
}

type OAuth2ClientPageReq struct {
	core.PageParam
	Name   string `form:"name"`
	Status *int   `form:"status"`
}
