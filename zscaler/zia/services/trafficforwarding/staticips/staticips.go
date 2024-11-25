package staticips

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	staticIPEndpoint = "/zia/api/v1/staticIP"
)

// Gets all provisioned static IP addresses.
type StaticIP struct {
	// The unique identifier for the static IP address
	ID int `json:"id,omitempty"`

	// The static IP address
	IpAddress string `json:"ipAddress"`

	// If not set, geographic coordinates and city are automatically determined from the IP address. Otherwise, the latitude and longitude coordinates must be provided.
	GeoOverride bool `json:"geoOverride"`

	// Required only if the geoOverride attribute is set. Latitude with 7 digit precision after decimal point, ranges between -90 and 90 degrees.
	Latitude float64 `json:"latitude,omitempty"`

	// Required only if the geoOverride attribute is set. Longitude with 7 digit precision after decimal point, ranges between -180 and 180 degrees.
	Longitude float64 `json:"longitude,omitempty"`

	// Indicates whether a non-RFC 1918 IP address is publicly routable. This attribute is ignored if there is no ZIA Private Service Edge associated to the organization.
	RoutableIP bool `json:"routableIP,omitempty"`

	City *City `json:"city,omitempty"`

	// When the static IP address was last modified
	LastModificationTime int `json:"lastModificationTime"`

	// Additional information about this static IP address
	Comment string `json:"comment,omitempty"`

	// SD-WAN Partner that manages the location. If a partner does not manage the location, this is set to Self.
	ManagedBy *ManagedBy `json:"managedBy,omitempty"` // Should probably move this to a common package. Used by multiple resources

	// Who modified the static IP address last
	LastModifiedBy *LastModifiedBy `json:"lastModifiedBy,omitempty"` // Should probably move this to a common package. Used by multiple resources
}

type ManagedBy struct {
	// Identifier that uniquely identifies an entity
	ID int `json:"id,omitempty"`

	// The configured name of the entity
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type LastModifiedBy struct {
	// Identifier that uniquely identifies an entity
	ID int `json:"id,omitempty"`

	// The configured name of the entity
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type City struct {
	// Identifier that uniquely identifies an entity
	ID int `json:"id,omitempty"`

	// The configured name of the entity
	Name string `json:"name,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, staticIpID int) (*StaticIP, error) {
	var staticIP StaticIP
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", staticIPEndpoint, staticIpID), &staticIP)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning static ip from Get: %d", staticIP.ID)
	return &staticIP, nil
}

func GetByIPAddress(ctx context.Context, service *zscaler.Service, address string) (*StaticIP, error) {
	var staticIPs []StaticIP
	err := common.ReadAllPages(ctx, service.Client, staticIPEndpoint, &staticIPs)
	if err != nil {
		return nil, err
	}
	for _, static := range staticIPs {
		if strings.EqualFold(static.IpAddress, address) {
			return &static, nil
		}
	}
	return nil, fmt.Errorf("no device group found with name: %s", address)
}

func Create(ctx context.Context, service *zscaler.Service, staticIpID *StaticIP) (*StaticIP, *http.Response, error) {
	resp, err := service.Client.Create(ctx, staticIPEndpoint, *staticIpID)
	if err != nil {
		return nil, nil, err
	}

	createdStaticIP, ok := resp.(*StaticIP)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a static ip pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning static ip from create: %d", createdStaticIP.ID)
	return createdStaticIP, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, staticIpID int, staticIP *StaticIP) (*StaticIP, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", staticIPEndpoint, staticIpID), *staticIP)
	if err != nil {
		return nil, nil, err
	}
	updatedStaticIP, _ := resp.(*StaticIP)

	service.Client.GetLogger().Printf("[DEBUG]returning static ip from update: %d", updatedStaticIP.ID)
	return updatedStaticIP, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, staticIpID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", staticIPEndpoint, staticIpID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]StaticIP, error) {
	var staticIPs []StaticIP
	err := common.ReadAllPages(ctx, service.Client, staticIPEndpoint, &staticIPs)
	return staticIPs, err
}
