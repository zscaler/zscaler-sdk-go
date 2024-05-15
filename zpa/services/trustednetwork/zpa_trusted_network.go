package trustednetwork

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfigV1           = "/mgmtconfig/v1/admin/customers/"
	mgmtConfigV2           = "/mgmtconfig/v2/admin/customers/"
	trustedNetworkEndpoint = "/network"
)

type TrustedNetwork struct {
	CreationTime     string `json:"creationTime,omitempty"`
	Domain           string `json:"domain,omitempty"`
	ID               string `json:"id,omitempty"`
	MasterCustomerID string `json:"masterCustomerId,omitempty"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedTime     string `json:"modifiedTime,omitempty"`
	Name             string `json:"name,omitempty"`
	NetworkID        string `json:"networkId,omitempty"`
	ZscalerCloud     string `json:"zscalerCloud,omitempty"`
}

func (service *Service) Get(networkID string) (*TrustedNetwork, *http.Response, error) {
	v := new(TrustedNetwork)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.Config.CustomerID+trustedNetworkEndpoint, networkID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByNetID(netID string) (*TrustedNetwork, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV2 + service.Client.Config.CustomerID + trustedNetworkEndpoint)
	list, resp, err := common.GetAllPagesGeneric[TrustedNetwork](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, trustedNetwork := range list {
		if trustedNetwork.NetworkID == netID {
			return &trustedNetwork, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no trusted network with NetworkID '%s' was found", netID)
}

func (service *Service) GetByName(trustedNetworkName string) (*TrustedNetwork, *http.Response, error) {
	adaptedtrustedNetworkName := common.RemoveCloudSuffix(trustedNetworkName)
	relativeURL := mgmtConfigV2 + service.Client.Config.CustomerID + trustedNetworkEndpoint

	// Set up custom filters for pagination
	filters := common.Filter{Search: adaptedtrustedNetworkName} // Using the adapted trusted Network Name for searching

	// Use the custom pagination function with custom filters
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[TrustedNetwork](service.Client, relativeURL, filters)
	if err != nil {
		return nil, nil, err
	}

	// Iterate through the list and find the trusted network by its name
	for _, trustedNetwork := range list {
		if strings.EqualFold(common.RemoveCloudSuffix(trustedNetwork.Name), adaptedtrustedNetworkName) {
			return &trustedNetwork, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no trusted network named '%s' was found", trustedNetworkName)
}

func (service *Service) GetAll() ([]TrustedNetwork, *http.Response, error) {
	relativeURL := mgmtConfigV2 + service.Client.Config.CustomerID + trustedNetworkEndpoint
	list, resp, err := common.GetAllPagesGeneric[TrustedNetwork](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
