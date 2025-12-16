// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/ftp_control_policy"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestFTPControlPolicy_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/ftpSettings"

	server.On("GET", path, common.SuccessResponse(ftp_control_policy.FTPControlPolicy{
		FtpOverHttpEnabled: true,
		FtpEnabled:         true,
		UrlCategories:      []string{"BUSINESS", "FINANCE"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ftp_control_policy.GetFTPControlPolicy(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.FtpOverHttpEnabled)
	assert.True(t, result.FtpEnabled)
}

func TestFTPControlPolicy_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/ftpSettings"

	server.On("PUT", path, common.SuccessResponse(ftp_control_policy.FTPControlPolicy{
		FtpOverHttpEnabled: false,
		FtpEnabled:         true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updatePolicy := &ftp_control_policy.FTPControlPolicy{
		FtpOverHttpEnabled: false,
		FtpEnabled:         true,
	}

	result, _, err := ftp_control_policy.UpdateFTPControlPolicy(context.Background(), service, updatePolicy)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.FtpOverHttpEnabled)
}

// =====================================================
// Structure Tests
// =====================================================

func TestFTPControlPolicy_Structure(t *testing.T) {
	t.Parallel()

	t.Run("FTPControlPolicy JSON marshaling", func(t *testing.T) {
		policy := ftp_control_policy.FTPControlPolicy{
			FtpOverHttpEnabled: true,
			FtpEnabled:         true,
			UrlCategories:      []string{"BUSINESS", "FINANCE"},
			Urls:               []string{"ftp.company.com", "*.ftp.example.com"},
		}

		data, err := json.Marshal(policy)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"ftpOverHttpEnabled":true`)
		assert.Contains(t, string(data), `"ftpEnabled":true`)
	})

	t.Run("FTPControlPolicy JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"ftpOverHttpEnabled": true,
			"ftpEnabled": false,
			"urlCategories": ["UNCATEGORIZED", "ADULT_CONTENT"],
			"urls": ["ftp.internal.com"]
		}`

		var policy ftp_control_policy.FTPControlPolicy
		err := json.Unmarshal([]byte(jsonData), &policy)
		require.NoError(t, err)

		assert.True(t, policy.FtpOverHttpEnabled)
		assert.False(t, policy.FtpEnabled)
		assert.Len(t, policy.UrlCategories, 2)
		assert.Len(t, policy.Urls, 1)
	})

	t.Run("FTPControlPolicy disabled", func(t *testing.T) {
		jsonData := `{
			"ftpOverHttpEnabled": false,
			"ftpEnabled": false
		}`

		var policy ftp_control_policy.FTPControlPolicy
		err := json.Unmarshal([]byte(jsonData), &policy)
		require.NoError(t, err)

		assert.False(t, policy.FtpOverHttpEnabled)
		assert.False(t, policy.FtpEnabled)
	})
}

func TestFTPControlPolicy_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse FTP settings", func(t *testing.T) {
		jsonResponse := `{
			"ftpOverHttpEnabled": true,
			"ftpEnabled": true,
			"urlCategories": ["ANY"],
			"urls": []
		}`

		var policy ftp_control_policy.FTPControlPolicy
		err := json.Unmarshal([]byte(jsonResponse), &policy)
		require.NoError(t, err)

		assert.True(t, policy.FtpOverHttpEnabled)
		assert.True(t, policy.FtpEnabled)
	})
}
