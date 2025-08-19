package private_cloud_group

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                          = "/zpa/mgmtconfig/v1/admin/customers/"
	privateCloudControllerGroupEndpoint = "/privateCloudControllerGroup"
)

type PrivateCloudGroup struct {
	ID                     string `json:"id,omitempty"`
	City                   string `json:"city,omitempty"`
	CityCountry            string `json:"cityCountry,omitempty"`
	CountryCode            string `json:"countryCode,omitempty"`
	Description            string `json:"description,omitempty"`
	Enabled                bool   `json:"enabled,omitempty"`
	GeoLocationID          string `json:"geoLocationId,omitempty"`
	IsPublic               string `json:"isPublic,omitempty"`
	Latitude               string `json:"latitude,omitempty"`
	Location               string `json:"location,omitempty"`
	Longitude              string `json:"longitude,omitempty"`
	Name                   string `json:"name,omitempty"`
	OverrideVersionProfile bool   `json:"overrideVersionProfile,omitempty"`
	ReadOnly               bool   `json:"readOnly,omitempty"`
	RestrictionType        string `json:"restrictionType,omitempty"`
	MicrotenantID          string `json:"microtenantId,omitempty"`
	MicrotenantName        string `json:"microtenantName,omitempty"`
	SiteID                 string `json:"siteId,omitempty"`
	SiteName               string `json:"siteName,omitempty"`
	UpgradeDay             string `json:"upgradeDay,omitempty"`
	UpgradeTimeInSecs      string `json:"upgradeTimeInSecs,omitempty"`
	VersionProfileID       string `json:"versionProfileId,omitempty"`
	VersionProfileName     string `json:"versionProfileName,omitempty"`
	ZscalerManaged         bool   `json:"zscalerManaged,omitempty"`
	CreationTime           string `json:"creationTime,omitempty"`
	ModifiedBy             string `json:"modifiedBy,omitempty"`
	ModifiedTime           string `json:"modifiedTime,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, groupID string) (*PrivateCloudGroup, *http.Response, error) {
	v := new(PrivateCloudGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+privateCloudControllerGroupEndpoint, groupID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, groupName string) (*PrivateCloudGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + privateCloudControllerGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PrivateCloudGroup](ctx, service.Client, relativeURL, common.Filter{Search: groupName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, groupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no private cloud group named '%s' was found", groupName)
}

func Create(ctx context.Context, service *zscaler.Service, controllerGroup PrivateCloudGroup) (*PrivateCloudGroup, *http.Response, error) {
	v := new(PrivateCloudGroup)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+privateCloudControllerGroupEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, controllerGroup, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, groupID string, controllerGroup *PrivateCloudGroup) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+privateCloudControllerGroupEndpoint, groupID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, controllerGroup, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, groupID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+privateCloudControllerGroupEndpoint, groupID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]PrivateCloudGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + privateCloudControllerGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PrivateCloudGroup](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetGroupSummary(ctx context.Context, service *zscaler.Service) ([]PrivateCloudGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + privateCloudControllerGroupEndpoint + "/summary"
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PrivateCloudGroup](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
