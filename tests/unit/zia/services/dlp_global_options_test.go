// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	ziacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_global_options"
)

const dlpGlobalOptionsPath = "/zia/api/v1/webDlpGlobalOptions"

// =====================================================
// GetDLPGlobalOptions
// =====================================================

func TestDLPGlobalOptions_GetDLPGlobalOptions_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", dlpGlobalOptionsPath, common.SuccessResponse(dlp_global_options.WebDlpGlobal{
		Applications:              []string{"APP1", "APP2"},
		Urls:                      []string{"https://example.com"},
		ExemptUrlEncodedData:      true,
		EnableInlineDlpOcr:        true,
		EnableEvaluateAllDlpRules: true,
		URLCategories: []ziacommon.IDNameExtensions{
			{ID: 1, Name: "SOCIAL_NETWORKING"},
		},
		HttpGetCustomUrlCategories: []string{"CUSTOM_01"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_global_options.GetDLPGlobalOptions(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, []string{"APP1", "APP2"}, result.Applications)
	assert.True(t, result.ExemptUrlEncodedData)
	assert.True(t, result.EnableInlineDlpOcr)
	assert.True(t, result.EnableEvaluateAllDlpRules)
	require.Len(t, result.URLCategories, 1)
	assert.Equal(t, "SOCIAL_NETWORKING", result.URLCategories[0].Name)
}

func TestDLPGlobalOptions_GetDLPGlobalOptions_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", dlpGlobalOptionsPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_global_options.GetDLPGlobalOptions(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// UpdateDLPGlobalOptions
// =====================================================

func TestDLPGlobalOptions_UpdateDLPGlobalOptions_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", dlpGlobalOptionsPath, common.SuccessResponse(dlp_global_options.WebDlpGlobal{
		ExemptUrlEncodedData:      true,
		EnableInlineDlpOcr:        true,
		EnableCasbOcr:             true,
		EnableEmailDlpOcr:         true,
		EnableEvaluateAllDlpRules: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	settings := dlp_global_options.WebDlpGlobal{
		ExemptUrlEncodedData:      true,
		EnableInlineDlpOcr:        true,
		EnableCasbOcr:             true,
		EnableEmailDlpOcr:         true,
		EnableEvaluateAllDlpRules: true,
	}

	result, _, err := dlp_global_options.UpdateDLPGlobalOptions(context.Background(), service, settings)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.ExemptUrlEncodedData)
	assert.True(t, result.EnableInlineDlpOcr)
	assert.True(t, result.EnableCasbOcr)
	assert.True(t, result.EnableEmailDlpOcr)
}

func TestDLPGlobalOptions_UpdateDLPGlobalOptions_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", dlpGlobalOptionsPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	settings := dlp_global_options.WebDlpGlobal{EnableInlineDlpOcr: true}

	result, _, err := dlp_global_options.UpdateDLPGlobalOptions(context.Background(), service, settings)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDLPGlobalOptions_UpdateDLPGlobalOptions_UnexpectedResponseType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Empty 204 body makes UpdateWithPut return (nil, nil); the service
	// type-asserts the response and surfaces "unexpected response type".
	server.On("PUT", dlpGlobalOptionsPath, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	settings := dlp_global_options.WebDlpGlobal{EnableInlineDlpOcr: true}

	result, _, err := dlp_global_options.UpdateDLPGlobalOptions(context.Background(), service, settings)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unexpected response type")
}

// =====================================================
// Structure Tests
// =====================================================

func TestDLPGlobalOptions_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WebDlpGlobal JSON marshaling", func(t *testing.T) {
		settings := dlp_global_options.WebDlpGlobal{
			Applications:                []string{"APP1"},
			Urls:                        []string{"https://corp.example.com"},
			ExemptUrlEncodedData:        true,
			EnableNpkEdmTemplates:       true,
			EnableNpkEdmTemplatesForOrg: true,
			EnableInlineDlpOcr:          true,
			EnableCasbOcr:               true,
			EnableEmailDlpOcr:           true,
			EnableEvaluateAllDlpRules:   true,
			EnableEdmPopularFormat:      true,
			HttpGetCustomUrlCategories:  []string{"CUSTOM_01"},
			URLCategories: []ziacommon.IDNameExtensions{
				{ID: 5, Name: "FINANCE"},
			},
		}

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"applications":["APP1"]`)
		assert.Contains(t, string(data), `"exemptUrlEncodedData":true`)
		assert.Contains(t, string(data), `"enableInlineDlpOcr":true`)
		assert.Contains(t, string(data), `"enableEvaluateAllDlpRules":true`)
		assert.Contains(t, string(data), `"httpGetCustomUrlCategories":["CUSTOM_01"]`)
		assert.Contains(t, string(data), `"urlCategories"`)
	})

	t.Run("WebDlpGlobal JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"applications": ["A", "B"],
			"urls": ["https://x.example.com"],
			"exemptUrlEncodedData": true,
			"enableInlineDlpOcr": false,
			"enableCasbOcr": true,
			"urlCategories": [{"id": 9, "name": "GAMBLING"}],
			"httpGetCustomUrlCategories": ["CUST"]
		}`

		var settings dlp_global_options.WebDlpGlobal
		err := json.Unmarshal([]byte(jsonData), &settings)
		require.NoError(t, err)

		assert.Equal(t, []string{"A", "B"}, settings.Applications)
		assert.True(t, settings.ExemptUrlEncodedData)
		assert.False(t, settings.EnableInlineDlpOcr)
		assert.True(t, settings.EnableCasbOcr)
		require.Len(t, settings.URLCategories, 1)
		assert.Equal(t, 9, settings.URLCategories[0].ID)
		assert.Equal(t, "GAMBLING", settings.URLCategories[0].Name)
	})
}
