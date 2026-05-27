// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/traffic_capture"
)

func TestTrafficCapture_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/trafficCaptureRules/12345"
	server.On("GET", path, common.SuccessResponse(traffic_capture.TrafficCaptureRules{
		ID: ruleID, Name: "tests-traffic-capture", Order: 1, Rank: 7,
		State: "ENABLED", Action: "CAPTURE",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_capture.Get(context.Background(), service, ruleID)
	require.NoError(t, err)
	assert.Equal(t, "CAPTURE", result.Action)
}

func TestTrafficCapture_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/trafficCaptureRules"
	server.On("GET", path, common.SuccessResponse([]traffic_capture.TrafficCaptureRules{
		{ID: 1, Name: "Rule 1", State: "ENABLED", Action: "CAPTURE"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_capture.GetAll(context.Background(), service, nil)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestTrafficCapture_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/trafficCaptureRules"
	server.On("POST", path, common.SuccessResponse(traffic_capture.TrafficCaptureRules{
		ID: 99999, Name: "tests-traffic-capture", Order: 1, Rank: 7, State: "ENABLED", Action: "CAPTURE",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &traffic_capture.TrafficCaptureRules{
		Name: "tests-traffic-capture", Order: 1, Rank: 7, State: "ENABLED", Action: "CAPTURE",
	}

	result, err := traffic_capture.Create(context.Background(), service, newRule)
	require.NoError(t, err)
	assert.Equal(t, 99999, result.ID)
}

func TestTrafficCapture_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/trafficCaptureRules/12345"
	server.On("PUT", path, common.SuccessResponse(traffic_capture.TrafficCaptureRules{
		ID: ruleID, Name: "tests-traffic-capture-updated",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	update := &traffic_capture.TrafficCaptureRules{ID: ruleID, Name: "tests-traffic-capture-updated"}
	result, err := traffic_capture.Update(context.Background(), service, ruleID, update)
	require.NoError(t, err)
	assert.Equal(t, "tests-traffic-capture-updated", result.Name)
}

func TestTrafficCapture_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/trafficCaptureRules/12345"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = traffic_capture.Delete(context.Background(), service, 12345)
	require.NoError(t, err)
}

func TestTrafficCapture_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	name := "tests-traffic-capture"
	path := "/zia/api/v1/trafficCaptureRules"
	server.On("GET", path, common.SuccessResponse([]traffic_capture.TrafficCaptureRules{
		{ID: 1, Name: name, Action: "CAPTURE"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_capture.GetByName(context.Background(), service, name)
	require.NoError(t, err)
	assert.Equal(t, name, result.Name)
}

func TestTrafficCapture_GetTrafficCaptureRuleCount_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/trafficCaptureRules/count"
	server.On("GET", path, common.SuccessResponse(5))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	count, err := traffic_capture.GetTrafficCaptureRuleCount(context.Background(), service, nil)
	require.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestTrafficCapture_GetTrafficCaptureRuleOrder_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/trafficCaptureRules/order"
	server.On("GET", path, common.SuccessResponse(traffic_capture.TrafficCaptureRuleOrderInfo{
		MaxOrderConfigured: 10,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_capture.GetTrafficCaptureRuleOrder(context.Background(), service)
	require.NoError(t, err)
	assert.Equal(t, 10, result.MaxOrderConfigured)
}

func TestTrafficCapture_GetTrafficCaptureRuleLabels_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/trafficCaptureRules/ruleLabels"
	server.On("GET", path, common.SuccessResponse([]traffic_capture.RuleLabelInfo{
		{ID: 1, Name: "Capture Label"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_capture.GetTrafficCaptureRuleLabels(context.Background(), service, nil)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestTrafficCapture_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/trafficCaptureRules/99999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_capture.Get(context.Background(), service, 99999)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestTrafficCapture_GetAll_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/trafficCaptureRules"
	server.On("GET", path, common.SuccessResponse([]traffic_capture.TrafficCaptureRules{
		{ID: 1, Name: "Filtered Rule", Action: "CAPTURE"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_capture.GetAll(context.Background(), service, &traffic_capture.GetAllFilterOptions{
		RuleName:   "Filtered",
		RuleAction: "CAPTURE",
		Location:   "HQ",
	})
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestTrafficCapture_GetTrafficCaptureRuleCount_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/trafficCaptureRules/count"
	server.On("GET", path, common.SuccessResponse(7))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	count, err := traffic_capture.GetTrafficCaptureRuleCount(context.Background(), service, &traffic_capture.TrafficCaptureRulesCountQuery{
		RuleName:            "Capture",
		PredefinedRuleCount: true,
		Department:          "Engineering",
	})
	require.NoError(t, err)
	assert.Equal(t, 7, count)
}

func TestTrafficCapture_GetTrafficCaptureRuleLabels_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/trafficCaptureRules/ruleLabels"
	server.On("GET", path, common.SuccessResponse([]traffic_capture.RuleLabelInfo{
		{ID: 1, Name: "Capture Label"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_capture.GetTrafficCaptureRuleLabels(context.Background(), service, &traffic_capture.GetTrafficCaptureRuleLabelsFilterOptions{
		SearchByField: "name",
		SearchByValue: "Capture",
	})
	require.NoError(t, err)
	assert.Len(t, result, 1)
}
