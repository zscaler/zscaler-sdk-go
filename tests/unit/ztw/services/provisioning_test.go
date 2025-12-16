// Package services provides unit tests for ZTW services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/provisioning/api_keys"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestProvisioningAPIKeys_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyID := 12345
	path := "/ztw/api/v1/apiKeys/12345"

	server.On("GET", path, common.SuccessResponse(api_keys.ProvisioningAPIKeys{
		ID:       keyID,
		KeyValue: "ABC123XYZ789",
		Enabled:  true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := api_keys.Get(context.Background(), service, keyID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, keyID, result.ID)
	assert.Equal(t, "ABC123XYZ789", result.KeyValue)
}

func TestProvisioningAPIKeys_GetPartnerAPIKey_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyValue := "PARTNER123KEY"
	path := "/ztw/api/v1/apiKeys"

	server.On("GET", path, common.SuccessResponse([]api_keys.ProvisioningAPIKeys{
		{ID: 1, KeyValue: "OTHER123KEY", Enabled: true},
		{ID: 2, KeyValue: keyValue, Enabled: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := api_keys.GetPartnerAPIKey(context.Background(), service, keyValue, true)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, keyValue, result.KeyValue)
}

func TestProvisioningAPIKeys_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/apiKeys"

	server.On("GET", path, common.SuccessResponse([]api_keys.ProvisioningAPIKeys{
		{ID: 1, KeyValue: "KEY1ABC", Enabled: true},
		{ID: 2, KeyValue: "KEY2XYZ", Enabled: false},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := api_keys.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestProvisioningAPIKeys_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/apiKeys"

	server.On("POST", path, common.SuccessResponse(api_keys.ProvisioningAPIKeys{
		ID:       99999,
		KeyValue: "NEW123KEY456",
		Enabled:  true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newKey := &api_keys.ProvisioningAPIKeys{
		Enabled: true,
	}

	result, err := api_keys.Create(context.Background(), service, newKey, false, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestProvisioningAPIKeys_CreateRegenerate_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyID := 12345
	path := "/ztw/api/v1/apiKeys/12345/ztw/api/v1/regenerate"

	server.On("POST", path, common.SuccessResponse(api_keys.ProvisioningAPIKeys{
		ID:       keyID,
		KeyValue: "REGEN123NEW",
		Enabled:  true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newKey := &api_keys.ProvisioningAPIKeys{}

	result, err := api_keys.Create(context.Background(), service, newKey, false, &keyID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "REGEN123NEW", result.KeyValue)
}

// =====================================================
// Structure Tests
// =====================================================

func TestProvisioningAPIKeys_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ProvisioningAPIKeys JSON marshaling", func(t *testing.T) {
		apiKey := api_keys.ProvisioningAPIKeys{
			ID:               12345,
			KeyValue:         "ABC123XYZ789",
			Permissions:      []string{"READ", "WRITE", "DELETE"},
			Enabled:          true,
			LastModifiedTime: 1699000000,
			PartnerUrl:       "https://partner.zscaler.com",
		}

		data, err := json.Marshal(apiKey)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"keyValue":"ABC123XYZ789"`)
		assert.Contains(t, string(data), `"enabled":true`)
		assert.Contains(t, string(data), `"permissions"`)
	})

	t.Run("ProvisioningAPIKeys JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"keyValue": "XYZ987ABC123",
			"permissions": ["READ", "PROVISION", "MANAGE"],
			"enabled": false,
			"lastModifiedTime": 1699500000,
			"lastModifiedBy": {
				"id": 100,
				"name": "admin@company.com"
			},
			"partnerUrl": "https://custom-partner.com"
		}`

		var apiKey api_keys.ProvisioningAPIKeys
		err := json.Unmarshal([]byte(jsonData), &apiKey)
		require.NoError(t, err)

		assert.Equal(t, 54321, apiKey.ID)
		assert.Equal(t, "XYZ987ABC123", apiKey.KeyValue)
		assert.False(t, apiKey.Enabled)
		assert.Len(t, apiKey.Permissions, 3)
		assert.NotNil(t, apiKey.LastModifiedBy)
		assert.Equal(t, "admin@company.com", apiKey.LastModifiedBy.Name)
	})
}

func TestProvisioning_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse API keys list", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"keyValue": "KEY1ABCDEFGH",
				"permissions": ["READ"],
				"enabled": true
			},
			{
				"id": 2,
				"keyValue": "KEY2IJKLMNOP",
				"permissions": ["READ", "WRITE"],
				"enabled": true
			},
			{
				"id": 3,
				"keyValue": "KEY3QRSTUVWX",
				"permissions": ["READ", "WRITE", "DELETE"],
				"enabled": false
			}
		]`

		var keys []api_keys.ProvisioningAPIKeys
		err := json.Unmarshal([]byte(jsonResponse), &keys)
		require.NoError(t, err)

		assert.Len(t, keys, 3)
		assert.True(t, keys[0].Enabled)
		assert.True(t, keys[1].Enabled)
		assert.False(t, keys[2].Enabled)
		assert.Len(t, keys[0].Permissions, 1)
		assert.Len(t, keys[2].Permissions, 3)
	})

	t.Run("Parse single API key", func(t *testing.T) {
		jsonResponse := `{
			"id": 99999,
			"keyValue": "PARTNERKEY123",
			"permissions": ["FULL_ACCESS"],
			"enabled": true,
			"partnerUrl": "https://partner-api.zscaler.com"
		}`

		var key api_keys.ProvisioningAPIKeys
		err := json.Unmarshal([]byte(jsonResponse), &key)
		require.NoError(t, err)

		assert.Equal(t, 99999, key.ID)
		assert.Equal(t, "PARTNERKEY123", key.KeyValue)
		assert.Contains(t, key.PartnerUrl, "partner-api")
	})
}

