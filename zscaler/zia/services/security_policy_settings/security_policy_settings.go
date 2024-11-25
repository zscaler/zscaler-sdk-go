package security_policy_settings

import (
	"context"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	securityEndpoint         = "/zia/api/v1/security"
	securityAdvancedEndpoint = "/zia/api/v1/security/advanced"
)

// TODO: because there isn't an endpoint to get all Urls, we need to have all action types here.
var AddRemoveURLFromList []string = []string{
	"ADD_TO_LIST",
	"REMOVE_FROM_LIST",
}

type ListUrls struct {
	// Allowlist URLs whose contents will not be scanned. Allows up to 255 URLs. There may be trusted websites the content of which might be blocked due to anti-virus, anti-spyware, or anti-malware policies. Enter the URLs of sites you do not want scanned. The service allows users to download content from these URLs without inspecting the traffic. The allowlist applies to the Malware Protection, Advanced Threats Protection, and Sandbox policies.
	White []string `json:"whitelistUrls,omitempty"`

	// URLs on the denylist for your organization. Allow up to 25000 URLs.
	Black []string `json:"blacklistUrls,omitempty"`
}

func GetListUrls(ctx context.Context, service *zscaler.Service) (*ListUrls, error) {
	whitelist, err := GetWhiteListUrls(ctx, service)
	if err != nil {
		return nil, err
	}
	blacklist, err := GetBlackListUrls(ctx, service)
	if err != nil {
		return nil, err
	}
	return &ListUrls{
		White: whitelist.White,
		Black: blacklist.Black,
	}, nil
}

func UpdateListUrls(ctx context.Context, service *zscaler.Service, listUrls ListUrls) (*ListUrls, error) {
	whitelist, err := UpdateWhiteListUrls(ctx, service, ListUrls{White: listUrls.White})
	if err != nil {
		return nil, err
	}
	blacklist, err := UpdateBlackListUrls(ctx, service, ListUrls{Black: listUrls.Black})
	if err != nil {
		return nil, err
	}
	return &ListUrls{
		White: whitelist.White,
		Black: blacklist.Black,
	}, nil
}

func UpdateWhiteListUrls(ctx context.Context, service *zscaler.Service, list ListUrls) (*ListUrls, error) {
	_, err := service.Client.UpdateWithPut(ctx, securityEndpoint, list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func UpdateBlackListUrls(ctx context.Context, service *zscaler.Service, list ListUrls) (*ListUrls, error) {
	_, err := service.Client.UpdateWithPut(ctx, securityAdvancedEndpoint, list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func GetWhiteListUrls(ctx context.Context, service *zscaler.Service) (*ListUrls, error) {
	var whitelist ListUrls
	err := service.Client.Read(ctx, securityEndpoint, &whitelist)
	if err != nil {
		return nil, err
	}
	return &whitelist, nil
}

func GetBlackListUrls(ctx context.Context, service *zscaler.Service) (*ListUrls, error) {
	var blacklist ListUrls
	err := service.Client.Read(ctx, securityAdvancedEndpoint, &blacklist)
	if err != nil {
		return nil, err
	}
	return &blacklist, nil
}
