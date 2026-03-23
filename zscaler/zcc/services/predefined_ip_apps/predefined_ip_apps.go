package predefined_ip_apps

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	predefinedIpAppsEndpoint = "/zcc/papi/public/v1/predefined-ip-based-apps"
)

type PredefinedIPAppsResponse struct {
	TotalCount          int              `json:"totalCount"`
	AppServiceContracts []PredefinedIPApp `json:"appServiceContracts"`
}

type PredefinedIPApp struct {
	ID              int           `json:"id,omitempty"`
	AppVersion      int           `json:"appVersion,omitempty"`
	AppSvcId        int           `json:"appSvcId,omitempty"`
	AppName         string        `json:"appName,omitempty"`
	Active          bool          `json:"active"`
	UID             string        `json:"uid,omitempty"`
	AppDataBlob     []AppDataBlob `json:"appDataBlob,omitempty"`
	AppDataBlobV6   []AppDataBlob `json:"appDataBlobV6,omitempty"`
	CreatedBy       string        `json:"createdBy,omitempty"`
	EditedBy        string        `json:"editedBy,omitempty"`
	EditedTimestamp string        `json:"editedTimestamp,omitempty"`
	ZappDataBlob    string        `json:"zappDataBlob,omitempty"`
	ZappDataBlobV6  string        `json:"zappDataBlobV6,omitempty"`
}

type AppDataBlob struct {
	Proto  string `json:"proto,omitempty"`
	Port   string `json:"port,omitempty"`
	Ipaddr string `json:"ipaddr,omitempty"`
	Fqdn   string `json:"fqdn,omitempty"`
}

func GetPredefinedIPApps(ctx context.Context, service *zscaler.Service, search string, page, pageSize *int) (*PredefinedIPAppsResponse, *http.Response, error) {
	params := common.QueryParams{
		Search: search,
	}
	if page != nil {
		params.Page = *page
	}
	if pageSize != nil {
		params.PageSize = *pageSize
	}

	var response PredefinedIPAppsResponse
	resp, err := service.Client.NewZccRequestDo(ctx, "GET", predefinedIpAppsEndpoint, params, nil, &response)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve predefined IP-based apps: %w", err)
	}
	return &response, resp, nil
}

func GetByAppID(ctx context.Context, service *zscaler.Service, appID string) (*PredefinedIPApp, *http.Response, error) {
	if appID == "" {
		return nil, nil, fmt.Errorf("appId is required")
	}
	endpoint := fmt.Sprintf("%s/%s", predefinedIpAppsEndpoint, appID)

	var app PredefinedIPApp
	resp, err := service.Client.NewZccRequestDo(ctx, "GET", endpoint, nil, nil, &app)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve predefined IP-based app %s: %w", appID, err)
	}
	return &app, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, name string) (*PredefinedIPApp, *http.Response, error) {
	pageSize := 1000
	page := 1

	for {
		res, resp, err := GetPredefinedIPApps(ctx, service, "", &page, &pageSize)
		if err != nil {
			return nil, resp, err
		}
		for _, a := range res.AppServiceContracts {
			if strings.EqualFold(a.AppName, name) {
				return &a, resp, nil
			}
		}
		if len(res.AppServiceContracts) < pageSize {
			break
		}
		page++
	}
	return nil, nil, fmt.Errorf("no predefined IP-based app found with name: %s", name)
}
