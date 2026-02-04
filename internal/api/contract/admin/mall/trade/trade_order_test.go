package trade

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTradeOrderBase_JSONContainsVipAndPointPrice(t *testing.T) {
	base := TradeOrderBase{
		VipPrice:   123,
		PointPrice: 456,
	}

	b, err := json.Marshal(base)
	require.NoError(t, err)

	var m map[string]any
	require.NoError(t, json.Unmarshal(b, &m))
	require.Equal(t, float64(123), m["vipPrice"])
	require.Equal(t, float64(456), m["pointPrice"])
}
