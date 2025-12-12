// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/rule_labels"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestRuleLabels_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	labelID := 12345
	path := "/zia/api/v1/ruleLabels/12345"

	server.On("GET", path, common.SuccessResponse(rule_labels.RuleLabels{
		ID:          labelID,
		Name:        "Test Label",
		Description: "Test description",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := rule_labels.Get(context.Background(), service, labelID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, labelID, result.ID)
	assert.Equal(t, "Test Label", result.Name)
}

func TestRuleLabels_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/ruleLabels"

	server.On("GET", path, common.SuccessResponse([]rule_labels.RuleLabels{
		{ID: 1, Name: "Label 1"},
		{ID: 2, Name: "Label 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := rule_labels.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestRuleLabels_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/ruleLabels"

	server.On("POST", path, common.SuccessResponse(rule_labels.RuleLabels{
		ID:   99999,
		Name: "New Label",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newLabel := rule_labels.RuleLabels{
		Name: "New Label",
	}

	result, _, err := rule_labels.Create(context.Background(), service, &newLabel)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestRuleLabels_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	labelID := 12345
	path := "/zia/api/v1/ruleLabels/12345"

	server.On("PUT", path, common.SuccessResponse(rule_labels.RuleLabels{
		ID:   labelID,
		Name: "Updated Label",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateLabel := rule_labels.RuleLabels{
		ID:   labelID,
		Name: "Updated Label",
	}

	result, _, err := rule_labels.Update(context.Background(), service, labelID, &updateLabel)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Label", result.Name)
}

func TestRuleLabels_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	labelID := 12345
	path := "/zia/api/v1/ruleLabels/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = rule_labels.Delete(context.Background(), service, labelID)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestRuleLabels_Structure(t *testing.T) {
	t.Parallel()

	t.Run("RuleLabels JSON marshaling", func(t *testing.T) {
		label := rule_labels.RuleLabels{
			ID:          12345,
			Name:        "Security Label",
			Description: "Label for security rules",
		}

		data, err := json.Marshal(label)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Security Label"`)
	})

	t.Run("RuleLabels JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Compliance Label",
			"description": "Label for compliance rules",
			"lastModifiedTime": 1699000000,
			"lastModifiedBy": {"id": 100, "name": "admin"}
		}`

		var label rule_labels.RuleLabels
		err := json.Unmarshal([]byte(jsonData), &label)
		require.NoError(t, err)

		assert.Equal(t, 54321, label.ID)
		assert.Equal(t, "Compliance Label", label.Name)
	})
}
