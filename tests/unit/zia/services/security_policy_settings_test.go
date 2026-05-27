// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/security_policy_settings"
)

const (
	securityWhitelistPath = "/zia/api/v1/security"
	securityBlacklistPath = "/zia/api/v1/security/advanced"
)

func TestSecurityPolicySettings_GetWhiteListUrls_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", securityWhitelistPath, common.SuccessResponse(security_policy_settings.ListUrls{
		White: []string{".example100.com", ".example200.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := security_policy_settings.GetWhiteListUrls(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result.White, 2)
}

func TestSecurityPolicySettings_GetBlackListUrls_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", securityBlacklistPath, common.SuccessResponse(security_policy_settings.ListUrls{
		Black: []string{".malware.example.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := security_policy_settings.GetBlackListUrls(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result.Black, 1)
}

func TestSecurityPolicySettings_GetListUrls_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", securityWhitelistPath, common.SuccessResponse(security_policy_settings.ListUrls{
		White: []string{"trusted.com"},
	}))
	server.On("GET", securityBlacklistPath, common.SuccessResponse(security_policy_settings.ListUrls{
		Black: []string{"malware.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := security_policy_settings.GetListUrls(context.Background(), service)
	require.NoError(t, err)
	assert.Equal(t, []string{"trusted.com"}, result.White)
	assert.Equal(t, []string{"malware.com"}, result.Black)
}

func TestSecurityPolicySettings_UpdateWhiteListUrls_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", securityWhitelistPath, common.SuccessResponse(security_policy_settings.ListUrls{
		White: []string{".example1.com", ".example2.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	list := security_policy_settings.ListUrls{White: []string{".example1.com", ".example2.com"}}
	result, err := security_policy_settings.UpdateWhiteListUrls(context.Background(), service, list)
	require.NoError(t, err)
	assert.Len(t, result.White, 2)
}

func TestSecurityPolicySettings_UpdateBlackListUrls_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", securityBlacklistPath, common.SuccessResponse(security_policy_settings.ListUrls{
		Black: []string{".bad.example.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	list := security_policy_settings.ListUrls{Black: []string{".bad.example.com"}}
	result, err := security_policy_settings.UpdateBlackListUrls(context.Background(), service, list)
	require.NoError(t, err)
	assert.Len(t, result.Black, 1)
}

func TestSecurityPolicySettings_UpdateListUrls_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", securityWhitelistPath, common.SuccessResponse(security_policy_settings.ListUrls{
		White: []string{".white.example.com"},
	}))
	server.On("PUT", securityBlacklistPath, common.SuccessResponse(security_policy_settings.ListUrls{
		Black: []string{".black.example.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	settings := security_policy_settings.ListUrls{
		White: []string{".white.example.com"},
		Black: []string{".black.example.com"},
	}
	result, err := security_policy_settings.UpdateListUrls(context.Background(), service, settings)
	require.NoError(t, err)
	assert.Equal(t, settings.White, result.White)
	assert.Equal(t, settings.Black, result.Black)
}

func TestSecurityPolicySettings_GetWhiteListUrls_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", securityWhitelistPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := security_policy_settings.GetWhiteListUrls(context.Background(), service)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestSecurityPolicySettings_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ListUrls JSON marshaling", func(t *testing.T) {
		list := security_policy_settings.ListUrls{
			White: []string{"trusted.com", "*.internal.com"},
			Black: []string{"malware.com", "phishing.com"},
		}

		data, err := json.Marshal(list)
		require.NoError(t, err)
		assert.Contains(t, string(data), `"whitelistUrls"`)
		assert.Contains(t, string(data), `"blacklistUrls"`)
	})
}
