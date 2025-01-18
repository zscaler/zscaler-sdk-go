package extranet

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	extranetEndpoint = "/zia/api/v1/extranet"
)

type Extranet struct {
	// The unique identifier for the extranet
	ID int `json:"id"`

	// The name of the extranet
	Name string `json:"name"`

	// The description of the extranet
	Description string `json:"description"`

	// Information about the DNS servers specified for the extranet
	ExtranetDNSList []ExtranetDNSList `json:"extranetDNSList"`

	// Information about the traffic selectors specified for the extranet. Type: Array
	ExtranetIpPoolList string `json:"extranetIpPoolList"` // Placeholder for "TBD" - refine as needed

	// The Unix timestamp when the extranet was created
	CreatedAt int `json:"createdAt"`

	// The Unix timestamp when the extranet was last modified
	ModifiedAt int `json:"modifiedAt"`
}

type ExtranetDNSList struct {
	// The ID generated for the DNS server configuration
	ID int `json:"id"`

	// The name of the DNS server
	Name string `json:"name"`

	// The IP address of the primary DNS server
	PrimaryDNSServer string `json:"primaryDNSServer"`

	// The IP address of the secondary DNS server
	SecondaryDNSServer string `json:"secondaryDNSServer"`

	// A Boolean value indicating that the DNS servers specified in the extranet are the designated default servers
	UseAsDefault bool `json:"useAsDefault"`
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

func GetExtranetByName(ctx context.Context, service *zscaler.Service, extranetName string) (*Extranet, error) {
	var extranets []Extranet
	err := common.ReadAllPages(ctx, service.Client, extranetEndpoint, &extranets)
	if err != nil {
		return nil, err
	}
	for _, extranet := range extranets {
		if strings.EqualFold(extranet.Name, extranetName) {
			return &extranet, nil
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
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", extranetEndpoint, extranetID), *extranet)
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

func GetAll(ctx context.Context, service *zscaler.Service) ([]Extranet, error) {
	var extranet []Extranet
	err := common.ReadAllPages(ctx, service.Client, extranetEndpoint, &extranet)
	return extranet, err
}

func GetLite(ctx context.Context, service *zscaler.Service) ([]common.IDNameExternalID, error) {
	var extranet []common.IDNameExternalID
	err := common.ReadAllPages(ctx, service.Client, fmt.Sprintf("%s/lite", extranetEndpoint), &extranet)
	return extranet, err
}
