package gretunnels

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	greTunnelsEndpoint = "/zia/api/v1/greTunnels"
	// IpGreTunnelInfoEndpoint = "/zia/api/v1/orgProvisioning/ipGreTunnelInfo".
)

type GreTunnels struct {
	// Unique identifier of the static IP address that is associated to a GRE tunnel
	ID int `json:"id,omitempty"`

	// The source IP address of the GRE tunnel. This is typically a static IP address in the organization or SD-WAN. This IP address must be provisioned within the Zscaler service using the /staticIP endpoint.
	SourceIP string `json:"sourceIp,omitempty"`

	// The start of the internal IP address in /29 CIDR range
	InternalIpRange string `json:"internalIpRange,omitempty"`

	// When the GRE tunnel information was last modified
	LastModificationTime int `json:"lastModificationTime,omitempty"`

	// Restrict the data center virtual IP addresses (VIPs) only to those within the same country as the source IP address
	WithinCountry *bool `json:"withinCountry"`

	// Additional information about this GRE tunnel
	Comment string `json:"comment,omitempty"`

	// This is required to support the automated SD-WAN provisioning of GRE tunnels, when set to true gre_tun_ip and gre_tun_id are set to null
	IPUnnumbered bool `json:"ipUnnumbered"`

	// Restrict the data center virtual IP addresses (VIPs) only to those part of the subcloud
	SubCloud string `json:"subcloud,omitempty"`

	// SD-WAN Partner that manages the location. If a partner does not manage the location, this is set to Self.
	ManagedBy *ManagedBy `json:"managedBy,omitempty"` // Should probably move this to a common package. Used by multiple resources

	// Who modified the GRE tunnel information last
	LastModifiedBy *LastModifiedBy `json:"lastModifiedBy,omitempty"` // Should probably move this to a common package. Used by multiple resources

	// The primary destination data center and virtual IP address (VIP) of the GRE tunnel
	PrimaryDestVip *PrimaryDestVip `json:"primaryDestVip,omitempty"`

	// The secondary destination data center and virtual IP address (VIP) of the GRE tunnel
	SecondaryDestVip *SecondaryDestVip `json:"secondaryDestVip,omitempty"`
}

type PrimaryDestVip struct {
	// Unique identifer of the GRE virtual IP address (VIP)
	ID int `json:"id,omitempty"`

	// GRE cluster virtual IP address (VIP)
	VirtualIP string `json:"virtualIp,omitempty"`

	// Set to true if the virtual IP address (VIP) is a ZIA Private Service Edge
	PrivateServiceEdge bool `json:"privateServiceEdge"`

	// Data center information
	Datacenter string `json:"datacenter,omitempty"`

	// Latitude with 7 digit precision after decimal point, ranges between -90 and 90 degrees.
	Latitude float64 `json:"latitude,omitempty"`

	// Longitude with 7 digit precision after decimal point, ranges between -180 and 180 degrees.
	Longitude float64 `json:"longitude,omitempty"`

	// City information
	City string `json:"city,omitempty"`

	// Country Code information
	CountryCode string `json:"countryCode,omitempty"`

	// Region information
	Region string `json:"region,omitempty"`
}

type SecondaryDestVip struct {
	// Unique identifer of the GRE virtual IP address (VIP)
	ID int `json:"id,omitempty"`

	// GRE cluster virtual IP address (VIP)
	VirtualIP string `json:"virtualIp,omitempty"`

	// Set to true if the virtual IP address (VIP) is a ZIA Private Service Edge
	PrivateServiceEdge bool `json:"privateServiceEdge"`

	// Data center information
	Datacenter string `json:"datacenter,omitempty"`

	// Latitude with 7 digit precision after decimal point, ranges between -90 and 90 degrees.
	Latitude float64 `json:"latitude,omitempty"`

	// Longitude with 7 digit precision after decimal point, ranges between -180 and 180 degrees.
	Longitude float64 `json:"longitude,omitempty"`

	// City information
	City string `json:"city,omitempty"`

	// Country Code information
	CountryCode string `json:"countryCode,omitempty"`

	// Region information
	Region string `json:"region,omitempty"`
}

type ManagedBy struct {
	// Identifier that uniquely identifies an entity
	ID int `json:"id,omitempty"`

	// The configured name of the entity
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

type LastModifiedBy struct {
	// Identifier that uniquely identifies an entity
	ID int `json:"id,omitempty"`

	// The configured name of the entity
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions"`
}

// Gets specific provisioned GRE tunnel information.
func GetGreTunnels(ctx context.Context, service *zscaler.Service, greTunnelID int) (*GreTunnels, error) {
	var greTunnels GreTunnels
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", greTunnelsEndpoint, greTunnelID), &greTunnels)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]returning gre tunnel from get: %d", greTunnels.ID)
	return &greTunnels, nil
}

// Gets specific provisioned GRE tunnel information by source IP address
func GetByIPAddress(ctx context.Context, service *zscaler.Service, sourceIP string) (*GreTunnels, error) {
	var sourceIPs []GreTunnels
	err := common.ReadAllPages(ctx, service.Client, greTunnelsEndpoint, &sourceIPs)
	if err != nil {
		return nil, err
	}
	for _, source := range sourceIPs {
		if strings.EqualFold(source.SourceIP, sourceIP) {
			return &source, nil
		}
	}
	return nil, fmt.Errorf("no device group found with name: %s", sourceIP)
}

// Adds a GRE tunnel configuration.
func CreateGreTunnels(ctx context.Context, service *zscaler.Service, greTunnelID *GreTunnels) (*GreTunnels, *http.Response, error) {
	resp, err := service.Client.Create(ctx, greTunnelsEndpoint, *greTunnelID)
	if err != nil {
		return nil, nil, err
	}

	createdGreTunnels, ok := resp.(*GreTunnels)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a gre tunnel pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning gre tunnels from create: %d", createdGreTunnels.ID)
	return createdGreTunnels, nil, nil
}

func UpdateGreTunnels(ctx context.Context, service *zscaler.Service, greTunnelID int, greTunnels *GreTunnels) (*GreTunnels, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", greTunnelsEndpoint, greTunnelID), *greTunnels)
	if err != nil {
		return nil, nil, err
	}
	updatedGreTunnels, _ := resp.(*GreTunnels)

	service.Client.GetLogger().Printf("[DEBUG]returning gre tunnels from update: %d", updatedGreTunnels.ID)
	return updatedGreTunnels, nil, nil
}

func DeleteGreTunnels(ctx context.Context, service *zscaler.Service, greTunnelID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", greTunnelsEndpoint, greTunnelID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]GreTunnels, error) {
	var greTunnels []GreTunnels
	err := common.ReadAllPages(ctx, service.Client, greTunnelsEndpoint, &greTunnels)
	return greTunnels, err
}
