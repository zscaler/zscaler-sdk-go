package virtualipaddress

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/staticips"
)

const (
	vipsEndpoint               = "/zia/api/v1/vips"
	vipRecommendedListEndpoint = "/zia/api/v1/vips/recommendedList"
	staticIPEndpoint           = "/zia/api/v1/staticIP"
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
func GetZscalerVIPs(ctx context.Context, service *zscaler.Service, datacenter string) (*ZscalerVIPs, error) {
	var zscalerVips []ZscalerVIPs

	err := common.ReadAllPages(ctx, service.Client, vipsEndpoint, &zscalerVips)
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
func GetZSGREVirtualIPList(ctx context.Context, service *zscaler.Service, sourceIP string, count int) (*[]GREVirtualIPList, error) {
	var zscalerVips []GREVirtualIPList
	err := common.ReadAllPages(ctx, service.Client, fmt.Sprintf("%s?sourceIp=%s", vipRecommendedListEndpoint, sourceIP), &zscalerVips)
	if err != nil {
		return nil, err
	}
	if len(zscalerVips) < count {
		return nil, fmt.Errorf("not enough vips, got %d vips, required: %d", len(zscalerVips), count)
	}
	return &zscalerVips, nil
}

// Gets a paginated list of the virtual IP addresses (VIPs) available in the Zscaler cloud by sourceIP within country.
func GetPairZSGREVirtualIPsWithinCountry(ctx context.Context, service *zscaler.Service, sourceIP, countryCode string) (*[]GREVirtualIPList, error) {
	var zscalerVips []GREVirtualIPList
	err := common.ReadAllPages(ctx, service.Client, fmt.Sprintf("%s?sourceIp=%s&withinCountryOnly=true", vipRecommendedListEndpoint, sourceIP), &zscalerVips)
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

// Gets a paginated list of the virtual IP addresses (VIPs) available in the Zscaler cloud based on optional parameters.
func GetVIPRecommendedList(ctx context.Context, service *zscaler.Service, options ...func(*url.Values)) (*[]GREVirtualIPList, error) {
	queryParams := url.Values{}

	// Apply any optional parameters passed via options
	for _, option := range options {
		option(&queryParams)
	}

	// Default to withinCountryOnly if no withinCountryOnly flag is provided
	if queryParams.Get("withinCountryOnly") == "" {
		queryParams.Set("withinCountryOnly", "true")
	}

	// Construct the full endpoint with the query parameters
	endpoint := fmt.Sprintf("%s?%s", vipRecommendedListEndpoint, queryParams.Encode())

	var zscalerVips []GREVirtualIPList
	err := common.ReadAllPages(ctx, service.Client, endpoint, &zscalerVips)
	if err != nil {
		return nil, err
	}

	// If less than 2 VIPs are found, ensure at least 2 are returned
	if len(zscalerVips) < 2 {
		for _, vip := range zscalerVips {
			if len(zscalerVips) >= 2 {
				break
			}
			if !containsVIP(zscalerVips, vip) {
				zscalerVips = append(zscalerVips, vip)
			}
		}
	}

	if len(zscalerVips) < 2 {
		return nil, fmt.Errorf("not enough vips, got %d vips, required: %d", len(zscalerVips), 2)
	}

	return &zscalerVips, nil
}

// Optional parameters as functions to be passed to GetVIPRecommendedList
func WithRoutableIP(routableIP bool) func(*url.Values) {
	return func(v *url.Values) {
		v.Set("routableIP", strconv.FormatBool(routableIP))
	}
}

func WithWithinCountryOnly(withinCountryOnly bool) func(*url.Values) {
	return func(v *url.Values) {
		v.Set("withinCountryOnly", strconv.FormatBool(withinCountryOnly))
	}
}

func WithIncludePrivateServiceEdge(includePrivateServiceEdge bool) func(*url.Values) {
	return func(v *url.Values) {
		v.Set("includePrivateServiceEdge", strconv.FormatBool(includePrivateServiceEdge))
	}
}

func WithIncludeCurrentVips(includeCurrentVips bool) func(*url.Values) {
	return func(v *url.Values) {
		v.Set("includeCurrentVips", strconv.FormatBool(includeCurrentVips))
	}
}

func WithSourceIP(sourceIp string) func(*url.Values) {
	return func(v *url.Values) {
		v.Set("sourceIp", sourceIp)
	}
}

func WithLatitude(latitude float64) func(*url.Values) {
	return func(v *url.Values) {
		v.Set("latitude", fmt.Sprintf("%f", latitude))
	}
}

func WithLongitude(longitude float64) func(*url.Values) {
	return func(v *url.Values) {
		v.Set("longitude", fmt.Sprintf("%f", longitude))
	}
}

func WithSubcloud(subcloud string) func(*url.Values) {
	return func(v *url.Values) {
		v.Set("subcloud", subcloud)
	}
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

func GetAll(ctx context.Context, service *zscaler.Service, sourceIP string) ([]GREVirtualIPList, error) {
	var zscalerVips []GREVirtualIPList
	err := common.ReadAllPages(ctx, service.Client, vipRecommendedListEndpoint+"?sourceIp="+sourceIP, &zscalerVips)
	return zscalerVips, err
}

func getAllStaticIPs(ctx context.Context, service *zscaler.Service) ([]staticips.StaticIP, error) {
	var staticIPs []staticips.StaticIP
	err := common.ReadAllPages(ctx, service.Client, staticIPEndpoint, &staticIPs)
	return staticIPs, err
}

// GetAllSourceIPs  gets all vips for all static ips
func GetAllSourceIPs(ctx context.Context, service *zscaler.Service) ([]GREVirtualIPList, error) {
	var zscalerVips []GREVirtualIPList
	ips, err := getAllStaticIPs(ctx, service)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		list, err := GetAll(ctx, service, ip.IpAddress)
		if err != nil {
			continue
		}
		zscalerVips = append(zscalerVips, list...)
	}
	return zscalerVips, nil
}
