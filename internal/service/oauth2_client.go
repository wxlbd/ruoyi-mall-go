package service

import (
	"context"
	"encoding/json"

	"backend-go/internal/api/req"
	"backend-go/internal/model"
	"backend-go/internal/pkg/core"

	"gorm.io/gorm"
)

type OAuth2ClientService struct {
	db *gorm.DB
}

func NewOAuth2ClientService(db *gorm.DB) *OAuth2ClientService {
	return &OAuth2ClientService{db: db}
}

func (s *OAuth2ClientService) CreateOAuth2Client(ctx context.Context, r *req.OAuth2ClientCreateReq) (int64, error) {
	RedirectUris, _ := json.Marshal(r.RedirectUris)
	AuthorizedGrantTypes, _ := json.Marshal(r.AuthorizedGrantTypes)
	Scopes, _ := json.Marshal(r.Scopes)
	AutoApproveScopes, _ := json.Marshal(r.AutoApproveScopes)
	Authorities, _ := json.Marshal(r.Authorities)
	ResourceIDs, _ := json.Marshal(r.ResourceIDs)

	c := &model.SystemOAuth2Client{
		ClientID:                    r.ClientID,
		ClientSecret:                r.ClientSecret,
		Name:                        r.Name,
		Logo:                        r.Logo,
		Description:                 r.Description,
		Status:                      r.Status,
		AccessTokenValiditySeconds:  r.AccessTokenValiditySeconds,
		RefreshTokenValiditySeconds: r.RefreshTokenValiditySeconds,
		RedirectUris:                string(RedirectUris),
		AuthorizedGrantTypes:        string(AuthorizedGrantTypes),
		Scopes:                      string(Scopes),
		AutoApproveScopes:           string(AutoApproveScopes),
		Authorities:                 string(Authorities),
		ResourceIDs:                 string(ResourceIDs),
		AdditionalInformation:       r.AdditionalInformation,
	}

	if err := s.db.WithContext(ctx).Create(c).Error; err != nil {
		return 0, err
	}
	return c.ID, nil
}

func (s *OAuth2ClientService) UpdateOAuth2Client(ctx context.Context, r *req.OAuth2ClientUpdateReq) error {
	RedirectUris, _ := json.Marshal(r.RedirectUris)
	AuthorizedGrantTypes, _ := json.Marshal(r.AuthorizedGrantTypes)
	Scopes, _ := json.Marshal(r.Scopes)
	AutoApproveScopes, _ := json.Marshal(r.AutoApproveScopes)
	Authorities, _ := json.Marshal(r.Authorities)
	ResourceIDs, _ := json.Marshal(r.ResourceIDs)

	c := &model.SystemOAuth2Client{
		ID:                          r.ID,
		ClientID:                    r.ClientID,
		ClientSecret:                r.ClientSecret,
		Name:                        r.Name,
		Logo:                        r.Logo,
		Description:                 r.Description,
		Status:                      r.Status,
		AccessTokenValiditySeconds:  r.AccessTokenValiditySeconds,
		RefreshTokenValiditySeconds: r.RefreshTokenValiditySeconds,
		RedirectUris:                string(RedirectUris),
		AuthorizedGrantTypes:        string(AuthorizedGrantTypes),
		Scopes:                      string(Scopes),
		AutoApproveScopes:           string(AutoApproveScopes),
		Authorities:                 string(Authorities),
		ResourceIDs:                 string(ResourceIDs),
		AdditionalInformation:       r.AdditionalInformation,
	}

	return s.db.WithContext(ctx).Updates(c).Error
}

func (s *OAuth2ClientService) DeleteOAuth2Client(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Delete(&model.SystemOAuth2Client{}, id).Error
}

func (s *OAuth2ClientService) GetOAuth2Client(ctx context.Context, id int64) (*model.SystemOAuth2Client, error) {
	var c model.SystemOAuth2Client
	if err := s.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *OAuth2ClientService) GetOAuth2ClientPage(ctx context.Context, r *req.OAuth2ClientPageReq) (*core.PageResult[*model.SystemOAuth2Client], error) {
	db := s.db.WithContext(ctx).Model(&model.SystemOAuth2Client{})
	if r.Name != "" {
		db = db.Where("name LIKE ?", "%"+r.Name+"%")
	}
	if r.Status != nil {
		db = db.Where("status = ?", *r.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	var list []*model.SystemOAuth2Client
	offset := (r.PageNo - 1) * r.PageSize
	if err := db.Order("id desc").Offset(offset).Limit(r.PageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return &core.PageResult[*model.SystemOAuth2Client]{List: list, Total: total}, nil
}
