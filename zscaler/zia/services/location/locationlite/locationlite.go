package locationlite

import (
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
	// Unique identifier for the location group
	ID int `json:"id"`

	// Location name
	Name string `json:"name,omitempty"`

	// Parent Location ID. If this ID does not exist or is 0, it is implied that it is a parent location. Otherwise, it is a sub-location whose parent has this ID. x-applicableTo: SUB
	ParentID int `json:"parentId,omitempty"`

	// Indicates the location group was deleted
	TZ string `json:"tz,omitempty"`

	// Enable XFF Forwarding for a location. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
	// Note: For sub-locations, this attribute is a read-only field as the value is inherited from the parent location.
	XFFForwardEnabled bool `json:"xffForwardEnabled,omitempty"`

	// Enable AUP. When set to true, AUP is enabled for the location. To learn more, see About End User Notifications
	AUPEnabled bool `json:"aupEnabled"`

	// Enable Caution. When set to true, a caution notifcation is enabled for the location
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
	ZappSSLScanEnabled bool `json:"zappSSLScanEnabled"`

	// If set to true, IPv6 is enabled for the location and IPv6 traffic from the location can be forwarded to the Zscaler service to enforce security policies.
	IPv6Enabled bool `json:"ipv6Enabled,omitempty"`

	// If set to true, IPv6 is enabled for the location and IPv6 traffic from the location can be forwarded to the Zscaler service to enforce security policies.
	ECLocation bool `json:"ecLocation,omitempty"`

	// If set to true, IPv6 is enabled for the location and IPv6 traffic from the location can be forwarded to the Zscaler service to enforce security policies.
	KerberosAuth bool `json:"kerberosAuth,omitempty"`

	// If set to true, IPv6 is enabled for the location and IPv6 traffic from the location can be forwarded to the Zscaler service to enforce security policies.
	DigestAuthEnabled bool `json:"digestAuthEnabled,omitempty"`
}

func GetLocationLiteID(service *zscaler.Service, locationID int) (*LocationLite, error) {
	var locationLite LocationLite
	err := service.Client.Read(fmt.Sprintf("%s/%d", locationLiteEndpoint, locationID), &locationLite)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]returning location group from Get: %d", locationLite.ID)
	return &locationLite, nil
}

func GetLocationLiteByName(service *zscaler.Service, locationLiteName string) (*LocationLite, error) {
	var locationsLite []LocationLite
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?name=%s", locationLiteEndpoint, url.QueryEscape(locationLiteName)), &locationsLite)
	if err != nil {
		return nil, err
	}
	for _, locationLite := range locationsLite {
		if strings.EqualFold(locationLite.Name, locationLiteName) {
			return &locationLite, nil
		}
	}
	return nil, fmt.Errorf("no location found with name: %s", locationLiteName)
}

func GetAll(service *zscaler.Service) ([]LocationLite, error) {
	var locations []LocationLite
	err := common.ReadAllPages(service.Client, locationLiteEndpoint, &locations)
	return locations, err
}
