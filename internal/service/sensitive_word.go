package service

import (
	"context"
	"sync"

	"backend-go/internal/api/req"
	"backend-go/internal/model"
	"backend-go/internal/pkg/core"
	"backend-go/internal/pkg/utils"

	"gorm.io/gorm"
)

type SensitiveWordService struct {
	db   *gorm.DB
	trie *utils.SensitiveTrie
	mu   sync.RWMutex
}

func NewSensitiveWordService(db *gorm.DB) *SensitiveWordService {
	s := &SensitiveWordService{
		db:   db,
		trie: utils.NewSensitiveTrie(),
	}
	// 初始化缓存
	s.RefreshCache()
	return s
}

// RefreshCache 刷新缓存
func (s *SensitiveWordService) RefreshCache() {
	// 查询所有开启的敏感词
	// 注意: model.SystemSensitiveWord 没有 'Tags' 过滤，全量加载
	var list []model.SystemSensitiveWord
	if err := s.db.Where("status = ?", 0).Select("name").Find(&list).Error; err != nil {
		return // Log error?
	}

	newTrie := utils.NewSensitiveTrie()
	for _, w := range list {
		newTrie.AddWord(w.Name)
	}

	s.mu.Lock()
	s.trie = newTrie
	s.mu.Unlock()
}

// CreateSensitiveWord 创建敏感词
func (s *SensitiveWordService) CreateSensitiveWord(ctx context.Context, r *req.SensitiveWordCreateReq) (int64, error) {
	// 校验是否存在
	var count int64
	s.db.WithContext(ctx).Model(&model.SystemSensitiveWord{}).Where("name = ?", r.Name).Count(&count)
	if count > 0 {
		return 0, core.NewBizError(1002008001, "敏感词已存在")
	}

	word := &model.SystemSensitiveWord{
		Name:        r.Name,
		Tags:        r.Tags,
		Status:      r.Status,
		Description: r.Description,
	}
	if err := s.db.WithContext(ctx).Create(word).Error; err != nil {
		return 0, err
	}
	s.RefreshCache()
	return word.ID, nil
}

// UpdateSensitiveWord 更新敏感词
func (s *SensitiveWordService) UpdateSensitiveWord(ctx context.Context, r *req.SensitiveWordUpdateReq) error {
	// 校验是否存在
	var count int64
	s.db.WithContext(ctx).Model(&model.SystemSensitiveWord{}).Where("name = ? AND id != ?", r.Name, r.ID).Count(&count)
	if count > 0 {
		return core.NewBizError(1002008001, "敏感词已存在")
	}

	word := &model.SystemSensitiveWord{
		ID:          r.ID,
		Name:        r.Name,
		Tags:        r.Tags,
		Status:      r.Status,
		Description: r.Description,
	}
	if err := s.db.WithContext(ctx).Updates(word).Error; err != nil {
		return err
	}
	s.RefreshCache()
	return nil
}

// DeleteSensitiveWord 删除敏感词
func (s *SensitiveWordService) DeleteSensitiveWord(ctx context.Context, id int64) error {
	if err := s.db.WithContext(ctx).Delete(&model.SystemSensitiveWord{}, id).Error; err != nil {
		return err
	}
	s.RefreshCache()
	return nil
}

// GetSensitiveWord 获得敏感词
func (s *SensitiveWordService) GetSensitiveWord(ctx context.Context, id int64) (*model.SystemSensitiveWord, error) {
	var word model.SystemSensitiveWord
	if err := s.db.WithContext(ctx).First(&word, id).Error; err != nil {
		return nil, err
	}
	return &word, nil
}

// GetSensitiveWordPage 获得敏感词分页
func (s *SensitiveWordService) GetSensitiveWordPage(ctx context.Context, r *req.SensitiveWordPageReq) (*core.PageResult[*model.SystemSensitiveWord], error) {
	db := s.db.WithContext(ctx).Model(&model.SystemSensitiveWord{})

	if r.Name != "" {
		db = db.Where("name LIKE ?", "%"+r.Name+"%")
	}
	if r.Tag != "" {
		// tags 是 json 数组，mysql 5.7+ 使用 JSON_CONTAINS
		// 为了兼容性，也可以用 like '%"tag"%' ?
		// 若 r.Tag 是 string
		// GORM serializer json stores as string usually `["a", "b"]`
		db = db.Where("tags LIKE ?", "%"+r.Tag+"%")
	}
	if r.Status != nil {
		db = db.Where("status = ?", *r.Status)
	}
	if r.StartDate != "" && r.EndDate != "" {
		db = db.Where("create_time BETWEEN ? AND ?", r.StartDate, r.EndDate)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []*model.SystemSensitiveWord
	offset := (r.PageNo - 1) * r.PageSize
	if err := db.Order("id desc").Offset(offset).Limit(r.PageSize).Find(&list).Error; err != nil {
		return nil, err
	}

	return &core.PageResult[*model.SystemSensitiveWord]{
		List:  list,
		Total: total,
	}, nil
}

// ValidateSensitiveWord 验证敏感词
func (s *SensitiveWordService) ValidateSensitiveWord(ctx context.Context, text string, tags []string) []string {
	if text == "" {
		return []string{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	// TODO: Support Tags filtering in Trie?
	// Currently Validator returns all.
	// If tags provided, implementing strict tag filtering in Trie is hard without node metadata.
	// RuoYi: Test API just calls validate logic.
	// We will validate all first.
	return s.trie.Validate(text)
}
