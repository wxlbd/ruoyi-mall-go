package service

import (
	"context"
	"errors"

	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model"
	"backend-go/internal/pkg/core"
	"backend-go/internal/pkg/utils"
	"backend-go/internal/repo/query"
)

type UserService struct {
	q *query.Query
}

func NewUserService(q *query.Query) *UserService {
	return &UserService{
		q: q,
	}
}

// GetSimpleUserList 获取用户精简列表 (只包含启用用户)
func (s *UserService) GetSimpleUserList(ctx context.Context) ([]resp.UserSimpleRespVO, error) {
	u := s.q.SystemUser
	list, err := u.WithContext(ctx).Where(u.Status.Eq(0)).Find() // 0 = Enabled
	if err != nil {
		return nil, err
	}

	result := make([]resp.UserSimpleRespVO, 0, len(list))
	for _, user := range list {
		result = append(result, resp.UserSimpleRespVO{
			ID:       user.ID,
			Nickname: user.Nickname,
		})
	}
	return result, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *req.UserSaveReq) (int64, error) {
	// 1. 校验唯一性
	if err := s.checkUsernameUnique(ctx, req.Username, 0); err != nil {
		return 0, err
	}
	if req.Mobile != "" {
		if err := s.checkMobileUnique(ctx, req.Mobile, 0); err != nil {
			return 0, err
		}
	}
	if req.Email != "" {
		if err := s.checkEmailUnique(ctx, req.Email, 0); err != nil {
			return 0, err
		}
	}

	// 2. 加密密码
	if req.Password == "" {
		req.Password = "123456" // Default password
	}
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return 0, err
	}

	// 3. 构造 User 对象
	user := &model.SystemUser{
		Username: req.Username,
		Password: hashedPwd,
		Nickname: req.Nickname,
		DeptID:   req.DeptID,
		PostIDs:  "", // TODO: Remove or usage? Old usage. Now using table.
		Email:    req.Email,
		Mobile:   req.Mobile,
		Sex:      req.Sex,
		Avatar:   req.Avatar,
		Status:   int32(req.Status),
		Remark:   req.Remark,
	}

	// 4. 事务执行
	err = s.q.Transaction(func(tx *query.Query) error {
		// 4.1 插入用户
		if err := tx.SystemUser.WithContext(ctx).Create(user); err != nil {
			return err
		}

		// 4.2 关联岗位
		if len(req.PostIDs) > 0 {
			var userPosts []*model.SystemUserPost
			for _, postId := range req.PostIDs {
				userPosts = append(userPosts, &model.SystemUserPost{
					UserID: user.ID,
					PostID: postId,
				})
			}
			if err := tx.SystemUserPost.WithContext(ctx).Create(userPosts...); err != nil {
				return err
			}
		}

		// 4.3 关联角色
		if len(req.RoleIDs) > 0 {
			var userRoles []*model.SystemUserRole
			for _, roleId := range req.RoleIDs {
				userRoles = append(userRoles, &model.SystemUserRole{
					UserID: user.ID,
					RoleID: roleId,
				})
			}
			if err := tx.SystemUserRole.WithContext(ctx).Create(userRoles...); err != nil {
				return err
			}
		}
		return nil
	})

	return user.ID, err
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, req *req.UserSaveReq) error {
	// 1. 校验存在
	u := s.q.SystemUser
	_, err := u.WithContext(ctx).Where(u.ID.Eq(req.ID)).First()
	if err != nil {
		return errors.New("用户不存在")
	}

	// 2. 校验唯一性
	if err := s.checkUsernameUnique(ctx, req.Username, req.ID); err != nil {
		return err
	}
	if req.Mobile != "" {
		if err := s.checkMobileUnique(ctx, req.Mobile, req.ID); err != nil {
			return err
		}
	}
	if req.Email != "" {
		if err := s.checkEmailUnique(ctx, req.Email, req.ID); err != nil {
			return err
		}
	}

	// 3. 事务更新
	return s.q.Transaction(func(tx *query.Query) error {
		// 3.1 更新基本信息
		_, err := tx.SystemUser.WithContext(ctx).Where(u.ID.Eq(req.ID)).Updates(&model.SystemUser{
			Nickname: req.Nickname,
			DeptID:   req.DeptID,
			Email:    req.Email,
			Mobile:   req.Mobile,
			Sex:      req.Sex,
			Avatar:   req.Avatar,
			Status:   int32(req.Status),
			Remark:   req.Remark,
		})
		if err != nil {
			return err
		}

		// 3.2 更新岗位 (Delete + Insert)
		if _, err := tx.SystemUserPost.WithContext(ctx).Where(tx.SystemUserPost.UserID.Eq(req.ID)).Delete(); err != nil {
			return err
		}
		if len(req.PostIDs) > 0 {
			var userPosts []*model.SystemUserPost
			for _, postId := range req.PostIDs {
				userPosts = append(userPosts, &model.SystemUserPost{
					UserID: req.ID,
					PostID: postId,
				})
			}
			if err := tx.SystemUserPost.WithContext(ctx).Create(userPosts...); err != nil {
				return err
			}
		}

		// 3.3 更新角色 (Delete + Insert)
		if _, err := tx.SystemUserRole.WithContext(ctx).Where(tx.SystemUserRole.UserID.Eq(req.ID)).Delete(); err != nil {
			return err
		}
		if len(req.RoleIDs) > 0 {
			var userRoles []*model.SystemUserRole
			for _, roleId := range req.RoleIDs {
				userRoles = append(userRoles, &model.SystemUserRole{
					UserID: req.ID,
					RoleID: roleId,
				})
			}
			if err := tx.SystemUserRole.WithContext(ctx).Create(userRoles...); err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	u := s.q.SystemUser
	_, err := u.WithContext(ctx).Where(u.ID.Eq(id)).Delete()
	// TODO: Delete UserRole and UserPost relations physically or logically?
	// Usually soft delete user is enough, but relations might stay.
	// For now, simple user delete.
	return err
}

// GetUser 获得用户详情
func (s *UserService) GetUser(ctx context.Context, id int64) (*resp.UserProfileRespVO, error) {
	u := s.q.SystemUser
	user, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}

	// Get Roles
	ur := s.q.SystemUserRole
	userRoles, _ := ur.WithContext(ctx).Where(ur.UserID.Eq(id)).Find()
	roleIds := make([]int64, len(userRoles))
	for i, r := range userRoles {
		roleIds[i] = r.RoleID
	}

	// Get Posts
	up := s.q.SystemUserPost
	userPosts, _ := up.WithContext(ctx).Where(up.UserID.Eq(id)).Find()
	postIds := make([]int64, len(userPosts))
	for i, p := range userPosts {
		postIds[i] = p.PostID
	}

	return &resp.UserProfileRespVO{
		UserRespVO: &resp.UserRespVO{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Remark:   user.Remark,
			DeptID:   user.DeptID,
			PostIDs:  postIds,
			RoleIDs:  roleIds,
			Email:    user.Email,
			Mobile:   user.Mobile,
			Sex:      user.Sex,
			Avatar:   user.Avatar,
			Status:   user.Status,
			LoginIP:  user.LoginIP,
			// LoginDate:  *user.LoginDate, // Handle nil
			CreateTime: user.CreatedAt,
		},
		Roles: nil, // Frontend usually doesn't need full role objects here for basic CRUD, but maybe for profile?
		Posts: nil,
	}, nil
}

// GetUserPage 获得用户分页
func (s *UserService) GetUserPage(ctx context.Context, req *req.UserPageReq) (*core.PageResult[*resp.UserRespVO], error) {
	u := s.q.SystemUser
	qb := u.WithContext(ctx)

	if req.Username != "" {
		qb = qb.Where(u.Username.Like("%" + req.Username + "%"))
	}
	if req.Mobile != "" {
		qb = qb.Where(u.Mobile.Like("%" + req.Mobile + "%"))
	}
	if req.Status != nil {
		qb = qb.Where(u.Status.Eq(int32(*req.Status)))
	}
	if req.DeptID > 0 {
		// TODO: Recursive Dept search? Usually yes.
		// For now simple match.
		qb = qb.Where(u.DeptID.Eq(req.DeptID))
	}
	if req.CreateTimeGe != nil {
		qb = qb.Where(u.CreatedAt.Gte(*req.CreateTimeGe))
	}
	if req.CreateTimeLe != nil {
		qb = qb.Where(u.CreatedAt.Lte(*req.CreateTimeLe))
	}

	total, err := qb.Count()
	if err != nil {
		return nil, err
	}

	list, err := qb.Order(u.ID).Offset(req.GetOffset()).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	var data []*resp.UserRespVO
	for _, item := range list {
		data = append(data, &resp.UserRespVO{
			ID:         item.ID,
			Username:   item.Username,
			Nickname:   item.Nickname,
			DeptID:     item.DeptID,
			Email:      item.Email,
			Mobile:     item.Mobile,
			Sex:        item.Sex,
			Avatar:     item.Avatar,
			Status:     item.Status,
			LoginIP:    item.LoginIP,
			CreateTime: item.CreatedAt,
		})
	}

	return &core.PageResult[*resp.UserRespVO]{
		List:  data,
		Total: total,
	}, nil
}

// UpdateUserStatus 修改用户状态
func (s *UserService) UpdateUserStatus(ctx context.Context, req *req.UserUpdateStatusReq) error {
	u := s.q.SystemUser
	_, err := u.WithContext(ctx).Where(u.ID.Eq(req.ID)).Update(u.Status, req.Status)
	return err
}

// UpdateUserPassword 修改用户密码
func (s *UserService) UpdateUserPassword(ctx context.Context, req *req.UserUpdatePasswordReq) error {
	u := s.q.SystemUser
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}
	_, err = u.WithContext(ctx).Where(u.ID.Eq(req.ID)).Update(u.Password, hashedPwd)
	return err
}

// ResetUserPassword 重置用户密码
func (s *UserService) ResetUserPassword(ctx context.Context, req *req.UserResetPasswordReq) error {
	u := s.q.SystemUser
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}
	_, err = u.WithContext(ctx).Where(u.ID.Eq(req.ID)).Update(u.Password, hashedPwd)
	return err
}

// GetUserList 获得用户列表 (用于导出)
func (s *UserService) GetUserList(ctx context.Context, req *req.UserExportReq) ([]*resp.UserRespVO, error) {
	u := s.q.SystemUser
	qb := u.WithContext(ctx)

	if req.Username != "" {
		qb = qb.Where(u.Username.Like("%" + req.Username + "%"))
	}
	if req.Mobile != "" {
		qb = qb.Where(u.Mobile.Like("%" + req.Mobile + "%"))
	}
	if req.Status != nil {
		qb = qb.Where(u.Status.Eq(int32(*req.Status)))
	}
	if req.DeptID > 0 {
		qb = qb.Where(u.DeptID.Eq(req.DeptID))
	}
	if req.CreateTimeGe != nil {
		qb = qb.Where(u.CreatedAt.Gte(*req.CreateTimeGe))
	}
	if req.CreateTimeLe != nil {
		qb = qb.Where(u.CreatedAt.Lte(*req.CreateTimeLe))
	}

	list, err := qb.Order(u.ID).Find()
	if err != nil {
		return nil, err
	}

	var data []*resp.UserRespVO
	for _, item := range list {
		data = append(data, &resp.UserRespVO{
			ID:         item.ID,
			Username:   item.Username,
			Nickname:   item.Nickname,
			DeptID:     item.DeptID,
			Email:      item.Email,
			Mobile:     item.Mobile,
			Sex:        item.Sex,
			Avatar:     item.Avatar,
			Status:     item.Status,
			LoginIP:    item.LoginIP,
			CreateTime: item.CreatedAt,
		})
	}
	return data, nil
}

// GetImportTemplate 获得导入模板
func (s *UserService) GetImportTemplate(ctx context.Context) ([]resp.UserImportExcelVO, error) {
	return []resp.UserImportExcelVO{
		{
			Username: "zhangsan",
			Nickname: "张三",
			Email:    "zhangsan@yudao.cn",
			Mobile:   "15601691300",
			Sex:      "1",
			Status:   "0",
			DeptID:   100,
		},
		{
			Username: "lisi",
			Nickname: "李四",
			Email:    "lisi@yudao.cn",
			Mobile:   "15601691301",
			Sex:      "2",
			Status:   "0",
			DeptID:   100,
		},
	}, nil
}

// Helpers

func (s *UserService) checkUsernameUnique(ctx context.Context, username string, excludeId int64) error {
	u := s.q.SystemUser
	qb := u.WithContext(ctx).Where(u.Username.Eq(username))
	if excludeId > 0 {
		qb = qb.Where(u.ID.Neq(excludeId))
	}
	count, err := qb.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}
	return nil
}

func (s *UserService) checkMobileUnique(ctx context.Context, mobile string, excludeId int64) error {
	u := s.q.SystemUser
	qb := u.WithContext(ctx).Where(u.Mobile.Eq(mobile))
	if excludeId > 0 {
		qb = qb.Where(u.ID.Neq(excludeId))
	}
	count, err := qb.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("手机号已存在")
	}
	return nil
}

func (s *UserService) checkEmailUnique(ctx context.Context, email string, excludeId int64) error {
	u := s.q.SystemUser
	qb := u.WithContext(ctx).Where(u.Email.Eq(email))
	if excludeId > 0 {
		qb = qb.Where(u.ID.Neq(excludeId))
	}
	count, err := qb.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("邮箱已存在")
	}
	return nil
}
