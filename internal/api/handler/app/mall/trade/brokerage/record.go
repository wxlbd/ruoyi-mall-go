package brokerage

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"
	tradeReq "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/app/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade/brokerage"
	brokerageSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type AppBrokerageRecordHandler struct {
	recordSvc *brokerageSvc.BrokerageRecordService
}

func NewAppBrokerageRecordHandler(recordSvc *brokerageSvc.BrokerageRecordService) *AppBrokerageRecordHandler {
	return &AppBrokerageRecordHandler{recordSvc: recordSvc}
}

// GetBrokerageRecordPage 获得分销记录分页
func (h *AppBrokerageRecordHandler) GetBrokerageRecordPage(c *gin.Context) {
	var reqVO tradeReq.AppBrokerageRecordPageReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}
	reqVO.CreateTime = resolveCreateTimeQuery(c, reqVO.CreateTime)

	userId := context.GetLoginUserID(c)
	pageReq := &trade.BrokerageRecordPageReq{
		PageParam:  reqVO.PageParam,
		UserID:     &userId,
		Status:     reqVO.Status,
		CreateTime: reqVO.CreateTime,
		BizType:    reqVO.BizType,
	}

	pageResult, err := h.recordSvc.GetBrokerageRecordPage(c, pageReq)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	writeResp := pagination.PageResult[*tradeReq.AppBrokerageRecordRespVO]{
		Total: pageResult.Total,
		List: lo.Map(pageResult.List, func(item *brokerage.BrokerageRecord, _ int) *tradeReq.AppBrokerageRecordRespVO {
			return &tradeReq.AppBrokerageRecordRespVO{
				ID:          item.ID,
				UserID:      item.UserID,
				BizType:     item.BizType,
				BizID:       item.BizID,
				Price:       item.Price,
				Title:       item.Title,
				Description: item.Description,
				Status:      item.Status,
				Total:       item.TotalPrice,
				CreateTime:  item.CreateTime,
				// StatusName: item.Status // TODO: Dict lookup
			}
		}),
	}
	response.WriteSuccess(c, writeResp)
}

// GetProductBrokeragePrice 获得商品的分销金额
func (h *AppBrokerageRecordHandler) GetProductBrokeragePrice(c *gin.Context) {
	spuIdStr := c.Query("spuId")
	spuId, err := strconv.ParseInt(spuIdStr, 10, 64)
	if err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}

	userId := context.GetLoginUserID(c)
	result, err := h.recordSvc.CalculateProductBrokeragePrice(c, userId, spuId)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, result)
}

func resolveCreateTimeQuery(c *gin.Context, createTime []string) []string {
	if len(createTime) == 2 {
		return createTime
	}
	createTime = c.QueryArray("createTime")
	if len(createTime) == 2 {
		return createTime
	}
	return c.QueryArray("createTime[]")
}
