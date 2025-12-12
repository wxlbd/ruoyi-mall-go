package promotion

import (
	"time"

	"gorm.io/gorm"
)

// PromotionRewardActivity 满减送活动
// Table: promotion_reward_activity
type PromotionRewardActivity struct {
	ID                 int64          `gorm:"primaryKey;autoIncrement;comment:活动编号" json:"id"`
	Name               string         `gorm:"column:name;type:varchar(64);not null;comment:活动名称" json:"name"`
	Status             int            `gorm:"column:status;type:int;not null;default:0;comment:状态" json:"status"` // 0: 开启, 1: 关闭
	StartTime          time.Time      `gorm:"column:start_time;not null;comment:开始时间" json:"startTime"`
	EndTime            time.Time      `gorm:"column:end_time;not null;comment:结束时间" json:"endTime"`
	ProductScope       int            `gorm:"column:product_scope;type:int;not null;default:1;comment:商品范围" json:"productScope"`   // 1: 全部商品, 2: 指定商品, 3: 指定分类
	ProductScopeValues string         `gorm:"column:product_scope_values;type:json;comment:商品范围值" json:"productScopeValues"`       // Array of IDs
	ConditionType      int            `gorm:"column:condition_type;type:int;not null;default:1;comment:条件类型" json:"conditionType"` // 10: 满N元, 20: 满N件
	Rules              string         `gorm:"column:rules;type:json;not null;comment:优惠规则" json:"rules"`                           // List<Rule>
	Sort               int            `gorm:"column:sort;type:int;not null;default:0;comment:排序" json:"sort"`
	Remark             string         `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`
	Creator            string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater            string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt          time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt          time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted            bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (PromotionRewardActivity) TableName() string {
	return "promotion_reward_activity"
}
