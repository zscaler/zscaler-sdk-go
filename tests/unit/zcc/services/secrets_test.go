// Package services provides unit tests for ZCC secrets services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/secrets/getotp"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/secrets/getpasswords"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestSecrets_GetOtp_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getOtp"

	server.On("GET", path, common.SuccessResponse(getotp.OtpResponse{
		Otp:                     "123456",
		ExitOtp:                 "exit-123",
		LogoutOtp:               "logout-456",
		UninstallOtp:            "uninstall-789",
		ZiaDisableOtp:           "zia-disable-001",
		ZpaDisableOtp:           "zpa-disable-002",
		ZdxDisableOtp:           "zdx-disable-003",
		ZdpDisableOtp:           "zdp-disable-004",
		RevertOtp:               "revert-005",
		DeceptionSettingsOtp:    "deception-006",
		AntiTemperingDisableOtp: "anti-temper-007",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := getotp.GetOtp(context.Background(), service, "device-udid-123")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "123456", result.Otp)
	assert.Equal(t, "exit-123", result.ExitOtp)
	assert.Equal(t, "logout-456", result.LogoutOtp)
}

func TestSecrets_GetOtp_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getOtp"

	server.On("GET", path, common.SuccessResponse(getotp.OtpResponse{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := getotp.GetOtp(context.Background(), service, "")

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestSecrets_GetPasswords_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getPasswords"

	server.On("GET", path, common.SuccessResponse(getpasswords.Passwords{
		ExitPass:             "exit-pass-123",
		LogoutPass:           "logout-pass-456",
		UninstallPass:        "uninstall-pass-789",
		ZiaDisablePass:       "zia-disable-pass",
		ZpaDisablePass:       "zpa-disable-pass",
		ZdxDisablePass:       "zdx-disable-pass",
		ZdSettingsAccessPass: "zd-settings-pass",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := getpasswords.GetPasswords(context.Background(), service, "user@example.com", "WINDOWS")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "exit-pass-123", result.ExitPass)
	assert.Equal(t, "logout-pass-456", result.LogoutPass)
	assert.Equal(t, "uninstall-pass-789", result.UninstallPass)
}

func TestSecrets_GetPasswords_NoParams_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getPasswords"

	server.On("GET", path, common.SuccessResponse(getpasswords.Passwords{
		ExitPass:   "default-exit",
		LogoutPass: "default-logout",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := getpasswords.GetPasswords(context.Background(), service, "", "")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "default-exit", result.ExitPass)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestSecrets_Structure(t *testing.T) {
	t.Parallel()

	t.Run("OtpResponse JSON marshaling", func(t *testing.T) {
		otp := getotp.OtpResponse{
			Otp:                     "123456",
			ExitOtp:                 "exit-otp",
			LogoutOtp:               "logout-otp",
			UninstallOtp:            "uninstall-otp",
			ZiaDisableOtp:           "zia-otp",
			ZpaDisableOtp:           "zpa-otp",
			ZdxDisableOtp:           "zdx-otp",
			ZdpDisableOtp:           "zdp-otp",
			RevertOtp:               "revert-otp",
			DeceptionSettingsOtp:    "deception-otp",
			AntiTemperingDisableOtp: "anti-temper-otp",
		}

		data, err := json.Marshal(otp)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"otp":"123456"`)
		assert.Contains(t, string(data), `"exitOtp":"exit-otp"`)
		assert.Contains(t, string(data), `"logoutOtp":"logout-otp"`)
		assert.Contains(t, string(data), `"uninstallOtp":"uninstall-otp"`)
	})

	t.Run("OtpResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"otp": "654321",
			"exitOtp": "exit-654",
			"logoutOtp": "logout-321",
			"uninstallOtp": "uninstall-987",
			"ziaDisableOtp": "zia-111",
			"zpaDisableOtp": "zpa-222",
			"zdxDisableOtp": "zdx-333",
			"zdpDisableOtp": "zdp-444",
			"revertOtp": "revert-555",
			"deceptionSettingsOtp": "deception-666",
			"antiTemperingDisableOtp": "anti-777"
		}`

		var otp getotp.OtpResponse
		err := json.Unmarshal([]byte(jsonData), &otp)
		require.NoError(t, err)

		assert.Equal(t, "654321", otp.Otp)
		assert.Equal(t, "exit-654", otp.ExitOtp)
		assert.Equal(t, "zia-111", otp.ZiaDisableOtp)
		assert.Equal(t, "anti-777", otp.AntiTemperingDisableOtp)
	})

	t.Run("Passwords JSON marshaling", func(t *testing.T) {
		passwords := getpasswords.Passwords{
			ExitPass:             "exit-pass",
			LogoutPass:           "logout-pass",
			UninstallPass:        "uninstall-pass",
			ZiaDisablePass:       "zia-pass",
			ZpaDisablePass:       "zpa-pass",
			ZdxDisablePass:       "zdx-pass",
			ZdSettingsAccessPass: "zd-settings-pass",
		}

		data, err := json.Marshal(passwords)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"exitPass":"exit-pass"`)
		assert.Contains(t, string(data), `"logoutPass":"logout-pass"`)
		assert.Contains(t, string(data), `"uninstallPass":"uninstall-pass"`)
	})

	t.Run("Passwords JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"exitPass": "exit-abc",
			"logoutPass": "logout-def",
			"uninstallPass": "uninstall-ghi",
			"ziaDisablePass": "zia-jkl",
			"zpaDisablePass": "zpa-mno",
			"zdxDisablePass": "zdx-pqr",
			"zdSettingsAccessPass": "zd-stu"
		}`

		var passwords getpasswords.Passwords
		err := json.Unmarshal([]byte(jsonData), &passwords)
		require.NoError(t, err)

		assert.Equal(t, "exit-abc", passwords.ExitPass)
		assert.Equal(t, "logout-def", passwords.LogoutPass)
		assert.Equal(t, "zia-jkl", passwords.ZiaDisablePass)
		assert.Equal(t, "zd-stu", passwords.ZdSettingsAccessPass)
	})

	t.Run("GetOtpQuery JSON marshaling", func(t *testing.T) {
		query := getotp.GetOtpQuery{
			Udid: "device-udid-123",
		}

		data, err := json.Marshal(query)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"udid":"device-udid-123"`)
	})

	t.Run("GetPasswordsQueryParams JSON marshaling", func(t *testing.T) {
		query := getpasswords.GetPasswordsQueryParams{
			Username: "user@example.com",
			OsType:   "WINDOWS",
		}

		data, err := json.Marshal(query)
		require.NoError(t, err)

		// The struct uses PascalCase for JSON tags via the `url` tag
		assert.Contains(t, string(data), `user@example.com`)
		assert.Contains(t, string(data), `WINDOWS`)
	})
}

func TestSecrets_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse OTP response", func(t *testing.T) {
		jsonResponse := `{
			"otp": "999888",
			"exitOtp": "exit-999",
			"logoutOtp": "logout-888",
			"uninstallOtp": "uninstall-777"
		}`

		var otp getotp.OtpResponse
		err := json.Unmarshal([]byte(jsonResponse), &otp)
		require.NoError(t, err)

		assert.Equal(t, "999888", otp.Otp)
		assert.Equal(t, "exit-999", otp.ExitOtp)
	})

	t.Run("Parse Passwords response", func(t *testing.T) {
		jsonResponse := `{
			"exitPass": "secret-exit",
			"logoutPass": "secret-logout",
			"uninstallPass": "secret-uninstall"
		}`

		var passwords getpasswords.Passwords
		err := json.Unmarshal([]byte(jsonResponse), &passwords)
		require.NoError(t, err)

		assert.Equal(t, "secret-exit", passwords.ExitPass)
		assert.Equal(t, "secret-logout", passwords.LogoutPass)
	})
}

