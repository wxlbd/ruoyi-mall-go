package brokerage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseTimeRange(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		begin, end, ok := parseTimeRange([]string{"2026-01-01 00:00:00", "2026-01-02 23:59:59"})
		require.True(t, ok)
		require.Equal(t, time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local), begin)
		require.Equal(t, time.Date(2026, 1, 2, 23, 59, 59, 0, time.Local), end)
	})

	t.Run("invalid format", func(t *testing.T) {
		_, _, ok := parseTimeRange([]string{"2026-01-01", "2026-01-02"})
		require.False(t, ok)
	})

	t.Run("invalid length", func(t *testing.T) {
		_, _, ok := parseTimeRange([]string{"2026-01-01 00:00:00"})
		require.False(t, ok)
	})
}
