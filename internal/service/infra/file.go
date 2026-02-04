package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/infra"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/file"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"github.com/samber/lo"
)

type FileService struct {
	q                 *query.Query
	fileConfigService *FileConfigService
}

func NewFileService(q *query.Query, fileConfigService *FileConfigService) *FileService {
	return &FileService{
		q:                 q,
		fileConfigService: fileConfigService,
	}
}

// GenerateUploadPath 生成上传路径（用于预签名上传）
func (s *FileService) GenerateUploadPath(name string, directory string) (string, error) {
	if err := s.validateFileType(name); err != nil {
		return "", err
	}
	if err := s.validatePath(directory); err != nil {
		return "", err
	}
	return s.generateSafePath(name, directory), nil
}

// CreateFile 上传/创建文件
func (s *FileService) CreateFile(ctx context.Context, name string, path string, content []byte, fileType string) (string, error) {
	// 1. 获取 Master 配置
	config, err := s.fileConfigService.GetMasterFileConfig(ctx)
	if err != nil {
		return "", errors.New("请先配置主文件存储")
	}

	// 2. 验证文件大小（最大 100MB）
	maxFileSize := int64(100 * 1024 * 1024)
	if int64(len(content)) > maxFileSize {
		return "", fmt.Errorf("文件大小超过限制: 最大 %d MB", maxFileSize/1024/1024)
	}

	// 3. 验证文件类型（白名单）
	if err := s.validateFileType(name); err != nil {
		return "", err
	}

	// 4. 验证路径安全性（防止路径遍历）
	if err := s.validatePath(path); err != nil {
		return "", err
	}

	// 5. 初始化客户端
	client, err := file.NewFileClient(config.ID, config.Storage, config.Config)
	if err != nil {
		return "", fmt.Errorf("初始化文件客户端失败: %v", err)
	}

	// 6. 生成安全路径（带防冲突时间戳）
	safePath := s.generateSafePath(name, path)

	// 7. 上传
	url, err := client.Upload(content, safePath)
	if err != nil {
		return "", err
	}

	// 8. 自动探测文件类型（作为兜底）
	if fileType == "" {
		fileType = http.DetectContentType(content)
	}

	// 9. 保存记录
	fileRecord := &model.InfraFile{
		ConfigId: config.ID,
		Name:     name,
		Path:     safePath,
		Url:      url,
		Type:     fileType,
		Size:     len(content),
	}
	err = s.q.InfraFile.WithContext(ctx).Create(fileRecord)
	if err != nil {
		return "", err
	}

	return url, nil
}

// validateFileType 验证文件类型（白名单）
func (s *FileService) validateFileType(name string) error {
	// 允许的文件类型
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".txt":  true,
		".zip":  true,
		".rar":  true,
	}

	ext := strings.ToLower(filepath.Ext(name))
	if ext == "" {
		return errors.New("文件必须包含扩展名")
	}

	if !allowedExtensions[ext] {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}

	return nil
}

// validatePath 验证路径安全性（防止路径遍历攻击）
func (s *FileService) validatePath(path string) error {
	if path == "" {
		return nil // 空路径允许，会使用默认路径
	}

	// 检查路径遍历攻击
	if strings.Contains(path, "..") {
		return errors.New("路径包含非法字符 '..'")
	}

	// 检查绝对路径
	if filepath.IsAbs(path) {
		return errors.New("不允许使用绝对路径")
	}

	// 检查特殊字符
	if strings.Contains(path, "\x00") {
		return errors.New("路径包含空字符")
	}

	return nil
}

// generateSafePath 生成安全路径（带防冲突时间戳）
func (s *FileService) generateSafePath(name string, directory string) string {
	// 1. 生成时间戳（纳秒级别，防止冲突）
	timestamp := time.Now().UnixNano()

	// 2. 分离文件名和扩展名
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	// 3. 清理基础名（移除特殊字符）
	base = strings.ReplaceAll(base, "/", "_")
	base = strings.ReplaceAll(base, "\\", "_")
	base = strings.ReplaceAll(base, "..", "_")

	// 4. 目录处理
	if directory == "" {
		directory = time.Now().Format("2006/01/02")
	} else {
		directory = strings.Trim(directory, "/")
	}

	// 5. 组合路径：目录/基础名_时间戳.扩展名
	// 示例：2024/12/18/avatar_1734523456789123456.jpg
	return fmt.Sprintf("%s/%s_%d%s", directory, base, timestamp, ext)
}

func (s *FileService) DeleteFile(ctx context.Context, id int64) error {
	f := s.q.InfraFile
	fileRecord, err := f.WithContext(ctx).Where(f.ID.Eq(id)).First()
	if err != nil {
		return errors.New("文件不存在")
	}

	// 获取配置
	config, err := s.fileConfigService.GetFileConfig(ctx, fileRecord.ConfigId)
	if err != nil {
		// 如果配置都不存在了，只删除数据库记录
		f.WithContext(ctx).Where(f.ID.Eq(id)).Delete()
		return nil
	}

	// 初始化客户端并删除物理文件
	if config.Config != nil {
		configBytes, _ := json.Marshal(config.Config)
		client, err := file.NewFileClient(config.ID, config.Storage, configBytes)
		if err == nil {
			_ = client.Delete(fileRecord.Path)
		}
	}

	_, err = f.WithContext(ctx).Where(f.ID.Eq(id)).Delete()
	return err
}

// DeleteFileList 批量删除文件
func (s *FileService) DeleteFileList(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if err := s.DeleteFile(ctx, id); err != nil {
			// 如果单个删除失败，记录日志并继续，对齐 Java 逻辑（Java 也是逐个删除并抛出异常，Go 这里选择继续处理以防部分删除后中断）
			// 实际上 Java 的 deleteFileList 是调用 fileService.deleteFileList(ids)，这里我们也保持 Service 层闭环
			continue
		}
	}
	return nil
}

// GetFile 获得文件详情
func (s *FileService) GetFile(ctx context.Context, id int64) (*infra.FileResp, error) {
	f := s.q.InfraFile
	fileRecord, err := f.WithContext(ctx).Where(f.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.New("文件不存在")
	}
	return s.convertResp(fileRecord), nil
}

// GetFileContent 获取文件内容
func (s *FileService) GetFileContent(ctx context.Context, configId int64, path string) ([]byte, error) {
	// 对齐 Java: URLUtil.decode(path, StandardCharsets.UTF_8, false);
	decodedPath, err := url.PathUnescape(path)
	if err == nil {
		path = decodedPath
	}

	config, err := s.fileConfigService.GetFileConfig(ctx, configId)
	if err != nil {
		return nil, errors.New("配置不存在")
	}
	if config.Config == nil {
		return nil, errors.New("配置内容为空")
	}

	configBytes, _ := json.Marshal(config.Config)
	client, err := file.NewFileClient(config.ID, config.Storage, configBytes)
	if err != nil {
		return nil, err
	}
	return client.GetContent(path)
}

// GetFilePage 获得文件分页
func (s *FileService) GetFilePage(ctx context.Context, req *infra.FilePageReq) (*pagination.PageResult[*infra.FileResp], error) {
	f := s.q.InfraFile
	qb := f.WithContext(ctx)

	if req.Path != "" {
		qb = qb.Where(f.Path.Like("%" + req.Path + "%"))
	}
	if req.Type != "" {
		qb = qb.Where(f.Type.Like("%" + req.Type + "%"))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(f.ID.Desc()).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*infra.FileResp]{
		List:  lo.Map(list, func(item *model.InfraFile, _ int) *infra.FileResp { return s.convertResp(item) }),
		Total: total,
	}, nil
}

func (s *FileService) convertResp(item *model.InfraFile) *infra.FileResp {
	return &infra.FileResp{
		ID:         item.ID,
		ConfigId:   item.ConfigId,
		Name:       item.Name,
		Path:       item.Path,
		Url:        item.Url,
		Type:       item.Type,
		Size:       item.Size,
		CreateTime: item.CreateTime,
	}
}

func (s *FileService) GetFilePresignedUrl(ctx context.Context, path string) (*infra.FilePresignedUrlResp, error) {
	config, err := s.fileConfigService.GetMasterFileConfig(ctx)
	if err != nil {
		return nil, errors.New("请先配置主文件存储")
	}

	configBytes, _ := json.Marshal(config.Config)
	client, err := file.NewFileClient(config.ID, config.Storage, configBytes)
	if err != nil {
		return nil, err
	}

	presignedUrl, err := client.GetPresignedURL(path)
	if err != nil {
		return nil, err
	}

	return &infra.FilePresignedUrlResp{
		ConfigID:  config.ID,
		UploadURL: presignedUrl,
		URL:       client.GetURL(path),
		Path:      path,
	}, nil
}

func (s *FileService) CreateFileCallback(ctx context.Context, req *infra.FileCreateReq) (int64, error) {
	// 验证配置是否存在
	_, err := s.fileConfigService.GetFileConfig(ctx, req.ConfigID)
	if err != nil {
		return 0, errors.New("配置不存在")
	}

	fileRecord := &model.InfraFile{
		ConfigId: req.ConfigID,
		Name:     req.Name,
		Path:     req.Path,
		Url:      req.URL,
		Type:     req.Type,
		Size:     req.Size,
	}
	err = s.q.InfraFile.WithContext(ctx).Create(fileRecord)
	if err != nil {
		return 0, err
	}
	return fileRecord.ID, nil
}
