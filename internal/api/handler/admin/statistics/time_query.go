package statistics

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const statisticsTimeLayout = "2006-01-02 15:04:05"

func resolveTimesQuery(c *gin.Context, times []time.Time) []time.Time {
	if len(times) == 2 {
		return times
	}

	if parsed := parseTimes(c.QueryArray("times")); len(parsed) == 2 {
		return parsed
	}
	if parsed := parseTimes(c.QueryArray("times[]")); len(parsed) == 2 {
		return parsed
	}

	start := c.Query("times[0]")
	end := c.Query("times[1]")
	if start == "" || end == "" {
		return nil
	}
	return parseTimes([]string{start, end})
}

func parseTimes(values []string) []time.Time {
	if len(values) != 2 {
		return nil
	}

	beginTime, err := time.ParseInLocation(statisticsTimeLayout, strings.TrimSpace(values[0]), time.Local)
	if err != nil {
		return nil
	}
	endTime, err := time.ParseInLocation(statisticsTimeLayout, strings.TrimSpace(values[1]), time.Local)
	if err != nil {
		return nil
	}
	return []time.Time{beginTime, endTime}
}
