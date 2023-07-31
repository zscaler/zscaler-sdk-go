package trustednetwork

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	mgmtConfig             = "/mgmtconfig/v2/admin/customers/"
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
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+trustedNetworkEndpoint, networkID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByNetID(netID string) (*TrustedNetwork, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + trustedNetworkEndpoint)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[TrustedNetwork](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
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
	adaptedTrustedNetworkName := common.RemoveCloudSuffix(trustedNetworkName)
	// to avoid such errors:
	// - {"params" : [ "-" ], "id" : "filtering.input.invalid.operand", "reason" : "Invalid operand:- in Filtering criteria."}
	// - search=Corp++Trusted++Networks {"params" : [ "" ],"id" : "filtering.input.invalid.operand","reason" : "Invalid operand: in Filtering criteria."}
	adaptedTrustedNetworkName = strings.ReplaceAll(adaptedTrustedNetworkName, "-", " ")
	adaptedTrustedNetworkName = strings.TrimSpace(adaptedTrustedNetworkName)
	adaptedTrustedNetworkName = strings.Split(adaptedTrustedNetworkName, " ")[0]
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + trustedNetworkEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[TrustedNetwork](service.Client, relativeURL, common.Filter{Search: adaptedTrustedNetworkName, MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	for _, trustedNetwork := range list {
		if strings.EqualFold(common.RemoveCloudSuffix(trustedNetwork.Name), common.RemoveCloudSuffix(trustedNetworkName)) {
			return &trustedNetwork, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no trusted network named '%s' was found", trustedNetworkName)
}

func (service *Service) GetAll() ([]TrustedNetwork, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + trustedNetworkEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[TrustedNetwork](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
