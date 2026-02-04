package infra

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/infra"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/file"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
)

type FileConfigService struct {
	q *query.Query
}

func NewFileConfigService(q *query.Query) *FileConfigService {
	return &FileConfigService{
		q: q,
	}
}

// CreateFileConfig 创建文件配置
func (s *FileConfigService) CreateFileConfig(ctx context.Context, req *infra.FileConfigSaveReq) (int64, error) {
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		return 0, err
	}

	config := &model.InfraFileConfig{
		Name:    req.Name,
		Storage: req.Storage,
		Config:  configBytes,
		Remark:  req.Remark,
		Master:  false, // 默认不为主配置
	}

	// 如果这是第一个配置，自动设为主配置
	count, _ := s.q.InfraFileConfig.WithContext(ctx).Count()
	if count == 0 {
		config.Master = true
	}

	err = s.q.InfraFileConfig.WithContext(ctx).Create(config)
	return config.ID, err
}

// UpdateFileConfig 更新文件配置
func (s *FileConfigService) UpdateFileConfig(ctx context.Context, req *infra.FileConfigSaveReq) error {
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		return err
	}

	_, err = s.q.InfraFileConfig.WithContext(ctx).Where(s.q.InfraFileConfig.ID.Eq(req.ID)).Updates(&model.InfraFileConfig{
		Name:    req.Name,
		Storage: req.Storage,
		Config:  configBytes,
		Remark:  req.Remark,
	})
	return err
}

// UpdateFileConfigMaster 更新主配置
func (s *FileConfigService) UpdateFileConfigMaster(ctx context.Context, id int64) error {
	return s.q.Transaction(func(tx *query.Query) error {
		c := tx.InfraFileConfig
		// 1. 将所有配置设为非主配置
		// Master 是 BitBool (field.Field)，使用 Eq
		if _, err := c.WithContext(ctx).Where(c.Master.Eq(model.BitBool(true))).Update(c.Master, false); err != nil {
			return err
		}
		// 2. 将当前配置设为主配置
		result, err := c.WithContext(ctx).Where(c.ID.Eq(id)).Update(c.Master, true)
		if err != nil {
			return err
		}
		if result.RowsAffected == 0 {
			return errors.New("配置不存在")
		}
		return nil
	})
}

// DeleteFileConfig 删除文件配置
func (s *FileConfigService) DeleteFileConfig(ctx context.Context, id int64) error {
	c := s.q.InfraFileConfig
	item, err := c.WithContext(ctx).Where(c.ID.Eq(id)).First()
	if err != nil {
		return err
	}
	if bool(item.Master) {
		return errors.New("不能删除主配置")
	}
	_, err = c.WithContext(ctx).Where(c.ID.Eq(id)).Delete()
	return err
}

// DeleteFileConfigList 批量删除文件配置
func (s *FileConfigService) DeleteFileConfigList(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	for _, id := range ids {
		if err := s.DeleteFileConfig(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// GetFileConfig 获得文件配置
func (s *FileConfigService) GetFileConfig(ctx context.Context, id int64) (*infra.FileConfigResp, error) {
	c := s.q.InfraFileConfig
	item, err := c.WithContext(ctx).Where(c.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return s.convertResp(item), nil
}

// GetMasterFileConfig 获得主配置
func (s *FileConfigService) GetMasterFileConfig(ctx context.Context) (*model.InfraFileConfig, error) {
	c := s.q.InfraFileConfig
	return c.WithContext(ctx).Where(c.Master.Eq(model.BitBool(true))).First()
}

// GetFileConfigPage 获得文件配置分页
func (s *FileConfigService) GetFileConfigPage(ctx context.Context, req *infra.FileConfigPageReq) (*pagination.PageResult[*infra.FileConfigResp], error) {
	c := s.q.InfraFileConfig
	qb := c.WithContext(ctx)

	if req.Name != "" {
		qb = qb.Where(c.Name.Like("%" + req.Name + "%"))
	}
	if req.Storage != nil {
		qb = qb.Where(c.Storage.Eq(*req.Storage))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(c.ID.Desc()).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*infra.FileConfigResp]{
		List:  lo.Map(list, func(item *model.InfraFileConfig, _ int) *infra.FileConfigResp { return s.convertResp(item) }),
		Total: total,
	}, nil
}

func (s *FileConfigService) convertResp(item *model.InfraFileConfig) *infra.FileConfigResp {
	var configMap map[string]interface{}
	_ = json.Unmarshal(item.Config, &configMap)
	return &infra.FileConfigResp{
		ID:         item.ID,
		Name:       item.Name,
		Storage:    item.Storage,
		Master:     bool(item.Master),
		Config:     &configMap,
		Remark:     item.Remark,
		CreateTime: item.CreateTime,
	}
}

func (s *FileConfigService) TestFileConfig(ctx context.Context, id int64) (string, error) {
	config, err := s.GetFileConfig(ctx, id)
	if err != nil {
		return "", errors.New("配置不存在")
	}
	if config.Config == nil {
		return "", errors.New("配置内容为空")
	}

	configBytes, _ := json.Marshal(config.Config)
	client, err := file.NewFileClient(config.ID, config.Storage, configBytes)
	if err != nil {
		return "", err
	}

	// 测试上传文件
	path := "test.txt"
	content := []byte("test")
	url, err := client.Upload(content, path)
	if err != nil {
		return "", errors.New("上传文件失败: " + err.Error())
	}

	// 测试删除文件
	if err := client.Delete(path); err != nil {
		return "", errors.New("删除文件失败: " + err.Error())
	}

	return url, nil
}
