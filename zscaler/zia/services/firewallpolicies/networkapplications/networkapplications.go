package networkapplications

import (
	"context"
	"fmt"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	networkApplicationsEndpoint = "/zia/api/v1/networkApplications"
)

type NetworkApplications struct {
	ID             string `json:"id"`
	ParentCategory string `json:"parentCategory,omitempty"`
	Description    string `json:"description,omitempty"`
	Deprecated     bool   `json:"deprecated"`
}

func GetNetworkApplication(ctx context.Context, service *zscaler.Service, id, locale string) (*NetworkApplications, error) {
	var networkApplications NetworkApplications
	url := fmt.Sprintf("%s/%s", networkApplicationsEndpoint, id)
	if locale != "" {
		url = fmt.Sprintf("%s?locale=%s", url, locale)
	}
	err := service.Client.Read(ctx, url, &networkApplications)
	if err != nil {
		return nil, err
	}
	return &networkApplications, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, nwApplicationName, locale string) (*NetworkApplications, error) {
	var networkApplications []NetworkApplications

	// Construct the URL with search and locale query parameters
	url := fmt.Sprintf("%s?search=%s&locale=%s", networkApplicationsEndpoint, url.QueryEscape(nwApplicationName), url.QueryEscape(locale))

	err := service.Client.Read(ctx, url, &networkApplications)
	if err != nil {
		return nil, err
	}

	// It's assumed that the API will return filtered results based on the search parameter.
	// Therefore, we should check if at least one result is returned.
	if len(networkApplications) > 0 {
		return &networkApplications[0], nil
	}

	return nil, fmt.Errorf("no network application found with name: %s", nwApplicationName)
}

func GetAll(ctx context.Context, service *zscaler.Service, locale string) ([]NetworkApplications, error) {
	var networkApplications []NetworkApplications
	endpoint := networkApplicationsEndpoint
	if locale != "" {
		// Properly escape the locale string and append it as a query parameter
		endpoint = fmt.Sprintf("%s?locale=%s", networkApplicationsEndpoint, url.QueryEscape(locale))
	}
	err := common.ReadAllPages(ctx, service.Client, endpoint, &networkApplications)
	return networkApplications, err
}
