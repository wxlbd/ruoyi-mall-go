package promotion

import (
	"context"
	"time"

	"backend-go/internal/api/req"
	"backend-go/internal/model/promotion"
	"backend-go/internal/pkg/core"
	"backend-go/internal/repo/query"

	"gorm.io/gorm/clause"
)

type BargainHelpService struct {
	q *query.Query
}

func NewBargainHelpService(q *query.Query) *BargainHelpService {
	return &BargainHelpService{q: q}
}

// GetBargainHelpUserCountMapByActivity 获得砍价活动的助力用户数量 Map
func (s *BargainHelpService) GetBargainHelpUserCountMapByActivity(ctx context.Context, activityIds []int64) (map[int64]int, error) {
	if len(activityIds) == 0 {
		return make(map[int64]int), nil
	}
	q := s.q.PromotionBargainHelp
	rows, err := q.WithContext(ctx).
		Select(q.ActivityID, q.UserID.Count()).
		Where(q.ActivityID.In(activityIds...)).
		Group(q.ActivityID).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]int)
	for rows.Next() {
		var activityID int64
		var count int
		if err := rows.Scan(&activityID, &count); err == nil {
			result[activityID] = count
		}
	}
	return result, nil
}

// GetBargainHelpList 获得砍价助力列表 (按 recordId)
func (s *BargainHelpService) GetBargainHelpList(ctx context.Context, recordID int64) ([]*promotion.PromotionBargainHelp, error) {
	q := s.q.PromotionBargainHelp
	return q.WithContext(ctx).Where(q.RecordID.Eq(recordID)).Order(q.CreatedAt.Desc()).Find()
}

// GetBargainHelp 获得指定记录和用户的助力记录
func (s *BargainHelpService) GetBargainHelp(ctx context.Context, recordID int64, userID int64) (*promotion.PromotionBargainHelp, error) {
	q := s.q.PromotionBargainHelp
	return q.WithContext(ctx).Where(q.RecordID.Eq(recordID), q.UserID.Eq(userID)).First()
}

// GetBargainHelpCountByActivity 获得用户在指定活动的助力次数
func (s *BargainHelpService) GetBargainHelpCountByActivity(ctx context.Context, activityID int64, userID int64) (int64, error) {
	q := s.q.PromotionBargainHelp
	return q.WithContext(ctx).Where(q.ActivityID.Eq(activityID), q.UserID.Eq(userID)).Count()
}

// CreateBargainHelp 砍价助力
func (s *BargainHelpService) CreateBargainHelp(ctx context.Context, userID int64, r *req.AppBargainHelpCreateReq) (*promotion.PromotionBargainHelp, error) {
	var help *promotion.PromotionBargainHelp
	err := s.q.Transaction(func(tx *query.Query) error {
		// 1. 校验砍价记录 (加锁)
		record, err := tx.PromotionBargainRecord.WithContext(ctx).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(tx.PromotionBargainRecord.ID.Eq(r.RecordID)).
			First()
		if err != nil {
			return core.NewBizError(1001007000, "砍价记录不存在")
		}
		if record.UserID == userID {
			return core.NewBizError(1001007001, "不能给自己砍价")
		}
		if record.Status != 1 { // 1: In Progress
			return core.NewBizError(1001007002, "砍价记录已结束")
		}

		// 2. 校验砍价活动
		activity, err := tx.PromotionBargainActivity.WithContext(ctx).Where(tx.PromotionBargainActivity.ID.Eq(record.ActivityID)).First()
		if err != nil {
			return core.NewBizError(1001004000, "砍价活动不存在")
		}
		if activity.Status != 1 { // 1: Open
			return core.NewBizError(1001004001, "砍价活动已结束")
		}
		now := time.Now()
		if now.Before(activity.StartTime) || now.After(activity.EndTime) {
			return core.NewBizError(1001004001, "砍价活动已结束")
		}

		// 3. 校验是否已经助力过
		count, err := tx.PromotionBargainHelp.WithContext(ctx).
			Where(tx.PromotionBargainHelp.UserID.Eq(userID), tx.PromotionBargainHelp.RecordID.Eq(r.RecordID)).
			Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return core.NewBizError(1001007003, "您已经助力过了")
		}

		// 4. 计算砍价金额
		leftPrice := record.BargainPrice - activity.BargainMinPrice
		if leftPrice <= 0 {
			return core.NewBizError(1001007004, "砍价已完成")
		}

		reducePrice := 0
		if activity.RandomMinPrice == activity.RandomMaxPrice {
			reducePrice = activity.RandomMinPrice
		} else {
			reducePrice = activity.RandomMinPrice + int(now.UnixNano()%int64(activity.RandomMaxPrice-activity.RandomMinPrice+1))
		}
		if reducePrice > leftPrice {
			reducePrice = leftPrice
		}
		if reducePrice <= 0 {
			reducePrice = 1
			if leftPrice < 1 {
				reducePrice = leftPrice
			}
		}

		// 5. 保存助力
		help = &promotion.PromotionBargainHelp{
			UserID:      userID,
			ActivityID:  record.ActivityID,
			RecordID:    record.ID,
			ReducePrice: reducePrice,
		}
		if err := tx.PromotionBargainHelp.WithContext(ctx).Create(help); err != nil {
			return err
		}

		// 6. 更新砍价记录
		newPrice := record.BargainPrice - reducePrice
		newStatus := record.Status
		if newPrice <= activity.BargainMinPrice {
			newPrice = activity.BargainMinPrice
			newStatus = 2 // Success
		}

		// 更新记录状态和金额
		updateData := &promotion.PromotionBargainRecord{
			BargainPrice: newPrice,
			Status:       newStatus,
		}
		if newStatus == 2 {
			updateData.EndTime = now // 成功时记录结束时间
		}

		if _, err := tx.PromotionBargainRecord.WithContext(ctx).
			Where(tx.PromotionBargainRecord.ID.Eq(record.ID)).
			Updates(updateData); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return help, nil
}

// GetBargainHelpPage 获得砍价助力分页 (Admin)
func (s *BargainHelpService) GetBargainHelpPage(ctx context.Context, req *req.BargainHelpPageReq) (*core.PageResult[*promotion.PromotionBargainHelp], error) {
	q := s.q.PromotionBargainHelp
	do := q.WithContext(ctx)

	if req.ActivityID > 0 {
		do = do.Where(q.ActivityID.Eq(req.ActivityID))
	}
	if req.RecordID > 0 {
		do = do.Where(q.RecordID.Eq(req.RecordID))
	}

	list, total, err := do.Order(q.CreatedAt.Desc()).FindByPage(int((req.PageNo-1)*req.PageSize), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	return &core.PageResult[*promotion.PromotionBargainHelp]{List: list, Total: total}, nil
}
