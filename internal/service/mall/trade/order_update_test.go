package trade

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	trade2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/trade"
	memberContract "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
)

type fakeMemberAddressService struct {
	expectedUserID    int64
	expectedAddressID int64
	calledUserID      int64
	calledAddressID   int64
	defaultAddress    *memberContract.AppAddressResp
	defaultCalled     bool
}

func (f *fakeMemberAddressService) GetAddress(_ context.Context, userID int64, addressID int64) (*memberContract.AppAddressResp, error) {
	f.calledUserID = userID
	f.calledAddressID = addressID
	if userID != f.expectedUserID || addressID != f.expectedAddressID {
		return nil, nil
	}
	return &memberContract.AppAddressResp{
		ID:            addressID,
		Name:          "张三",
		Mobile:        "13800138000",
		AreaID:        110101,
		DetailAddress: "中关村大街 1 号",
	}, nil
}

func (f *fakeMemberAddressService) GetDefaultAddress(context.Context, int64) (*memberContract.AppAddressResp, error) {
	f.defaultCalled = true
	return f.defaultAddress, nil
}

func TestBuildTradeOrder_ExpressDeliveryLoadsReceiverFromAddress(t *testing.T) {
	const userID int64 = 20001
	addressID := int64(30001)

	addrSvc := &fakeMemberAddressService{
		expectedUserID:    userID,
		expectedAddressID: addressID,
	}
	svc := &TradeOrderUpdateService{
		addressSvc: addrSvc,
	}

	req := &trade2.AppTradeOrderCreateReq{
		AppTradeOrderSettlementReq: trade2.AppTradeOrderSettlementReq{
			DeliveryType: consts.DeliveryTypeExpress,
			AddressID:    &addressID,
		},
	}
	priceResp := &TradePriceCalculateRespBO{
		Type: 1,
		Price: TradePriceCalculatePriceBO{
			TotalPrice: 1000,
			PayPrice:   900,
		},
		Items: []TradePriceCalculateItemRespBO{
			{Count: 1},
		},
	}

	order := svc.buildTradeOrder(context.Background(), userID, "127.0.0.1", 10, req, priceResp)

	require.Equal(t, userID, addrSvc.calledUserID)
	require.Equal(t, addressID, addrSvc.calledAddressID)
	require.Equal(t, "张三", order.ReceiverName)
	require.Equal(t, "13800138000", order.ReceiverMobile)
	require.Equal(t, 110101, order.ReceiverAreaID)
	require.Equal(t, "中关村大街 1 号", order.ReceiverDetailAddress)
}

func TestBuildTradeOrder_ExpressDeliveryFallbackToDefaultAddress(t *testing.T) {
	const userID int64 = 20001

	addrSvc := &fakeMemberAddressService{
		defaultAddress: &memberContract.AppAddressResp{
			ID:            30002,
			Name:          "李四",
			Mobile:        "13900139000",
			AreaID:        310101,
			DetailAddress: "南京东路 1 号",
		},
	}
	svc := &TradeOrderUpdateService{
		addressSvc: addrSvc,
	}

	req := &trade2.AppTradeOrderCreateReq{
		AppTradeOrderSettlementReq: trade2.AppTradeOrderSettlementReq{
			DeliveryType: consts.DeliveryTypeExpress,
		},
	}
	priceResp := &TradePriceCalculateRespBO{
		Type: 1,
		Price: TradePriceCalculatePriceBO{
			TotalPrice: 1000,
			PayPrice:   900,
		},
		Items: []TradePriceCalculateItemRespBO{
			{Count: 1},
		},
	}

	order := svc.buildTradeOrder(context.Background(), userID, "127.0.0.1", 10, req, priceResp)

	require.True(t, addrSvc.defaultCalled)
	require.Equal(t, "李四", order.ReceiverName)
	require.Equal(t, "13900139000", order.ReceiverMobile)
	require.Equal(t, 310101, order.ReceiverAreaID)
	require.Equal(t, "南京东路 1 号", order.ReceiverDetailAddress)
}

func TestBuildTradeOrder_AndItemsPersistVipPrice(t *testing.T) {
	svc := &TradeOrderUpdateService{}
	req := &trade2.AppTradeOrderCreateReq{
		AppTradeOrderSettlementReq: trade2.AppTradeOrderSettlementReq{
			DeliveryType: consts.DeliveryTypePickUp,
		},
	}
	priceResp := &TradePriceCalculateRespBO{
		Type: 1,
		Price: TradePriceCalculatePriceBO{
			TotalPrice: 1500,
			VipPrice:   300,
			PayPrice:   1200,
		},
		Items: []TradePriceCalculateItemRespBO{
			{
				SpuID:    10,
				SkuID:    11,
				Count:    1,
				Price:    1500,
				VipPrice: 300,
				PayPrice: 1200,
			},
		},
	}

	order := svc.buildTradeOrder(context.Background(), 1, "127.0.0.1", 10, req, priceResp)
	items := svc.buildTradeOrderItems(order, priceResp)

	require.Equal(t, 300, order.VipPrice)
	require.Len(t, items, 1)
	require.Equal(t, 300, items[0].VipPrice)
}
