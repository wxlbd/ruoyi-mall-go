package trade

import (
	"time"

	"gorm.io/gorm"
)

// TradeOrder 交易订单
type TradeOrder struct {
	ID                       int64          `gorm:"primaryKey;autoIncrement;comment:订单编号"`
	No                       string         `gorm:"column:no;type:varchar(32);not null;comment:订单流水号"`
	Type                     int            `gorm:"column:type;type:int;not null;comment:订单类型"`
	Terminal                 int            `gorm:"column:terminal;type:int;not null;comment:订单来源"`
	UserID                   int64          `gorm:"column:user_id;type:bigint;not null;comment:用户编号"`
	UserIP                   string         `gorm:"column:user_ip;type:varchar(50);not null;comment:用户 IP"`
	UserRemark               string         `gorm:"column:user_remark;type:varchar(255);comment:用户备注"`
	Status                   int            `gorm:"column:status;type:int;not null;comment:订单状态"`
	ProductCount             int            `gorm:"column:product_count;type:int;not null;comment:购买的商品数量"`
	FinishTime               *time.Time     `gorm:"column:finish_time;comment:订单完成时间"`
	CancelTime               *time.Time     `gorm:"column:cancel_time;comment:订单取消时间"`
	CancelType               int            `gorm:"column:cancel_type;type:int;comment:取消类型"`
	Remark                   string         `gorm:"column:remark;type:varchar(255);comment:商家备注"`
	CommentStatus            bool           `gorm:"column:comment_status;type:tinyint(1);not null;default:0;comment:是否评价"`
	BrokerageUserID          *int64         `gorm:"column:brokerage_user_id;type:bigint;comment:推广人编号"`
	PayOrderID               *int64         `gorm:"column:pay_order_id;type:bigint;comment:支付订单编号"`
	PayStatus                bool           `gorm:"column:pay_status;type:tinyint(1);not null;default:0;comment:是否已支付"`
	PayTime                  *time.Time     `gorm:"column:pay_time;comment:付款时间"`
	PayChannelCode           string         `gorm:"column:pay_channel_code;type:varchar(16);comment:支付渠道"`
	TotalPrice               int            `gorm:"column:total_price;type:int;not null;comment:商品原价"`
	DiscountPrice            int            `gorm:"column:discount_price;type:int;not null;comment:优惠金额"`
	DeliveryPrice            int            `gorm:"column:delivery_price;type:int;not null;comment:运费金额"`
	AdjustPrice              int            `gorm:"column:adjust_price;type:int;not null;comment:订单调价"`
	PayPrice                 int            `gorm:"column:pay_price;type:int;not null;comment:应付金额"`
	DeliveryType             int            `gorm:"column:delivery_type;type:int;not null;comment:配送方式"`
	LogisticsID              int64          `gorm:"column:logistics_id;type:bigint;comment:发货物流公司编号"`
	LogisticsNo              string         `gorm:"column:logistics_no;type:varchar(64);comment:发货物流单号"`
	DeliveryTime             *time.Time     `gorm:"column:delivery_time;comment:发货时间"`
	ReceiveTime              *time.Time     `gorm:"column:receive_time;comment:收货时间"`
	ReceiverName             string         `gorm:"column:receiver_name;type:varchar(30);not null;comment:收件人名称"`
	ReceiverMobile           string         `gorm:"column:receiver_mobile;type:varchar(20);not null;comment:收件人手机"`
	ReceiverAreaID           int            `gorm:"column:receiver_area_id;type:int;not null;comment:收件人地区编号"`
	ReceiverDetailAddress    string         `gorm:"column:receiver_detail_address;type:varchar(255);not null;comment:收件人详细地址"`
	PickUpStoreID            int64          `gorm:"column:pick_up_store_id;type:bigint;comment:自提门店编号"`
	PickUpVerifyCode         string         `gorm:"column:pick_up_verify_code;type:varchar(64);comment:自提核销码"`
	RefundStatus             int            `gorm:"column:refund_status;type:int;not null;comment:售后状态"`
	RefundPrice              int            `gorm:"column:refund_price;type:int;not null;comment:退款金额"`
	CouponID                 int64          `gorm:"column:coupon_id;type:bigint;not null;default:0;comment:优惠劵编号"`
	CouponPrice              int            `gorm:"column:coupon_price;type:int;not null;comment:优惠劵减免金额"`
	UsePoint                 int            `gorm:"column:use_point;type:int;not null;default:0;comment:使用的积分"`
	PointPrice               int            `gorm:"column:point_price;type:int;not null;default:0;comment:积分抵扣的金额"`
	GivePoint                int            `gorm:"column:give_point;type:int;not null;default:0;comment:赠送的积分"`
	RefundPoint              int            `gorm:"column:refund_point;type:int;not null;default:0;comment:退还的使用的积分"`
	VipPrice                 int            `gorm:"column:vip_price;type:int;not null;default:0;comment:VIP 减免金额"`
	GiveCouponTemplateCounts map[int64]int  `gorm:"column:give_coupon_template_counts;type:json;serializer:json;comment:赠送的优惠劵"`
	GiveCouponIDs            []int64        `gorm:"column:give_coupon_ids;type:json;serializer:json;comment:赠送的优惠劵编号"`
	SeckillActivityID        int64          `gorm:"column:seckill_activity_id;type:bigint;comment:秒杀活动编号"`
	BargainActivityID        int64          `gorm:"column:bargain_activity_id;type:bigint;comment:砍价活动编号"`
	BargainRecordID          int64          `gorm:"column:bargain_record_id;type:bigint;comment:砍价记录编号"`
	CombinationActivityID    int64          `gorm:"column:combination_activity_id;type:bigint;comment:拼团活动编号"`
	CombinationHeadID        int64          `gorm:"column:combination_head_id;type:bigint;comment:拼团团长编号"`
	CombinationRecordID      int64          `gorm:"column:combination_record_id;type:bigint;comment:拼团记录编号"`
	PointActivityID          int64          `gorm:"column:point_activity_id;type:bigint;comment:积分商城活动的编号"`
	Creator                  string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater                  string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt                time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt                time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt                gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted                  bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (TradeOrder) TableName() string {
	return "trade_order"
}

// TradeOrderItem 交易订单项
type TradeOrderItem struct {
	ID              int64                    `gorm:"primaryKey;autoIncrement;comment:订单项编号"`
	UserID          int64                    `gorm:"column:user_id;type:bigint;not null;comment:用户编号"`
	OrderID         int64                    `gorm:"column:order_id;type:bigint;not null;comment:订单编号"`
	CartID          int64                    `gorm:"column:cart_id;type:bigint;not null;comment:购物车项编号"`
	SpuID           int64                    `gorm:"column:spu_id;type:bigint;not null;comment:商品 SPU 编号"`
	SpuName         string                   `gorm:"column:spu_name;type:varchar(255);not null;comment:商品 SPU 名称"`
	SkuID           int64                    `gorm:"column:sku_id;type:bigint;not null;comment:商品 SKU 编号"`
	Properties      []TradeOrderItemProperty `gorm:"column:properties;type:json;serializer:json;comment:属性数组"`
	PicURL          string                   `gorm:"column:pic_url;type:varchar(255);comment:商品图片"`
	Count           int                      `gorm:"column:count;type:int;not null;comment:购买数量"`
	CommentStatus   bool                     `gorm:"column:comment_status;type:tinyint(1);not null;default:0;comment:是否评价"`
	Price           int                      `gorm:"column:price;type:int;not null;comment:商品原价"`
	DiscountPrice   int                      `gorm:"column:discount_price;type:int;not null;comment:优惠金额"`
	DeliveryPrice   int                      `gorm:"column:delivery_price;type:int;not null;comment:运费金额"`
	AdjustPrice     int                      `gorm:"column:adjust_price;type:int;not null;comment:订单调价"`
	PayPrice        int                      `gorm:"column:pay_price;type:int;not null;comment:应付金额"`
	CouponPrice     int                      `gorm:"column:coupon_price;type:int;not null;comment:优惠劵减免金额"`
	PointPrice      int                      `gorm:"column:point_price;type:int;not null;comment:积分抵扣的金额"`
	UsePoint        int                      `gorm:"column:use_point;type:int;not null;default:0;comment:使用的积分"`
	GivePoint       int                      `gorm:"column:give_point;type:int;not null;default:0;comment:赠送的积分"`
	VipPrice        int                      `gorm:"column:vip_price;type:int;not null;default:0;comment:VIP 减免金额"`
	AfterSaleID     int64                    `gorm:"column:after_sale_id;type:bigint;comment:售后单编号"`
	AfterSaleStatus int                      `gorm:"column:after_sale_status;type:int;not null;comment:售后状态"`
	Creator         string                   `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater         string                   `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt       time.Time                `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt       time.Time                `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt       gorm.DeletedAt           `gorm:"column:deleted;index;comment:删除时间"`
	Deleted         bool                     `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (TradeOrderItem) TableName() string {
	return "trade_order_item"
}

type TradeOrderItemProperty struct {
	PropertyID   int64  `json:"propertyId"`
	PropertyName string `json:"propertyName"`
	ValueID      int64  `json:"valueId"`
	ValueName    string `json:"valueName"`
}

// TradeOrderLog 订单日志
type TradeOrderLog struct {
	ID           int64          `gorm:"primaryKey;autoIncrement;comment:日志编号"`
	UserID       int64          `gorm:"column:user_id;type:bigint;not null;comment:用户编号"`
	UserType     int            `gorm:"column:user_type;type:tinyint;not null;comment:用户类型"`
	OrderID      int64          `gorm:"column:order_id;type:bigint;not null;comment:订单号"`
	BeforeStatus int            `gorm:"column:before_status;type:int;comment:操作前状态"`
	AfterStatus  int            `gorm:"column:after_status;type:int;comment:操作后状态"`
	OperateType  int            `gorm:"column:operate_type;type:int;not null;comment:操作类型"`
	Content      string         `gorm:"column:content;type:varchar(2000);not null;comment:订单日志信息"`
	Creator      string         `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater      string         `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreatedAt    time.Time      `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdatedAt    time.Time      `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted;index;comment:删除时间"`
	Deleted      bool           `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
}

func (TradeOrderLog) TableName() string {
	return "trade_order_log"
}
