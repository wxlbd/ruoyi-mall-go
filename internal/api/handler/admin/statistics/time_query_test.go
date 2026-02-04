package statistics

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestResolveTimesQuery(t *testing.T) {
	t.Run("supports times[]", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?times%5B%5D=2026-01-28+00:00:00&times%5B%5D=2026-02-03+23:59:59", nil)

		got := resolveTimesQuery(c, nil)
		require.Len(t, got, 2)
		require.Equal(t, time.Date(2026, 1, 28, 0, 0, 0, 0, time.Local), got[0])
		require.Equal(t, time.Date(2026, 2, 3, 23, 59, 59, 0, time.Local), got[1])
	})

	t.Run("supports indexed times[0] times[1]", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?times%5B0%5D=2026-01-28+00:00:00&times%5B1%5D=2026-02-03+23:59:59", nil)

		got := resolveTimesQuery(c, nil)
		require.Len(t, got, 2)
		require.Equal(t, time.Date(2026, 1, 28, 0, 0, 0, 0, time.Local), got[0])
		require.Equal(t, time.Date(2026, 2, 3, 23, 59, 59, 0, time.Local), got[1])
	})

	t.Run("returns nil on invalid time format", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?times%5B0%5D=2026-01-28&times%5B1%5D=2026-02-03", nil)

		got := resolveTimesQuery(c, nil)
		require.Nil(t, got)
	})
}
