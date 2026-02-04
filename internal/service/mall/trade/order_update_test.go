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
	return nil, nil
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
