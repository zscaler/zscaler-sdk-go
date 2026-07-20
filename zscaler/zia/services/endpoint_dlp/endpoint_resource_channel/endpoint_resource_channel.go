package endpoint_resource_channel

import (
	"context"
	"fmt"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource"
)

const (
	dlpEndpointResourceEndpoint = "/zia/api/v1/dlpEndpointResource"
)

// Channel identifies a DLP endpoint resource channel. The dlpEndpointResource
// endpoints support Printing, Removable Storage, Network Share, and Personal
// Cloud Storage channels.
type Channel string

const (
	ChannelPrinting               Channel = "PRINTING"
	ChannelRemovableDriveTransfer Channel = "REMOVABLE_DRIVE_TRANSFER"
	ChannelNetworkDriveTransfer   Channel = "NETWORK_DRIVE_TRANSFER"
	ChannelPersonalCloudStorage   Channel = "PERSONAL_CLOUD_STORAGE"
)

// GetChannelFilterOptions holds the optional query parameters supported by
// GetChannelList.
type GetChannelFilterOptions struct {
	// SortOrder specifies the sorting order for the list by ascending or
	// descending order of the DLP resource names. Optional.
	SortOrder string

	// Name is the search string used to filter the list by DLP resource name
	// and other fields. Optional.
	Name string
}

// GetChannelList retrieves the list of DLP resources configured for the
// specified channel.
//
// The channel parameter is required and is part of the endpoint path. The
// sortOrder and name parameters are optional. The endpoint supports pagination,
// so common.ReadAllPages is used to aggregate all pages; the page and pageSize
// query parameters are handled internally by the pagination helper.
func GetChannelList(ctx context.Context, service *zscaler.Service, channel Channel, opts *GetChannelFilterOptions) ([]endpoint_resource.EndpointResource, error) {
	var channels []endpoint_resource.EndpointResource
	endpoint := fmt.Sprintf("%s/%s", dlpEndpointResourceEndpoint, url.PathEscape(string(channel)))

	queryParams := url.Values{}
	if opts != nil {
		if opts.SortOrder != "" {
			queryParams.Set("sortOrder", opts.SortOrder)
		}
		if opts.Name != "" {
			queryParams.Set("name", opts.Name)
		}
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &channels)
	return channels, err
}

// GetChannel retrieves information about a single DLP resource with the
// specified ID for the given channel. Both the channel and id parameters are
// required and are part of the endpoint path.
func GetChannel(ctx context.Context, service *zscaler.Service, channel Channel, id int) (*endpoint_resource.EndpointResource, error) {
	var resource endpoint_resource.EndpointResource
	endpoint := fmt.Sprintf("%s/%s/%d", dlpEndpointResourceEndpoint, url.PathEscape(string(channel)), id)
	err := service.Client.Read(ctx, endpoint, &resource)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning dlp endpoint resource from Get: %d", resource.ID)
	return &resource, nil
}
