package extranet

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
	extranetEndpoint = "/zia/api/v1/extranet"
)

type Extranet struct {
	// The unique identifier for the extranet
	ID int `json:"id,omitempty"`

	// The name of the extranet
	Name string `json:"name,omitempty"`

	// The description of the extranet
	Description string `json:"description,omitempty"`

	// Information about the DNS servers specified for the extranet
	ExtranetDNSList []ExtranetDNSList `json:"extranetDNSList,omitempty"`

	// Information about the traffic selectors specified for the extranet (API returns an array).
	ExtranetIpPoolList []ExtranetPoolList `json:"extranetIpPoolList,omitempty"`

	// The Unix timestamp when the extranet was created
	CreatedAt int `json:"createdAt,omitempty"`

	// The Unix timestamp when the extranet was last modified
	ModifiedAt int `json:"modifiedAt,omitempty"`
}

type ExtranetDNSList struct {
	// The ID generated for the DNS server configuration
	ID int `json:"id,omitempty"`

	// The name of the DNS server
	Name string `json:"name,omitempty"`

	// The IP address of the primary DNS server
	PrimaryDNSServer string `json:"primaryDNSServer,omitempty"`

	// The IP address of the secondary DNS server
	SecondaryDNSServer string `json:"secondaryDNSServer,omitempty"`

	// A Boolean value indicating that the DNS servers specified in the extranet are the designated default servers
	UseAsDefault bool `json:"useAsDefault,omitempty"`
}

type ExtranetPoolList struct {
	// The ID generated for the DNS server configuration
	ID int `json:"id,omitempty"`

	// The name of the DNS server
	Name string `json:"name,omitempty"`

	IPStart string `json:"ipStart,omitempty"`

	IPEnd string `json:"ipEnd,omitempty"`

	// A Boolean value indicating that the DNS servers specified in the extranet are the designated default servers
	UseAsDefault bool `json:"useAsDefault,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, extranetID int) (*Extranet, error) {
	var extranet Extranet
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", extranetEndpoint, extranetID), &extranet)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning extranet from Get: %d", extranet.ID)
	return &extranet, nil
}

// GetAllOptions holds optional query parameters supported by the extranet API.
// The API does not support page or pageSize; only orderBy, order, and search are supported.
type GetAllOptions struct {
	OrderBy string // e.g. "id", "name"
	Order   string // e.g. "asc", "desc"
	Search  string
}

// GetExtranetByName returns an extranet by name (case-insensitive). It uses GetAll internally
// so that the non-paginated API is used and page/pageSize are not sent.
func GetExtranetByName(ctx context.Context, service *zscaler.Service, extranetName string) (*Extranet, error) {
	extranets, err := GetAll(ctx, service, nil)
	if err != nil {
		return nil, err
	}
	for i := range extranets {
		if strings.EqualFold(extranets[i].Name, extranetName) {
			return &extranets[i], nil
		}
	}
	return nil, fmt.Errorf("no extranet found with name: %s", extranetName)
}

func Create(ctx context.Context, service *zscaler.Service, extranet *Extranet) (*Extranet, *http.Response, error) {
	resp, err := service.Client.Create(ctx, extranetEndpoint, *extranet)
	if err != nil {
		return nil, nil, err
	}

	createdExtranet, ok := resp.(*Extranet)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a extranet pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new extranet from create: %d", createdExtranet.ID)
	return createdExtranet, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, extranetID int, extranet *Extranet) (*Extranet, error) {
	extranet.ID = extranetID
	resp, err := service.Client.UpdateWithPut(ctx, extranetEndpoint, *extranet)
	if err != nil {
		return nil, err
	}
	updatedExtranet, _ := resp.(*Extranet)
	service.Client.GetLogger().Printf("[DEBUG]returning extranet from update: %d", updatedExtranet.ID)
	return updatedExtranet, nil
}

func Delete(ctx context.Context, service *zscaler.Service, extranetID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", extranetEndpoint, extranetID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAll returns all extranets. The extranet API does not support pagination (page, pageSize);
// it only supports orderBy, order, and search. opts may be nil to use no query parameters.
func GetAll(ctx context.Context, service *zscaler.Service, opts *GetAllOptions) ([]Extranet, error) {
	endpoint := extranetEndpoint
	if opts != nil {
		params := url.Values{}
		if opts.OrderBy != "" {
			params.Set("orderBy", opts.OrderBy)
		}
		if opts.Order != "" {
			params.Set("order", opts.Order)
		}
		if opts.Search != "" {
			params.Set("search", opts.Search)
		}
		if len(params) > 0 {
			endpoint = endpoint + "?" + params.Encode()
		}
	}
	var extranets []Extranet
	err := service.Client.Read(ctx, endpoint, &extranets)
	return extranets, err
}

// GetLite returns lite extranet list. Uses Read directly; the extranet API does not support page/pageSize.
func GetLite(ctx context.Context, service *zscaler.Service) ([]common.IDNameExternalID, error) {
	var list []common.IDNameExternalID
	err := service.Client.Read(ctx, fmt.Sprintf("%s/lite", extranetEndpoint), &list)
	return list, err
}
