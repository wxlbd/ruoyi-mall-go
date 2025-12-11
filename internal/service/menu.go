package service

import (
	"context"

	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/repo/query"
)

type MenuService struct {
	q *query.Query
}

func NewMenuService(q *query.Query) *MenuService {
	return &MenuService{
		q: q,
	}
}

// GetMenuList 获取菜单列表
func (s *MenuService) GetMenuList(ctx context.Context, req *req.MenuListReq) ([]*resp.MenuResp, error) {
	m := s.q.SystemMenu
	qb := m.WithContext(ctx)

	// 动态条件
	if req.Name != "" {
		qb = qb.Where(m.Name.Like("%" + req.Name + "%"))
	}
	if req.Status != nil {
		qb = qb.Where(m.Status.Eq(*req.Status))
	}

	// 排序
	qb = qb.Order(m.Sort, m.ID)

	list, err := qb.Find()
	if err != nil {
		return nil, err
	}

	// DO -> DTO
	var res []*resp.MenuResp
	for _, item := range list {
		res = append(res, &resp.MenuResp{
			ID:            item.ID,
			ParentID:      item.ParentID,
			Name:          item.Name,
			Type:          item.Type,
			Sort:          item.Sort,
			Path:          item.Path,
			Icon:          item.Icon,
			Component:     item.Component,
			ComponentName: item.ComponentName,
			Permission:    item.Permission,
			Status:        item.Status,
			Visible:       bool(item.Visible),
			KeepAlive:     bool(item.KeepAlive),
			AlwaysShow:    bool(item.AlwaysShow),
			CreateTime:    item.CreatedAt,
		})
	}
	return res, nil
}

// GetSimpleMenuList 获取精简菜单列表 (仅返回开启状态的菜单)
func (s *MenuService) GetSimpleMenuList(ctx context.Context) ([]*resp.MenuSimpleResp, error) {
	m := s.q.SystemMenu
	// 这里硬编码 Status=0 (CommonStatusEnum.ENABLE)
	list, err := m.WithContext(ctx).Where(m.Status.Eq(0)).Order(m.Sort, m.ID).Find()
	if err != nil {
		return nil, err
	}

	var res []*resp.MenuSimpleResp
	for _, item := range list {
		res = append(res, &resp.MenuSimpleResp{
			ID:       item.ID,
			ParentID: item.ParentID,
			Name:     item.Name,
			Type:     item.Type,
		})
	}
	return res, nil
}

// GetMenu 获取菜单详情
func (s *MenuService) GetMenu(ctx context.Context, id int64) (*resp.MenuResp, error) {
	m := s.q.SystemMenu
	item, err := m.WithContext(ctx).Where(m.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}

	return &resp.MenuResp{
		ID:            item.ID,
		ParentID:      item.ParentID,
		Name:          item.Name,
		Type:          item.Type,
		Sort:          item.Sort,
		Path:          item.Path,
		Icon:          item.Icon,
		Component:     item.Component,
		ComponentName: item.ComponentName,
		Permission:    item.Permission,
		Status:        item.Status,
		Visible:       bool(item.Visible),
		KeepAlive:     bool(item.KeepAlive),
		AlwaysShow:    bool(item.AlwaysShow),
		CreateTime:    item.CreatedAt,
	}, nil
}

// BuildMenuTree 构建菜单树
func (s *MenuService) BuildMenuTree(menus []*resp.MenuResp) []resp.MenuVO {
	if len(menus) == 0 {
		return []resp.MenuVO{}
	}

	// 1. 构建 Map 和 根节点列表
	menuMap := make(map[int64]*resp.MenuVO)
	var roots []resp.MenuVO

	// 先把所有 MenuResp 转为 MenuVO 并存入 Map
	for _, m := range menus {
		vo := resp.MenuVO{
			ID:            m.ID,
			ParentID:      m.ParentID,
			Name:          m.Name,
			Path:          m.Path,
			Component:     m.Component,
			ComponentName: m.ComponentName,
			Icon:          m.Icon,
			Visible:       m.Visible,
			KeepAlive:     m.KeepAlive,
			AlwaysShow:    m.AlwaysShow,
			Children:      make([]resp.MenuVO, 0),
		}
		menuMap[m.ID] = &vo
	}

	// 2. 再次遍历，组装父子关系
	// 注意：这里是有序的，前提是传入的 menus 已经按 Sort 排序
	for _, m := range menus {
		node := menuMap[m.ID]
		if m.ParentID == 0 {
			// 根节点
			roots = append(roots, *node)
		} else {
			// 子节点，挂载到父节点
			if parent, ok := menuMap[m.ParentID]; ok {
				parent.Children = append(parent.Children, *node)
			} else {
				// 如果父节点找不到（可能是被禁用了，或者数据不一致），也可以考虑作为根节点，或者忽略
				// 这里选择作为根节点兜底，或者根据业务需求丢弃
				// 暂时忽略
			}
		}
	}

	// 3. 这里的 roots 是值拷贝，但是 Children 是指针引用（Slice），所以 append 到 parent.Children 的修改是生效的吗？
	// 不，menuMap 存的是 *MenuVO。
	// parent.Children = append(parent.Children, *node) 这里存的是 Value。
	// 如果 *node 后来又有 Children了，parent.Children 里的那个 copy 并不会更新。
	// 所以必须用 Pointer 这种递归构建，或者 Two-Pass Pointer 链接。

	// 修正逻辑：使用递归或者 Pointer Map
	return s.buildTreeRecursive(menus, 0)
}

func (s *MenuService) buildTreeRecursive(list []*resp.MenuResp, parentId int64) []resp.MenuVO {
	var tree []resp.MenuVO
	for _, item := range list {
		// 过滤按钮类型 (type=3), 只保留目录(type=1)和菜单(type=2)
		// Java: MenuTypeEnum.BUTTON = 3, 在 convert 时被过滤
		if item.Type == 3 {
			continue
		}

		if item.ParentID == parentId {
			node := resp.MenuVO{
				ID:            item.ID,
				ParentID:      item.ParentID,
				Name:          item.Name,
				Path:          item.Path,
				Component:     item.Component,
				ComponentName: item.ComponentName,
				Icon:          item.Icon,
				Visible:       item.Visible,
				KeepAlive:     item.KeepAlive,
				AlwaysShow:    item.AlwaysShow,
			}
			children := s.buildTreeRecursive(list, item.ID)
			if len(children) > 0 {
				node.Children = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// GetMenuListByIds 根据ID列表获取菜单
func (s *MenuService) GetMenuListByIds(ctx context.Context, ids []int64) ([]*resp.MenuResp, error) {
	if len(ids) == 0 {
		return []*resp.MenuResp{}, nil
	}
	m := s.q.SystemMenu
	list, err := m.WithContext(ctx).Where(m.ID.In(ids...)).Order(m.Sort, m.ID).Find()
	if err != nil {
		return nil, err
	}

	var res []*resp.MenuResp
	for _, item := range list {
		modelPerm := item.Permission // Capture value
		res = append(res, &resp.MenuResp{
			ID:            item.ID,
			ParentID:      item.ParentID,
			Name:          item.Name,
			Type:          item.Type,
			Sort:          item.Sort,
			Path:          item.Path,
			Icon:          item.Icon,
			Component:     item.Component,
			ComponentName: item.ComponentName,
			Permission:    modelPerm,
			Status:        item.Status,
			Visible:       bool(item.Visible),
			KeepAlive:     bool(item.KeepAlive),
			AlwaysShow:    bool(item.AlwaysShow),
			CreateTime:    item.CreatedAt,
		})
	}
	return res, nil
}
