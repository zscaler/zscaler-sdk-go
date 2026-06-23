package votiro_cdr

import (
	"context"
	"net/http"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	votiroCdrEndpoint            = "/zia/api/v1/isolationVotiroCdr"
	votiroCdrTokenConfigEndpoint = votiroCdrEndpoint + "/tokenConfig"
)

type VotiroCDRPolicies struct {
	ID   int `json:"id,omitempty"`
	Name int `json:"name,omitempty"`
}

func GetCDRPolicies(ctx context.Context, service *zscaler.Service) ([]VotiroCDRPolicies, error) {
	var cdrPolicies []VotiroCDRPolicies
	err := service.Client.Read(ctx, votiroCdrEndpoint, &cdrPolicies)
	return cdrPolicies, err
}

type VotiroCDRTokenConfig struct {
	Token       string `json:"token,omitempty"`
	TokenEnding string `json:"tokenEnding,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	TokenExp    int    `json:"tokenExp,omitempty"`
}

// UpdateCDRTokenConfig configures or updates the Votiro CDR authentication
// credentials.
//
// The action query parameter selects which field is updated: when set to "token"
// or "hostname", only the respective field in the VotiroCdrTokenConfig is
// updated. The token configuration details are supplied in the request body.
func UpdateCDRTokenConfig(ctx context.Context, service *zscaler.Service, action string, tokenConfig *VotiroCDRTokenConfig) (*VotiroCDRTokenConfig, *http.Response, error) {
	endpoint := votiroCdrTokenConfigEndpoint
	if action != "" {
		queryParams := url.Values{}
		queryParams.Set("action", action)
		endpoint += "?" + queryParams.Encode()
	}

	resp, err := service.Client.UpdateWithPut(ctx, endpoint, *tokenConfig)
	if err != nil {
		return nil, nil, err
	}
	updatedTokenConfig, _ := resp.(*VotiroCDRTokenConfig)

	service.Client.GetLogger().Printf("[DEBUG] returning updated votiro cdr token config from update")
	return updatedTokenConfig, nil, nil
}

func DeleteCDRTokenConfig(ctx context.Context, service *zscaler.Service) (*http.Response, error) {
	err := service.Client.Delete(ctx, votiroCdrTokenConfigEndpoint)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetCDRTokenConfig(ctx context.Context, service *zscaler.Service) (*VotiroCDRTokenConfig, error) {
	var cdrTokenConfig VotiroCDRTokenConfig
	err := service.Client.Read(ctx, votiroCdrTokenConfigEndpoint, &cdrTokenConfig)
	if err != nil {
		return nil, err
	}
	return &cdrTokenConfig, nil
}
