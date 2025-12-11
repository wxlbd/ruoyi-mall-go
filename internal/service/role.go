package service

import (
	"context"
	"errors"

	"github.com/samber/lo"

	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"
)

type RoleService struct {
	q *query.Query
}

func NewRoleService(q *query.Query) *RoleService {
	return &RoleService{
		q: q,
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, req *req.RoleSaveReq) (int64, error) {
	if err := s.checkDuplicate(ctx, req.Name, req.Code, 0); err != nil {
		return 0, err
	}

	role := &model.SystemRole{
		Name:      req.Name,
		Code:      req.Code,
		Sort:      req.Sort,
		Status:    int32(req.Status),
		Remark:    req.Remark,
		Type:      2, // Default Custom
		DataScope: 1, // Default All
	}

	err := s.q.SystemRole.WithContext(ctx).Create(role)
	return role.ID, err
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, req *req.RoleSaveReq) error {
	r := s.q.SystemRole
	role, err := r.WithContext(ctx).Where(r.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.New("角色不存在")
	}
	if role.Type == 1 {
		// Allow updating basic info even for system roles, but maybe restricted in some systems.
		// For now allow it.
	}

	if err := s.checkDuplicate(ctx, req.Name, req.Code, req.ID); err != nil {
		return err
	}

	_, err = r.WithContext(ctx).Where(r.ID.Eq(req.ID)).Updates(&model.SystemRole{
		Name:   req.Name,
		Code:   req.Code,
		Sort:   req.Sort,
		Status: int32(req.Status),
		Remark: req.Remark,
	})
	return err
}

// UpdateRoleStatus 更新角色状态
func (s *RoleService) UpdateRoleStatus(ctx context.Context, req *req.RoleUpdateStatusReq) error {
	r := s.q.SystemRole
	role, err := r.WithContext(ctx).Where(r.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.New("角色不存在")
	}
	if role.Type == 1 {
		return errors.New("内置角色不能修改状态")
	}
	_, err = r.WithContext(ctx).Where(r.ID.Eq(req.ID)).Update(r.Status, req.Status)
	return err
}

// UpdateRoleDataScope 更新数据权限
func (s *RoleService) UpdateRoleDataScope(ctx context.Context, roleId int64, dataScope int, deptIds []int64) error {
	r := s.q.SystemRole
	_, err := r.WithContext(ctx).Where(r.ID.Eq(roleId)).First()
	if err != nil {
		return errors.New("角色不存在")
	}

	_, err = r.WithContext(ctx).Where(r.ID.Eq(roleId)).Updates(&model.SystemRole{
		DataScope:        int32(dataScope),
		DataScopeDeptIds: deptIds, // Handled by serializer:json
	})
	return err
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, id int64) error {
	r := s.q.SystemRole
	role, err := r.WithContext(ctx).Where(r.ID.Eq(id)).First()
	if err != nil {
		return errors.New("角色不存在")
	}
	if role.Type == 1 {
		return errors.New("内置角色不能删除")
	}
	// Check assigned users count
	userRoleCount, err := s.q.SystemUserRole.WithContext(ctx).Where(s.q.SystemUserRole.RoleID.Eq(id)).Count()
	if err != nil {
		return err
	}
	if userRoleCount > 0 {
		return errors.New("角色已分配给用户，无法删除")
	}
	_, err = r.WithContext(ctx).Where(r.ID.Eq(id)).Delete()
	return err
}

// GetRole 获得角色
func (s *RoleService) GetRole(ctx context.Context, id int64) (*resp.RoleRespVO, error) {
	r := s.q.SystemRole
	item, err := r.WithContext(ctx).Where(r.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return s.convertResp(item), nil
}

// GetRolePage 分页
func (s *RoleService) GetRolePage(ctx context.Context, req *req.RolePageReq) (*core.PageResult[*resp.RoleRespVO], error) {
	r := s.q.SystemRole
	qb := r.WithContext(ctx)

	if req.Name != "" {
		qb = qb.Where(r.Name.Like("%" + req.Name + "%"))
	}
	if req.Code != "" {
		qb = qb.Where(r.Code.Like("%" + req.Code + "%"))
	}
	if req.Status != nil {
		qb = qb.Where(r.Status.Eq(int32(*req.Status)))
	}
	if req.CreateTimeGe != nil {
		qb = qb.Where(r.CreatedAt.Gte(*req.CreateTimeGe))
	}
	if req.CreateTimeLe != nil {
		qb = qb.Where(r.CreatedAt.Lte(*req.CreateTimeLe))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}
	list, err := qb.Order(r.Sort, r.ID).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	return &core.PageResult[*resp.RoleRespVO]{
		List:  lo.Map(list, func(item *model.SystemRole, _ int) *resp.RoleRespVO { return s.convertResp(item) }),
		Total: total,
	}, nil
}

// GetRoleListByStatus 获取全列表
func (s *RoleService) GetRoleListByStatus(ctx context.Context, status int) ([]*resp.RoleRespVO, error) {
	r := s.q.SystemRole
	list, err := r.WithContext(ctx).Where(r.Status.Eq(int32(status))).Order(r.Sort, r.ID).Find()
	if err != nil {
		return nil, err
	}
	return lo.Map(list, func(item *model.SystemRole, _ int) *resp.RoleRespVO { return s.convertResp(item) }), nil
}

// GetRoleList IDs (Already existed, keep it)
func (s *RoleService) GetRoleList(ctx context.Context, ids []int64) ([]*model.SystemRole, error) {
	if len(ids) == 0 {
		return []*model.SystemRole{}, nil
	}
	r := s.q.SystemRole
	return r.WithContext(ctx).Where(r.ID.In(ids...)).Find()
}

// Helpers

func (s *RoleService) checkDuplicate(ctx context.Context, name, code string, excludeId int64) error {
	r := s.q.SystemRole
	// Name unique
	qb := r.WithContext(ctx).Where(r.Name.Eq(name))
	if excludeId > 0 {
		qb = qb.Where(r.ID.Neq(excludeId))
	}
	count, err := qb.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色名称已存在")
	}

	// Code unique
	qb = r.WithContext(ctx).Where(r.Code.Eq(code))
	if excludeId > 0 {
		qb = qb.Where(r.ID.Neq(excludeId))
	}
	count, err = qb.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色编码已存在")
	}
	return nil
}

func (s *RoleService) convertResp(item *model.SystemRole) *resp.RoleRespVO {
	return &resp.RoleRespVO{
		ID:               item.ID,
		Name:             item.Name,
		Code:             item.Code,
		Sort:             item.Sort,
		Status:           item.Status,
		Type:             item.Type,
		Remark:           item.Remark,
		DataScope:        item.DataScope,
		DataScopeDeptIDs: item.DataScopeDeptIds,
		CreateTime:       item.CreatedAt,
	}
}
