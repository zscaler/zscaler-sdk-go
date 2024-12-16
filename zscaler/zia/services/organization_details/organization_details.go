package organization_details

import (
	"context"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	subscriptionsEndpoint      = "/zia/api/v1/subscriptions"
	orgInformationEndpoint     = "/zia/api/v1/orgInformation"
	orgInformationLiteEndpoint = "/zia/api/v1/orgInformation/lite"
)

type Subscription struct {
	// The unique identifier of the subscription
	ID string `json:"id"`

	// The status of the subscription indicating whether it is on trial, expired and disabled, not subscribed, subscribed, or temporarily extended
	Status string `json:"status"`

	// The current state of the subscription indicating whether it is active, expired, disabled, or not started.
	State string `json:"state"`

	// The number of licenses owned by the tenant under the subscription. The license could be for users, Virtual Service Edges, proxy ports, NSS feeds, etc.
	Licenses int `json:"licenses"`

	// The subscription start time in Unix time format
	StartDate int `json:"startDate"`

	// The subscription start date in MM/DD/YYYY format
	StrStartDate string `json:"strStartDate"`

	// The subscription end date in MM/DD/YYYY format
	StrEndDate string `json:"strEndDate"`

	// The subscription end time in Unix time format
	EndDate int `json:"endDate"`

	// The ID of the service enabled through the subscription
	SKU string `json:"sku"`

	CellCount string `json:"cellCount"`

	// The timestamp in Unix time format when the subscription was last updated
	UpdatedAtTimestamp int `json:"updatedAtTimestamp"`

	// A Boolean value indicating that the subscription is active or is not started
	Subscribed bool `json:"subscribed"`
}

type Organization struct {
	// Organization identifier
	OrgID int `json:"orgId"`

	// Name of the organization
	Name string `json:"name"`

	// Headquarter location
	HQLocation string `json:"hqLocation"`

	// Domain names
	Domains []string `json:"domains"`

	// Geographic Region Group
	GeoLocation string `json:"geoLocation"`

	// Industry Vertical Group
	IndustryVertical string `json:"industryVertical"`

	// Address Line 1
	AddrLine1 string `json:"addrLine1"`

	// Address Line 2
	AddrLine2 string `json:"addrLine2"`

	// City
	City string `json:"city"`

	// State
	State string `json:"state"`

	// Zip Code
	ZipCode string `json:"zipcode"`

	// Country Code
	Country string `json:"country"`

	// Number of employees
	EmployeeCount string `json:"employeeCount"`

	// Primary Language. If not set, set to NONE.
	Language string `json:"language"`

	// Primary Time Zone. If not set, set to GMT.
	Timezone string `json:"timezone"`

	// Time after which the alert is resent if the condition is not cleared.
	AlertTimer string `json:"alertTimer"`

	// Organization pseudo domain
	PDomain string `json:"pdomain"`

	// Internal Company
	InternalCompany bool `json:"internalCompany"`

	// Primary technical contact type
	PrimaryTechnicalContactType string `json:"primaryTechnicalContactcontactType"`

	// Primary technical contact name
	PrimaryTechnicalContactName string `json:"primaryTechnicalContactName"`

	// Primary technical contact title
	PrimaryTechnicalContactTitle string `json:"primaryTechnicalContactTitle"`

	// Primary technical contact email address
	PrimaryTechnicalContactEmail string `json:"primaryTechnicalContactEmail"`

	// Primary technical contact phone number
	PrimaryTechnicalContactPhone string `json:"primaryTechnicalContactPhone"`

	// Primary technical contact alternate phone number
	PrimaryTechnicalContactAltPhone string `json:"primaryTechnicalContactAltPhone"`

	// Contains href to the most recent insights newsletter when available
	PrimaryTechnicalContactInsightsHref string `json:"primaryTechnicalContactInsightsHref"`

	// Secondary technical contact type
	SecondaryTechnicalContactType string `json:"secondaryTechnicalContactcontactType"`

	// Secondary technical contact name
	SecondaryTechnicalContactName string `json:"secondaryTechnicalContactName"`

	// Secondary technical contact title
	SecondaryTechnicalContactTitle string `json:"secondaryTechnicalContactTitle"`

	// Secondary technical contact email address
	SecondaryTechnicalContactEmail string `json:"secondaryTechnicalContactEmail"`

	// Secondary technical contact phone number
	SecondaryTechnicalContactPhone string `json:"secondaryTechnicalContactPhone"`

	// Secondary technical contact alternate phone number
	SecondaryTechnicalContactAltPhone string `json:"secondaryTechnicalContactAltPhone"`

	// Contains href to the most recent insights newsletter when available
	SecondaryTechnicalContactInsightsHref string `json:"secondaryTechnicalContactInsightsHref"`

	// Primary billing contact type
	PrimaryBillingContactType string `json:"primaryBillingContactcontactType"`

	// Primary billing contact name
	PrimaryBillingContactName string `json:"primaryBillingContactName"`

	// Primary billing contact title
	PrimaryBillingContactTitle string `json:"primaryBillingContactTitle"`

	// Primary billing contact email address
	PrimaryBillingContactEmail string `json:"primaryBillingContactEmail"`

	// Primary billing contact phone number
	PrimaryBillingContactPhone string `json:"primaryBillingContactPhone"`

	// Primary billing contact alternate phone number
	PrimaryBillingContactAltPhone string `json:"primaryBillingContactAltPhone"`

	// Contains href to the most recent insights newsletter when available
	PrimaryBillingContactInsightsHref string `json:"primaryBillingContactInsightsHref"`

	// Secondary billing contact type
	SecondaryBillingContactType string `json:"secondaryBillingContactcontactType"`

	// Secondary billing contact name
	SecondaryBillingContactName string `json:"secondaryBillingContactName"`

	// Secondary billing contact title
	SecondaryBillingContactTitle string `json:"secondaryBillingContactTitle"`

	// Secondary billing contact email address
	SecondaryBillingContactEmail string `json:"secondaryBillingContactEmail"`

	// Secondary billing contact phone number
	SecondaryBillingContactPhone string `json:"secondaryBillingContactPhone"`

	// Secondary billing contact alternate phone number
	SecondaryBillingContactAltPhone string `json:"secondaryBillingContactAltPhone"`

	// Contains href to the most recent insights newsletter when available
	SecondaryBillingContactInsightsHref string `json:"secondaryBillingContactInsightsHref"`

	// Primary business contact type
	PrimaryBusinessContactType string `json:"primaryBusinessContactcontactType"`

	// Primary business contact name
	PrimaryBusinessContactName string `json:"primaryBusinessContactName"`

	// Primary business contact title
	PrimaryBusinessContactTitle string `json:"primaryBusinessContactTitle"`

	// Primary business contact email address
	PrimaryBusinessContactEmail string `json:"primaryBusinessContactEmail"`

	// Primary business contact phone number
	PrimaryBusinessContactPhone string `json:"primaryBusinessContactPhone"`

	// Primary business contact alternate phone number
	PrimaryBusinessContactAltPhone string `json:"primaryBusinessContactAltPhone"`

	// Contains href to the most recent insights newsletter when available
	PrimaryBusinessContactInsightsHref string `json:"primaryBusinessContactInsightsHref"`

	// Secondary business contact type
	SecondaryBusinessContactType string `json:"secondaryBusinessContactcontactType"`

	// Secondary business contact name
	SecondaryBusinessContactName string `json:"secondaryBusinessContactName"`

	// Secondary business contact title
	SecondaryBusinessContactTitle string `json:"secondaryBusinessContactTitle"`

	// Secondary business contact email address
	SecondaryBusinessContactEmail string `json:"secondaryBusinessContactEmail"`

	// Secondary business contact phone number
	SecondaryBusinessContactPhone string `json:"secondaryBusinessContactPhone"`

	// Secondary business contact alternate phone number
	SecondaryBusinessContactAltPhone string `json:"secondaryBusinessContactAltPhone"`

	// Contains href to the most recent insights newsletter when available
	SecondaryBusinessContactInsightsHref string `json:"secondaryBusinessContactInsightsHref"`

	// Contains href to the most recent insights newsletter when available
	ExecInsightsHref string `json:"execInsightsHref"`

	// Enable the executive report
	LegacyInsightsReportWasEnabled bool `json:"legacyInsightsReportWasEnabled"`

	// Organization logo data. Base64 encoded. Must use uploadLogo API to update the logo.
	LogoBase64Data string `json:"logoBase64Data"`

	// Organization logo file MIME type
	LogoMimeType string `json:"logoMimeType"`

	// Cloud name
	CloudName string `json:"cloudName"`

	// Email portal is an external application
	ExternalEmailPortal bool `json:"externalEmailPortal"`

	// ZPA tenant ID
	ZpaTenantID int64 `json:"zpaTenantId"`

	// ZPA tenant cloud
	ZpaTenantCloud string `json:"zpaTenantCloud"`

	// Customer contact inherit flag
	CustomerContactInherit bool `json:"customerContactInherit"`
}

type OrganizationInfoLite struct {
	// Organization identifier
	OrgID int `json:"orgId"`

	// Name of the organization
	Name string `json:"name"`

	// Cloud Name
	CloudName string `json:"cloudName"`

	// Domain Names
	Domains []string `json:"domains"`

	// Primary Language. If not set, set to NONE.
	Language string `json:"language"`

	// Primary Time Zone. If not set, set to GMT.
	Timezone string `json:"timezone"`

	// Primary Time Zone. If not set, set to GMT.
	OrgDisabled bool `json:"orgDisabled"`
}

func GetSubscriptions(ctx context.Context, service *zscaler.Service) ([]Subscription, error) {
	var orgs []Subscription // Use a slice to match the array response
	err := service.Client.Read(ctx, subscriptionsEndpoint, &orgs)
	if err != nil {
		return nil, err
	}
	service.Client.GetLogger().Printf("[DEBUG] Returning subscription information from Get: %v", orgs)
	return orgs, nil
}

func GetOrgInformation(ctx context.Context, service *zscaler.Service) (*Organization, error) {
	var org Organization
	err := service.Client.Read(ctx, orgInformationEndpoint, &org)
	if err != nil {
		return nil, err
	}
	service.Client.GetLogger().Printf("[DEBUG] Returning organization information from Get: %v", org)
	return &org, nil
}

func GetOrgInformationLite(ctx context.Context, service *zscaler.Service) (*OrganizationInfoLite, error) {
	var org OrganizationInfoLite
	err := service.Client.Read(ctx, orgInformationLiteEndpoint, &org)
	if err != nil {
		return nil, err
	}
	service.Client.GetLogger().Printf("[DEBUG] Returning organization information from Get: %v", org)
	return &org, nil
}
