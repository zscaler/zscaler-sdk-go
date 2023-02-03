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
	ID                                  int              `json:"id,omitempty"`
	Name                                string           `json:"name,omitempty"`
	ParentID                            int              `json:"parentId,omitempty"`
	UpBandwidth                         int              `json:"upBandwidth,omitempty"`
	DnBandwidth                         int              `json:"dnBandwidth,omitempty"`
	Country                             string           `json:"country,omitempty"`
	TZ                                  string           `json:"tz,omitempty"`
	IPAddresses                         []string         `json:"ipAddresses,omitempty"`
	Ports                               string           `json:"ports,omitempty"`
	VPNCredentials                      []VPNCredentials `json:"vpnCredentials,omitempty"`
	AuthRequired                        bool             `json:"authRequired"`
	SSLScanEnabled                      bool             `json:"sslScanEnabled"`
	ZappSSLScanEnabled                  bool             `json:"zappSSLScanEnabled"`
	XFFForwardEnabled                   bool             `json:"xffForwardEnabled"`
	SurrogateIP                         bool             `json:"surrogateIP"`
	IdleTimeInMinutes                   int              `json:"idleTimeInMinutes,omitempty"`
	DisplayTimeUnit                     string           `json:"displayTimeUnit,omitempty"`
	SurrogateIPEnforcedForKnownBrowsers bool             `json:"surrogateIPEnforcedForKnownBrowsers"`
	SurrogateRefreshTimeInMinutes       int              `json:"surrogateRefreshTimeInMinutes,omitempty"`
	SurrogateRefreshTimeUnit            string           `json:"surrogateRefreshTimeUnit,omitempty"`
	OFWEnabled                          bool             `json:"ofwEnabled"`
	IPSControl                          bool             `json:"ipsControl"`
	AUPEnabled                          bool             `json:"aupEnabled"`
	CautionEnabled                      bool             `json:"cautionEnabled"`
	AUPBlockInternetUntilAccepted       bool             `json:"aupBlockInternetUntilAccepted"`
	AUPForceSSLInspection               bool             `json:"aupForceSslInspection"`
	AUPTimeoutInDays                    int              `json:"aupTimeoutInDays,omitempty"`
	Profile                             string           `json:"profile,omitempty"`
	Description                         string           `json:"description,omitempty"`
	OtherSubLocation                    bool             `json:"otherSubLocation,omitempty"`
	Other6SubLocation                   bool             `json:"other6SubLocation,omitempty"`
	IPv6Enabled                         bool             `json:"ipv6Enabled,omitempty"`
	IPv6Dns64Prefix                     bool             `json:"ipv6Dns64Prefix,omitempty"`
}

type Location struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type ManagedBy struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type VPNCredentials struct {
	ID           int         `json:"id,omitempty"`
	Type         string      `json:"type,omitempty"`
	FQDN         string      `json:"fqdn,omitempty"`
	IPAddress    string      `json:"ipAddress"`
	PreSharedKey string      `json:"preSharedKey,omitempty"`
	Comments     string      `json:"comments,omitempty"`
	Location     []Location  `json:"location,omitempty"`
	ManagedBy    []ManagedBy `json:"managedBy,omitempty"`
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
