// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/web_policy"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestWebPolicy_GetListByCompanyID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/web/policy/listByCompany"

	server.On("GET", path, common.SuccessResponse([]map[string]interface{}{
		{"id": 1, "name": "Default Policy", "active": "true"},
		{"id": 2, "name": "Custom Policy", "active": "false"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := web_policy.GetPolicyListByCompanyID(context.Background(), service, nil, nil, nil, nil, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestWebPolicy_Activate_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/web/policy/activate"

	server.On("PUT", path, common.SuccessResponse(web_policy.WebPolicyActivation{
		PolicyId:   123,
		DeviceType: 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	activation := &web_policy.WebPolicyActivation{
		PolicyId:   123,
		DeviceType: 1,
	}

	result, err := web_policy.ActivateWebPolicy(context.Background(), service, activation)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 123, result.PolicyId)
}

func TestWebPolicy_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/web/policy/edit"

	server.On("PUT", path, common.SuccessResponse(web_policy.WebPolicy{
		ID:          "123",
		Name:        "Updated Policy",
		Active:      "true",
		Description: "Updated policy description",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updatePolicy := &web_policy.WebPolicy{
		ID:     "123",
		Name:   "Updated Policy",
		Active: "true",
	}

	result, err := web_policy.UpdateWebPolicy(context.Background(), service, updatePolicy)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Policy", result.Name)
}

func TestWebPolicy_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/web/policy/123/delete"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = web_policy.DeleteWebPolicy(context.Background(), service, 123)

	require.NoError(t, err)
}

func TestWebPolicy_GetAppService_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webAppService/listByCompany"

	server.On("GET", path, common.SuccessResponse([]map[string]interface{}{
		{"id": 1, "serviceName": "Service 1", "active": true},
		{"id": 2, "serviceName": "Service 2", "active": true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := web_policy.GetWebAppServiceInfoByCompanyID(context.Background(), service, nil, nil, nil, nil, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestWebPolicy_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WebPolicy JSON marshaling", func(t *testing.T) {
		policy := web_policy.WebPolicy{
			ID:          "123",
			Name:        "Enterprise Policy",
			Description: "Main enterprise web policy",
			Active:      "true",
			DeviceType:  "WINDOWS",
		}

		data, err := json.Marshal(policy)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"123"`)
		assert.Contains(t, string(data), `"name":"Enterprise Policy"`)
		assert.Contains(t, string(data), `"active":"true"`)
	})

	t.Run("WebPolicy JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "456",
			"name": "Branch Policy",
			"description": "Policy for branch offices",
			"active": "false",
			"device_type": "MAC"
		}`

		var policy web_policy.WebPolicy
		err := json.Unmarshal([]byte(jsonData), &policy)
		require.NoError(t, err)

		assert.Equal(t, "456", policy.ID)
		assert.Equal(t, "Branch Policy", policy.Name)
		assert.Equal(t, "false", policy.Active)
	})

	t.Run("WebPolicyActivation JSON marshaling", func(t *testing.T) {
		activation := web_policy.WebPolicyActivation{
			PolicyId:   123,
			DeviceType: 1,
		}

		data, err := json.Marshal(activation)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"policyId":123`)
		assert.Contains(t, string(data), `"deviceType":1`)
	})
}
