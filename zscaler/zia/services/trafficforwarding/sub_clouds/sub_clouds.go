package sub_clouds

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	subCloudsEndpoint = "/zia/api/v1/subclouds"
)

type SubClouds struct {
	ID         int          `json:"id,omitempty"`
	Name       string       `json:"name,omitempty"`
	Dcs        []DCs        `json:"dcs,omitempty"`
	Exclusions []Exclusions `json:"exclusions,omitempty"`
}

type DCs struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Country string `json:"country,omitempty"`
}

type Exclusions struct {
	Datacenter       *common.IDNameExtensions `json:"datacenter,omitempty"`
	LastModifiedUser *common.IDNameExtensions `json:"lastModifiedUser,omitempty"`
	Country          string                   `json:"country,omitempty"`
	Expired          bool                     `json:"expired,omitempty"`
	DisabledByOps    bool                     `json:"disabledByOps,omitempty"`
	CreateTime       int                      `json:"createTime,omitempty"`
	StartTime        int                      `json:"startTime,omitempty"`
	EndTime          int                      `json:"endTime,omitempty"`
	LastModifiedTime int                      `json:"lastModifiedTime,omitempty"`
}

type SubCloudCountryDCExclusionInfo struct {
	ID              int    `json:"id,omitempty"`
	DcIDs           []int  `json:"dcIds,omitempty"`
	Country         string `json:"country,omitempty"`
	LastDCExclusion bool   `json:"lastDCExclusion,omitempty"`
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]SubClouds, error) {
	var subClouds []SubClouds
	err := common.ReadAllPages(ctx, service.Client, subCloudsEndpoint, &subClouds)
	return subClouds, err
}

func GetByName(ctx context.Context, service *zscaler.Service, subCloudName string) (*SubClouds, error) {
	subClouds, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	for i := range subClouds {
		if strings.EqualFold(subClouds[i].Name, subCloudName) {
			return &subClouds[i], nil
		}
	}
	return nil, fmt.Errorf("no subcloud found with name: %s", subCloudName)
}

func Update(ctx context.Context, service *zscaler.Service, cloudID int, subClouds *SubClouds) (*SubClouds, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", subCloudsEndpoint, cloudID), *subClouds)
	if err != nil {
		return nil, nil, err
	}
	updatedSubCloud, _ := resp.(*SubClouds)

	service.Client.GetLogger().Printf("[DEBUG]returning updates subclouds from update: %d", updatedSubCloud.ID)
	return updatedSubCloud, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, cloudID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", subCloudsEndpoint, cloudID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func Get(ctx context.Context, service *zscaler.Service, cloudID int) (*SubCloudCountryDCExclusionInfo, error) {
	var subClouds SubCloudCountryDCExclusionInfo
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", subCloudsEndpoint+"/isLastDcInCountry/", cloudID), &subClouds)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning subcloud from Get: %d", subClouds.ID)
	return &subClouds, nil
}
