// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_classes"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_control_rules"
)

// =====================================================
// SDK Function Tests - Bandwidth Control Rules
// =====================================================

func TestBandwidthControlRules_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/bandwidthControlRules/12345"

	server.On("GET", path, common.SuccessResponse(bandwidth_control_rules.BandwidthControlRules{
		ID:           ruleID,
		Name:         "Limit Streaming",
		MaxBandwidth: 5000000,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := bandwidth_control_rules.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
}

func TestBandwidthControlRules_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleName := "Limit Streaming"
	path := "/zia/api/v1/bandwidthControlRules"

	server.On("GET", path, common.SuccessResponse([]bandwidth_control_rules.BandwidthControlRules{
		{ID: 1, Name: "Other Rule", MaxBandwidth: 1000000},
		{ID: 2, Name: ruleName, MaxBandwidth: 5000000},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := bandwidth_control_rules.GetByName(context.Background(), service, ruleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
}

func TestBandwidthControlRules_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/bandwidthControlRules"

	server.On("POST", path, common.SuccessResponse(bandwidth_control_rules.BandwidthControlRules{
		ID:           100,
		Name:         "New Rule",
		MaxBandwidth: 10000000,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &bandwidth_control_rules.BandwidthControlRules{
		Name:         "New Rule",
		MaxBandwidth: 10000000,
		State:        "ENABLED",
		Order:        1,
	}

	result, err := bandwidth_control_rules.Create(context.Background(), service, newRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 100, result.ID)
}

func TestBandwidthControlRules_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/bandwidthControlRules/12345"

	server.On("PUT", path, common.SuccessResponse(bandwidth_control_rules.BandwidthControlRules{
		ID:           ruleID,
		Name:         "Updated Rule",
		MaxBandwidth: 20000000,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRule := &bandwidth_control_rules.BandwidthControlRules{
		ID:           ruleID,
		Name:         "Updated Rule",
		MaxBandwidth: 20000000,
	}

	result, err := bandwidth_control_rules.Update(context.Background(), service, ruleID, updateRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Rule", result.Name)
}

func TestBandwidthControlRules_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/bandwidthControlRules/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = bandwidth_control_rules.Delete(context.Background(), service, ruleID)

	require.NoError(t, err)
}

func TestBandwidthControlRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/bandwidthControlRules"

	server.On("GET", path, common.SuccessResponse([]bandwidth_control_rules.BandwidthControlRules{
		{ID: 1, Name: "Rule 1", MaxBandwidth: 5000000},
		{ID: 2, Name: "Rule 2", MaxBandwidth: 10000000},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := bandwidth_control_rules.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestBandwidthControlRules_GetAllLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/bandwidthControlRules/lite"

	server.On("GET", path, common.SuccessResponse([]bandwidth_control_rules.BandwidthControlRules{
		{ID: 1, Name: "Rule 1"},
		{ID: 2, Name: "Rule 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := bandwidth_control_rules.GetAllLite(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// SDK Function Tests - Bandwidth Classes
// =====================================================

func TestBandwidthClasses_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	classID := 12345
	path := "/zia/api/v1/bandwidthClasses/12345"

	server.On("GET", path, common.SuccessResponse(bandwidth_classes.BandwidthClasses{
		ID:   classID,
		Name: "Streaming Media Class",
		Type: "CUSTOM",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := bandwidth_classes.Get(context.Background(), service, classID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, classID, result.ID)
}

func TestBandwidthClasses_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	className := "Streaming Media Class"
	path := "/zia/api/v1/bandwidthClasses"

	server.On("GET", path, common.SuccessResponse([]bandwidth_classes.BandwidthClasses{
		{ID: 1, Name: "Other Class", Type: "PREDEFINED"},
		{ID: 2, Name: className, Type: "CUSTOM"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := bandwidth_classes.GetByName(context.Background(), service, className)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
}

func TestBandwidthClasses_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/bandwidthClasses"

	server.On("GET", path, common.SuccessResponse([]bandwidth_classes.BandwidthClasses{
		{ID: 1, Name: "Class 1", Type: "PREDEFINED"},
		{ID: 2, Name: "Class 2", Type: "CUSTOM"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := bandwidth_classes.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestBandwidthClasses_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/bandwidthClasses"

	server.On("POST", path, common.SuccessResponse(bandwidth_classes.BandwidthClasses{
		ID:   100,
		Name: "New Class",
		Type: "CUSTOM",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newClass := &bandwidth_classes.BandwidthClasses{
		Name: "New Class",
		Type: "CUSTOM",
	}

	result, _, err := bandwidth_classes.Create(context.Background(), service, newClass)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 100, result.ID)
}

func TestBandwidthClasses_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	classID := 12345
	path := "/zia/api/v1/bandwidthClasses/12345"

	server.On("PUT", path, common.SuccessResponse(bandwidth_classes.BandwidthClasses{
		ID:   classID,
		Name: "Updated Class",
		Type: "CUSTOM",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateClass := &bandwidth_classes.BandwidthClasses{
		ID:   classID,
		Name: "Updated Class",
		Type: "CUSTOM",
	}

	result, _, err := bandwidth_classes.Update(context.Background(), service, classID, updateClass)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Class", result.Name)
}

func TestBandwidthClasses_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	classID := 12345
	path := "/zia/api/v1/bandwidthClasses/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = bandwidth_classes.Delete(context.Background(), service, classID)

	require.NoError(t, err)
}

func TestBandwidthClasses_GetAllLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/bandwidthClasses/lite"

	server.On("GET", path, common.SuccessResponse([]bandwidth_classes.BandwidthClasses{
		{ID: 1, Name: "Class 1"},
		{ID: 2, Name: "Class 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := bandwidth_classes.GetAllLite(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// Structure Tests
// =====================================================

func TestBandwidthClasses_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BandwidthClasses JSON marshaling", func(t *testing.T) {
		class := bandwidth_classes.BandwidthClasses{
			ID:                       12345,
			Name:                     "Streaming Media Class",
			Type:                     "CUSTOM",
			IsNameL10nTag:            false,
			WebApplications:          []string{"YOUTUBE", "NETFLIX"},
			UrlCategories:            []string{"STREAMING_MEDIA"},
			NetworkApplications:      []string{"HTTP", "HTTPS"},
			NetworkServices:          []string{"TCP_443", "TCP_80"},
			ApplicationServiceGroups: []string{"VIDEO_STREAMING"},
		}

		data, err := json.Marshal(class)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"CUSTOM"`)
	})

	t.Run("BandwidthClasses JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Cloud Storage Class",
			"type": "PREDEFINED",
			"isNameL10nTag": true,
			"webApplications": ["GOOGLE_DRIVE", "DROPBOX", "BOX"],
			"urlCategories": ["CLOUD_STORAGE"],
			"urls": ["*.googleapis.com", "*.dropbox.com"],
			"applications": ["GOOGLE_DRIVE", "DROPBOX"]
		}`

		var class bandwidth_classes.BandwidthClasses
		err := json.Unmarshal([]byte(jsonData), &class)
		require.NoError(t, err)

		assert.Equal(t, 54321, class.ID)
		assert.True(t, class.IsNameL10nTag)
		assert.Len(t, class.WebApplications, 3)
	})
}

func TestBandwidthControlRules_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BandwidthControlRules JSON marshaling", func(t *testing.T) {
		rule := bandwidth_control_rules.BandwidthControlRules{
			ID:                12345,
			Name:              "Limit Streaming",
			Order:             1,
			State:             "ENABLED",
			Description:       "Limit streaming media bandwidth",
			MaxBandwidth:      5000000, // 5 Mbps
			MinBandwidth:      1000000, // 1 Mbps
			Rank:              7,
			AccessControl:     "READ_WRITE",
			DefaultRule:       false,
			Protocols:         []string{"HTTPS_RULE", "HTTP_RULE"},
			DeviceTrustLevels: []string{"HIGH_TRUST"},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"maxBandwidth":5000000`)
		assert.Contains(t, string(data), `"minBandwidth":1000000`)
	})

	t.Run("BandwidthControlRules JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Priority Traffic",
			"order": 1,
			"state": "ENABLED",
			"description": "Prioritize business traffic",
			"maxBandwidth": 0,
			"minBandwidth": 10000000,
			"rank": 5,
			"accessControl": "READ_ONLY",
			"defaultRule": false,
			"protocols": ["HTTPS_RULE"],
			"lastModifiedTime": 1699000000,
			"lastModifiedBy": {
				"id": 100,
				"name": "admin@company.com"
			},
			"bandwidthClasses": [
				{"id": 200, "name": "Business Apps"}
			],
			"locations": [
				{"id": 300, "name": "HQ"}
			],
			"locationGroups": [
				{"id": 400, "name": "North America"}
			]
		}`

		var rule bandwidth_control_rules.BandwidthControlRules
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.Equal(t, 54321, rule.ID)
		assert.Equal(t, 10000000, rule.MinBandwidth)
		assert.NotNil(t, rule.LastModifiedBy)
		assert.Len(t, rule.BandwidthClasses, 1)
	})

	t.Run("BandwidthControlRules default rule", func(t *testing.T) {
		jsonData := `{
			"id": 1,
			"name": "Default Rule",
			"order": 100,
			"state": "ENABLED",
			"defaultRule": true,
			"maxBandwidth": 100000000
		}`

		var rule bandwidth_control_rules.BandwidthControlRules
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.True(t, rule.DefaultRule)
	})
}

func TestBandwidthControl_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse bandwidth classes list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Video Streaming", "type": "PREDEFINED"},
			{"id": 2, "name": "Cloud Storage", "type": "PREDEFINED"},
			{"id": 3, "name": "Custom Class", "type": "CUSTOM"}
		]`

		var classes []bandwidth_classes.BandwidthClasses
		err := json.Unmarshal([]byte(jsonResponse), &classes)
		require.NoError(t, err)

		assert.Len(t, classes, 3)
		assert.Equal(t, "CUSTOM", classes[2].Type)
	})

	t.Run("Parse bandwidth control rules list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Rule 1", "state": "ENABLED", "maxBandwidth": 10000000},
			{"id": 2, "name": "Rule 2", "state": "ENABLED", "maxBandwidth": 5000000},
			{"id": 3, "name": "Default", "state": "ENABLED", "defaultRule": true}
		]`

		var rules []bandwidth_control_rules.BandwidthControlRules
		err := json.Unmarshal([]byte(jsonResponse), &rules)
		require.NoError(t, err)

		assert.Len(t, rules, 3)
		assert.True(t, rules[2].DefaultRule)
	})
}

