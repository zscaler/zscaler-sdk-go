package gretunnelinfo

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	ipGreTunnelInfoEndpoint = "/zia/api/v1/orgProvisioning/ipGreTunnelInfo"
)

type GRETunnelInfo struct {
	TunID             int    `json:"tunID,omitempty"`
	IPaddress         string `json:"ipAddress,omitempty"`
	GREEnabled        bool   `json:"greEnabled,omitempty"`
	GREtunnelIP       string `json:"greTunnelIP,omitempty"`
	PrimaryGW         string `json:"primaryGW,omitempty"`
	SecondaryGW       string `json:"secondaryGW,omitempty"`
	GRERangePrimary   string `json:"greRangePrimary,omitempty"`
	GRERangeSecondary string `json:"greRangeSecondary,omitempty"`
}

// Gets a list of IP addresses with GRE tunnel details.
func GetGRETunnelInfo(ctx context.Context, service *zscaler.Service, ipAddress string) (*GRETunnelInfo, error) {
	var greTunnelInfo []GRETunnelInfo
	err := service.Client.Read(ctx, fmt.Sprintf("%s?ipAddress=%s", ipGreTunnelInfoEndpoint, url.QueryEscape(ipAddress)), &greTunnelInfo)
	if err != nil {
		return nil, err
	}
	for _, greIP := range greTunnelInfo {
		if strings.EqualFold(greIP.IPaddress, ipAddress) {
			return &greIP, nil
		}
	}
	return nil, fmt.Errorf("no information for gre tunnel ip address: %s", ipAddress)
}
