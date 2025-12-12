package trade

import "time"

// AfterSale 售后
type AfterSale struct {
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Creator          string    `gorm:"column:creator" json:"creator"`
	Updater          string    `gorm:"column:updater" json:"updater"`
	CreatedAt        time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt        time.Time `gorm:"column:update_time" json:"updateTime"`
	DeletedAt        time.Time `gorm:"column:deleted" json:"deletedTime"`
	Deleted          bool      `gorm:"column:deleted" json:"deleted"`
	No               string    `gorm:"column:no" json:"no"`
	Status           int       `gorm:"column:status" json:"status"`
	Way              int       `gorm:"column:way" json:"way"`
	Type             int       `gorm:"column:type" json:"type"`
	UserID           int64     `gorm:"column:user_id" json:"userId"`
	ApplyReason      string    `gorm:"column:apply_reason" json:"applyReason"`
	ApplyDescription string    `gorm:"column:apply_description" json:"applyDescription"`
	ApplyPicURLs     string    `gorm:"column:apply_pic_urls" json:"applyPicUrls"` // serialized
	OrderID          int64     `gorm:"column:order_id" json:"orderId"`
	OrderNo          string    `gorm:"column:order_no" json:"orderNo"`
	OrderItemID      int64     `gorm:"column:order_item_id" json:"orderItemId"`
	SpuID            int64     `gorm:"column:spu_id" json:"spuId"`
	SpuName          string    `gorm:"column:spu_name" json:"spuName"`
	SkuID            int64     `gorm:"column:sku_id" json:"skuId"`
	Properties       string    `gorm:"column:properties" json:"properties"` // serialized
	PicURL           string    `gorm:"column:pic_url" json:"picUrl"`
	Count            int       `gorm:"column:count" json:"count"`
	AuditTime        time.Time `gorm:"column:audit_time" json:"auditTime"`
	AuditUserID      int64     `gorm:"column:audit_user_id" json:"auditUserId"`
	AuditReason      string    `gorm:"column:audit_reason" json:"auditReason"`
	RefundPrice      int       `gorm:"column:refund_price" json:"refundPrice"`
	PayRefundID      int64     `gorm:"column:pay_refund_id" json:"payRefundId"`
	RefundTime       time.Time `gorm:"column:refund_time" json:"refundTime"`
	LogisticsID      int64     `gorm:"column:logistics_id" json:"logisticsId"`
	LogisticsNo      string    `gorm:"column:logistics_no" json:"logisticsNo"`
	DeliveryTime     time.Time `gorm:"column:delivery_time" json:"deliveryTime"`
	ReceiveTime      time.Time `gorm:"column:receive_time" json:"receiveTime"`
	ReceiveReason    string    `gorm:"column:receive_reason" json:"receiveReason"`
}

func (AfterSale) TableName() string {
	return "trade_after_sale"
}
