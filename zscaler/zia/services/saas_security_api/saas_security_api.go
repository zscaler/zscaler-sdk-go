package saas_security_api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	// Endpoint misses important attributes
	// domainProfilesLiteEndpoint          = "/zia/api/v1/domainProfiles/lite"

	// Endpoint is publicly available but not present in the swagger
	domainProfilesEndpoint              = "/zia/api/v1/domainProfiles"
	quarantineTombstoneTemplateEndpoint = "/zia/api/v1/quarantineTombstoneTemplate/lite"
	casbEmailLabelEndpoint              = "/zia/api/v1/casbEmailLabel/lite"
	casbTenantEndpoint                  = "/zia/api/v1/casbTenant"
)

type DomainProfiles struct {
	// Domain profile ID
	ProfileID int `json:"profileId,omitempty"`

	// Domain profile name
	ProfileName string `json:"profileName,omitempty"`

	// A Boolean flag to determine if the organizational domains have to be included in the domain profile
	IncludeCompanyDomains bool `json:"includeCompanyDomains,omitempty"`

	// A Boolean flag to determine whether or not to include subdomains
	IncludeSubdomains bool `json:"includeSubdomains,omitempty"`

	// Additional notes or information about the domain profile
	Description string `json:"description,omitempty"`

	// List of custom domains for the domain profile. There can be one or more custom domains.
	CustomDomains []string `json:"customDomains,omitempty"`

	// List of predefined email service provider domains for the domain profile
	// See SaaS Security API for supported values: https://help.zscaler.com/zia/saas-security-api#/domainProfiles/lite-get
	PredefinedEmailDomains []string `json:"predefinedEmailDomains"`
}

type QuarantineTombstoneLite struct {
	// Tombstone file template ID
	ID int `json:"id,omitempty"`

	// Tombstone file templat
	Name string `json:"name,omitempty"`

	// The text that is included in the tombstone file
	Description string `json:"description,omitempty"`
}

type CasbEmailLabel struct {
	// SaaS Security API email label ID
	ID int `json:"id,omitempty"`

	// SaaS Security API email label name
	Name string `json:"name,omitempty"`

	// Description of the email label
	LabelDesc string `json:"labelDesc,omitempty"`

	// Color to apply to the email label
	LabelColor string `json:"labelColor,omitempty"`

	// A Boolean value that indicates whether or not the email label is deleted
	LabelDeleted bool `json:"labelDeleted,omitempty"`
}

type CasbTenantTags struct {
	// System-generated tag ID
	TagID int `json:"tagId,omitempty"`

	// Tenant ID to which the tag belongs
	TenantID int `json:"tenantId,omitempty"`

	// Universally Unique Identifier (UUID) of the tag
	TagUUID string `json:"tagUUID,omitempty"`

	// Tag name
	TagName string `json:"tagName,omitempty"`

	// A Boolean value that indicates whether or not a tag is deleted
	Deleted bool `json:"deleted,omitempty"`
}

type CasbTenants struct {
	TenantID                 int            `json:"tenantId,omitempty"`
	ModifiedTime             int            `json:"modifiedTime,omitempty"`
	LastTenantValidationTime int            `json:"lastTenantValidationTime,omitempty"`
	TenantDeleted            bool           `json:"tenantDeleted,omitempty"`
	TenantWebhookEnabled     bool           `json:"tenantWebhookEnabled,omitempty"`
	ReAuth                   bool           `json:"reAuth,omitempty"`
	FeaturesSupported        []string       `json:"featuresSupported,omitempty"`
	Status                   []string       `json:"status,omitempty"`
	EnterpriseTenantID       string         `json:"enterpriseTenantId,omitempty"`
	TenantName               string         `json:"tenantName,omitempty"`
	SaaSApplication          string         `json:"saasApplication,omitempty"`
	ZscalerAppTenantID       *common.IDName `json:"zscalerAppTenantId,omitempty"`
}

type CasbTenantScanInfo struct {
	TenantName      string   `json:"tenantName,omitempty"`
	TenantID        int      `json:"tenantId,omitempty"`
	SaasApplication string   `json:"saasApplication,omitempty"`
	ScanInfo        ScanInfo `json:"scanInfo,omitempty"`
	ScanAction      int      `json:"scanAction,omitempty"`
}

type ScanInfo struct {
	CurScanStartTime int `json:"cur_scan_start_time,omitempty"`
	PrevScanEndTime  int `json:"prev_scan_end_time,omitempty"`
	ScanResetNum     int `json:"scan_reset_num,omitempty"`
}

func GetDomainProfiles(ctx context.Context, service *zscaler.Service) ([]DomainProfiles, error) {
	var profiles []DomainProfiles
	err := common.ReadAllPages(ctx, service.Client, domainProfilesEndpoint, &profiles)
	return profiles, err
}

func GetQuarantineTombstoneLite(ctx context.Context, service *zscaler.Service) ([]QuarantineTombstoneLite, error) {
	var templates []QuarantineTombstoneLite
	err := common.ReadAllPages(ctx, service.Client, quarantineTombstoneTemplateEndpoint, &templates)
	return templates, err
}

func GetCasbEmailLabelLite(ctx context.Context, service *zscaler.Service) ([]CasbEmailLabel, error) {
	var labels []CasbEmailLabel
	err := common.ReadAllPages(ctx, service.Client, casbEmailLabelEndpoint, &labels)
	return labels, err
}

func GetCasbTenantTagPolicy(ctx context.Context, service *zscaler.Service, tenantID int) ([]CasbTenantTags, error) {
	var tags []CasbTenantTags
	endpoint := fmt.Sprintf("%s/%d/tags/policy", casbTenantEndpoint, tenantID)
	err := common.ReadAllPages(ctx, service.Client, endpoint, &tags)
	return tags, err
}

func GetCasbTenantLite(ctx context.Context, service *zscaler.Service, queryParams map[string]interface{}) ([]CasbTenants, error) {
	var tenants []CasbTenants
	baseEndpoint := fmt.Sprintf("%s/lite", casbTenantEndpoint)
	queryString := ""
	if len(queryParams) > 0 {
		values := url.Values{}
		for k, v := range queryParams {
			switch val := v.(type) {
			case string:
				values.Set(k, val)
			case bool:
				values.Set(k, strconv.FormatBool(val))
			case []string:
				for _, item := range val {
					values.Add(k, item)
				}
			}
		}
		queryString = "?" + values.Encode()
	}

	endpoint := baseEndpoint + queryString

	err := common.ReadAllPages(ctx, service.Client, endpoint, &tenants)
	return tenants, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]CasbTenantScanInfo, error) {
	var scanInfos []CasbTenantScanInfo
	err := common.ReadAllPages(ctx, service.Client, casbTenantEndpoint+"/scanInfo", &scanInfos)
	return scanInfos, err
}
