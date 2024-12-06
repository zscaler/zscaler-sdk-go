package datacenter

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	vipGroupByDatacenterEndpoint = "/zia/api/v1/vips/groupByDatacenter"
)

type DatacenterVIPS struct {
	Datacenter struct {
		Name        string  `json:"datacenter"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		City        string  `json:"city"`
		CountryCode string  `json:"countryCode"`
		Region      string  `json:"region"`
	} `json:"datacenter"`
	GreVIP []GreVIP `json:"greVips"`
}

type GreVIP struct {
	ID                 int    `json:"id,omitempty"`
	VirtualIp          string `json:"virtualIp,omitempty"`
	PrivateServiceEdge bool   `json:"privateServiceEdge,omitempty"`
	Datacenter         string `json:"datacenter,omitempty"`
}

func SearchByDatacenters(ctx context.Context, service *zscaler.Service, params common.DatacenterSearchParameters) ([]DatacenterVIPS, error) {
	var zscalerVips []DatacenterVIPS
	var queryParams []string

	if params.RoutableIP {
		queryParams = append(queryParams, "routableIP=true")
	}
	if params.WithinCountryOnly {
		queryParams = append(queryParams, "withinCountryOnly=true")
	}
	if params.IncludePrivateServiceEdge {
		queryParams = append(queryParams, "includePrivateServiceEdge=true")
	}
	if params.IncludeCurrentVips {
		queryParams = append(queryParams, "includeCurrentVips=true")
	}
	if params.SourceIp != "" {
		queryParams = append(queryParams, "sourceIp="+url.QueryEscape(params.SourceIp))
	}
	if params.Latitude != 0 {
		latitudeStr := strconv.FormatFloat(params.Latitude, 'f', -1, 64)
		queryParams = append(queryParams, "latitude="+latitudeStr)
	}
	if params.Longitude != 0 {
		longitudeStr := strconv.FormatFloat(params.Longitude, 'f', -1, 64)
		queryParams = append(queryParams, "longitude="+longitudeStr)
	}
	if params.Subcloud != "" {
		queryParams = append(queryParams, "subcloud="+url.QueryEscape(params.Subcloud))
	}

	endpoint := vipGroupByDatacenterEndpoint
	if len(queryParams) > 0 {
		endpoint += "?" + strings.Join(queryParams, "&")
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &zscalerVips)
	if err != nil {
		return nil, err
	}
	return zscalerVips, nil
}
