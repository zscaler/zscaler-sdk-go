// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestDLPWebRules_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/webDlpRules/12345"

	server.On("GET", path, common.SuccessResponse(dlp_web_rules.WebDLPRules{
		ID:          ruleID,
		Name:        "Block SSN Uploads",
		Description: "Block uploads containing SSN",
		Action:      "BLOCK",
		State:       "ENABLED",
		Order:       1,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_web_rules.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
	assert.Equal(t, "Block SSN Uploads", result.Name)
}

func TestDLPWebRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/webDlpRules"

	server.On("GET", path, common.SuccessResponse([]dlp_web_rules.WebDLPRules{
		{ID: 1, Name: "Rule 1", Action: "BLOCK", State: "ENABLED"},
		{ID: 2, Name: "Rule 2", Action: "ALLOW", State: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_web_rules.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDLPWebRules_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/webDlpRules"

	server.On("POST", path, common.SuccessResponse(dlp_web_rules.WebDLPRules{
		ID:     99999,
		Name:   "New DLP Rule",
		Action: "BLOCK",
		State:  "ENABLED",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &dlp_web_rules.WebDLPRules{
		Name:   "New DLP Rule",
		Action: "BLOCK",
		State:  "ENABLED",
	}

	result, err := dlp_web_rules.Create(context.Background(), service, newRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestDLPWebRules_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/webDlpRules/12345"

	server.On("PUT", path, common.SuccessResponse(dlp_web_rules.WebDLPRules{
		ID:   ruleID,
		Name: "Updated DLP Rule",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRule := &dlp_web_rules.WebDLPRules{
		ID:   ruleID,
		Name: "Updated DLP Rule",
	}

	result, err := dlp_web_rules.Update(context.Background(), service, ruleID, updateRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated DLP Rule", result.Name)
}

func TestDLPWebRules_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/webDlpRules/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = dlp_web_rules.Delete(context.Background(), service, ruleID)

	require.NoError(t, err)
}

func TestDLPWebRules_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleName := "Block SSN Uploads"
	path := "/zia/api/v1/webDlpRules"

	server.On("GET", path, common.SuccessResponse([]dlp_web_rules.WebDLPRules{
		{ID: 1, Name: "Other Rule", Action: "ALLOW"},
		{ID: 2, Name: ruleName, Action: "BLOCK"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_web_rules.GetByName(context.Background(), service, ruleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleName, result.Name)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestDLPWebRules_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WebDLPRules JSON marshaling", func(t *testing.T) {
		rule := dlp_web_rules.WebDLPRules{
			ID:          12345,
			Name:        "Block Credit Cards",
			Description: "Block uploads containing credit card numbers",
			Action:      "BLOCK",
			State:       "ENABLED",
			Order:       1,
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Block Credit Cards"`)
		assert.Contains(t, string(data), `"action":"BLOCK"`)
	})

	t.Run("WebDLPRules JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "SSN Detection Rule",
			"action": "BLOCK",
			"state": "ENABLED",
			"order": 2,
			"dlpEngines": [{"id": 100, "name": "SSN Engine"}],
			"matchOnly": false
		}`

		var rule dlp_web_rules.WebDLPRules
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.Equal(t, 54321, rule.ID)
		assert.Equal(t, "BLOCK", rule.Action)
	})
}
