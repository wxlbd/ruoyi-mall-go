package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"backend-go/internal/api/req"
	"backend-go/internal/model"
	"backend-go/internal/pkg/core"

	"gorm.io/gorm"
)

type NotifyService struct {
	db *gorm.DB
	// Cache
	templateCache map[string]*model.SystemNotifyTemplate
	mu            sync.RWMutex
}

func NewNotifyService(db *gorm.DB) *NotifyService {
	s := &NotifyService{
		db: db,
	}
	s.RefreshCache()
	return s
}

func (s *NotifyService) RefreshCache() {
	var list []model.SystemNotifyTemplate
	s.db.Find(&list)
	m := make(map[string]*model.SystemNotifyTemplate)
	for i := range list {
		m[list[i].Code] = &list[i]
	}
	s.mu.Lock()
	s.templateCache = m
	s.mu.Unlock()
}

// ================= Template CRUD =================

func (s *NotifyService) CreateNotifyTemplate(ctx context.Context, r *req.NotifyTemplateCreateReq) (int64, error) {
	t := &model.SystemNotifyTemplate{
		Name:     r.Name,
		Code:     r.Code,
		Nickname: r.Nickname,
		Content:  r.Content,
		Type:     r.Type,
		Status:   r.Status,
		Remark:   r.Remark,
	}
	if err := s.db.WithContext(ctx).Create(t).Error; err != nil {
		return 0, err
	}
	s.RefreshCache()
	return t.ID, nil
}

func (s *NotifyService) UpdateNotifyTemplate(ctx context.Context, r *req.NotifyTemplateUpdateReq) error {
	t := &model.SystemNotifyTemplate{
		ID:       r.ID,
		Name:     r.Name,
		Code:     r.Code,
		Nickname: r.Nickname,
		Content:  r.Content,
		Type:     r.Type,
		Status:   r.Status,
		Remark:   r.Remark,
	}
	if err := s.db.WithContext(ctx).Updates(t).Error; err != nil {
		return err
	}
	s.RefreshCache()
	return nil
}

func (s *NotifyService) DeleteNotifyTemplate(ctx context.Context, id int64) error {
	if err := s.db.WithContext(ctx).Delete(&model.SystemNotifyTemplate{}, id).Error; err != nil {
		return err
	}
	s.RefreshCache()
	return nil
}

func (s *NotifyService) GetNotifyTemplate(ctx context.Context, id int64) (*model.SystemNotifyTemplate, error) {
	var t model.SystemNotifyTemplate
	if err := s.db.WithContext(ctx).First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *NotifyService) GetNotifyTemplatePage(ctx context.Context, r *req.NotifyTemplatePageReq) (*core.PageResult[*model.SystemNotifyTemplate], error) {
	db := s.db.WithContext(ctx).Model(&model.SystemNotifyTemplate{})
	if r.Name != "" {
		db = db.Where("name LIKE ?", "%"+r.Name+"%")
	}
	if r.Code != "" {
		db = db.Where("code LIKE ?", "%"+r.Code+"%")
	}
	if r.Status != nil {
		db = db.Where("status = ?", *r.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	var list []*model.SystemNotifyTemplate
	offset := (r.PageNo - 1) * r.PageSize
	if err := db.Order("id desc").Offset(offset).Limit(r.PageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return &core.PageResult[*model.SystemNotifyTemplate]{List: list, Total: total}, nil
}

// ================= Message Logic =================

func (s *NotifyService) SendNotify(ctx context.Context, userID int64, userType int, templateCode string, params map[string]interface{}) (int64, error) {
	s.mu.RLock()
	template, ok := s.templateCache[templateCode]
	s.mu.RUnlock()
	if !ok || template == nil {
		return 0, core.NewBizError(1002006001, "站内信模板不存在")
	}

	content := template.Content
	for k, v := range params {
		content = strings.ReplaceAll(content, "{"+k+"}", fmt.Sprintf("%v", v))
	}

	paramsStr, _ := json.Marshal(params)
	msg := &model.SystemNotifyMessage{
		UserID:           userID,
		UserType:         userType,
		TemplateID:       template.ID,
		TemplateCode:     template.Code,
		TemplateNickname: template.Nickname,
		TemplateContent:  content,
		TemplateType:     template.Type,
		TemplateParams:   string(paramsStr),
		ReadStatus:       false,
	}
	if err := s.db.WithContext(ctx).Create(msg).Error; err != nil {
		return 0, err
	}
	return msg.ID, nil
}

func (s *NotifyService) GetNotifyMessagePage(ctx context.Context, r *req.NotifyMessagePageReq) (*core.PageResult[*model.SystemNotifyMessage], error) {
	db := s.db.WithContext(ctx).Model(&model.SystemNotifyMessage{})
	if r.UserID != 0 {
		db = db.Where("user_id = ?", r.UserID)
	}
	if r.UserType != 0 {
		db = db.Where("user_type = ?", r.UserType)
	}
	if r.TemplateCode != "" {
		db = db.Where("template_code LIKE ?", "%"+r.TemplateCode+"%")
	}
	if r.TemplateType != nil {
		db = db.Where("template_type = ?", *r.TemplateType)
	}
	if r.ReadStatus != nil {
		db = db.Where("read_status = ?", *r.ReadStatus)
	}
	if r.StartDate != "" && r.EndDate != "" {
		db = db.Where("create_time BETWEEN ? AND ?", r.StartDate, r.EndDate)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	var list []*model.SystemNotifyMessage
	offset := (r.PageNo - 1) * r.PageSize
	if err := db.Order("id desc").Offset(offset).Limit(r.PageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return &core.PageResult[*model.SystemNotifyMessage]{List: list, Total: total}, nil
}

func (s *NotifyService) GetMyNotifyMessagePage(ctx context.Context, userID int64, userType int, r *req.MyNotifyMessagePageReq) (*core.PageResult[*model.SystemNotifyMessage], error) {
	db := s.db.WithContext(ctx).Model(&model.SystemNotifyMessage{}).Where("user_id = ? AND user_type = ?", userID, userType)
	if r.ReadStatus != nil {
		db = db.Where("read_status = ?", *r.ReadStatus)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	var list []*model.SystemNotifyMessage
	offset := (r.PageNo - 1) * r.PageSize
	if err := db.Order("id desc").Offset(offset).Limit(r.PageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return &core.PageResult[*model.SystemNotifyMessage]{List: list, Total: total}, nil
}

func (s *NotifyService) UpdateNotifyMessageRead(ctx context.Context, userID int64, userType int, ids []int64) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&model.SystemNotifyMessage{}).
		Where("id IN ? AND user_id = ? AND user_type = ?", ids, userID, userType).
		Updates(map[string]interface{}{
			"read_status": true,
			"read_time":   &now,
		}).Error
}

func (s *NotifyService) UpdateAllNotifyMessageRead(ctx context.Context, userID int64, userType int) error {
	now := time.Now()
	return s.db.WithContext(ctx).Model(&model.SystemNotifyMessage{}).
		Where("user_id = ? AND user_type = ? AND read_status = ?", userID, userType, false).
		Updates(map[string]interface{}{
			"read_status": true,
			"read_time":   &now,
		}).Error
}

func (s *NotifyService) GetUnreadNotifyMessageCount(ctx context.Context, userID int64, userType int) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&model.SystemNotifyMessage{}).
		Where("user_id = ? AND user_type = ? AND read_status = ?", userID, userType, false).
		Count(&count).Error
	return count, err
}
