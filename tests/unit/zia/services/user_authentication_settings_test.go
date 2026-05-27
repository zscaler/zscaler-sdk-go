// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/user_authentication_settings"
)

const exemptedUrlsPath = "/zia/api/v1/authSettings/exemptedUrls"

func TestUserAuthenticationSettings_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", exemptedUrlsPath, common.SuccessResponse(user_authentication_settings.ExemptedUrls{
		URLs: []string{"site100.example.com", "site200.test.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := user_authentication_settings.Get(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result.URLs, 2)
}

func TestUserAuthenticationSettings_Update_AddURLs_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", exemptedUrlsPath, common.SuccessResponse(user_authentication_settings.ExemptedUrls{
		URLs: []string{"existing.example.com"},
	}))
	server.On("POST", exemptedUrlsPath, common.SuccessResponse(user_authentication_settings.ExemptedUrls{
		URLs: []string{"site123.example.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newURLs := user_authentication_settings.ExemptedUrls{
		URLs: []string{"existing.example.com", "site123.example.com"},
	}
	result, err := user_authentication_settings.Update(context.Background(), service, newURLs)
	require.NoError(t, err)
	assert.Contains(t, result.URLs, "site123.example.com")
}

func TestUserAuthenticationSettings_Update_RemoveURLs_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", exemptedUrlsPath, common.SuccessResponse(user_authentication_settings.ExemptedUrls{
		URLs: []string{"keep.example.com", "remove.example.com"},
	}))
	server.On("POST", exemptedUrlsPath, common.SuccessResponse(user_authentication_settings.ExemptedUrls{
		URLs: []string{"remove.example.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	target := user_authentication_settings.ExemptedUrls{URLs: []string{"keep.example.com"}}
	result, err := user_authentication_settings.Update(context.Background(), service, target)
	require.NoError(t, err)
	assert.Equal(t, []string{"keep.example.com"}, result.URLs)
}

func TestUserAuthenticationSettings_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", exemptedUrlsPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := user_authentication_settings.Get(context.Background(), service)
	require.Error(t, err)
	assert.Nil(t, result)
}
