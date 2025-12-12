// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/timewindow"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/time_intervals"
)

func TestTimeIntervals_Structure(t *testing.T) {
	t.Parallel()

	t.Run("TimeInterval JSON marshaling", func(t *testing.T) {
		interval := time_intervals.TimeInterval{
			ID:         12345,
			Name:       "Business Hours",
			StartTime:  540,  // 9:00 AM in minutes
			EndTime:    1080, // 6:00 PM in minutes
			DaysOfWeek: []string{"MON", "TUE", "WED", "THU", "FRI"},
		}

		data, err := json.Marshal(interval)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"startTime":540`)
		assert.Contains(t, string(data), `"endTime":1080`)
	})

	t.Run("TimeInterval JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "After Hours",
			"startTime": 1080,
			"endTime": 540,
			"daysOfWeek": ["MON", "TUE", "WED", "THU", "FRI"]
		}`

		var interval time_intervals.TimeInterval
		err := json.Unmarshal([]byte(jsonData), &interval)
		require.NoError(t, err)

		assert.Equal(t, 54321, interval.ID)
		assert.Equal(t, 1080, interval.StartTime)
		assert.Len(t, interval.DaysOfWeek, 5)
	})

	t.Run("TimeInterval weekend config", func(t *testing.T) {
		interval := time_intervals.TimeInterval{
			ID:         12346,
			Name:       "Weekend",
			StartTime:  0,
			EndTime:    1440, // Full day
			DaysOfWeek: []string{"SAT", "SUN"},
		}

		data, err := json.Marshal(interval)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"daysOfWeek":["SAT","SUN"]`)
	})

	t.Run("TimeInterval everyday config", func(t *testing.T) {
		interval := time_intervals.TimeInterval{
			ID:         12347,
			Name:       "All Day Every Day",
			StartTime:  0,
			EndTime:    1440,
			DaysOfWeek: []string{"EVERYDAY"},
		}

		data, err := json.Marshal(interval)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"EVERYDAY"`)
	})
}

func TestTimeWindow_Structure(t *testing.T) {
	t.Parallel()

	t.Run("TimeWindow JSON marshaling", func(t *testing.T) {
		tw := timewindow.TimeWindow{
			ID:        12345,
			Name:      "Work Hours",
			StartTime: 540,
			EndTime:   1080,
			DayOfWeek: []string{"MON", "TUE", "WED", "THU", "FRI"},
		}

		data, err := json.Marshal(tw)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Work Hours"`)
	})

	t.Run("TimeWindow JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "Night Shift",
			"startTime": 1320,
			"endTime": 360,
			"dayOfWeek": ["MON", "TUE", "WED", "THU", "FRI"]
		}`

		var tw timewindow.TimeWindow
		err := json.Unmarshal([]byte(jsonData), &tw)
		require.NoError(t, err)

		assert.Equal(t, 67890, tw.ID)
		assert.Equal(t, int32(1320), tw.StartTime)
	})
}

func TestTimeIntervals_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse time intervals list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Morning", "startTime": 480, "endTime": 720},
			{"id": 2, "name": "Afternoon", "startTime": 720, "endTime": 1020},
			{"id": 3, "name": "Evening", "startTime": 1020, "endTime": 1320}
		]`

		var intervals []time_intervals.TimeInterval
		err := json.Unmarshal([]byte(jsonResponse), &intervals)
		require.NoError(t, err)

		assert.Len(t, intervals, 3)
		assert.Equal(t, "Afternoon", intervals[1].Name)
	})
}

