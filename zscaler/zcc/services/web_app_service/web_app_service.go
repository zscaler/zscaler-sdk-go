package web_app_service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	baseWebAppServiceEndpoint = "/zcc/papi/public/v1/webAppService"
)

type WebAppService struct {
	ID             int           `json:"id,omitempty"`
	AppVersion     int           `json:"appVersion,omitempty"`
	AppSvcId       int           `json:"appSvcId,omitempty"`
	AppName        string        `json:"appName,omitempty"`
	Active         bool          `json:"active"`
	UID            string        `json:"uid,omitempty"`
	AppDataBlob    []AppDataBlob `json:"appDataBlob,omitempty"`
	AppDataBlobV6  []AppDataBlob `json:"appDataBlobV6,omitempty"`
	CreatedBy      string        `json:"createdBy,omitempty"`
	EditedBy       string        `json:"editedBy,omitempty"`
	EditedTimestamp string        `json:"editedTimestamp,omitempty"`
	ZappDataBlob   string        `json:"zappDataBlob,omitempty"`
	ZappDataBlobV6 string        `json:"zappDataBlobV6,omitempty"`
	Version        int           `json:"version,omitempty"`
}

type AppDataBlob struct {
	Proto  string `json:"proto,omitempty"`
	Port   string `json:"port,omitempty"`
	Ipaddr string `json:"ipaddr,omitempty"`
	Fqdn   string `json:"fqdn,omitempty"`
}

func GetWebAppServices(ctx context.Context, service *zscaler.Service, search string, page, pageSize *int) ([]WebAppService, error) {
	endpoint := fmt.Sprintf("%s/listByCompany", baseWebAppServiceEndpoint)

	params := common.QueryParams{
		Search: search,
	}
	if page != nil {
		params.Page = *page
	}
	if pageSize != nil {
		params.PageSize = *pageSize
	}

	var result []WebAppService
	_, err := service.Client.NewZccRequestDo(ctx, "GET", endpoint, params, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve web app services: %w", err)
	}
	return result, nil
}

func GetByAppID(ctx context.Context, service *zscaler.Service, appID string) (*WebAppService, error) {
	if appID == "" {
		return nil, fmt.Errorf("appId is required")
	}

	idInt, err := strconv.Atoi(appID)
	if err != nil {
		return nil, fmt.Errorf("invalid appId %q: %w", appID, err)
	}

	pageSize := 1000
	page := 1

	for {
		items, err := GetWebAppServices(ctx, service, "", &page, &pageSize)
		if err != nil {
			return nil, err
		}
		for i := range items {
			if items[i].ID == idInt {
				return &items[i], nil
			}
		}
		if len(items) < pageSize {
			break
		}
		page++
	}
	return nil, fmt.Errorf("web app service with ID %s not found", appID)
}

func GetByName(ctx context.Context, service *zscaler.Service, name string) (*WebAppService, error) {
	pageSize := 1000
	page := 1

	for {
		items, err := GetWebAppServices(ctx, service, "", &page, &pageSize)
		if err != nil {
			return nil, err
		}
		for i := range items {
			if strings.EqualFold(items[i].AppName, name) {
				return &items[i], nil
			}
		}
		if len(items) < pageSize {
			break
		}
		page++
	}
	return nil, fmt.Errorf("web app service with name %q not found", name)
}

func UpdateWebAppService(ctx context.Context, service *zscaler.Service, app *WebAppService) (*WebAppService, error) {
	if app == nil {
		return nil, errors.New("web app service payload is required")
	}

	url := fmt.Sprintf("%s/edit", baseWebAppServiceEndpoint)

	var updated WebAppService
	_, err := service.Client.NewZccRequestDo(ctx, "PUT", url, nil, app, &updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update web app service: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning web app service from update: %+v", updated)
	return &updated, nil
}
