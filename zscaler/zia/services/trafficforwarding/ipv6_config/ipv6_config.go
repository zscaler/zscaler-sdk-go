package ipv6_config

import (
	"context"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	ipv6configEndpoint = "/zia/api/v1/ipv6config"
)

type IPv6Config struct {
	IpV6Enabled bool               `json:"ipV6Enabled,omitempty"`
	NatPrefixes []IPv6ConfigPrefix `json:"natPrefixes,omitempty"`
	DnsPrefix   string             `json:"dnsPrefix,omitempty"`
}

type IPv6ConfigPrefix struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	PrefixMask  string `json:"prefixMask,omitempty"`
	DnsPrefix   bool   `json:"dnsPrefix,omitempty"`
	NonEditable bool   `json:"nonEditable,omitempty"`
}

func GetIPv6Config(ctx context.Context, service *zscaler.Service) (*IPv6Config, error) {
	var config IPv6Config
	err := service.Client.Read(ctx, ipv6configEndpoint, &config)
	return &config, err
}

func GetDns64Prefix(ctx context.Context, service *zscaler.Service, search ...string) ([]IPv6ConfigPrefix, error) {
	var prefix []IPv6ConfigPrefix
	endpoint := ipv6configEndpoint + "/dns64prefix"

	if len(search) > 0 && strings.TrimSpace(search[0]) != "" {
		endpoint += "?search=" + url.QueryEscape(search[0])
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &prefix)
	return prefix, err
}

func GetNat64Prefix(ctx context.Context, service *zscaler.Service, search ...string) ([]IPv6ConfigPrefix, error) {
	var prefix []IPv6ConfigPrefix
	endpoint := ipv6configEndpoint + "/nat64prefix"

	if len(search) > 0 && strings.TrimSpace(search[0]) != "" {
		endpoint += "?search=" + url.QueryEscape(search[0])
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &prefix)
	return prefix, err
}
