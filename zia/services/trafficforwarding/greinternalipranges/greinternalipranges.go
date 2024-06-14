package greinternalipranges

import (
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
)

const (
	greTunnelIPRangeEndpoint = "/greTunnels/availableInternalIpRanges"
)

type GREInternalIPRange struct {
	// Starting IP address in the range
	StartIPAddress string `json:"startIPAddress,omitempty"`

	// Ending IP address in the range
	EndIPAddress string `json:"endIPAddress,omitempty"`
}

func GetGREInternalIPRange(service *services.Service, count int) (*[]GREInternalIPRange, error) {
	var greInternalIPRanges []GREInternalIPRange
	err := service.Client.Read(fmt.Sprintf("%s?limit=%d", greTunnelIPRangeEndpoint, count), &greInternalIPRanges)
	if err != nil {
		return nil, err
	}
	if len(greInternalIPRanges) < count {
		return nil, fmt.Errorf("not enough internal IP range available, got %d internal IP range, required: %d", len(greInternalIPRanges), count)
	}
	service.Client.Logger.Printf("[DEBUG]Returning internal IP range: %s", greInternalIPRanges)
	return &greInternalIPRanges, nil
}
