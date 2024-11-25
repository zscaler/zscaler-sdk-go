package greinternalipranges

import (
	"context"
	"fmt"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	greTunnelIPRangeEndpoint = "/zia/api/v1/greTunnels/availableInternalIpRanges"
)

type GREInternalIPRange struct {
	// Starting IP address in the range
	StartIPAddress string `json:"startIPAddress,omitempty"`

	// Ending IP address in the range
	EndIPAddress string `json:"endIPAddress,omitempty"`
}

func GetGREInternalIPRange(ctx context.Context, service *zscaler.Service, count int) (*[]GREInternalIPRange, error) {
	var greInternalIPRanges []GREInternalIPRange
	err := service.Client.Read(ctx, fmt.Sprintf("%s?limit=%d", greTunnelIPRangeEndpoint, count), &greInternalIPRanges)
	if err != nil {
		return nil, err
	}
	if len(greInternalIPRanges) < count {
		return nil, fmt.Errorf("not enough internal IP range available, got %d internal IP range, required: %d", len(greInternalIPRanges), count)
	}
	service.Client.GetLogger().Printf("[DEBUG]Returning internal IP range: %s", greInternalIPRanges)
	return &greInternalIPRanges, nil
}
