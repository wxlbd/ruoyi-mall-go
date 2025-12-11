package service

import (
	"context"
	"strings"

	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/pkg/utils"
	"backend-go/internal/repo/query"
)

type AuthService struct {
	repo     *query.Query
	permSvc  *PermissionService
	roleSvc  *RoleService
	menuSvc  *MenuService
	tokenSvc *OAuth2TokenService
}

func NewAuthService(
	repo *query.Query,
	permSvc *PermissionService,
	roleSvc *RoleService,
	menuSvc *MenuService,
	tokenSvc *OAuth2TokenService,
) *AuthService {
	return &AuthService{
		repo:     repo,
		tokenSvc: tokenSvc,
		permSvc:  permSvc,
		roleSvc:  roleSvc,
		menuSvc:  menuSvc,
	}
}

// GetPermissionInfo 获取登录用户的权限信息
func (s *AuthService) GetPermissionInfo(ctx context.Context) (*resp.AuthPermissionInfoResp, error) {
	// 1. 获取当前用户 ID (从 Context)
	userIdVal := ctx.Value(core.CtxUserIDKey)
	if userIdVal == nil {
		return nil, core.NewBizError(401, "未登录")
	}
	userId, ok := userIdVal.(int64)
	if !ok {
		return nil, core.NewBizError(401, "用户标识无效")
	}

	// 2. 获取用户信息
	uRepo := s.repo.SystemUser
	user, err := uRepo.WithContext(ctx).Where(uRepo.ID.Eq(userId)).First()
	if err != nil {
		// 区分是否是 RecordNotFound (已在 InitDB 中忽略日志，这里 err 可能为 gorm.ErrRecordNotFound)
		// 但为了安全，通常模糊返回
		return nil, core.NewBizError(1002000002, "账号或密码不正确")
	}

	// 3. 获取用户角色
	roleIds, err := s.permSvc.GetUserRoleIdListByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	// 获取角色Code列表
	rolesData, err := s.roleSvc.GetRoleList(ctx, roleIds)
	if err != nil {
		return nil, err
	}
	// 过滤禁用的角色 (Java: roles.removeIf(role -> !CommonStatusEnum.ENABLE.getStatus().equals(role.getStatus())))
	var roles []string
	for _, r := range rolesData {
		if r.Status == 0 { // 0 = ENABLE
			roles = append(roles, r.Code)
		}
	}

	// 4. 获取角色菜单
	menuIds, err := s.permSvc.GetRoleMenuListByRoleId(ctx, roleIds)
	if err != nil {
		return nil, err
	}

	// 5. 获取菜单列表
	menus, err := s.menuSvc.GetMenuListByIds(ctx, menuIds)
	if err != nil {
		return nil, err
	}

	// 5.1 过滤禁用的菜单 (Java: menuList = menuService.filterDisableMenus(menuList))
	var enabledMenus []*resp.MenuResp
	for _, m := range menus {
		if m.Status == 0 { // 0 = ENABLE
			enabledMenus = append(enabledMenus, m)
		}
	}

	// 6. 获取角色权限 (从菜单中提取)
	permissions := make([]string, 0)
	for _, m := range enabledMenus {
		if m.Permission != "" {
			permissions = append(permissions, m.Permission)
		}
	}

	// 7. 构建菜单树
	menuTree := s.menuSvc.BuildMenuTree(enabledMenus)

	return &resp.AuthPermissionInfoResp{
		User: resp.UserVO{
			ID:       user.ID,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			DeptID:   user.DeptID,
			Username: user.Username,
			Email:    user.Email,
		},
		Roles:       roles,
		Permissions: permissions,
		Menus:       menuTree,
	}, nil
}

// Login 登录业务
func (s *AuthService) Login(ctx context.Context, req *req.AuthLoginReq) (*resp.AuthLoginResp, error) {
	// 0. 解析租户
	var tenantId int64 = 1 // 默认租户ID
	if req.TenantName != "" {
		tenantRepo := s.repo.SystemTenant
		tenant, err := tenantRepo.WithContext(ctx).Where(tenantRepo.Name.Eq(req.TenantName)).First()
		if err != nil {
			return nil, core.NewBizError(1002000003, "租户不存在")
		}
		if tenant.Status != 0 {
			return nil, core.NewBizError(1002000004, "租户已被禁用")
		}
		tenantId = tenant.ID
	}

	// 1. 查询用户
	userRepo := s.repo.SystemUser
	user, err := userRepo.WithContext(ctx).Where(userRepo.Username.Eq(req.Username), userRepo.TenantID.Eq(tenantId)).First()
	if err != nil {
		// 区分是否是 RecordNotFound (已在 InitDB 中忽略日志，这里 err 可能为 gorm.ErrRecordNotFound)
		// 但为了安全，通常模糊返回
		return nil, core.NewBizError(1002000002, "账号或密码不正确")
	}

	// 2. 校验状态
	if user.Status != 0 { // 假设 0 是开启
		return nil, core.NewBizError(1002000001, "用户已被禁用")
	}

	// 3. 校验密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, core.NewBizError(1002000002, "账号或密码不正确")
	}

	// 4. 构建用户信息
	userInfo := map[string]string{
		"nickname": user.Nickname,
	}
	if user.DeptID != 0 {
		userInfo["deptId"] = string(rune(user.DeptID))
	}

	// 5. 创建访问令牌（使用 OAuth2TokenService，与 Java 对齐）
	tokenDO, err := s.tokenSvc.CreateAccessToken(ctx, user.ID, UserTypeAdmin, tenantId, userInfo)
	if err != nil {
		return nil, core.ErrUnknown
	}

	// 6. 返回结果
	return &resp.AuthLoginResp{
		UserId:       user.ID,
		AccessToken:  tokenDO.AccessToken,
		RefreshToken: tokenDO.RefreshToken,
		ExpiresTime:  tokenDO.ExpiresTime,
	}, nil
}

// Logout 登出
func (s *AuthService) Logout(ctx context.Context, token string) error {
	// 1. 处理 token，移除 Bearer 前缀
	if strings.HasPrefix(strings.ToUpper(token), "BEARER ") {
		token = token[7:]
	}
	if token == "" {
		return nil
	}

	// 2. 使用 OAuth2TokenService 删除访问令牌
	_, _ = s.tokenSvc.RemoveAccessToken(ctx, token)

	// 3. 记录登出日志（可选，这里简化处理）
	return nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*resp.AuthLoginResp, error) {
	// 1. 验证 refreshToken（从 Redis 获取原令牌信息）
	oldToken, err := s.tokenSvc.GetAccessToken(ctx, refreshToken)
	if err != nil || oldToken == nil {
		return nil, core.NewBizError(1002000005, "刷新令牌无效或已过期")
	}

	// 2. 获取用户信息
	userRepo := s.repo.SystemUser
	user, err := userRepo.WithContext(ctx).Where(userRepo.ID.Eq(oldToken.UserID)).First()
	if err != nil {
		return nil, core.NewBizError(1002000002, "用户不存在")
	}

	// 3. 校验状态
	if user.Status != 0 {
		return nil, core.NewBizError(1002000001, "用户已被禁用")
	}

	// 4. 构建用户信息
	userInfo := map[string]string{
		"nickname": user.Nickname,
	}

	// 5. 创建新的访问令牌
	tokenDO, err := s.tokenSvc.CreateAccessToken(ctx, user.ID, oldToken.UserType, oldToken.TenantID, userInfo)
	if err != nil {
		return nil, core.ErrUnknown
	}

	// 6. 返回结果
	return &resp.AuthLoginResp{
		UserId:       user.ID,
		AccessToken:  tokenDO.AccessToken,
		RefreshToken: tokenDO.RefreshToken,
		ExpiresTime:  tokenDO.ExpiresTime,
	}, nil
}

// SmsLogin 短信登录
func (s *AuthService) SmsLogin(ctx context.Context, req *req.AuthSmsLoginReq) (*resp.AuthLoginResp, error) {
	// TODO: 验证短信验证码
	// 这里简化实现，实际需要调用短信服务验证

	// 1. 根据手机号查询用户
	userRepo := s.repo.SystemUser
	user, err := userRepo.WithContext(ctx).Where(userRepo.Mobile.Eq(req.Mobile)).First()
	if err != nil {
		return nil, core.NewBizError(1002000002, "用户不存在")
	}

	// 2. 校验状态
	if user.Status != 0 {
		return nil, core.NewBizError(1002000001, "用户已被禁用")
	}

	// 3. 构建用户信息
	userInfo := map[string]string{
		"nickname": user.Nickname,
	}

	// 4. 创建访问令牌（使用 OAuth2TokenService）
	tokenDO, err := s.tokenSvc.CreateAccessToken(ctx, user.ID, UserTypeAdmin, user.TenantID, userInfo)
	if err != nil {
		return nil, core.ErrUnknown
	}

	return &resp.AuthLoginResp{
		UserId:       user.ID,
		AccessToken:  tokenDO.AccessToken,
		RefreshToken: tokenDO.RefreshToken,
		ExpiresTime:  tokenDO.ExpiresTime,
	}, nil
}

// SendSmsCode 发送短信验证码
func (s *AuthService) SendSmsCode(ctx context.Context, req *req.AuthSmsSendReq) error {
	// TODO: 调用短信服务发送验证码
	// 这里简化实现，实际需要调用短信服务
	return nil
}

// Register 注册
func (s *AuthService) Register(ctx context.Context, req *req.AuthRegisterReq) (*resp.AuthLoginResp, error) {
	// 1. 检查用户名是否已存在
	userRepo := s.repo.SystemUser
	existUser, _ := userRepo.WithContext(ctx).Where(userRepo.Username.Eq(req.Username)).First()
	if existUser != nil {
		return nil, core.NewBizError(1002000006, "用户名已存在")
	}

	// 2. 创建用户
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, core.ErrUnknown
	}

	// TODO: 实际创建用户逻辑
	_ = hashedPassword

	return nil, core.NewBizError(1002000007, "注册功能暂未开放")
}

// ResetPassword 重置密码
func (s *AuthService) ResetPassword(ctx context.Context, req *req.AuthResetPasswordReq) error {
	// TODO: 验证短信验证码
	// 这里简化实现，实际需要调用短信服务验证

	// 1. 根据手机号查询用户
	userRepo := s.repo.SystemUser
	user, err := userRepo.WithContext(ctx).Where(userRepo.Mobile.Eq(req.Mobile)).First()
	if err != nil {
		return core.NewBizError(1002000002, "用户不存在")
	}

	// 2. 更新密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return core.ErrUnknown
	}

	_, err = userRepo.WithContext(ctx).Where(userRepo.ID.Eq(user.ID)).Update(userRepo.Password, hashedPassword)
	return err
}

// SocialAuthRedirect 社交授权跳转
func (s *AuthService) SocialAuthRedirect(ctx context.Context, socialType int, redirectUri string) (string, error) {
	// TODO: 调用社交服务获取授权 URL
	// 这里简化实现，返回空字符串
	return "", core.NewBizError(1002000008, "社交登录功能暂未开放")
}

// SocialLogin 社交登录
func (s *AuthService) SocialLogin(ctx context.Context, req *req.AuthSocialLoginReq) (*resp.AuthLoginResp, error) {
	// TODO: 调用社交服务验证并获取用户信息
	// 这里简化实现
	return nil, core.NewBizError(1002000008, "社交登录功能暂未开放")
}
