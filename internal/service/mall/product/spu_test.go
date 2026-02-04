package product

import (
	"testing"

	"github.com/stretchr/testify/require"
	modelProduct "github.com/wxlbd/ruoyi-mall-go/internal/model/product"
)

func TestMergeVirtualSalesCount(t *testing.T) {
	spu := &modelProduct.ProductSpu{
		SalesCount:        120,
		VirtualSalesCount: 30,
	}

	mergeVirtualSalesCount(spu)

	require.Equal(t, 150, spu.SalesCount)
}
