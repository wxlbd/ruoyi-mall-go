package trade

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
)

// 交易模块错误码定义 (对齐 Java ErrorCodeConstants)
// 错误码格式: 1004XXXYYY
// - 1004: 业务模块前缀
// - XXX: 子模块编号 (003=价格计算, 004=订单操作, 005=售后, 006=购物车)
// - YYY: 具体错误编号

const (
	// ========== 价格计算相关错误码 (1004003xxx) ==========

	// 价格计算基础错误 (1004003000-1004003099)
	ErrorCodePriceCalculateError     = 1004003000 // 价格计算失败
	ErrorCodePriceCalculateItemEmpty = 1004003001 // 价格计算商品为空
	ErrorCodePriceCalculateItemError = 1004003002 // 价格计算商品错误
	ErrorCodePriceCalculateUserError = 1004003003 // 价格计算用户错误

	// 商品相关错误 (1004003100-1004003199)
	ErrorCodeProductNotExists      = 1004003100 // 商品不存在
	ErrorCodeProductNotEnable      = 1004003101 // 商品未启用
	ErrorCodeProductStockNotEnough = 1004003102 // 商品库存不足
	ErrorCodeProductSkuNotExists   = 1004003103 // 商品SKU不存在
	ErrorCodeProductSkuNotEnable   = 1004003104 // 商品SKU未启用

	// 优惠券相关错误 (1004003200-1004003299)
	ErrorCodeCouponNotExists = 1004003200 // 优惠券不存在
	ErrorCodeCouponNotMatch  = 1004003201 // 优惠券不匹配
	ErrorCodeCouponUsed      = 1004003202 // 优惠券已使用
	ErrorCodeCouponExpired   = 1004003203 // 优惠券已过期
	ErrorCodeCouponNotStart  = 1004003204 // 优惠券未开始
	ErrorCodeCouponNotEnough = 1004003205 // 优惠券数量不足

	// 积分相关错误 (1004003300-1004003399)
	ErrorCodePointNotEnough      = 1004003300 // 积分不足
	ErrorCodePointCalculateError = 1004003301 // 积分计算错误

	// 活动相关错误 (1004003400-1004003499)
	ErrorCodeActivityNotExists      = 1004003400 // 活动不存在
	ErrorCodeActivityNotStart       = 1004003401 // 活动未开始
	ErrorCodeActivityExpired        = 1004003402 // 活动已结束
	ErrorCodeActivityNotMatch       = 1004003403 // 活动不匹配
	ErrorCodeActivityStockNotEnough = 1004003404 // 活动库存不足

	// 运费相关错误 (1004003500-1004003599)
	ErrorCodeDeliveryNotSupport        = 1004003500 // 不支持配送
	ErrorCodeDeliveryTemplateNotExists = 1004003501 // 运费模板不存在
	ErrorCodeDeliveryCalculateError    = 1004003502 // 运费计算错误

	// ========== 订单操作相关错误码 (1004004xxx) ==========

	// 订单基础错误 (1004004000-1004004099)
	ErrorCodeOrderNotExists    = 1004004000 // 订单不存在
	ErrorCodeOrderStatusError  = 1004004001 // 订单状态错误
	ErrorCodeOrderUserNotMatch = 1004004002 // 订单用户不匹配
	ErrorCodeOrderCreateError  = 1004004003 // 订单创建失败
	ErrorCodeOrderUpdateError  = 1004004004 // 订单更新失败
	ErrorCodeOrderDeleteError  = 1004004005 // 订单删除失败

	// 订单支付相关错误 (1004004100-1004004199)
	ErrorCodeOrderNotPaid        = 1004004100 // 订单未支付
	ErrorCodeOrderAlreadyPaid    = 1004004101 // 订单已支付
	ErrorCodeOrderPayError       = 1004004102 // 订单支付失败
	ErrorCodeOrderPayTimeout     = 1004004103 // 订单支付超时
	ErrorCodeOrderPayAmountError = 1004004104 // 订单支付金额错误

	// 订单发货相关错误 (1004004200-1004004299)
	ErrorCodeOrderNotDelivered     = 1004004200 // 订单未发货
	ErrorCodeOrderAlreadyDelivered = 1004004201 // 订单已发货
	ErrorCodeOrderDeliveryError    = 1004004202 // 订单发货失败
	ErrorCodeOrderLogisticsError   = 1004004203 // 物流信息错误

	// 订单收货相关错误 (1004004300-1004004399)
	ErrorCodeOrderNotReceived     = 1004004300 // 订单未收货
	ErrorCodeOrderAlreadyReceived = 1004004301 // 订单已收货
	ErrorCodeOrderReceiveError    = 1004004302 // 订单收货失败

	// 订单取消相关错误 (1004004400-1004004499)
	ErrorCodeOrderNotCanceled     = 1004004400 // 订单未取消
	ErrorCodeOrderAlreadyCanceled = 1004004401 // 订单已取消
	ErrorCodeOrderCancelError     = 1004004402 // 订单取消失败
	ErrorCodeOrderCancelNotAllow  = 1004004403 // 订单不允许取消

	// 订单退款相关错误 (1004004500-1004004599)
	ErrorCodeOrderRefundError       = 1004004500 // 订单退款失败
	ErrorCodeOrderRefundAmountError = 1004004501 // 退款金额错误
	ErrorCodeOrderRefundNotAllow    = 1004004502 // 订单不允许退款

	// 订单核销相关错误 (1004004600-1004004699)
	ErrorCodeOrderPickUpError     = 1004004600 // 订单核销失败
	ErrorCodeOrderNotPickUp       = 1004004601 // 非自提订单
	ErrorCodeOrderPickUpCodeError = 1004004602 // 核销码错误
	ErrorCodeOrderAlreadyPickUp   = 1004004603 // 订单已核销

	// 订单评价相关错误 (1004004700-1004004799)
	ErrorCodeOrderCommentError    = 1004004700 // 订单评价失败
	ErrorCodeOrderAlreadyComment  = 1004004701 // 订单已评价
	ErrorCodeOrderCommentNotAllow = 1004004702 // 订单不允许评价

	// ========== 售后相关错误码 (1004005xxx) ==========

	// 售后基础错误 (1004005000-1004005099)
	ErrorCodeAfterSaleNotExists   = 1004005000 // 售后单不存在
	ErrorCodeAfterSaleStatusError = 1004005001 // 售后单状态错误
	ErrorCodeAfterSaleCreateError = 1004005002 // 售后单创建失败
	ErrorCodeAfterSaleUpdateError = 1004005003 // 售后单更新失败

	// ========== 购物车相关错误码 (1004006xxx) ==========

	// 购物车基础错误 (1004006000-1004006099)
	ErrorCodeCartNotExists   = 1004006000 // 购物车项不存在
	ErrorCodeCartAddError    = 1004006001 // 添加购物车失败
	ErrorCodeCartUpdateError = 1004006002 // 更新购物车失败
	ErrorCodeCartDeleteError = 1004006003 // 删除购物车失败
	ErrorCodeCartCountError  = 1004006004 // 购物车数量错误

	// ========== 分销相关错误码 (1011007xxx) ==========
	// 对齐 Java: ErrorCodeConstants.BROKERAGE_* (1_011_007_xxx)

	// 分销用户相关错误 (1011007000-1011007099)
	ErrorCodeBrokerageUserNotExists            = 1011007000 // 分销用户不存在
	ErrorCodeBrokerageUserFrozenPriceNotEnough = 1011007001 // 用户冻结佣金数量不足
	ErrorCodeBrokerageBindSelf                 = 1011007002 // 不能绑定自己
	ErrorCodeBrokerageBindUserNotEnabled       = 1011007003 // 绑定用户没有推广资格
	ErrorCodeBrokerageBindConditionAdmin       = 1011007004 // 仅可在后台绑定推广员
	ErrorCodeBrokerageBindModeRegister         = 1011007005 // 只有在注册时可以绑定
	ErrorCodeBrokerageBindOverride             = 1011007006 // 已绑定了推广人
	ErrorCodeBrokerageBindLoop                 = 1011007007 // 下级不能绑定自己的上级
	ErrorCodeBrokerageUserLevelNotSupport      = 1011007008 // 目前只支持 level 小于等于 2
	ErrorCodeBrokerageUserExists               = 1011007009 // 分销用户已存在

	// 分销提现相关错误 (1011008000-1011008099)
	ErrorCodeBrokerageWithdrawNotExists = 1011008000 // 佣金提现记录不存在
)

// 错误消息映射表 (对齐 Java 版本的错误消息)
var errorMessages = map[int]string{
	// 价格计算相关错误消息
	ErrorCodePriceCalculateError:     "价格计算失败",
	ErrorCodePriceCalculateItemEmpty: "价格计算商品为空",
	ErrorCodePriceCalculateItemError: "价格计算商品错误",
	ErrorCodePriceCalculateUserError: "价格计算用户错误",

	ErrorCodeProductNotExists:      "商品不存在",
	ErrorCodeProductNotEnable:      "商品未启用",
	ErrorCodeProductStockNotEnough: "商品库存不足",
	ErrorCodeProductSkuNotExists:   "商品SKU不存在",
	ErrorCodeProductSkuNotEnable:   "商品SKU未启用",

	ErrorCodeCouponNotExists: "优惠券不存在",
	ErrorCodeCouponNotMatch:  "优惠券不匹配",
	ErrorCodeCouponUsed:      "优惠券已使用",
	ErrorCodeCouponExpired:   "优惠券已过期",
	ErrorCodeCouponNotStart:  "优惠券未开始",
	ErrorCodeCouponNotEnough: "优惠券数量不足",

	ErrorCodePointNotEnough:      "积分不足",
	ErrorCodePointCalculateError: "积分计算错误",

	ErrorCodeActivityNotExists:      "活动不存在",
	ErrorCodeActivityNotStart:       "活动未开始",
	ErrorCodeActivityExpired:        "活动已结束",
	ErrorCodeActivityNotMatch:       "活动不匹配",
	ErrorCodeActivityStockNotEnough: "活动库存不足",

	ErrorCodeDeliveryNotSupport:        "不支持配送",
	ErrorCodeDeliveryTemplateNotExists: "运费模板不存在",
	ErrorCodeDeliveryCalculateError:    "运费计算错误",

	// 订单操作相关错误消息
	ErrorCodeOrderNotExists:    "订单不存在",
	ErrorCodeOrderStatusError:  "订单状态错误",
	ErrorCodeOrderUserNotMatch: "订单用户不匹配",
	ErrorCodeOrderCreateError:  "订单创建失败",
	ErrorCodeOrderUpdateError:  "订单更新失败",
	ErrorCodeOrderDeleteError:  "订单删除失败",

	ErrorCodeOrderNotPaid:        "订单未支付",
	ErrorCodeOrderAlreadyPaid:    "订单已支付",
	ErrorCodeOrderPayError:       "订单支付失败",
	ErrorCodeOrderPayTimeout:     "订单支付超时",
	ErrorCodeOrderPayAmountError: "订单支付金额错误",

	ErrorCodeOrderNotDelivered:     "订单未发货",
	ErrorCodeOrderAlreadyDelivered: "订单已发货",
	ErrorCodeOrderDeliveryError:    "订单发货失败",
	ErrorCodeOrderLogisticsError:   "物流信息错误",

	ErrorCodeOrderNotReceived:     "订单未收货",
	ErrorCodeOrderAlreadyReceived: "订单已收货",
	ErrorCodeOrderReceiveError:    "订单收货失败",

	ErrorCodeOrderNotCanceled:     "订单未取消",
	ErrorCodeOrderAlreadyCanceled: "订单已取消",
	ErrorCodeOrderCancelError:     "订单取消失败",
	ErrorCodeOrderCancelNotAllow:  "订单不允许取消",

	ErrorCodeOrderRefundError:       "订单退款失败",
	ErrorCodeOrderRefundAmountError: "退款金额错误",
	ErrorCodeOrderRefundNotAllow:    "订单不允许退款",

	ErrorCodeOrderPickUpError:     "订单核销失败",
	ErrorCodeOrderNotPickUp:       "非自提订单",
	ErrorCodeOrderPickUpCodeError: "核销码错误",
	ErrorCodeOrderAlreadyPickUp:   "订单已核销",

	ErrorCodeOrderCommentError:    "订单评价失败",
	ErrorCodeOrderAlreadyComment:  "订单已评价",
	ErrorCodeOrderCommentNotAllow: "订单不允许评价",

	// 售后相关错误消息
	ErrorCodeAfterSaleNotExists:   "售后单不存在",
	ErrorCodeAfterSaleStatusError: "售后单状态错误",
	ErrorCodeAfterSaleCreateError: "售后单创建失败",
	ErrorCodeAfterSaleUpdateError: "售后单更新失败",

	// 购物车相关错误消息
	ErrorCodeCartNotExists:   "购物车项不存在",
	ErrorCodeCartAddError:    "添加购物车失败",
	ErrorCodeCartUpdateError: "更新购物车失败",
	ErrorCodeCartDeleteError: "删除购物车失败",
	ErrorCodeCartCountError:  "购物车数量错误",

	// 分销相关错误消息
	ErrorCodeBrokerageUserNotExists:            "分销用户不存在",
	ErrorCodeBrokerageUserFrozenPriceNotEnough: "用户冻结佣金数量不足",
	ErrorCodeBrokerageBindSelf:                 "不能绑定自己",
	ErrorCodeBrokerageBindUserNotEnabled:       "绑定用户没有推广资格",
	ErrorCodeBrokerageBindConditionAdmin:       "仅可在后台绑定推广员",
	ErrorCodeBrokerageBindModeRegister:         "只有在注册时可以绑定",
	ErrorCodeBrokerageBindOverride:             "已绑定了推广人",
	ErrorCodeBrokerageBindLoop:                 "下级不能绑定自己的上级",
	ErrorCodeBrokerageUserLevelNotSupport:      "目前只支持 level 小于等于 2",
	ErrorCodeBrokerageUserExists:               "分销用户已存在",
	ErrorCodeBrokerageWithdrawNotExists:        "佣金提现记录不存在",
}

// NewTradeError 创建交易模块业务错误
func NewTradeError(code int) error {
	if msg, exists := errorMessages[code]; exists {
		return errors.NewBizError(code, msg)
	}
	return errors.NewBizError(code, "未知错误")
}

// NewTradeErrorWithMsg 创建交易模块业务错误（自定义消息）
func NewTradeErrorWithMsg(code int, msg string) error {
	return errors.NewBizError(code, msg)
}

// 便捷的错误创建函数

// 价格计算相关错误
func ErrPriceCalculateError() error     { return NewTradeError(ErrorCodePriceCalculateError) }
func ErrPriceCalculateItemEmpty() error { return NewTradeError(ErrorCodePriceCalculateItemEmpty) }
func ErrProductNotExists() error        { return NewTradeError(ErrorCodeProductNotExists) }
func ErrProductStockNotEnough() error   { return NewTradeError(ErrorCodeProductStockNotEnough) }
func ErrCouponNotExists() error         { return NewTradeError(ErrorCodeCouponNotExists) }
func ErrCouponNotMatch() error          { return NewTradeError(ErrorCodeCouponNotMatch) }
func ErrPointNotEnough() error          { return NewTradeError(ErrorCodePointNotEnough) }
func ErrActivityNotExists() error       { return NewTradeError(ErrorCodeActivityNotExists) }
func ErrDeliveryNotSupport() error      { return NewTradeError(ErrorCodeDeliveryNotSupport) }

// 订单操作相关错误
func ErrOrderNotExists() error        { return NewTradeError(ErrorCodeOrderNotExists) }
func ErrOrderStatusError() error      { return NewTradeError(ErrorCodeOrderStatusError) }
func ErrOrderUserNotMatch() error     { return NewTradeError(ErrorCodeOrderUserNotMatch) }
func ErrOrderAlreadyPaid() error      { return NewTradeError(ErrorCodeOrderAlreadyPaid) }
func ErrOrderNotDelivered() error     { return NewTradeError(ErrorCodeOrderNotDelivered) }
func ErrOrderAlreadyDelivered() error { return NewTradeError(ErrorCodeOrderAlreadyDelivered) }
func ErrOrderNotReceived() error      { return NewTradeError(ErrorCodeOrderNotReceived) }
func ErrOrderAlreadyReceived() error  { return NewTradeError(ErrorCodeOrderAlreadyReceived) }
func ErrOrderAlreadyCanceled() error  { return NewTradeError(ErrorCodeOrderAlreadyCanceled) }
func ErrOrderCancelNotAllow() error   { return NewTradeError(ErrorCodeOrderCancelNotAllow) }
func ErrOrderRefundNotAllow() error   { return NewTradeError(ErrorCodeOrderRefundNotAllow) }
func ErrOrderNotPickUp() error        { return NewTradeError(ErrorCodeOrderNotPickUp) }
func ErrOrderPickUpCodeError() error  { return NewTradeError(ErrorCodeOrderPickUpCodeError) }
func ErrOrderAlreadyPickUp() error    { return NewTradeError(ErrorCodeOrderAlreadyPickUp) }

// 分销相关错误
func ErrBrokerageUserNotExists() error { return NewTradeError(ErrorCodeBrokerageUserNotExists) }
func ErrBrokerageUserFrozenPriceNotEnough() error {
	return NewTradeError(ErrorCodeBrokerageUserFrozenPriceNotEnough)
}
func ErrBrokerageBindSelf() error { return NewTradeError(ErrorCodeBrokerageBindSelf) }
func ErrBrokerageBindUserNotEnabled() error {
	return NewTradeError(ErrorCodeBrokerageBindUserNotEnabled)
}
func ErrBrokerageBindConditionAdmin() error {
	return NewTradeError(ErrorCodeBrokerageBindConditionAdmin)
}
func ErrBrokerageBindModeRegister() error { return NewTradeError(ErrorCodeBrokerageBindModeRegister) }
func ErrBrokerageBindOverride() error     { return NewTradeError(ErrorCodeBrokerageBindOverride) }
func ErrBrokerageBindLoop() error         { return NewTradeError(ErrorCodeBrokerageBindLoop) }
func ErrBrokerageUserLevelNotSupport() error {
	return NewTradeError(ErrorCodeBrokerageUserLevelNotSupport)
}
func ErrBrokerageUserExists() error        { return NewTradeError(ErrorCodeBrokerageUserExists) }
func ErrBrokerageWithdrawNotExists() error { return NewTradeError(ErrorCodeBrokerageWithdrawNotExists) }
