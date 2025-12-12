// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/ftp_control_policy"
)

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
