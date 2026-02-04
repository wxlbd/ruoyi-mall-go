package brokerage

import (
	"context"
	"errors"
	"strings"
	"time"

	trade2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/app/mall/trade"
	tradeModel "github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/product"
	tradeSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"go.uber.org/zap"
)

type BrokerageRecordService struct {
	q              *query.Query
	logger         *zap.Logger
	tradeConfigSvc *tradeSvc.TradeConfigService
	spuSvc         *product.ProductSpuService
	skuSvc         *product.ProductSkuService
}

func NewBrokerageRecordService(q *query.Query, logger *zap.Logger, tradeConfigSvc *tradeSvc.TradeConfigService, spuSvc *product.ProductSpuService, skuSvc *product.ProductSkuService) *BrokerageRecordService {
	return &BrokerageRecordService{
		q:              q,
		logger:         logger,
		tradeConfigSvc: tradeConfigSvc,
		spuSvc:         spuSvc,
		skuSvc:         skuSvc,
	}
}

// GetSummaryPriceByUserId 获得分销佣金统计
func (s *BrokerageRecordService) GetSummaryPriceByUserId(ctx context.Context, userId int64, bizType int, status int, beginTime, endTime time.Time) (int, error) {
	q := s.q.BrokerageRecord.WithContext(ctx).
		Where(s.q.BrokerageRecord.UserID.Eq(userId))

	if bizType > 0 {
		q = q.Where(s.q.BrokerageRecord.BizType.Eq(bizType))
	}
	if status > 0 {
		q = q.Where(s.q.BrokerageRecord.Status.Eq(status))
	}
	if !beginTime.IsZero() && !endTime.IsZero() {
		q = q.Where(s.q.BrokerageRecord.CreateTime.Between(beginTime, endTime))
	}

	var sum int
	err := q.Select(s.q.BrokerageRecord.Price.Sum()).Scan(&sum)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// GetBrokerageRecord 获得分销记录
func (s *BrokerageRecordService) GetBrokerageRecord(ctx context.Context, id int64) (*brokerage.BrokerageRecord, error) {
	return s.q.BrokerageRecord.WithContext(ctx).Where(s.q.BrokerageRecord.ID.Eq(id)).First()
}

// GetBrokerageRecordPage 获得分销记录分页
func (s *BrokerageRecordService) GetBrokerageRecordPage(ctx context.Context, r *trade2.BrokerageRecordPageReq) (*pagination.PageResult[*brokerage.BrokerageRecord], error) {
	q := s.q.BrokerageRecord.WithContext(ctx)

	if r.UserID != nil && *r.UserID > 0 {
		q = q.Where(s.q.BrokerageRecord.UserID.Eq(*r.UserID))
	}
	if r.BizType != nil {
		q = q.Where(s.q.BrokerageRecord.BizType.Eq(*r.BizType))
	}
	if r.Status != nil {
		q = q.Where(s.q.BrokerageRecord.Status.Eq(*r.Status))
	}
	if r.SourceUserLevel != nil {
		q = q.Where(s.q.BrokerageRecord.SourceUserLevel.Eq(*r.SourceUserLevel))
	}
	if r.BizID != "" {
		q = q.Where(s.q.BrokerageRecord.BizID.Eq(r.BizID))
	}
	if begin, end, ok := parseTimeRange(r.CreateTime); ok {
		q = q.Where(s.q.BrokerageRecord.CreateTime.Between(begin, end))
	}

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	offset := (r.PageNo - 1) * r.PageSize
	list, err := q.Limit(r.PageSize).Offset(offset).Order(s.q.BrokerageRecord.ID.Desc()).Find()
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[*brokerage.BrokerageRecord]{
		List:  list,
		Total: total,
	}, nil
}

// mapBizTypeToInt 将业务类型字符串映射为整数
func (s *BrokerageRecordService) mapBizTypeToInt(bizType string) int {
	switch strings.TrimSpace(bizType) {
	case "1", "order":
		return tradeModel.BrokerageRecordBizTypeOrder
	case "2", "withdraw":
		return tradeModel.BrokerageRecordBizTypeWithdraw
	case "3", "withdraw_reject":
		return tradeModel.BrokerageRecordBizTypeWithdrawReject
	default:
		return 0
	}
}

func parseTimeRange(values []string) (time.Time, time.Time, bool) {
	if len(values) != 2 {
		return time.Time{}, time.Time{}, false
	}
	begin, err := time.ParseInLocation(time.DateTime, strings.TrimSpace(values[0]), time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, false
	}
	end, err := time.ParseInLocation(time.DateTime, strings.TrimSpace(values[1]), time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, false
	}
	return begin, end, true
}

// AddBrokerage 添加分销记录（增加佣金）
func (s *BrokerageRecordService) AddBrokerage(ctx context.Context, userID int64, bizType string, bizID string, price int, title string) error {
	// 1. 更新用户佣金余额
	u := s.q.BrokerageUser
	if _, err := u.WithContext(ctx).Where(u.ID.Eq(userID)).UpdateSimple(u.BrokeragePrice.Add(price)); err != nil {
		return err
	}

	// 2. 创建佣金记录
	bizTypeInt := s.mapBizTypeToInt(bizType)
	record := &brokerage.BrokerageRecord{
		UserID:      userID,
		BizType:     bizTypeInt,
		BizID:       bizID,
		Price:       price,
		Title:       title,
		Description: title,
		Status:      tradeModel.BrokerageRecordStatusSettlement,
		TotalPrice:  price,
	}
	return s.q.BrokerageRecord.WithContext(ctx).Create(record)
}

// ReduceBrokerageForWithdraw 提现扣减佣金
func (s *BrokerageRecordService) ReduceBrokerageForWithdraw(ctx context.Context, userID int64, bizID string, price int) error {
	// 1. 校验佣金余额并更新（原子操作）
	u := s.q.BrokerageUser
	info, err := u.WithContext(ctx).Where(u.ID.Eq(userID)).UpdateSimple(
		u.BrokeragePrice.Sub(price),
		u.FrozenPrice.Add(price),
	)
	if err != nil {
		return err
	}
	if info.RowsAffected == 0 {
		return errors.New("佣金不足")
	}

	// 2. 创建佣金记录
	record := &brokerage.BrokerageRecord{
		UserID:      userID,
		BizType:     tradeModel.BrokerageRecordBizTypeWithdraw,
		BizID:       bizID,
		Price:       -price, // 负数表示扣减
		Title:       "佣金提现",
		Description: "佣金提现",
		Status:      tradeModel.BrokerageRecordStatusSettlement,
		TotalPrice:  price,
	}
	return s.q.BrokerageRecord.WithContext(ctx).Create(record)
}

// CalculateProductBrokeragePrice 计算商品佣金
func (s *BrokerageRecordService) CalculateProductBrokeragePrice(ctx context.Context, userId int64, spuId int64) (*trade.AppBrokerageProductPriceRespVO, error) {
	resp := &trade.AppBrokerageProductPriceRespVO{
		BrokerageEnabled: false,
		BrokeragePrice:   0,
	}

	// 1. 校验分销功能是否开启
	config, err := s.tradeConfigSvc.GetTradeConfig(ctx)
	if err != nil || config == nil || !config.BrokerageEnabled {
		return resp, nil
	}

	// 2. 校验用户是否有分销资格
	user, err := s.q.BrokerageUser.WithContext(ctx).Where(s.q.BrokerageUser.ID.Eq(userId)).First()
	if err != nil || user == nil || !user.BrokerageEnabled {
		return resp, nil
	}
	resp.BrokerageEnabled = true

	// 3. 校验商品是否存在
	spu, err := s.spuSvc.GetSpu(ctx, spuId)
	if err != nil || spu == nil {
		return resp, nil
	}

	// 4. 计算佣金
	skus, err := s.skuSvc.GetSkuListBySpuId(ctx, spuId)
	if err != nil {
		return resp, nil
	}

	minPrice := 0
	maxPrice := 0
	percent := config.BrokerageFirstPercent

	for _, sku := range skus {
		var brokeragePrice int
		if spu.SubCommissionType {
			// 商品单独分佣模式：使用 SKU 的固定佣金
			brokeragePrice = sku.FirstBrokeragePrice
		} else {
			// 全局分佣模式：根据商品价格比例计算
			brokeragePrice = sku.Price * percent / 100
		}

		if minPrice == 0 || brokeragePrice < minPrice {
			minPrice = brokeragePrice
		}
		if brokeragePrice > maxPrice {
			maxPrice = brokeragePrice
		}
	}

	// 使用最大佣金作为展示值
	resp.BrokeragePrice = maxPrice

	return resp, nil
}

// GetBrokerageUserRankPageByPrice 获得分销用户排行分页（基于佣金）
func (s *BrokerageRecordService) GetBrokerageUserRankPageByPrice(ctx context.Context, r *trade2.AppBrokerageUserRankPageReq) (*pagination.PageResult[*trade.AppBrokerageUserRankByPriceRespVO], error) {
	// 解析时间范围
	var beginTime, endTime time.Time
	if len(r.Times) >= 2 {
		beginTime = parseTime(r.Times[0])
		endTime = parseTime(r.Times[1])
	}

	// 使用 Gen 生成的字段和表名
	br := s.q.BrokerageRecord
	tableName := br.TableName()
	userIDCol := br.UserID.ColumnName().String()
	priceCol := br.Price.ColumnName().String()
	bizTypeCol := br.BizType.ColumnName().String()
	statusCol := br.Status.ColumnName().String()
	createTimeCol := br.CreateTime.ColumnName().String()

	// 构建基础查询条件统计总数
	q := br.WithContext(ctx).
		Where(br.BizType.Eq(tradeModel.BrokerageRecordBizTypeOrder)).
		Where(br.Status.Eq(tradeModel.BrokerageRecordStatusSettlement))

	if !beginTime.IsZero() && !endTime.IsZero() {
		q = q.Where(br.CreateTime.Between(beginTime, endTime))
	}

	// 获取总数 (不同用户数)
	total, err := q.Distinct(br.UserID).Count()
	if err != nil {
		return nil, err
	}

	// 分组查询需要使用原生 GORM，因为 Gen 不支持 GROUP BY 聚合
	db := br.WithContext(ctx).UnderlyingDB()
	offset := (r.PageNo - 1) * r.PageSize

	type RankResult struct {
		UserID int64 `gorm:"column:user_id"`
		Price  int   `gorm:"column:price"`
	}

	selectClause := userIDCol + ", SUM(" + priceCol + ") as price"
	query := db.Table(tableName).
		Select(selectClause).
		Where(bizTypeCol+" = ? AND "+statusCol+" = ?", tradeModel.BrokerageRecordBizTypeOrder, tradeModel.BrokerageRecordStatusSettlement).
		Where("deleted = 0")

	if !beginTime.IsZero() && !endTime.IsZero() {
		query = query.Where(createTimeCol+" BETWEEN ? AND ?", beginTime, endTime)
	}

	var results []RankResult
	query.Group(userIDCol).
		Order("price DESC").
		Limit(r.PageSize).
		Offset(offset).
		Scan(&results)

	// 转换为 VO
	list := make([]*trade.AppBrokerageUserRankByPriceRespVO, len(results))
	for i, item := range results {
		list[i] = &trade.AppBrokerageUserRankByPriceRespVO{
			ID:             item.UserID,
			BrokeragePrice: item.Price,
		}
	}

	return &pagination.PageResult[*trade.AppBrokerageUserRankByPriceRespVO]{
		List:  list,
		Total: total,
	}, nil
}

// GetUserRankByPrice 获得分销用户排行（基于佣金）
func (s *BrokerageRecordService) GetUserRankByPrice(ctx context.Context, userId int64, times []time.Time) (int, error) {
	var beginTime, endTime time.Time
	if len(times) >= 2 {
		beginTime = times[0]
		endTime = times[1]
	}

	// 1. 获取用户的分销佣金总额
	userPrice, err := s.GetSummaryPriceByUserId(ctx, userId, tradeModel.BrokerageRecordBizTypeOrder, tradeModel.BrokerageRecordStatusSettlement, beginTime, endTime)
	if err != nil {
		return 0, err
	}

	// 使用 Gen 生成的字段和表名
	br := s.q.BrokerageRecord
	tableName := br.TableName()
	userIDCol := br.UserID.ColumnName().String()
	priceCol := br.Price.ColumnName().String()
	bizTypeCol := br.BizType.ColumnName().String()
	statusCol := br.Status.ColumnName().String()
	createTimeCol := br.CreateTime.ColumnName().String()

	// 2. 获取比用户佣金高的用户数量
	db := br.WithContext(ctx).UnderlyingDB()

	// 子查询：获取每个用户的佣金总额大于当前用户的数量
	subQuery := db.Table(tableName).
		Select(userIDCol+", SUM("+priceCol+") as total_price").
		Where(bizTypeCol+" = ? AND "+statusCol+" = ?", tradeModel.BrokerageRecordBizTypeOrder, tradeModel.BrokerageRecordStatusSettlement).
		Where("deleted = 0")
	if !beginTime.IsZero() && !endTime.IsZero() {
		subQuery = subQuery.Where(createTimeCol+" BETWEEN ? AND ?", beginTime, endTime)
	}
	subQuery = subQuery.Group(userIDCol).Having("SUM("+priceCol+") > ?", userPrice)

	var greaterCount int64
	db.Table("(?) as ranked", subQuery).Count(&greaterCount)

	// 3. 返回排名 (比自己高的人数 + 1)
	return int(greaterCount) + 1, nil
}

// BrokerageAddReqBO 分佣请求参数
// 对齐 Java: cn.iocoder.yudao.module.trade.service.brokerage.bo.BrokerageAddReqBO
type BrokerageAddReqBO struct {
	BizID            string // 业务编号
	BasePrice        int    // 分佣基础价格
	FirstFixedPrice  int    // 一级固定佣金
	SecondFixedPrice int    // 二级固定佣金
	Title            string // 标题
	SourceUserId     int64  // 来源用户编号（下单用户）
}

// AddBrokerageMultiLevel 添加多级分佣（一级/二级推广人分佣）
// 对齐 Java: BrokerageRecordServiceImpl.addBrokerage(userId, bizType, list)
func (s *BrokerageRecordService) AddBrokerageMultiLevel(ctx context.Context, sourceUserID int64, bizType int, list []*BrokerageAddReqBO, brokerageUserSvc *BrokerageUserService) error {
	// 0. 校验分销功能是否启用
	config, err := s.tradeConfigSvc.GetTradeConfig(ctx)
	if err != nil || config == nil || !config.BrokerageEnabled {
		s.logger.Warn("[AddBrokerageMultiLevel] 分销功能未开启")
		return nil
	}

	// 1.1 获得一级推广人
	firstUser, err := brokerageUserSvc.GetBindBrokerageUser(ctx, sourceUserID)
	if err != nil || firstUser == nil || !firstUser.BrokerageEnabled {
		return nil
	}

	// 1.2 计算一级分佣
	if err := s.addBrokerageForLevel(ctx, firstUser, list, config.BrokerageFrozenDays, config.BrokerageFirstPercent, bizType, 1, sourceUserID, brokerageUserSvc); err != nil {
		s.logger.Error("[AddBrokerageMultiLevel] 一级分佣失败", zap.Error(err))
	}

	// 2.1 获得二级推广员
	if firstUser.BindUserID == 0 {
		return nil
	}
	secondUser, err := brokerageUserSvc.GetBrokerageUser(ctx, firstUser.BindUserID)
	if err != nil || secondUser == nil || !secondUser.BrokerageEnabled {
		return nil
	}

	// 2.2 计算二级分佣
	if err := s.addBrokerageForLevel(ctx, secondUser, list, config.BrokerageFrozenDays, config.BrokerageSecondPercent, bizType, 2, sourceUserID, brokerageUserSvc); err != nil {
		s.logger.Error("[AddBrokerageMultiLevel] 二级分佣失败", zap.Error(err))
	}

	return nil
}

// addBrokerageForLevel 为指定级别用户添加佣金
func (s *BrokerageRecordService) addBrokerageForLevel(ctx context.Context, user *brokerage.BrokerageUser, list []*BrokerageAddReqBO, frozenDays int, percent int, bizType int, level int, sourceUserID int64, brokerageUserSvc *BrokerageUserService) error {
	// 计算解冻时间
	var unfreezeTime *time.Time
	if frozenDays > 0 {
		t := time.Now().AddDate(0, 0, frozenDays)
		unfreezeTime = &t
	}

	totalBrokerage := 0
	var records []*brokerage.BrokerageRecord

	for _, item := range list {
		// 计算佣金金额
		var fixedPrice int
		if level == 1 {
			fixedPrice = item.FirstFixedPrice
		} else {
			fixedPrice = item.SecondFixedPrice
		}

		brokeragePrice := s.calculatePrice(item.BasePrice, percent, fixedPrice)
		if brokeragePrice <= 0 {
			continue
		}
		totalBrokerage += brokeragePrice

		// 确定状态
		status := tradeModel.BrokerageRecordStatusSettlement
		if frozenDays > 0 {
			status = tradeModel.BrokerageRecordStatusWait
		}

		records = append(records, &brokerage.BrokerageRecord{
			UserID:          user.ID,
			BizType:         bizType,
			BizID:           item.BizID,
			Price:           brokeragePrice,
			TotalPrice:      brokeragePrice,
			Title:           item.Title,
			Description:     item.Title,
			Status:          status,
			FrozenDays:      frozenDays,
			UnfreezeTime:    unfreezeTime,
			SourceUserLevel: level,
			SourceUserID:    item.SourceUserId, // 对齐 Java: 使用 item 中的 sourceUserId
		})
	}

	if len(records) == 0 {
		return nil
	}

	// 批量插入记录
	if err := s.q.BrokerageRecord.WithContext(ctx).CreateInBatches(records, 100); err != nil {
		return err
	}

	// 更新用户佣金（冻结或可用）
	if frozenDays > 0 {
		return brokerageUserSvc.UpdateUserFrozenPrice(ctx, user.ID, totalBrokerage)
	}
	return brokerageUserSvc.UpdateUserPrice(ctx, user.ID, totalBrokerage)
}

// calculatePrice 计算佣金价格
func (s *BrokerageRecordService) calculatePrice(basePrice int, percent int, fixedPrice int) int {
	if fixedPrice > 0 {
		return fixedPrice
	}
	return basePrice * percent / 100
}

// CancelBrokerage 取消佣金（用于退款/订单取消）
// 对齐 Java: BrokerageRecordServiceImpl.cancelBrokerage
func (s *BrokerageRecordService) CancelBrokerage(ctx context.Context, bizType int, bizID string, brokerageUserSvc *BrokerageUserService) error {
	// 1. 查询佣金记录
	records, err := s.q.BrokerageRecord.WithContext(ctx).
		Where(s.q.BrokerageRecord.BizType.Eq(bizType)).
		Where(s.q.BrokerageRecord.BizID.Eq(bizID)).
		Find()
	if err != nil {
		return err
	}
	if len(records) == 0 {
		s.logger.Warn("[CancelBrokerage] 记录不存在", zap.Int("bizType", bizType), zap.String("bizID", bizID))
		return nil
	}

	// 2. 遍历更新
	for _, record := range records {
		// 2.1 更新记录状态为已失效
		info, err := s.q.BrokerageRecord.WithContext(ctx).
			Where(s.q.BrokerageRecord.ID.Eq(record.ID)).
			Where(s.q.BrokerageRecord.Status.Eq(record.Status)).
			Update(s.q.BrokerageRecord.Status, tradeModel.BrokerageRecordStatusCancel)
		if err != nil {
			s.logger.Error("[CancelBrokerage] 更新记录失败", zap.Int64("id", record.ID), zap.Error(err))
			continue
		}
		if info.RowsAffected == 0 {
			s.logger.Warn("[CancelBrokerage] 更新记录失败（乐观锁）", zap.Int64("id", record.ID))
			continue
		}

		// 2.2 更新用户佣金（回滚）
		if record.Status == tradeModel.BrokerageRecordStatusWait {
			// 待结算状态：减少冻结佣金
			_ = brokerageUserSvc.UpdateUserFrozenPrice(ctx, record.UserID, -record.Price)
		} else if record.Status == tradeModel.BrokerageRecordStatusSettlement {
			// 已结算状态：减少可用佣金
			_ = brokerageUserSvc.UpdateUserPrice(ctx, record.UserID, -record.Price)
		}
	}
	return nil
}

// UnfreezeRecord 解冻佣金（定时任务调用）
// 对齐 Java: BrokerageRecordServiceImpl.unfreezeRecord
func (s *BrokerageRecordService) UnfreezeRecord(ctx context.Context, brokerageUserSvc *BrokerageUserService) (int, error) {
	// 1. 查询待结算且已到解冻时间的记录
	now := time.Now()
	records, err := s.q.BrokerageRecord.WithContext(ctx).
		Where(s.q.BrokerageRecord.Status.Eq(tradeModel.BrokerageRecordStatusWait)).
		Where(s.q.BrokerageRecord.UnfreezeTime.Lt(now)).
		Find()
	if err != nil {
		return 0, err
	}
	if len(records) == 0 {
		return 0, nil
	}

	// 2. 遍历执行解冻
	count := 0
	for _, record := range records {
		success, err := s.unfreezeSingleRecord(ctx, record, brokerageUserSvc)
		if err != nil {
			s.logger.Error("[UnfreezeRecord] 解冻失败", zap.Int64("id", record.ID), zap.Error(err))
			continue
		}
		if success {
			count++
		}
	}
	return count, nil
}

// unfreezeSingleRecord 解冻单条记录
func (s *BrokerageRecordService) unfreezeSingleRecord(ctx context.Context, record *brokerage.BrokerageRecord, brokerageUserSvc *BrokerageUserService) (bool, error) {
	// 1. 更新记录状态为已结算
	now := time.Now()
	info, err := s.q.BrokerageRecord.WithContext(ctx).
		Where(s.q.BrokerageRecord.ID.Eq(record.ID)).
		Where(s.q.BrokerageRecord.Status.Eq(record.Status)).
		Updates(map[string]interface{}{
			"status":        tradeModel.BrokerageRecordStatusSettlement,
			"unfreeze_time": now,
		})
	if err != nil {
		return false, err
	}
	if info.RowsAffected == 0 {
		s.logger.Warn("[unfreezeSingleRecord] 更新失败（乐观锁）", zap.Int64("id", record.ID))
		return false, nil
	}

	// 2. 更新用户冻结佣金减少，可用佣金增加
	if err := brokerageUserSvc.UpdateFrozenPriceDecrAndPriceIncr(ctx, record.UserID, -record.Price); err != nil {
		s.logger.Error("[unfreezeSingleRecord] 更新用户佣金失败", zap.Int64("userId", record.UserID), zap.Error(err))
		return false, err
	}

	s.logger.Info("[unfreezeSingleRecord] 解冻成功", zap.Int64("id", record.ID))
	return true, nil
}
