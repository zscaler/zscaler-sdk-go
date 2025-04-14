package zparesources

import (
	"context"

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

	ZpaID string `json:"zpaId,omitempty"`
}

func GetZPAApplicationSegments(ctx context.Context, service *zscaler.Service) ([]ZPAApplicationSegment, error) {
	var zpaAppSegments []ZPAApplicationSegment
	err := common.ReadAllPages(ctx, service.Client, zpaResourcesEndpoint, &zpaAppSegments)
	return zpaAppSegments, err
}
