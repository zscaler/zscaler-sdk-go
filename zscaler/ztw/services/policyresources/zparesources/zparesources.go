package zparesources

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	zpaResourcesEndpoint = "/ztw/api/v1/zpaResources/applicationSegments"
)

type ZPAApplicationSegment struct {
	// ID of the ZPA application segment.
	ID int `json:"id"`

	// Name of the ZPA application segment.
	Name string `json:"name,omitempty"`

	// Description of the ZPA application segment.
	Description string `json:"description,omitempty"`

	ZpaID int `json:"zpaId,omitempty"`

	Deleted bool `json:"deleted,omitempty"`
}

func GetZPAApplicationSegments(ctx context.Context, service *zscaler.Service) ([]ZPAApplicationSegment, error) {
	var zpaAppSegments []ZPAApplicationSegment
	err := common.ReadAllPages(ctx, service.Client, zpaResourcesEndpoint, &zpaAppSegments)
	return zpaAppSegments, err
}

func GetByName(ctx context.Context, service *zscaler.Service, targetGroup string) (*ZPAApplicationSegment, error) {
	var groups []ZPAApplicationSegment
	page := 1

	// Construct the endpoint with the search parameter
	endpointWithSearch := fmt.Sprintf("%s?search=%s&%s", zpaResourcesEndpoint, url.QueryEscape(targetGroup), common.GetSortParams(common.SortField(service.SortBy), common.SortOrder(service.SortOrder)))

	for {
		err := common.ReadPage(ctx, service.Client, endpointWithSearch, page, &groups)
		if err != nil {
			return nil, err
		}

		for _, group := range groups {
			if strings.EqualFold(group.Name, targetGroup) {
				return &group, nil
			}
		}

		// Break the loop if there are no more pages
		if len(groups) < common.GetPageSize() {
			break
		}
		page++
	}

	return nil, fmt.Errorf("no zpa resource with name: %s", targetGroup)
}
