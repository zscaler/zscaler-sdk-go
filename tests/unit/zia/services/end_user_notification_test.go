// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/end_user_notification"
)

const eunPath = "/zia/api/v1/eun"

// sampleEUNSettings mirrors the integration test update payload in end_user_notification_test.go.
func sampleEUNSettings() end_user_notification.UserNotificationSettings {
	return end_user_notification.UserNotificationSettings{
		AUPFrequency:       "NEVER",
		AUPDayOffset:       1,
		AUPMessage:         "",
		NotificationType:   "CUSTOM",
		DisplayReason:      false,
		DisplayCompName:    false,
		DisplayCompLogo:    false,
		CustomText:         "",
		URLCatReviewEnabled:               true,
		URLCatReviewSubmitToSecurityCloud: true,
		URLCatReviewText:                  "If you believe you received this message in error, please click here to request a review of this site.",
		SecurityReviewEnabled:               true,
		SecurityReviewSubmitToSecurityCloud: true,
		SecurityReviewText:                  "Click to request security review.",
		WebDLPReviewEnabled:                 true,
		WebDLPReviewCustomLocation:          "https://redirect.acme.com",
		WebDLPReviewText:                    "Click to request policy review.",
		RedirectURL:                         "https://redirect.acme.com",
		SupportEmail:                        "support@8061240.zscalerbeta.net",
		SupportPhone:                        "+91-9000000000",
		OrgPolicyLink:                       "http://8061240.zscalerbeta.net/policy.html",
		CautionAgainAfter:                   300,
		CautionPerDomain:                    true,
		CautionCustomText:                   "Proceeding to visit the site may violate your company policy. Press the \"Continue\" button to access the site anyway or press the \"Back\" button on your browser to go back",
		IDPProxyNotificationText:            "",
		QuarantineCustomNotificationText:  "We are checking this file for a potential security risk.The file you attempted to download is being analyzed for your protection. \n        \tIt is not blocked. The analysis can take up to 10 minutes, depending on the size and type of the file.If safe, your file downloads automatically. \n        \tIf unsafe, the file will be blocked.",
	}
}

// =====================================================
// SDK Function Tests
// =====================================================

func TestEndUserNotification_GetUserNotificationSettings_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", eunPath, common.SuccessResponse(sampleEUNSettings()))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := end_user_notification.GetUserNotificationSettings(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "NEVER", result.AUPFrequency)
	assert.Equal(t, "CUSTOM", result.NotificationType)
	assert.Equal(t, "https://redirect.acme.com", result.RedirectURL)
	assert.Equal(t, "support@8061240.zscalerbeta.net", result.SupportEmail)
	assert.Equal(t, "+91-9000000000", result.SupportPhone)
	assert.Equal(t, 300, result.CautionAgainAfter)
	assert.True(t, result.CautionPerDomain)
	assert.True(t, result.URLCatReviewEnabled)
	assert.True(t, result.SecurityReviewEnabled)
	assert.True(t, result.WebDLPReviewEnabled)
}

func TestEndUserNotification_GetUserNotificationSettings_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", eunPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := end_user_notification.GetUserNotificationSettings(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestEndUserNotification_UpdateUserNotificationSettings_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	settings := sampleEUNSettings()
	server.On("PUT", eunPath, common.SuccessResponse(settings))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := end_user_notification.UpdateUserNotificationSettings(context.Background(), service, settings)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "NEVER", result.AUPFrequency)
	assert.Equal(t, "CUSTOM", result.NotificationType)
	assert.Equal(t, "https://redirect.acme.com", result.WebDLPReviewCustomLocation)
	assert.Contains(t, result.URLCatReviewText, "request a review of this site")
	assert.Contains(t, result.QuarantineCustomNotificationText, "potential security risk")
}

func TestEndUserNotification_UpdateUserNotificationSettings_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", eunPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := end_user_notification.UpdateUserNotificationSettings(context.Background(), service, sampleEUNSettings())

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestEndUserNotification_UpdateUserNotificationSettings_UnexpectedResponseType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", eunPath, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := end_user_notification.UpdateUserNotificationSettings(context.Background(), service, sampleEUNSettings())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unexpected response type")
}

// =====================================================
// Structure Tests
// =====================================================

func TestEndUserNotification_Structure(t *testing.T) {
	t.Parallel()

	t.Run("UserNotificationSettings JSON marshaling", func(t *testing.T) {
		settings := sampleEUNSettings()

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"aupFrequency":"NEVER"`)
		assert.Contains(t, string(data), `"notificationType":"CUSTOM"`)
		assert.Contains(t, string(data), `"redirectUrl":"https://redirect.acme.com"`)
		assert.Contains(t, string(data), `"supportEmail":"support@8061240.zscalerbeta.net"`)
		assert.Contains(t, string(data), `"cautionAgainAfter":300`)
		assert.Contains(t, string(data), `"cautionPerDomain":true`)
		assert.Contains(t, string(data), `"urlCatReviewEnabled":true`)
		assert.Contains(t, string(data), `"webDlpReviewEnabled":true`)
	})

	t.Run("UserNotificationSettings JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"aupFrequency": "NEVER",
			"aupDayOffset": 1,
			"notificationType": "CUSTOM",
			"displayReason": false,
			"displayCompName": false,
			"displayCompLogo": false,
			"urlCatReviewEnabled": true,
			"urlCatReviewSubmitToSecurityCloud": true,
			"urlCatReviewText": "If you believe you received this message in error, please click here to request a review of this site.",
			"securityReviewEnabled": true,
			"securityReviewSubmitToSecurityCloud": true,
			"securityReviewText": "Click to request security review.",
			"webDlpReviewEnabled": true,
			"webDlpReviewCustomLocation": "https://redirect.acme.com",
			"webDlpReviewText": "Click to request policy review.",
			"redirectUrl": "https://redirect.acme.com",
			"supportEmail": "support@8061240.zscalerbeta.net",
			"supportPhone": "+91-9000000000",
			"orgPolicyLink": "http://8061240.zscalerbeta.net/policy.html",
			"cautionAgainAfter": 300,
			"cautionPerDomain": true,
			"cautionCustomText": "Proceeding to visit the site may violate your company policy.",
			"quarantineCustomNotificationText": "We are checking this file for a potential security risk."
		}`

		var settings end_user_notification.UserNotificationSettings
		err := json.Unmarshal([]byte(jsonData), &settings)
		require.NoError(t, err)

		assert.Equal(t, "NEVER", settings.AUPFrequency)
		assert.Equal(t, "CUSTOM", settings.NotificationType)
		assert.Equal(t, "https://redirect.acme.com", settings.RedirectURL)
		assert.Equal(t, 300, settings.CautionAgainAfter)
		assert.True(t, settings.CautionPerDomain)
		assert.True(t, settings.URLCatReviewSubmitToSecurityCloud)
	})
}

func TestEndUserNotification_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse end user notification settings response", func(t *testing.T) {
		jsonResponse := `{
			"aupFrequency": "NEVER",
			"notificationType": "CUSTOM",
			"redirectUrl": "https://redirect.acme.com",
			"supportEmail": "support@8061240.zscalerbeta.net"
		}`

		var settings end_user_notification.UserNotificationSettings
		err := json.Unmarshal([]byte(jsonResponse), &settings)
		require.NoError(t, err)

		assert.Equal(t, "NEVER", settings.AUPFrequency)
		assert.Equal(t, "support@8061240.zscalerbeta.net", settings.SupportEmail)
	})
}
