package end_user_notification

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	userNotificationEndpoint = "/zia/api/v1/eun"
	eunTemplateEndpoint      = "/zia/api/v1/eunTemplate"
	userConfirmationEndpoint = "/zia/api/v1/userConfirmation"
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

type EunTemplateProduct struct {
	ID                  int                 `json:"id,omitempty"`
	Name                string              `json:"name,omitempty"`
	Channel             string              `json:"channel,omitempty"`
	Product             string              `json:"product,omitempty"`
	Type                string              `json:"type,omitempty"`
	NotificationDetails []string            `json:"notificationDetails,omitempty"`
	CautionInterval     string              `json:"cautionInterval,omitempty"`
	Default             bool                `json:"default,omitempty"`
	RecommendedCloudApp RecommendedCloudApp `json:"recommendedCloudApp,omitempty"`
	LanguageTemplates   []LanguageTemplates `json:"languageTemplates,omitempty"`
}

type RecommendedCloudApp struct {
	Val                 int    `json:"val,omitempty"`
	Name                string `json:"name,omitempty"`
	Channel             string `json:"channel,omitempty"`
	Product             string `json:"product,omitempty"`
	Type                string `json:"type,omitempty"`
	CautionMiscInterval string `json:"misc,omitempty"`
	AppNotReady         bool   `json:"appNotReady,omitempty"`
	UnderMigration      bool   `json:"underMigration,omitempty"`
	AppCatModified      bool   `json:"appCatModified,omitempty"`
	Deprecated          bool   `json:"deprecated,omitempty"`
}

type LanguageTemplates struct {
	Language                string `json:"language,omitempty"`
	AllowMessage            string `json:"allowMessage,omitempty"`
	BlockMessage            string `json:"blockMessage,omitempty"`
	EncryptMessage          string `json:"encryptMessage,omitempty"`
	ReadonlyMessage         string `json:"readonlyMessage,omitempty"`
	CautionMessage          string `json:"cautionMessage,omitempty"`
	RedirectResponseMessage string `json:"redirectResponseMessage,omitempty"`
	Default                 bool   `json:"default,omitempty"`
}

func GetEunTemplateBrowserBasedZCC(ctx context.Context, service *zscaler.Service, templateType, product string) ([]EunTemplateProduct, error) {
	var eunTemplates []EunTemplateProduct
	endpoint := fmt.Sprintf("%s/%s/product/%s", eunTemplateEndpoint, templateType, product)
	err := service.Client.Read(ctx, endpoint, &eunTemplates)
	if err != nil {
		return nil, err
	}
	return eunTemplates, nil
}

type EunEnablementStatus struct {
	InlineDlpStatus         map[string]string `json:"inlineDlpStatus,omitempty"`
	EptDlpStatus            map[string]string `json:"eptDlpStatus,omitempty"`
	CloudAppStatus          map[string]string `json:"cloudAppStatus,omitempty"`
	UrlFilteringStatus      map[string]string `json:"urlFilteringStatus,omitempty"`
	DnsRuleStatus           map[string]string `json:"dnsRuleStatus,omitempty"`
	FirewallFilteringStatus map[string]string `json:"firewallFilteringStatus,omitempty"`
	IpsControlStatus        map[string]string `json:"ipsControlStatus,omitempty"`
	FileTypeFilteringStatus map[string]string `json:"fileTypeFilteringStatus,omitempty"`
}

func GetEunEnablementStatus(ctx context.Context, service *zscaler.Service, templateType string, productType *string) (*EunEnablementStatus, error) {
	var status EunEnablementStatus
	endpoint := fmt.Sprintf("%s/%s/featureEnablementStatus", eunTemplateEndpoint, templateType)

	if productType != nil && *productType != "" {
		queryParams := url.Values{}
		queryParams.Set("productType", *productType)
		endpoint += "?" + queryParams.Encode()
	}

	err := service.Client.Read(ctx, endpoint, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

type UserConfirmationByPolicyType struct {
	ID                int                     `json:"id,omitempty"`
	Name              string                  `json:"name,omitempty"`
	Channel           string                  `json:"channel,omitempty"`
	Product           string                  `json:"product,omitempty"`
	Default           bool                    `json:"default"`
	LanguageTemplates []LanguageTemplatesLite `json:"languageTemplates,omitempty"`
}

type LanguageTemplatesLite struct {
	Language string `json:"language,omitempty"`
	Message  string `json:"message,omitempty"`
	Default  bool   `json:"default"`
}

// GetEunTemplateByPolicy retrieves the list of user confirmation notification
// templates by policy type.
//
// The product path parameter is required and identifies the policy type
// associated with the notification template (INLINE, ENDPOINT_DLP, CLOUDAPP,
// URL, FILE_TYPE, FIREWALL, DNS, IPS).
func GetEunTemplateByPolicy(ctx context.Context, service *zscaler.Service, product string) ([]UserConfirmationByPolicyType, error) {
	var templates []UserConfirmationByPolicyType
	endpoint := fmt.Sprintf("%s/product/%s", userConfirmationEndpoint, product)
	err := service.Client.Read(ctx, endpoint, &templates)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// GetGlobalDefaultTemplates retrieves the global default user confirmation
// templates for all policy types and policy channels where applicable. This
// endpoint takes no parameters.
func GetGlobalDefaultTemplates(ctx context.Context, service *zscaler.Service) ([]UserConfirmationByPolicyType, error) {
	var templates []UserConfirmationByPolicyType
	err := service.Client.Read(ctx, userConfirmationEndpoint+"/globalDefaultTemplates", &templates)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

type NotificationEnablementFeatureStatus struct {
	InlineDlpConfirmStatus   string `json:"inlineDlpConfirmStatus,omitempty"`
	EndpointDlpConfirmStatus string `json:"endpointDlpConfirmStatus,omitempty"`
}

// GetNotificationEnablementFeatureStatus retrieves the notification enablement
// status for user confirmation notifications.
//
// The templateType path parameter is required (ZCC, BROWSER). User confirmation
// notifications are supported using Zscaler Client Connector only, so the type
// is typically ZCC for this request. The productType query parameter is optional
// and filters by policy type (INLINE, ENDPOINT_DLP). Pass nil to omit it.
func GetNotificationEnablementFeatureStatus(ctx context.Context, service *zscaler.Service, templateType string, productType *string) (*NotificationEnablementFeatureStatus, error) {
	var status NotificationEnablementFeatureStatus
	endpoint := fmt.Sprintf("%s/%s/featureEnablementStatus", userConfirmationEndpoint, templateType)

	if productType != nil && *productType != "" {
		queryParams := url.Values{}
		queryParams.Set("productType", *productType)
		endpoint += "?" + queryParams.Encode()
	}

	err := service.Client.Read(ctx, endpoint, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}
