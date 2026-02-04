package statistics

import (
	trade2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/excel"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

// TradeStatisticsHandler 交易统计处理器
type TradeStatisticsHandler struct {
	tradeStatisticsService      trade.TradeStatisticsService
	tradeOrderStatisticsService trade.TradeOrderStatisticsService
	afterSaleStatisticsService  trade.AfterSaleStatisticsService
	brokerageStatisticsService  trade.BrokerageStatisticsService
}

// NewTradeStatisticsHandler 创建交易统计处理器
func NewTradeStatisticsHandler(
	tradeStatisticsService trade.TradeStatisticsService,
	tradeOrderStatisticsService trade.TradeOrderStatisticsService,
	afterSaleStatisticsService trade.AfterSaleStatisticsService,
	brokerageStatisticsService trade.BrokerageStatisticsService,
) *TradeStatisticsHandler {
	return &TradeStatisticsHandler{
		tradeStatisticsService:      tradeStatisticsService,
		tradeOrderStatisticsService: tradeOrderStatisticsService,
		afterSaleStatisticsService:  afterSaleStatisticsService,
		brokerageStatisticsService:  brokerageStatisticsService,
	}
}

// GetTradeSummaryComparison 获得交易统计对比
// GET /statistics/trade/summary
func (h *TradeStatisticsHandler) GetTradeSummaryComparison(c *gin.Context) {
	// 1.1 昨天的数据
	yesterdayData, err := h.tradeStatisticsService.GetTradeSummaryByDays(c, -1)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 1.2 前天的数据（用于对照昨天的数据）
	beforeYesterdayData, err := h.tradeStatisticsService.GetTradeSummaryByDays(c, -2)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 2.1 本月数据
	monthData, err := h.tradeStatisticsService.GetTradeSummaryByMonths(c, 0)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 2.2 上月数据（用于对照本月的数据）
	lastMonthData, err := h.tradeStatisticsService.GetTradeSummaryByMonths(c, -1)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 拼接数据
	result := &trade2.DataComparisonRespVO[trade2.TradeSummaryRespVO]{
		Summary: &trade2.TradeSummaryRespVO{
			Yesterday: yesterdayData,
			Month:     monthData,
		},
		Comparison: &trade2.TradeSummaryRespVO{
			Yesterday: beforeYesterdayData,
			Month:     lastMonthData,
		},
	}

	response.WriteSuccess(c, result)
}

// GetTradeStatisticsAnalyse 获得交易状况统计
// GET /statistics/trade/analyse
func (h *TradeStatisticsHandler) GetTradeStatisticsAnalyse(c *gin.Context) {
	var reqVO trade2.TradeStatisticsReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil && len(resolveTimesQuery(c, reqVO.Times)) != 2 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	reqVO.Times = resolveTimesQuery(c, reqVO.Times)

	if len(reqVO.Times) != 2 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	result, err := h.tradeStatisticsService.GetTradeStatisticsAnalyse(c, reqVO.Times[0], reqVO.Times[1])
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, result)
}

// GetTradeStatisticsList 获得交易状况明细
// GET /statistics/trade/list
func (h *TradeStatisticsHandler) GetTradeStatisticsList(c *gin.Context) {
	var reqVO trade2.TradeStatisticsReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil && len(resolveTimesQuery(c, reqVO.Times)) != 2 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	reqVO.Times = resolveTimesQuery(c, reqVO.Times)

	if len(reqVO.Times) != 2 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	result, err := h.tradeStatisticsService.GetTradeStatisticsList(c, reqVO.Times[0], reqVO.Times[1])
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, result)
}

// GetOrderCount 获得交易订单数量
// GET /statistics/trade/order-count
func (h *TradeStatisticsHandler) GetOrderCount(c *gin.Context) {
	// 订单统计
	// 待发货：Status=Undelivered (10), DeliveryType=Express (1)
	undeliveredCount, err := h.tradeOrderStatisticsService.GetCountByStatusAndDeliveryType(c, consts.TradeOrderStatusUndelivered, consts.DeliveryTypeExpress)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 待自提：Status=Delivered (20), DeliveryType=PickUp (2)
	pickUpCount, err := h.tradeOrderStatisticsService.GetCountByStatusAndDeliveryType(c, consts.TradeOrderStatusDelivered, consts.DeliveryTypePickUp)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 售后统计: Status=Applied (10)
	// TODO: Replace magic number 1 with AfterSaleStatus enum constant if available
	afterSaleApplyCount, err := h.afterSaleStatisticsService.GetCountByStatus(c, 1)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 提现统计: Status=Auditing (1)
	// TODO: Replace magic number 1 with BrokerageWithdrawStatus enum constant if available
	auditingWithdrawCount, err := h.brokerageStatisticsService.GetWithdrawCountByStatus(c, 1)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	result := &trade2.TradeOrderCountRespVO{
		UndeliveredCount:      undeliveredCount,
		PickUpCount:           pickUpCount,
		AfterSaleApplyCount:   afterSaleApplyCount,
		AuditingWithdrawCount: auditingWithdrawCount,
	}

	response.WriteSuccess(c, result)
}

// GetOrderComparison 获得交易订单数量对比
// GET /statistics/trade/order-comparison
func (h *TradeStatisticsHandler) GetOrderComparison(c *gin.Context) {
	result, err := h.tradeOrderStatisticsService.GetOrderComparison(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, result)
}

// GetOrderCountTrendComparison 获得订单量趋势统计
// GET /statistics/trade/order-count-trend
func (h *TradeStatisticsHandler) GetOrderCountTrendComparison(c *gin.Context) {
	result, err := h.tradeOrderStatisticsService.GetOrderCountTrendComparison(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	response.WriteSuccess(c, result)
}

// ExportTradeStatisticsExcel 导出交易统计 Excel
// GET /statistics/trade/export-excel
func (h *TradeStatisticsHandler) ExportTradeStatisticsExcel(c *gin.Context) {
	var reqVO trade2.TradeStatisticsReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil && len(resolveTimesQuery(c, reqVO.Times)) != 2 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	reqVO.Times = resolveTimesQuery(c, reqVO.Times)

	if len(reqVO.Times) != 2 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	// 查询数据
	data, err := h.tradeStatisticsService.GetTradeStatisticsList(c, reqVO.Times[0], reqVO.Times[1])
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 导出 Excel
	if err = excel.WriteExcel(c, "交易状况.xlsx", "数据", data); err != nil {
		response.WriteBizError(c, err)
		return
	}
}
