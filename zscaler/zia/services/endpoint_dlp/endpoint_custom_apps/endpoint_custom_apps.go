package endpoint_custom_apps

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource"
)

const (
	endPointApplicationsEndpoint = "/zia/api/v1/endPointApplications"
)

type EndpointApplications struct {
	ResourceID       int        `json:"resourceId,omitempty"`
	Description      string     `json:"description,omitempty"`
	OsType           string     `json:"osType,omitempty"`
	ApplicationName  string     `json:"applicationName,omitempty"`
	Bundle           string     `json:"bundleID,omitempty"`
	Filename         string     `json:"filename,omitempty"`
	OriginalFileName string     `json:"originalFileName,omitempty"`
	DigitallySigned  bool       `json:"digitallySigned,omitempty"`
	ModUId           int        `json:"modUId,omitempty"`
	LastModifiedTime int        `json:"lastModifiedTime,omitempty"`
	ApplicationType  string     `json:"applicationType,omitempty"`
	Version          Version    `json:"version,omitempty"`
	Versions         []Versions `json:"versions,omitempty"`
	ZappID           string     `json:"zappId,omitempty"`
	Deleted          bool       `json:"deleted,omitempty"`
}

type Version struct {
	Version                      string `json:"version,omitempty"`
	ZverIDMD32                   int    `json:"z_ver_id_md32,omitempty"`
	ThreatType                   int    `json:"threat_type,omitempty"`
	ThreatLevel                  string `json:"threat_level,omitempty"`
	Bundle                       string `json:"bundleID,omitempty"`
	CodeSigningCertificateStatus int    `json:"code_signing_certificate_status,omitempty"`
	ThreatLevelUpdated           bool   `json:"threatLevelUpdated,omitempty"`
}

type Versions struct {
	Version                      string `json:"version,omitempty"`
	ZverIDMD32                   int    `json:"z_ver_id_md32,omitempty"`
	ThreatType                   int    `json:"threat_type,omitempty"`
	ThreatLevel                  string `json:"threat_level,omitempty"`
	Bundle                       string `json:"bundleID,omitempty"`
	CodeSigningCertificateStatus int    `json:"code_signing_certificate_status,omitempty"`
	ThreatLevelUpdated           bool   `json:"threatLevelUpdated,omitempty"`
}

// GetCustomAppsFilterOptions holds the optional query parameters supported by
// GetCustomApps.
type GetCustomAppsFilterOptions struct {
	// Search is the search string used to match against application names. Optional.
	Search string

	// OsType filters the results by operating system (e.g., Windows OS and Mac OS). Optional.
	OsType string
}

// GetCustomApps retrieves the list of custom applications.
//
// The search and osType parameters are optional. The endpoint supports
// pagination, so common.ReadAllPages is used to aggregate all pages; the page
// and pageSize query parameters are handled internally by the pagination helper.
func GetCustomApps(ctx context.Context, service *zscaler.Service, opts *GetCustomAppsFilterOptions) ([]EndpointApplications, error) {
	var applications []EndpointApplications
	endpoint := endPointApplicationsEndpoint + "/customApps"

	queryParams := url.Values{}
	if opts != nil {
		if opts.Search != "" {
			queryParams.Set("search", opts.Search)
		}
		if opts.OsType != "" {
			queryParams.Set("osType", opts.OsType)
		}
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &applications)
	return applications, err
}

// GetCustomApp retrieves information about the custom endpoint application with
// the specified ID. The id parameter is required.
func GetCustomApp(ctx context.Context, service *zscaler.Service, id int) (*EndpointApplications, error) {
	var customApp EndpointApplications
	err := service.Client.Read(ctx, fmt.Sprintf("%s/customApp/%d", endPointApplicationsEndpoint, id), &customApp)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning custom endpoint application from Get: %d", customApp.ResourceID)
	return &customApp, nil
}

func Create(ctx context.Context, service *zscaler.Service, customApp *endpoint_resource.EndpointResource) (*endpoint_resource.EndpointResource, *http.Response, error) {
	resp, err := service.Client.Create(ctx, endPointApplicationsEndpoint+"/customApp", *customApp)
	if err != nil {
		return nil, nil, err
	}

	createdCustomApp, ok := resp.(*endpoint_resource.EndpointResource)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a endpoint custom app pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new endpoint custom app from create: %d", createdCustomApp.ID)
	return createdCustomApp, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, appID int, customApp *endpoint_resource.EndpointResource) (*endpoint_resource.EndpointResource, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/customApp/%d", endPointApplicationsEndpoint, appID), *customApp)
	if err != nil {
		return nil, nil, err
	}
	updatedCustomApp, _ := resp.(*endpoint_resource.EndpointResource)

	service.Client.GetLogger().Printf("[DEBUG]returning updated endpoint custom app from update: %d", updatedCustomApp.ID)
	return updatedCustomApp, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, appID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/customApp/%d", endPointApplicationsEndpoint, appID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
