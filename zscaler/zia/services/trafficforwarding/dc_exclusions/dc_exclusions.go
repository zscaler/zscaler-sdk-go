package dc_exclusions

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	dcExclusionsEndpoint = "/zia/api/v1/dcExclusions"
	datacentersEndpoint  = "/zia/api/v1/datacenters"
)

type DCExclusions struct {
	DcID        int                      `json:"dcid,omitempty"`
	Expired     bool                     `json:"expired,omitempty"`
	StartTime   int                      `json:"startTime,omitempty"`
	EndTime     int                      `json:"endTime,omitempty"`
	Description string                   `json:"description,omitempty"`
	DcName      *common.IDNameExtensions `json:"dcName,omitempty"`
}

type Datacenter struct {
	Datacenter string `json:"datacenter,omitempty"`
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]DCExclusions, error) {
	var gws []DCExclusions
	err := common.ReadAllPages(ctx, service.Client, dcExclusionsEndpoint, &gws)
	return gws, err
}

func Create(ctx context.Context, service *zscaler.Service, dcID *DCExclusions) (*DCExclusions, *http.Response, error) {
	resp, err := service.Client.Create(ctx, dcExclusionsEndpoint, *dcID)
	if err != nil {
		return nil, nil, err
	}

	createdExclusions, ok := resp.(*DCExclusions)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a dc exclusion pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new dc exclusion from create: %d", createdExclusions.DcID)
	return createdExclusions, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, dcID int, dcExclusions *DCExclusions) (*DCExclusions, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", dcExclusionsEndpoint, dcID), *dcExclusions)
	if err != nil {
		return nil, nil, err
	}
	updatedExclusions, _ := resp.(*DCExclusions)

	service.Client.GetLogger().Printf("[DEBUG]returning updates dc exclusion from update: %d", updatedExclusions.DcID)
	return updatedExclusions, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, dcID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", dcExclusionsEndpoint, dcID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetDatacenters(ctx context.Context, service *zscaler.Service) ([]Datacenter, error) {
	var datacenters []Datacenter
	err := common.ReadAllPages(ctx, service.Client, datacentersEndpoint, &datacenters)
	return datacenters, err
}
