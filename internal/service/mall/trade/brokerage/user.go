package brokerage

import (
	"context"
	"time"

	trade2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"
	tradeReq "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/app/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"

	"go.uber.org/zap"
)

type BrokerageUserService struct {
	q         *query.Query
	logger    *zap.Logger
	memberSvc *member.MemberUserService
	configSvc *trade.TradeConfigService
}

func NewBrokerageUserService(q *query.Query, logger *zap.Logger, memberSvc *member.MemberUserService, configSvc *trade.TradeConfigService) *BrokerageUserService {
	return &BrokerageUserService{
		q:         q,
		logger:    logger,
		memberSvc: memberSvc,
		configSvc: configSvc,
	}
}

// GetBrokerageUser 获得分销用户
func (s *BrokerageUserService) GetBrokerageUser(ctx context.Context, id int64) (*brokerage.BrokerageUser, error) {
	return s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(id)).First()
}

// GetBindBrokerageUser 获得用户的推广人
// 对齐 Java: BrokerageUserServiceImpl.getBindBrokerageUser
func (s *BrokerageUserService) GetBindBrokerageUser(ctx context.Context, id int64) (*brokerage.BrokerageUser, error) {
	user, err := s.GetBrokerageUser(ctx, id)
	if err != nil || user == nil {
		return nil, err
	}
	if user.BindUserID == 0 {
		return nil, nil
	}
	return s.GetBrokerageUser(ctx, user.BindUserID)
}

// GetUserBrokerageEnabled 获得用户分销资格
// 对齐 Java: BrokerageUserServiceImpl.getUserBrokerageEnabled
func (s *BrokerageUserService) GetUserBrokerageEnabled(ctx context.Context, userId int64) (bool, error) {
	user, err := s.GetBrokerageUser(ctx, userId)
	if err != nil || user == nil {
		return false, err
	}
	return bool(user.BrokerageEnabled), nil
}

// UpdateUserPrice 更新用户可用佣金
// 对齐 Java: BrokerageUserServiceImpl.updateUserPrice
func (s *BrokerageUserService) UpdateUserPrice(ctx context.Context, id int64, price int) error {
	u := s.q.BrokerageUser
	_, err := u.WithContext(ctx).Where(u.ID.Eq(id)).UpdateSimple(u.BrokeragePrice.Add(price))
	return err
}

// UpdateUserFrozenPrice 更新用户冻结佣金
// 对齐 Java: BrokerageUserServiceImpl.updateUserFrozenPrice
func (s *BrokerageUserService) UpdateUserFrozenPrice(ctx context.Context, id int64, frozenPrice int) error {
	u := s.q.BrokerageUser
	_, err := u.WithContext(ctx).Where(u.ID.Eq(id)).UpdateSimple(u.FrozenPrice.Add(frozenPrice))
	return err
}

// UpdateFrozenPriceDecrAndPriceIncr 冻结佣金减少，可用佣金增加
// 对齐 Java: BrokerageUserServiceImpl.updateFrozenPriceDecrAndPriceIncr
func (s *BrokerageUserService) UpdateFrozenPriceDecrAndPriceIncr(ctx context.Context, id int64, frozenPrice int) error {
	u := s.q.BrokerageUser
	// frozenPrice 是负数（减少冻结），所以可用佣金增加的是 -frozenPrice
	_, err := u.WithContext(ctx).Where(u.ID.Eq(id)).UpdateSimple(
		u.FrozenPrice.Add(frozenPrice),
		u.BrokeragePrice.Add(-frozenPrice),
	)
	return err
}

func parseTime(t string) time.Time {
	res, _ := time.ParseInLocation(time.DateTime, t, time.Local)
	return res
}

// GetBrokerageUserPage 获得分销用户分页
func (s *BrokerageUserService) GetBrokerageUserPage(ctx context.Context, r *trade2.BrokerageUserPageReq) (*pagination.PageResult[*brokerage.BrokerageUser], error) {
	q := s.q.BrokerageUser.WithContext(ctx)

	// Filter by BindUserId and Level
	if r.BindUserID > 0 {
		childIDs, err := s.GetChildUserIdsByLevel(ctx, r.BindUserID, r.Level)
		if err != nil {
			return nil, err
		}
		if len(childIDs) == 0 {
			return &pagination.PageResult[*brokerage.BrokerageUser]{List: []*brokerage.BrokerageUser{}, Total: 0}, nil
		}
		q = q.Where(s.q.BrokerageUser.ID.In(childIDs...))
	}

	if r.BrokerageEnabled != nil {
		q = q.Where(s.q.BrokerageUser.BrokerageEnabled.Eq(model.BitBool(*r.BrokerageEnabled)))
	}

	// Time ranges
	if len(r.CreateTime) == 2 {
		q = q.Where(s.q.BrokerageUser.CreateTime.Between(parseTime(r.CreateTime[0]), parseTime(r.CreateTime[1])))
	}
	// BindUserTime?

	page := r.PageNo
	size := r.PageSize
	offset := (page - 1) * size

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Limit(size).Offset(offset).Order(s.q.BrokerageUser.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*brokerage.BrokerageUser]{
		List:  list,
		Total: total,
	}, nil
}

// CreateBrokerageUser 创建分销用户（Admin 手动创建）
func (s *BrokerageUserService) CreateBrokerageUser(ctx context.Context, r *trade2.BrokerageUserCreateReq) (int64, error) {
	// Convert FlexInt64 to int64
	userID := int64(r.UserID)
	bindUserID := int64(r.BindUserID)

	// 1.1 校验分销用户是否已存在
	exists, _ := s.GetBrokerageUser(ctx, userID)
	if exists != nil {
		return 0, trade.ErrBrokerageUserExists()
	}

	// 1.2 校验是否能绑定用户
	user := &brokerage.BrokerageUser{ID: userID}
	if err := s.validateCanBindUser(ctx, user, bindUserID); err != nil {
		return 0, err
	}

	// 2. 创建分销人
	now := time.Now()
	newUser := &brokerage.BrokerageUser{
		ID:            userID,
		BindUserID:    bindUserID,
		BindUserTime:  &now,
		BrokerageTime: &now,
	}
	if err := s.q.BrokerageUser.WithContext(ctx).Create(newUser); err != nil {
		return 0, err
	}
	return newUser.ID, nil
}

// GetOrCreateBrokerageUser 获得或创建分销用户
// 特殊：人人分销的情况下，如果分销人为空则创建分销人
func (s *BrokerageUserService) GetOrCreateBrokerageUser(ctx context.Context, id int64) (*brokerage.BrokerageUser, error) {
	user, _ := s.GetBrokerageUser(ctx, id)
	if user != nil {
		return user, nil
	}

	// 获取分销配置
	config, err := s.configSvc.GetTradeConfig(ctx)
	if err != nil {
		return nil, err
	}

	// 人人分销（BrokerageEnabledCondition = 1）的情况下才自动创建
	if config == nil || config.BrokerageEnabledCondition != 1 {
		return nil, nil
	}

	now := time.Now()
	newUser := &brokerage.BrokerageUser{
		ID:               id,
		BrokerageEnabled: true,
		BrokeragePrice:   0,
		FrozenPrice:      0,
		BrokerageTime:    &now,
	}
	if err := s.q.BrokerageUser.WithContext(ctx).Create(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

// BindBrokerageUser 绑定推广员
func (s *BrokerageUserService) BindBrokerageUser(ctx context.Context, userId int64, bindUserId int64) (bool, error) {
	// 1. 获得分销用户
	isNewBrokerageUser := false
	user, _ := s.GetBrokerageUser(ctx, userId)
	if user == nil {
		// 分销用户不存在的情况：1. 新注册；2. 旧数据；3. 分销功能关闭后又打开
		isNewBrokerageUser = true
		user = &brokerage.BrokerageUser{
			ID:               userId,
			BrokerageEnabled: false,
			BrokeragePrice:   0,
			FrozenPrice:      0,
		}
	}

	// 2.1 校验是否能绑定用户
	canBind, err := s.isUserCanBind(ctx, user)
	if err != nil {
		return false, err
	}
	if !canBind {
		return false, nil
	}

	// 2.2 校验能否绑定
	if err := s.validateCanBindUser(ctx, user, bindUserId); err != nil {
		return false, err
	}

	// 2.3 绑定用户
	now := time.Now()
	if isNewBrokerageUser {
		// 获取分销配置，判断是否人人分销
		config, _ := s.configSvc.GetTradeConfig(ctx)
		if config != nil && config.BrokerageEnabledCondition == consts.BrokerageEnabledConditionAll {
			// 人人分销：用户默认就有分销资格
			user.BrokerageEnabled = true
			user.BrokerageTime = &now
		} else {
			user.BrokerageEnabled = false
			user.BrokerageTime = &now
		}
		user.BindUserID = bindUserId
		user.BindUserTime = &now
		if err := s.q.BrokerageUser.WithContext(ctx).Create(user); err != nil {
			return false, err
		}
	} else {
		_, err := s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(userId)).Updates(map[string]interface{}{
			"bind_user_id":   bindUserId,
			"bind_user_time": &now,
		})
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// UpdateBrokerageUserEnabled 修改推广资格
func (s *BrokerageUserService) UpdateBrokerageUserEnabled(ctx context.Context, id int64, enabled bool) error {
	u, err := s.GetBrokerageUser(ctx, id)
	if err != nil {
		return err
	}
	if bool(u.BrokerageEnabled) == enabled {
		return nil
	}

	updates := map[string]interface{}{
		"brokerage_enabled": enabled,
	}
	if enabled {
		now := time.Now()
		updates["brokerage_time"] = &now
	} else {
		updates["brokerage_time"] = nil
	}
	_, err = s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(id)).Updates(updates)
	return err
}

// UpdateBrokerageUserId 修改推广员
func (s *BrokerageUserService) UpdateBrokerageUserId(ctx context.Context, id int64, bindUserId int64) error {
	user, err := s.GetBrokerageUser(ctx, id)
	if err != nil {
		return err
	}
	if user.BindUserID == bindUserId {
		return nil
	}

	// Clear
	if bindUserId == 0 {
		_, err = s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(id)).Updates(map[string]interface{}{
			"bind_user_id":   0, // Or null
			"bind_user_time": nil,
		})
		return err
	}

	// Validate
	if err := s.validateCanBindUser(ctx, user, bindUserId); err != nil {
		return err
	}

	now := time.Now()
	_, err = s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(id)).Updates(map[string]interface{}{
		"bind_user_id":   bindUserId,
		"bind_user_time": &now,
	})
	return err
}

// isUserCanBind 校验是否能绑定用户
func (s *BrokerageUserService) isUserCanBind(ctx context.Context, user *brokerage.BrokerageUser) (bool, error) {
	// 校验分销功能是否启用
	config, err := s.configSvc.GetTradeConfig(ctx)
	if err != nil {
		return false, err
	}
	if config == nil || !config.BrokerageEnabled {
		return false, nil
	}

	// 校验分销关系绑定模式
	// BrokerageBindMode: 1=首次绑定 2=注册绑定 3=覆盖绑定
	switch config.BrokerageBindMode {
	case 2:
		// 注册绑定模式：判断是否为新用户（注册时间在 30 秒内）
		memberUser, err := s.memberSvc.GetUser(ctx, user.ID)
		if err != nil || memberUser == nil {
			return false, nil
		}
		if time.Since(memberUser.CreateTime) > 30*time.Second {
			return false, trade.ErrBrokerageBindModeRegister()
		}
	case 3:
		// 覆盖绑定模式：不允许
		if user.BindUserID > 0 {
			return false, trade.ErrBrokerageBindOverride()
		}
	}
	// 首次绑定模式 (默认)：如果已绑定则返回 false
	if user.BindUserID > 0 {
		return false, nil
	}
	return true, nil
}

// validateCanBindUser 校验能否绑定
func (s *BrokerageUserService) validateCanBindUser(ctx context.Context, user *brokerage.BrokerageUser, bindUserId int64) error {
	if bindUserId == 0 {
		return nil
	}
	// 1.1 校验推广人是否存在
	bindUser, err := s.memberSvc.GetUser(ctx, bindUserId)
	if err != nil || bindUser == nil {
		return trade.ErrBrokerageUserNotExists()
	}

	// 1.2 校验要绑定的用户有无推广资格
	brokerageBindUser, _ := s.GetOrCreateBrokerageUser(ctx, bindUserId)
	if brokerageBindUser == nil || !brokerageBindUser.BrokerageEnabled {
		return trade.ErrBrokerageBindUserNotEnabled()
	}

	// 2. 校验绑定自己
	if user.ID == bindUserId {
		return trade.ErrBrokerageBindSelf()
	}

	// 3. 下级不能绑定自己的上级（循环绑定检查）
	currentBindId := brokerageBindUser.BindUserID
	for i := 0; i < 32767; i++ { // Short.MAX_VALUE
		if currentBindId == 0 {
			break
		}
		if currentBindId == user.ID {
			return trade.ErrBrokerageBindLoop()
		}
		next, _ := s.GetBrokerageUser(ctx, currentBindId)
		if err != nil || next == nil {
			break
		}
		currentBindId = next.BindUserID
	}
	return nil
}

// GetChildUserIdsByLevel 获得下级用户编号列表
func (s *BrokerageUserService) GetChildUserIdsByLevel(ctx context.Context, bindUserId int64, level int) ([]int64, error) {
	// Level 1
	var level1Ids []int64
	err := s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.BindUserID.Eq(bindUserId)).Select(s.q.BrokerageUser.ID).Scan(&level1Ids)
	if err != nil {
		return nil, err
	}
	if len(level1Ids) == 0 {
		return []int64{}, nil
	}

	if level == consts.BrokerageUserLevelOne {
		return level1Ids, nil
	}

	// Level 2
	var level2Ids []int64
	if len(level1Ids) > 0 {
		err = s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.BindUserID.In(level1Ids...)).Select(s.q.BrokerageUser.ID).Scan(&level2Ids)
		if err != nil {
			return nil, err
		}
	}

	if level == consts.BrokerageUserLevelTwo {
		return level2Ids, nil
	}

	// All (1 + 2)
	return append(level1Ids, level2Ids...), nil
}

// GetBrokerageUserCountByBindUserId 获得推广用户数量
func (s *BrokerageUserService) GetBrokerageUserCountByBindUserId(ctx context.Context, bindUserId int64, level int) (int64, error) {
	ids, err := s.GetChildUserIdsByLevel(ctx, bindUserId, level)
	if err != nil {
		return 0, err
	}
	return int64(len(ids)), nil
}

// GetBrokerageUserChildSummaryPage 获得下级分销统计分页
func (s *BrokerageUserService) GetBrokerageUserChildSummaryPage(ctx context.Context, r *tradeReq.AppBrokerageUserChildSummaryPageReqVO, userId int64) (*pagination.PageResult[*brokerage.BrokerageUser], error) {
	childIDs, err := s.GetChildUserIdsByLevel(ctx, userId, r.Level)
	if err != nil {
		return nil, err
	}
	if len(childIDs) == 0 {
		return &pagination.PageResult[*brokerage.BrokerageUser]{List: []*brokerage.BrokerageUser{}, Total: 0}, nil
	}

	q := s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.In(childIDs...))

	// 注：昵称过滤需要与 MemberUser 表关联，此处简化实现，昵称信息由调用层补充
	// Java 实现通过 join 或 map 方式获取用户信息
	// 排序逻辑:
	switch r.Sorting {
	case "userCount":
		// 按用户数量排序：复杂查询，暂用默认排序
	case "brokeragePrice":
		q = q.Order(s.q.BrokerageUser.BrokeragePrice.Desc())
	default:
		q = q.Order(s.q.BrokerageUser.BrokerageTime.Desc())
	}

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	list, err := q.Limit(r.PageSize).Offset((r.PageNo - 1) * r.PageSize).Find()
	if err != nil {
		return nil, err
	}
	return &pagination.PageResult[*brokerage.BrokerageUser]{List: list, Total: total}, nil
}

// BrokerageUserRankByUserCountResult 分销用户排行（基于用户量）结果
type BrokerageUserRankByUserCountResult struct {
	ID                 int64 `gorm:"column:id"`
	BrokerageUserCount int   `gorm:"column:brokerage_user_count"`
}

// GetBrokerageUserRankPageByUserCount 获得分销用户排行分页（基于用户量）
func (s *BrokerageUserService) GetBrokerageUserRankPageByUserCount(ctx context.Context, r *tradeReq.AppBrokerageUserRankPageReqVO) (*pagination.PageResult[*BrokerageUserRankByUserCountResult], error) {
	// 解析时间范围
	var beginTime, endTime time.Time
	if len(r.Times) >= 2 {
		beginTime = parseTime(r.Times[0])
		endTime = parseTime(r.Times[1])
	}

	// 使用 Gen 生成的字段和表名构建查询
	bu := s.q.BrokerageUser
	tableName := bu.TableName()
	bindUserIDCol := bu.BindUserID.ColumnName().String()
	bindUserTimeCol := bu.BindUserTime.ColumnName().String()

	// 构建基础查询条件
	q := bu.WithContext(ctx).
		Where(bu.BindUserID.Gt(0))

	if !beginTime.IsZero() && !endTime.IsZero() {
		q = q.Where(bu.BindUserTime.Between(beginTime, endTime))
	}

	// 获取总数（不同绑定用户数）
	// 使用 GORM Gen 的 Distinct 进行统计
	total, err := q.Distinct(bu.BindUserID).Count()
	if err != nil {
		return nil, err
	}

	// 分组查询需要使用原生 GORM，因为 Gen 不支持 GROUP BY 聚合
	// 使用 Gen 生成的字段名确保类型安全
	db := bu.WithContext(ctx).UnderlyingDB()
	offset := (r.PageNo - 1) * r.PageSize

	var results []*BrokerageUserRankByUserCountResult
	selectClause := bindUserIDCol + " as id, COUNT(*) as brokerage_user_count"
	query := db.Table(tableName).
		Select(selectClause).
		Where(bindUserIDCol+" > ?", 0).
		Where("deleted = 0")

	if !beginTime.IsZero() && !endTime.IsZero() {
		query = query.Where(bindUserTimeCol+" BETWEEN ? AND ?", beginTime, endTime)
	}

	query.Group(bindUserIDCol).
		Order("brokerage_user_count DESC").
		Limit(r.PageSize).
		Offset(offset).
		Scan(&results)

	return &pagination.PageResult[*BrokerageUserRankByUserCountResult]{
		List:  results,
		Total: total,
	}, nil
}
