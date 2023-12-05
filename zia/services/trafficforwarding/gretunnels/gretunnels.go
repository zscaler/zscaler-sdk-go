package gretunnels

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	greTunnelsEndpoint = "/greTunnels"
	// IpGreTunnelInfoEndpoint = "/orgProvisioning/ipGreTunnelInfo".
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
	WithinCountry bool `json:"withinCountry"`

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
func (service *Service) GetGreTunnels(greTunnelID int) (*GreTunnels, error) {
	var greTunnels GreTunnels
	err := service.Client.Read(fmt.Sprintf("%s/%d", greTunnelsEndpoint, greTunnelID), &greTunnels)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]returning gre tunnel from get: %d", greTunnels.ID)
	return &greTunnels, nil
}

// Gets specific provisioned GRE tunnel information by source IP address
func (service *Service) GetByIPAddress(sourceIP string) (*GreTunnels, error) {
	var sourceIPs []GreTunnels
	err := common.ReadAllPages(service.Client, greTunnelsEndpoint, &sourceIPs)
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
func (service *Service) CreateGreTunnels(greTunnelID *GreTunnels) (*GreTunnels, *http.Response, error) {
	resp, err := service.Client.Create(greTunnelsEndpoint, *greTunnelID)
	if err != nil {
		return nil, nil, err
	}

	createdGreTunnels, ok := resp.(*GreTunnels)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a gre tunnel pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning gre tunnels from create: %d", createdGreTunnels.ID)
	return createdGreTunnels, nil, nil
}

func (service *Service) UpdateGreTunnels(greTunnelID int, greTunnels *GreTunnels) (*GreTunnels, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", greTunnelsEndpoint, greTunnelID), *greTunnels)
	if err != nil {
		return nil, nil, err
	}
	updatedGreTunnels, _ := resp.(*GreTunnels)

	service.Client.Logger.Printf("[DEBUG]returning gre tunnels from update: %d", updatedGreTunnels.ID)
	return updatedGreTunnels, nil, nil
}

func (service *Service) DeleteGreTunnels(greTunnelID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", greTunnelsEndpoint, greTunnelID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (service *Service) GetAll() ([]GreTunnels, error) {
	var greTunnels []GreTunnels
	err := common.ReadAllPages(service.Client, greTunnelsEndpoint, &greTunnels)
	return greTunnels, err
}
