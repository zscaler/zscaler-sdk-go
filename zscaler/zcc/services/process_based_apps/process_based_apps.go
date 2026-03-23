package process_based_apps

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	processBasedAppsEndpoint = "/zcc/papi/public/v1/process-based-apps"
)

type ProcessBasedAppsResponse struct {
	TotalCount    int               `json:"totalCount"`
	AppIdentities []ProcessBasedApp `json:"appIdentities"`
}

type ProcessBasedApp struct {
	ID                 int      `json:"id,omitempty"`
	AppName            string   `json:"appName,omitempty"`
	FileNames          []string `json:"fileNames,omitempty"`
	FilePaths          []string `json:"filePaths,omitempty"`
	MatchingCriteria   int      `json:"matchingCriteria,omitempty"`
	SignaturePayload   string   `json:"signaturePayload,omitempty"`
	CertificatePayload string   `json:"certificatePayload,omitempty"`
	CreatedBy          string   `json:"createdBy,omitempty"`
	EditedBy           string   `json:"editedBy,omitempty"`
	EditedTimestamp    string   `json:"editedTimestamp,omitempty"`
}

func GetProcessBasedApps(ctx context.Context, service *zscaler.Service, search string, page, pageSize *int) (*ProcessBasedAppsResponse, *http.Response, error) {
	params := common.QueryParams{
		Search: search,
	}
	if page != nil {
		params.Page = *page
	}
	if pageSize != nil {
		params.PageSize = *pageSize
	}

	var response ProcessBasedAppsResponse
	resp, err := service.Client.NewZccRequestDo(ctx, "GET", processBasedAppsEndpoint, params, nil, &response)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve process-based apps: %w", err)
	}
	return &response, resp, nil
}

func GetByAppID(ctx context.Context, service *zscaler.Service, appID string) (*ProcessBasedApp, *http.Response, error) {
	if appID == "" {
		return nil, nil, fmt.Errorf("appId is required")
	}
	endpoint := fmt.Sprintf("%s/%s", processBasedAppsEndpoint, appID)

	var app ProcessBasedApp
	resp, err := service.Client.NewZccRequestDo(ctx, "GET", endpoint, nil, nil, &app)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve process-based app %s: %w", appID, err)
	}
	return &app, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, name string) (*ProcessBasedApp, *http.Response, error) {
	pageSize := 1000
	page := 1

	for {
		res, resp, err := GetProcessBasedApps(ctx, service, "", &page, &pageSize)
		if err != nil {
			return nil, resp, err
		}
		for _, a := range res.AppIdentities {
			if strings.EqualFold(a.AppName, name) {
				return &a, resp, nil
			}
		}
		if len(res.AppIdentities) < pageSize {
			break
		}
		page++
	}
	return nil, nil, fmt.Errorf("no process-based app found with name: %s", name)
}
