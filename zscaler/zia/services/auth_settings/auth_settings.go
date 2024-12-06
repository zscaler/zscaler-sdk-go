package auth_settings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	authSettingsEndpoint     = "/zia/api/v1/authSettings"
	authSettingsLiteEndpoint = "/zia/api/v1/authSettings/lite"
)

type AuthenticationSettings struct {
	// User authentication type. Setting it to an LDAP-based authentication requires a complete LdapProperties configuration.
	OrgAuthType string `json:"orgAuthType"`

	// When the orgAuthType is NONE, administrators must manually provide the password to new end users.
	OneTimeAuth string `json:"oneTimeAuth"`

	// Whether or not to authenticate users using SAML Single Sign-On. Enabling SAML requires complete SamlSettings.
	SamlEnabled bool `json:"samlEnabled"`

	// Whether or not to authenticate users using Kerberos
	KerberosEnabled bool `json:"kerberosEnabled"`

	// Read Only. Kerberos password can only be set through generateKerberosPassword api.
	KerberosPwd string `json:"kerberosPwd"`

	// How frequently the users are required to authenticate (i.e., cookie expiration duration after a user is first authenticated). This field is not applicable to the Lite API.
	AuthFrequency string `json:"authFrequency"`

	// How frequently the users are required to authenticate. This field is customized to set the value in days. Valid range is 1–180.
	AuthCustomFrequency int `json:"authCustomFrequency"`

	// Password strength required for form-based authentication of hosted DB users. Not applicable for other authentication types (e.g. SAML SSO or Directory).
	PasswordStrength string `json:"passwordStrength"`

	// Password expiration required for form-based authentication of hosted DB users. Not applicable for other authentication types (e.g. SAML SSO or Directory).
	PasswordExpiry string `json:"passwordExpiry"`

	// Timestamp (epoch time in seconds) corresponding to the start of the last LDAP sync.
	LastSyncStartTime int64 `json:"lastSyncStartTime"`

	// Timestamp (epoch time in seconds) corresponding to the end of the last LDAP sync.
	LastSyncEndTime int64 `json:"lastSyncEndTime"`

	// Indicate the use of Mobile Admin as IdP
	MobileAdminSamlIdpEnabled bool `json:"mobileAdminSamlIdpEnabled"`

	// Enable SAML Auto-Provisioning
	AutoProvision bool `json:"autoProvision"`

	// Enable to disable directory synchronization for this user repository type so you can enable SCIM provisioning or SAML auto-provisioning.
	DirectorySyncMigrateToScimEnabled bool `json:"directorySyncMigrateToScimEnabled"`
}

func Get(ctx context.Context, service *zscaler.Service) (*AuthenticationSettings, error) {
	var auth AuthenticationSettings
	err := service.Client.Read(ctx, authSettingsEndpoint, &auth)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning authentication settings from Get: %v", auth)
	return &auth, nil
}

func GetLite(ctx context.Context, service *zscaler.Service) (*AuthenticationSettings, error) {
	var auth AuthenticationSettings
	err := service.Client.Read(ctx, authSettingsLiteEndpoint, &auth)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning authentication settings from Get: %v", auth)
	return &auth, nil
}

func UpdateAuthSettings(ctx context.Context, service *zscaler.Service, authSettings AuthenticationSettings) (*AuthenticationSettings, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, authSettingsEndpoint, authSettings)
	if err != nil {
		return nil, nil, err
	}

	updatedAuthSettings, ok := resp.(*AuthenticationSettings)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected response type")
	}
	service.Client.GetLogger().Printf("[DEBUG] Updated Auth Settings Settings: %+v", updatedAuthSettings)
	return updatedAuthSettings, nil, nil
}
