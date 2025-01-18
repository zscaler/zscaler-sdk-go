package end_user_notification

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	userNotificationEndpoint = "/zia/api/v1/eun"
)

type UserNotificationSettings struct {
	AUPFrequency                        string `json:"aupFrequency"`
	AUPCustomFrequency                  int    `json:"aupCustomFrequency"`
	AUPDayOffset                        int    `json:"aupDayOffset"`
	AUPMessage                          string `json:"aupMessage"`
	NotificationType                    string `json:"notificationType"`
	DisplayReason                       bool   `json:"displayReason"`
	DisplayCompName                     bool   `json:"displayCompName"`
	DisplayCompLogo                     bool   `json:"displayCompLogo"`
	CustomText                          string `json:"customText"`
	URLCatReviewEnabled                 bool   `json:"urlCatReviewEnabled"`
	URLCatReviewSubmitToSecurityCloud   bool   `json:"urlCatReviewSubmitToSecurityCloud"`
	URLCatReviewCustomLocation          string `json:"urlCatReviewCustomLocation"`
	URLCatReviewText                    string `json:"urlCatReviewText"`
	SecurityReviewEnabled               bool   `json:"securityReviewEnabled"`
	SecurityReviewSubmitToSecurityCloud bool   `json:"securityReviewSubmitToSecurityCloud"`
	SecurityReviewCustomLocation        string `json:"securityReviewCustomLocation"`
	SecurityReviewText                  string `json:"securityReviewText"`
	WebDLPReviewEnabled                 bool   `json:"webDlpReviewEnabled"`
	WebDLPReviewSubmitToSecurityCloud   bool   `json:"webDlpReviewSubmitToSecurityCloud"`
	WebDLPReviewCustomLocation          string `json:"webDlpReviewCustomLocation"`
	WebDLPReviewText                    string `json:"webDlpReviewText"`
	RedirectURL                         string `json:"redirectUrl,omitempty"`
	SupportEmail                        string `json:"supportEmail"`
	SupportPhone                        string `json:"supportPhone"`
	OrgPolicyLink                       string `json:"orgPolicyLink"`
	CautionAgainAfter                   int    `json:"cautionAgainAfter"`
	CautionPerDomain                    bool   `json:"cautionPerDomain"`
	CautionCustomText                   string `json:"cautionCustomText"`
	IDPProxyNotificationText            string `json:"idpProxyNotificationText"`
	QuarantineCustomNotificationText    string `json:"quarantineCustomNotificationText"`
}

func GetUserNotificationSettings(ctx context.Context, service *zscaler.Service) (*UserNotificationSettings, error) {
	var notificationSettings UserNotificationSettings
	err := service.Client.Read(ctx, userNotificationEndpoint, &notificationSettings)
	if err != nil {
		return nil, err
	}
	return &notificationSettings, nil
}

func UpdateUserNotificationSettings(ctx context.Context, service *zscaler.Service, settings UserNotificationSettings) (*UserNotificationSettings, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, userNotificationEndpoint, settings)
	if err != nil {
		return nil, nil, err
	}

	notificationSettings, ok := resp.(*UserNotificationSettings)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected response type")
	}
	service.Client.GetLogger().Printf("[DEBUG] Updated End User Notification Settings : %+v", notificationSettings)
	return notificationSettings, nil, nil
}
