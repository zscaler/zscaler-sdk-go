package locationlite

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	locationLiteEndpoint = "/zia/api/v1/locations/lite"
)

type LocationLite struct {
	// Unique identifier for the location
	ID int `json:"id"`

	// Location name
	Name string `json:"name,omitempty"`

	// Parent Location ID. If this ID does not exist or is 0, it is implied that it is a parent location. Otherwise, it is a sub-location whose parent has this ID. x-applicableTo: SUB
	ParentID int `json:"parentId,omitempty"`

	// Timezone of the location. If not specified, it defaults to GMT.
	TZ string `json:"tz,omitempty"`

	// Enable XFF Forwarding for a location. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
	// Note: For sub-locations, this attribute is a read-only field as the value is inherited from the parent location.
	XFFForwardEnabled bool `json:"xffForwardEnabled,omitempty"`

	// Enable AUP. When set to true, AUP is enabled for the location. To learn more, see About End User Notifications
	AUPEnabled bool `json:"aupEnabled"`

	// Enable Caution. When set to true, a caution notification is enabled for the location
	CautionEnabled bool `json:"cautionEnabled"`

	// For First Time AUP Behavior, Block Internet Access. When set, all internet access (including non-HTTP traffic) is disabled until the user accepts the AUP.
	AUPBlockInternetUntilAccepted bool `json:"aupBlockInternetUntilAccepted"`

	// For First Time AUP Behavior, Force SSL Inspection. When set, Zscaler forces SSL Inspection in order to enforce AUP for HTTPS traffic.
	AUPForceSSLInspection bool `json:"aupForceSslInspection"`

	// Enable Surrogate IP. When set to true, users are mapped to internal device IP addresses
	SurrogateIP bool `json:"surrogateIP"`

	// Enforce Surrogate IP for Known Browsers. When set to true, IP Surrogate is enforced for all known browsers
	SurrogateIPEnforcedForKnownBrowsers bool `json:"surrogateIPEnforcedForKnownBrowsers"`

	// If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv4 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other and it can be renamed, if required.
	OtherSubLocation bool `json:"otherSubLocation,omitempty"`

	// If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv6 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other6 and it can be renamed, if required. This field is applicable only if ipv6Enabled is set is true.
	Other6SubLocation bool `json:"other6SubLocation,omitempty"`

	// Enable Firewall. When set to true, Firewall is enabled for the location.
	OFWEnabled bool `json:"ofwEnabled"`

	// Enable IPS Control. When set to true, IPS Control is enabled for the location if Firewall is enabled.
	IPSControl bool `json:"ipsControl"`

	// This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future.
	// Enable Zscaler App SSL Setting. When set to true, the Zscaler App SSL Scan Setting takes effect, irrespective of the SSL policy that is configured for the location.
	ZappSSLScanEnabled bool `json:"zappSslScanEnabled"`

	// If set to true, IPv6 is enabled for the location and IPv6 traffic from the location can be forwarded to the Zscaler service to enforce security policies.
	IPv6Enabled bool `json:"ipv6Enabled,omitempty"`

	// Indicates whether defining scopes is allowed for this sublocation. Sublocation scopes are available only for the Workload traffic type sublocations whose parent locations are associated with Amazon Web Services (AWS) Cloud Connector groups.
	SubLocScopeEnabled bool `json:"subLocScopeEnabled,omitempty"`

	// Defines a scope for the sublocation from the available types to segregate workload traffic from a single sublocation to apply different Cloud Connector and ZIA security policies. This field is only available for the Workload traffic type sublocations whose parent locations are associated with Amazon Web Services (AWS) Cloud Connector groups.
	SubLocScope string `json:"subLocScope,omitempty"`

	// Specifies values for the selected sublocation scope type
	SubLocScopeValues []string `json:"subLocScopeValues,omitempty"`

	// Specifies values for the selected sublocation scope type
	SubLocAccIDs []string `json:"subLocAccIds,omitempty"`
}

func GetLocationLiteID(ctx context.Context, service *zscaler.Service, locationID int) (*LocationLite, error) {
	var locationLite LocationLite
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", locationLiteEndpoint, locationID), &locationLite)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]returning location group from Get: %d", locationLite.ID)
	return &locationLite, nil
}

// GetAllFilterOptions represents optional filter parameters for GetAll
type GetAllFilterOptions struct {
	// If set to true sub-locations are included in the response otherwise they are excluded
	IncludeSubLocations *bool
	// If set to true locations with sub locations are included in the response, otherwise only locations without sub-locations are included
	IncludeParentLocations *bool
	// This parameter was deprecated. Filter based on whether the Enable SSL Scanning setting is enabled or disabled for a location.
	SslScanEnabled *bool
	// The search string used to partially match against a location's name and port attributes.
	Search *string
	// If set to true, the city field (containing IoT-enabled location IDs, names, latitudes, and longitudes) and the iotDiscoveryEnabled filter are included in the response. Otherwise, they are not included.
	EnableIOT *bool
}

// GetAll retrieves all location lite entries with optional filters.
// The API supports a maximum page size of 1000.
func GetAll(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) ([]LocationLite, error) {
	var locations []LocationLite
	endpoint := locationLiteEndpoint

	// Build query parameters
	queryParams := url.Values{}
	if opts != nil {
		if opts.IncludeSubLocations != nil {
			queryParams.Add("includeSubLocations", fmt.Sprintf("%t", *opts.IncludeSubLocations))
		}
		if opts.IncludeParentLocations != nil {
			queryParams.Add("includeParentLocations", fmt.Sprintf("%t", *opts.IncludeParentLocations))
		}
		if opts.SslScanEnabled != nil {
			queryParams.Add("sslScanEnabled", fmt.Sprintf("%t", *opts.SslScanEnabled))
		}
		if opts.Search != nil {
			queryParams.Add("search", *opts.Search)
		}
		if opts.EnableIOT != nil {
			queryParams.Add("enableIOT", fmt.Sprintf("%t", *opts.EnableIOT))
		}
	}

	// Build base endpoint with query parameters
	baseQuery := queryParams.Encode()
	if baseQuery != "" {
		endpoint += "?" + baseQuery
	}

	// Use common.ReadAllPages with default page size 1000 (API maximum)
	err := common.ReadAllPages(ctx, service.Client, endpoint, &locations)
	return locations, err
}

func GetLocationLiteByName(ctx context.Context, service *zscaler.Service, locationLiteName string) (*LocationLite, error) {
	// Use GetAll with search filter to leverage API filtering
	opts := &GetAllFilterOptions{
		Search: &locationLiteName,
	}
	locationsLite, err := GetAll(ctx, service, opts)
	if err != nil {
		return nil, err
	}
	// API may do partial matching, so verify exact match (case-insensitive)
	for _, locationLite := range locationsLite {
		if strings.EqualFold(locationLite.Name, locationLiteName) {
			return &locationLite, nil
		}
	}
	return nil, fmt.Errorf("no location found with name: %s", locationLiteName)
}
