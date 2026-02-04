package brokerage

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestResolveCreateTimeQuery(t *testing.T) {
	t.Run("prefer bound createTime", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest("GET", "/?createTime=2026-01-01+00:00:00&createTime=2026-01-02+23:59:59", nil)
		c.Request = req

		got := resolveCreateTimeQuery(c, []string{"2026-01-03 00:00:00", "2026-01-04 23:59:59"})
		require.Equal(t, []string{"2026-01-03 00:00:00", "2026-01-04 23:59:59"}, got)
	})

	t.Run("fallback createTime", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest("GET", "/?createTime=2026-01-01+00:00:00&createTime=2026-01-02+23:59:59", nil)
		c.Request = req

		got := resolveCreateTimeQuery(c, nil)
		require.Equal(t, []string{"2026-01-01 00:00:00", "2026-01-02 23:59:59"}, got)
	})

	t.Run("fallback createTime bracket", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest("GET", "/?createTime%5B%5D=2026-01-01+00:00:00&createTime%5B%5D=2026-01-02+23:59:59", nil)
		c.Request = req

		got := resolveCreateTimeQuery(c, nil)
		require.Equal(t, []string{"2026-01-01 00:00:00", "2026-01-02 23:59:59"}, got)
	})
}
