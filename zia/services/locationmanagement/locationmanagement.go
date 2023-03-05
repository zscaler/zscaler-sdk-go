package locationmanagement

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
)

const (
	locationsEndpoint               = "/locations"
	locationIPv6NAT64PrefixEndpoint = "/ipv6config/nat64prefix"
)

// Gets locations only, not sub-locations. When a location matches the given search parameter criteria only its parent location is included in the result set, not its sub-locations.
type Locations struct {
	// Location ID
	ID int `json:"id,omitempty"`

	// Location Name
	Name string `json:"name,omitempty"`

	// Parent Location ID. If this ID does not exist or is 0, it is implied that it is a parent location. Otherwise, it is a sub-location whose parent has this ID. x-applicableTo: SUB
	ParentID int `json:"parentId,omitempty"`

	// Upload bandwidth in kbps. The value 0 implies no Bandwidth Control enforcement
	UpBandwidth int `json:"upBandwidth,omitempty"`

	// Download bandwidth in kbps. The value 0 implies no Bandwidth Control enforcement
	DnBandwidth int `json:"dnBandwidth,omitempty"`

	// Country
	Country string `json:"country,omitempty"`

	// Timezone of the location. If not specified, it defaults to GMT.
	TZ string `json:"tz,omitempty"`

	// For locations: IP addresses of the egress points that are provisioned in the Zscaler Cloud. Each entry is a single IP address (e.g., 238.10.33.9).
	// For sub-locations: Egress, internal, or GRE tunnel IP addresses. Each entry is either a single IP address, CIDR (e.g., 10.10.33.0/24), or range (e.g., 10.10.33.1-10.10.33.10)).
	IPAddresses []string `json:"ipAddresses,omitempty"`

	// IP ports that are associated with the location
	Ports string `json:"ports,omitempty"`

	// VPN User Credentials that are associated with the location.
	VPNCredentials []VPNCredentials `json:"vpnCredentials,omitempty"`

	// Enable Basic Authentication at the location
	BasicAuthEnabled bool `json:"basicAuthEnabled"`

	// Enforce Authentication. Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled.
	AuthRequired bool `json:"authRequired"`

	// This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future.
	// Enable SSL Inspection. Set to true in order to apply your SSL Inspection policy to HTTPS traffic in the location and inspect HTTPS transactions for data leakage, malicious content, and viruses.
	SSLScanEnabled bool `json:"sslScanEnabled"`

	// This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future.
	// Enable Zscaler App SSL Setting. When set to true, the Zscaler App SSL Scan Setting takes effect, irrespective of the SSL policy that is configured for the location.
	ZappSSLScanEnabled bool `json:"zappSSLScanEnabled"`

	// Enable XFF Forwarding for a location. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
	// Note: For sub-locations, this attribute is a read-only field as the value is inherited from the parent location.
	XFFForwardEnabled bool `json:"xffForwardEnabled"`

	//Enable Surrogate IP. When set to true, users are mapped to internal device IP addresses
	SurrogateIP bool `json:"surrogateIP"`

	// Idle Time to Disassociation. The user mapping idle time (in minutes) is required if a Surrogate IP is enabled
	IdleTimeInMinutes int `json:"idleTimeInMinutes,omitempty"`

	// Display Time Unit. The time unit to display for IP Surrogate idle time to disassociation
	DisplayTimeUnit string `json:"displayTimeUnit,omitempty"`

	// Enforce Surrogate IP for Known Browsers. When set to true, IP Surrogate is enforced for all known browsers
	SurrogateIPEnforcedForKnownBrowsers bool `json:"surrogateIPEnforcedForKnownBrowsers"`

	// Refresh Time for re-validation of Surrogacy. The surrogate refresh time (in minutes) to re-validate the IP surrogates
	SurrogateRefreshTimeInMinutes int `json:"surrogateRefreshTimeInMinutes,omitempty"`

	// Display Refresh Time Unit. The time unit to display for refresh time for re-validation of surrogacy
	SurrogateRefreshTimeUnit string `json:"surrogateRefreshTimeUnit,omitempty"`

	// Enable Firewall. When set to true, Firewall is enabled for the location.
	OFWEnabled bool `json:"ofwEnabled"`

	// Enable IPS Control. When set to true, IPS Control is enabled for the location if Firewall is enabled.
	IPSControl bool `json:"ipsControl"`

	// Enable AUP. When set to true, AUP is enabled for the location
	AUPEnabled bool `json:"aupEnabled"`

	// Enable Caution. When set to true, a caution notifcation is enabled for the location
	CautionEnabled bool `json:"cautionEnabled"`

	// For First Time AUP Behavior, Block Internet Access. When set, all internet access (including non-HTTP traffic) is disabled until the user accepts the AUP.
	AUPBlockInternetUntilAccepted bool `json:"aupBlockInternetUntilAccepted"`

	// For First Time AUP Behavior, Force SSL Inspection. When set, Zscaler forces SSL Inspection in order to enforce AUP for HTTPS traffic.
	AUPForceSSLInspection bool `json:"aupForceSslInspection"`

	// Custom AUP Frequency. Refresh time (in days) to re-validate the AUP.
	AUPTimeoutInDays int `json:"aupTimeoutInDays,omitempty"`

	// Profile tag that specifies the location traffic type. If not specified, this tag defaults to "Unassigned".
	Profile string `json:"profile,omitempty"`

	// Additional notes or information regarding the location or sub-location. The description cannot exceed 1024 characters.
	Description string `json:"description,omitempty"`

	// If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv4 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other and it can be renamed, if required.
	OtherSubLocation bool `json:"otherSubLocation,omitempty"`

	// If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv6 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other6 and it can be renamed, if required. This field is applicable only if ipv6Enabled is set is true.
	Other6SubLocation bool `json:"other6SubLocation,omitempty"`

	// If set to true, IPv6 is enabled for the location and IPv6 traffic from the location can be forwarded to the Zscaler service to enforce security policies.
	IPv6Enabled bool `json:"ipv6Enabled,omitempty"`

	// (Optional) Name-ID pair of the NAT64 prefix configured as the DNS64 prefix for the location. If specified, the DNS64 prefix is used for the IP addresses that reside in this location. If not specified, a prefix is selected from the set of supported prefixes. This field is applicable only if ipv6Enabled is set is true.
	// Before you can configure a DNS64 prefix, you must send a GET request to /ipv6config/nat64prefix to retrieve the IDs of NAT64 prefixes, which can be configured as the DNS64 prefix.
	IPv6Dns64Prefix bool `json:"ipv6Dns64Prefix,omitempty"`
}

type Location struct {
	// Identifier that uniquely identifies an entity
	ID int `json:"id,omitempty"`

	// The configured name of the entity
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type ManagedBy struct {
	// Identifier that uniquely identifies an entity
	ID int `json:"id,omitempty"`

	// The configured name of the entity
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type VPNCredentials struct {
	// VPN credential id
	ID int `json:"id,omitempty"`

	// VPN authentication type (i.e., how the VPN credential is sent to the server). It is not modifiable after VpnCredential is created.
	// Note: Zscaler no longer supports adding a new XAUTH VPN credential, but existing entries can be edited or deleted using the respective endpoints.
	Type string `json:"type,omitempty"`

	// Fully Qualified Domain Name. Applicable only to UFQDN or XAUTH (or HOSTED_MOBILE_USERS) auth type.
	FQDN string `json:"fqdn,omitempty"`

	// Static IP address for VPN that is self-provisioned or provisioned by Zscaler. This is a required field for IP auth type and is not applicable to other auth types.
	// Note: If you want Zscaler to provision static IP addresses for your organization, contact Zscaler Support.
	IPAddress string `json:"ipAddress"`

	// Pre-shared key. This is a required field for UFQDN and IP auth type.
	PreSharedKey string `json:"preSharedKey,omitempty"`

	// Additional information about this VPN credential.
	Comments string `json:"comments,omitempty"`

	// Location that is associated to this VPN credential. Non-existence means not associated to any location.
	Location []Location `json:"location,omitempty"`

	// SD-WAN Partner that manages the location. If a partner does not manage the location, this is set to Self.
	ManagedBy []ManagedBy `json:"managedBy,omitempty"`
}

// Gets locations only, not sub-locations. When a location matches the given search parameter criteria only its parent location is included in the result set, not its sub-locations.
func (service *Service) GetLocation(locationID int) (*Locations, error) {
	var location Locations
	err := service.Client.Read(fmt.Sprintf("%s/%d", locationsEndpoint, locationID), &location)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning Location from Get: %d", location.ID)
	return &location, nil
}

func (service *Service) GetLocationByName(locationName string) (*Locations, error) {
	var locations []Locations
	// We are assuming this location name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(service.Client, locationsEndpoint, &locations)
	if err != nil {
		return nil, err
	}
	for _, location := range locations {
		if strings.EqualFold(location.Name, locationName) {
			return &location, nil
		}
	}
	return nil, fmt.Errorf("no location found with name: %s", locationName)
}

func (service *Service) Create(locations *Locations) (*Locations, error) {
	resp, err := service.Client.Create(locationsEndpoint, *locations)
	if err != nil {
		return nil, err
	}

	createdLocations, ok := resp.(*Locations)
	if !ok {
		return nil, errors.New("object returned from api was not a location pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning locations from create: %d", createdLocations.ID)
	return createdLocations, nil
}

func (service *Service) Update(locationID int, locations *Locations) (*Locations, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", locationsEndpoint, locationID), *locations)
	if err != nil {
		return nil, nil, err
	}
	updatedLocations, _ := resp.(*Locations)

	service.Client.Logger.Printf("[DEBUG]returning locations from Update: %d", updatedLocations.ID)
	return updatedLocations, nil, nil
}

func (service *Service) Delete(locationID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", locationsEndpoint, locationID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Gets a name and ID of locations.
func (service *Service) GetSublocations(sublocations string) (*Locations, error) {
	var subLocations Locations
	err := service.Client.Read(locationsEndpoint, "/"+"%s"+"/sublocations")
	if err != nil {
		return nil, err
	}

	return &subLocations, nil
}

func (service *Service) GetAll() ([]Locations, error) {
	var locations []Locations
	// We are assuming this location name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(service.Client, locationsEndpoint, &locations)
	return locations, err
}
