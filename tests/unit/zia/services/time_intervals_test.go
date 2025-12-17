// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/timewindow"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/time_intervals"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestTimeWindow_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/timeWindows"

	server.On("GET", path, common.SuccessResponse([]timewindow.TimeWindow{
		{ID: 1, Name: "Business Hours", StartTime: 540, EndTime: 1080},
		{ID: 2, Name: "After Hours", StartTime: 1080, EndTime: 540},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := timewindow.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestTimeWindow_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	windowName := "Business Hours"
	path := "/zia/api/v1/timeWindows"

	server.On("GET", path, common.SuccessResponse([]timewindow.TimeWindow{
		{ID: 1, Name: windowName, StartTime: 540, EndTime: 1080},
		{ID: 2, Name: "After Hours", StartTime: 1080, EndTime: 540},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := timewindow.GetTimeWindowByName(context.Background(), service, windowName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, windowName, result.Name)
}

// =====================================================
// Time Intervals SDK Function Tests
// =====================================================

func TestTimeIntervals_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	intervalID := 12345
	path := "/zia/api/v1/timeIntervals/12345"

	server.On("GET", path, common.SuccessResponse(time_intervals.TimeInterval{
		ID:         intervalID,
		Name:       "Business Hours",
		StartTime:  540,
		EndTime:    1080,
		DaysOfWeek: []string{"MON", "TUE", "WED", "THU", "FRI"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := time_intervals.Get(context.Background(), service, intervalID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, intervalID, result.ID)
}

func TestTimeIntervals_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/timeIntervals"

	server.On("GET", path, common.SuccessResponse([]time_intervals.TimeInterval{
		{ID: 1, Name: "Business Hours", StartTime: 540, EndTime: 1080},
		{ID: 2, Name: "After Hours", StartTime: 1080, EndTime: 540},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := time_intervals.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestTimeIntervals_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	intervalName := "Business Hours"
	path := "/zia/api/v1/timeIntervals"

	server.On("GET", path, common.SuccessResponse([]time_intervals.TimeInterval{
		{ID: 1, Name: intervalName, StartTime: 540, EndTime: 1080},
		{ID: 2, Name: "After Hours", StartTime: 1080, EndTime: 540},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := time_intervals.GetTimeIntervalByName(context.Background(), service, intervalName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, intervalName, result.Name)
}

func TestTimeIntervals_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/timeIntervals"

	server.On("POST", path, common.SuccessResponse(time_intervals.TimeInterval{
		ID:        99999,
		Name:      "New Interval",
		StartTime: 600,
		EndTime:   1200,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newInterval := &time_intervals.TimeInterval{
		Name:      "New Interval",
		StartTime: 600,
		EndTime:   1200,
	}

	result, _, err := time_intervals.Create(context.Background(), service, newInterval)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestTimeIntervals_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	intervalID := 12345
	path := "/zia/api/v1/timeIntervals/12345"

	server.On("PUT", path, common.SuccessResponse(time_intervals.TimeInterval{
		ID:   intervalID,
		Name: "Updated Interval",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateInterval := &time_intervals.TimeInterval{
		ID:   intervalID,
		Name: "Updated Interval",
	}

	result, _, err := time_intervals.Update(context.Background(), service, intervalID, updateInterval)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Interval", result.Name)
}

func TestTimeIntervals_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	intervalID := 12345
	path := "/zia/api/v1/timeIntervals/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = time_intervals.Delete(context.Background(), service, intervalID)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

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

