package virtualipaddress

import (
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
)

const (
	vipsEndpoint               = "/vips"
	vipRecommendedListEndpoint = "/vips/recommendedList"
)

type ZscalerVIPs struct {
	CloudName     string   `json:"cloudName"`
	Region        string   `json:"region"`
	City          string   `json:"city"`
	DataCenter    string   `json:"dataCenter"`
	Location      string   `json:"location"`
	VPNIPs        []string `json:"vpnIps"`
	VPNDomainName string   `json:"vpnDomainName"`
	GREIPs        []string `json:"greIps"`
	GREDomainName string   `json:"greDomainName"`
	PACIPs        []string `json:"pacIps"`
	PACDomainName string   `json:"pacDomainName"`
}

type GREVirtualIPList struct {
	// Unique identifer of the GRE virtual IP address (VIP)
	ID int `json:"id"`

	// GRE cluster virtual IP address (VIP)
	VirtualIp string `json:"virtualIp,omitempty"`

	// Set to true if the virtual IP address (VIP) is a ZIA Private Service Edge
	PrivateServiceEdge bool `json:"privateServiceEdge,omitempty"`

	// Data center information
	DataCenter string `json:"dataCenter,omitempty"`

	// Country code information
	CountryCode string `json:"countryCode,omitempty"`

	City      string  `json:"city,omitempty"`
	Region    string  `json:"region,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Gets a paginated list of the virtual IP addresses (VIPs) available in the Zscaler cloud, including region and data center information. By default, the request gets all public VIPs in the cloud, but you can also include private or all VIPs in the request, if necessary.
func (service *Service) GetZscalerVIPs(datacenter string) (*ZscalerVIPs, error) {
	var zscalerVips []ZscalerVIPs

	err := common.ReadAllPages(service.Client, vipsEndpoint, &zscalerVips)
	if err != nil {
		return nil, err
	}
	for _, vips := range zscalerVips {
		if strings.EqualFold(vips.DataCenter, datacenter) {
			return &vips, nil
		}
	}
	return nil, fmt.Errorf("no datacenter found with name: %s", datacenter)
}

// Gets a paginated list of the virtual IP addresses (VIPs) available in the Zscaler cloud by sourceIP.
func (service *Service) GetZSGREVirtualIPList(sourceIP string, count int) (*[]GREVirtualIPList, error) {
	var zscalerVips []GREVirtualIPList
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?sourceIp=%s", vipRecommendedListEndpoint, sourceIP), &zscalerVips)
	if err != nil {
		return nil, err
	}
	if len(zscalerVips) < count {
		return nil, fmt.Errorf("not enough vips, got %d vips, required: %d", len(zscalerVips), count)
	}
	return &zscalerVips, nil
}

// Gets a paginated list of the virtual IP addresses (VIPs) available in the Zscaler cloud by sourceIP within country.
func (service *Service) GetPairZSGREVirtualIPsWithinCountry(sourceIP, countryCode string) (*[]GREVirtualIPList, error) {
	var zscalerVips []GREVirtualIPList
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?sourceIp=%s&withinCountryOnly=true", vipRecommendedListEndpoint, sourceIP), &zscalerVips)
	if err != nil {
		return nil, err
	}
	var pairVips []GREVirtualIPList
	for _, vip := range zscalerVips {
		if strings.EqualFold(vip.CountryCode, countryCode) {
			pairVips = append(pairVips, vip)
		}
	}
	// If not enough VIPs in the specified country, add any VIPs until there are at least two.
	if len(pairVips) < 2 {
		for _, vip := range zscalerVips {
			if len(pairVips) >= 2 {
				break
			}
			if !containsVIP(pairVips, vip) {
				pairVips = append(pairVips, vip)
			}
		}
	}
	if len(pairVips) < 2 {
		return nil, fmt.Errorf("not enough vips, got %d vips, required: %d", len(pairVips), 2)
	}
	return &pairVips, nil
}

// Helper function to check if a VIP is already in the list
func containsVIP(vips []GREVirtualIPList, vip GREVirtualIPList) bool {
	for _, v := range vips {
		if v.ID == vip.ID {
			return true
		}
	}
	return false
}

func (service *Service) GetAll(sourceIP string) ([]GREVirtualIPList, error) {
	var zscalerVips []GREVirtualIPList
	err := common.ReadAllPages(service.Client, vipRecommendedListEndpoint+"?sourceIp="+sourceIP, &zscalerVips)
	return zscalerVips, err
}

func (service *Service) getAllStaticIPs() ([]staticips.StaticIP, error) {
	var staticIPs []staticips.StaticIP
	err := common.ReadAllPages(service.Client, "/staticIP", &staticIPs)
	return staticIPs, err
}

// GetAllSourceIPs  gets all vips for all static ips
func (service *Service) GetAllSourceIPs() ([]GREVirtualIPList, error) {
	var zscalerVips []GREVirtualIPList
	ips, err := service.getAllStaticIPs()
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		list, err := service.GetAll(ip.IpAddress)
		if err != nil {
			continue
		}
		zscalerVips = append(zscalerVips, list...)
	}
	return zscalerVips, nil
}
