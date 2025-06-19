package tenancy_restriction

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	tenantRestrictionEndpoint = "/zia/api/v1/tenancyRestrictionProfile"
)

type TenancyRestrictionProfile struct {
	ID                          int      `json:"id,omitempty"`
	Name                        string   `json:"name,omitempty"`
	AppType                     string   `json:"appType,omitempty"`
	Description                 string   `json:"description,omitempty"`
	ItemTypePrimary             string   `json:"itemTypePrimary,omitempty"`
	ItemTypeSecondary           string   `json:"itemTypeSecondary,omitempty"`
	RestrictPersonalO365Domains bool     `json:"restrictPersonalO365Domains,omitempty"`
	AllowGoogleConsumers        bool     `json:"allowGoogleConsumers,omitempty"`
	MsLoginServicesTrV2         bool     `json:"msLoginServicesTrV2,omitempty"`
	AllowGoogleVisitors         bool     `json:"allowGoogleVisitors,omitempty"`
	AllowGcpCloudStorageRead    bool     `json:"allowGcpCloudStorageRead,omitempty"`
	ItemDataPrimary             []string `json:"itemDataPrimary,omitempty"`
	ItemDataSecondary           []string `json:"itemDataSecondary,omitempty"`
	ItemValue                   []string `json:"itemValue,omitempty"`
	LastModifiedTime            int      `json:"lastModifiedTime,omitempty"`
	LastModifiedUserID          int      `json:"lastModifiedUserId,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, profileID int) (*TenancyRestrictionProfile, error) {
	var profile TenancyRestrictionProfile
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", tenantRestrictionEndpoint, profileID), &profile)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning tenant restriction profile from Get: %d", profile.ID)
	return &profile, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, profileName string) (*TenancyRestrictionProfile, error) {
	var profiles []TenancyRestrictionProfile
	err := common.ReadAllPages(ctx, service.Client, tenantRestrictionEndpoint, &profiles)
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		if strings.EqualFold(profile.Name, profileName) {
			return &profile, nil
		}
	}
	return nil, fmt.Errorf("no tenant restriction profile found with name: %s", profileName)
}

func Create(ctx context.Context, service *zscaler.Service, instanceID *TenancyRestrictionProfile) (*TenancyRestrictionProfile, *http.Response, error) {
	resp, err := service.Client.Create(ctx, tenantRestrictionEndpoint, *instanceID)
	if err != nil {
		return nil, nil, err
	}

	createdProfile, ok := resp.(*TenancyRestrictionProfile)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a tenant restriction profile pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new tenant restriction profile from create: %d", createdProfile.ID)
	return createdProfile, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, instanceID int, cloudInstance *TenancyRestrictionProfile) (*TenancyRestrictionProfile, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", tenantRestrictionEndpoint, instanceID), *cloudInstance)
	if err != nil {
		return nil, nil, err
	}
	updatedProfile, _ := resp.(*TenancyRestrictionProfile)

	service.Client.GetLogger().Printf("[DEBUG]returning updates tenant restriction profile  from update: %d", updatedProfile.ID)
	return updatedProfile, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, profileID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", tenantRestrictionEndpoint, profileID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]TenancyRestrictionProfile, error) {
	var profiles []TenancyRestrictionProfile
	err := common.ReadAllPages(ctx, service.Client, tenantRestrictionEndpoint, &profiles)
	return profiles, err
}

func GetAppItemCount(ctx context.Context, service *zscaler.Service, appType, itemType string, excludeProfile ...int) (map[string]int, error) {
	endpoint := fmt.Sprintf("%s/app-item-count/%s/%s", tenantRestrictionEndpoint, url.PathEscape(appType), url.PathEscape(itemType))

	if len(excludeProfile) > 0 {
		endpoint += fmt.Sprintf("?excludeProfile=%d", excludeProfile[0])
	}

	var result map[string]int
	err := service.Client.Read(ctx, endpoint, &result)
	return result, err
}
