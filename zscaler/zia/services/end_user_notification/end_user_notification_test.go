package end_user_notification

import (
	"context"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestUserNotificationSettings(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	ctx := context.Background()

	t.Run("RetrieveUserNotificationSettings", func(t *testing.T) {
		settings, err := GetUserNotificationSettings(ctx, service)
		if err != nil {
			t.Fatalf("Error retrieving user notification settings: %v", err)
		}
		t.Logf("Successfully retrieved user notification settings: %+v", settings)
	})

	t.Run("UpdateUserNotificationSettings", func(t *testing.T) {
		NotificationSettings, err := GetUserNotificationSettings(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching end user notification settings: %v", err)
		}

		updatedSettings := *NotificationSettings
		updatedSettings.AUPFrequency = "NEVER"
		// updatedSettings.AUPCustomFrequency = true
		updatedSettings.AUPDayOffset = 1
		updatedSettings.AUPMessage = ""
		updatedSettings.NotificationType = "CUSTOM"
		updatedSettings.DisplayReason = false
		updatedSettings.DisplayCompName = false
		updatedSettings.DisplayCompLogo = false
		updatedSettings.CustomText = ""
		updatedSettings.URLCatReviewEnabled = true
		updatedSettings.URLCatReviewSubmitToSecurityCloud = true
		// updatedSettings.URLCatReviewCustomLocation = true
		updatedSettings.URLCatReviewText = "If you believe you received this message in error, please click here to request a review of this site."
		updatedSettings.SecurityReviewEnabled = true
		updatedSettings.SecurityReviewSubmitToSecurityCloud = true
		// updatedSettings.SecurityReviewCustomLocation = true
		updatedSettings.SecurityReviewText = "Click to request security review."
		updatedSettings.WebDLPReviewEnabled = true
		// updatedSettings.WebDLPReviewSubmitToSecurityCloud = true
		updatedSettings.WebDLPReviewCustomLocation = "https://redirect.acme.com"
		updatedSettings.WebDLPReviewText = "Click to request policy review."
		updatedSettings.RedirectURL = "https://redirect.acme.com"
		updatedSettings.SupportEmail = "support@8061240.zscalerbeta.net"
		updatedSettings.SupportPhone = "+91-9000000000"
		updatedSettings.OrgPolicyLink = "http://8061240.zscalerbeta.net/policy.html"
		updatedSettings.CautionAgainAfter = 300
		updatedSettings.CautionPerDomain = true
		updatedSettings.CautionCustomText = "Proceeding to visit the site may violate your company policy. Press the \"Continue\" button to access the site anyway or press the \"Back\" button on your browser to go back"
		updatedSettings.IDPProxyNotificationText = ""
		updatedSettings.QuarantineCustomNotificationText = "We are checking this file for a potential security risk.The file you attempted to download is being analyzed for your protection. \n        \tIt is not blocked. The analysis can take up to 10 minutes, depending on the size and type of the file.If safe, your file downloads automatically. \n        \tIf unsafe, the file will be blocked."

		result, _, err := UpdateUserNotificationSettings(ctx, service, updatedSettings)
		if err != nil {
			t.Fatalf("Error updating end user notification settings: %v", err)
		}
		t.Logf("Successfully updated end user notification settings: %+v", result)
	})
}
