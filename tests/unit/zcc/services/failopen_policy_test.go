// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/failopen_policy"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestFailOpenPolicy_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webFailOpenPolicy"

	server.On("GET", path, common.SuccessResponse([]failopen_policy.WebFailOpenPolicy{
		{ID: "policy-001", Active: "true", CompanyID: "company-123", EnableFailOpen: 1},
		{ID: "policy-002", Active: "false", CompanyID: "company-123", EnableFailOpen: 0},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := failopen_policy.GetFailOpenPolicy(context.Background(), service, 100)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "policy-001", result[0].ID)
	assert.Equal(t, "true", result[0].Active)
}

func TestFailOpenPolicy_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webFailOpenPolicy"

	server.On("PUT", path, common.SuccessResponse(failopen_policy.WebFailOpenPolicy{
		ID:             "policy-001",
		Active:         "false",
		CompanyID:      "company-123",
		EnableFailOpen: 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updatePolicy := &failopen_policy.WebFailOpenPolicy{
		ID:             "policy-001",
		Active:         "false",
		EnableFailOpen: 0,
	}

	result, err := failopen_policy.UpdateFailOpenPolicy(context.Background(), service, updatePolicy)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "false", result.Active)
	assert.Equal(t, 0, result.EnableFailOpen)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestFailOpenPolicy_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WebFailOpenPolicy JSON marshaling", func(t *testing.T) {
		policy := failopen_policy.WebFailOpenPolicy{
			ID:                                "policy-123",
			Active:                            "true",
			CompanyID:                         "company-456",
			EnableFailOpen:                    1,
			EnableCaptivePortalDetection:      1,
			CaptivePortalWebSecDisableMinutes: 5,
			EnableStrictEnforcementPrompt:     1,
			StrictEnforcementPromptDelayMins:  10,
			TunnelFailureRetryCount:           3,
			EnableWebSecOnTunnelFailure:       "true",
			EnableWebSecOnProxyUnreachable:    "false",
			CreatedBy:                         "admin@example.com",
			EditedBy:                          "admin@example.com",
		}

		data, err := json.Marshal(policy)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"policy-123"`)
		assert.Contains(t, string(data), `"active":"true"`)
		assert.Contains(t, string(data), `"enableFailOpen":1`)
	})

	t.Run("WebFailOpenPolicy JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "policy-789",
			"active": "false",
			"companyId": "company-123",
			"enableFailOpen": 0,
			"enableCaptivePortalDetection": 1,
			"captivePortalWebSecDisableMinutes": 10,
			"enableStrictEnforcementPrompt": 0,
			"tunnelFailureRetryCount": 5,
			"enableWebSecOnTunnelFailure": "false",
			"enableWebSecOnProxyUnreachable": "true"
		}`

		var policy failopen_policy.WebFailOpenPolicy
		err := json.Unmarshal([]byte(jsonData), &policy)
		require.NoError(t, err)

		assert.Equal(t, "policy-789", policy.ID)
		assert.Equal(t, "false", policy.Active)
		assert.Equal(t, 0, policy.EnableFailOpen)
		assert.Equal(t, 5, policy.TunnelFailureRetryCount)
	})
}

func TestFailOpenPolicy_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse fail open policies list", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": "policy-001",
				"active": "true",
				"companyId": "123",
				"enableFailOpen": 1
			},
			{
				"id": "policy-002",
				"active": "false",
				"companyId": "123",
				"enableFailOpen": 0
			}
		]`

		var policies []failopen_policy.WebFailOpenPolicy
		err := json.Unmarshal([]byte(jsonResponse), &policies)
		require.NoError(t, err)

		assert.Len(t, policies, 2)
		assert.Equal(t, "policy-001", policies[0].ID)
		assert.Equal(t, "true", policies[0].Active)
		assert.Equal(t, 1, policies[0].EnableFailOpen)
		assert.Equal(t, "false", policies[1].Active)
	})
}
