package networkapplications

import (
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	networkApplicationsEndpoint = "/networkApplications"
)

type NetworkApplications struct {
	ID             string `json:"id"`
	ParentCategory string `json:"parentCategory,omitempty"`
	Description    string `json:"description,omitempty"`
	Deprecated     bool   `json:"deprecated"`
}

func (service *Service) GetNetworkApplication(id, locale string) (*NetworkApplications, error) {
	var networkApplications NetworkApplications
	url := fmt.Sprintf("%s/%s", networkApplicationsEndpoint, id)
	if locale != "" {
		url = fmt.Sprintf("%s?locale=%s", url, locale)
	}
	err := service.Client.Read(url, &networkApplications)
	if err != nil {
		return nil, err
	}
	return &networkApplications, nil
}

func (service *Service) GetAll() ([]NetworkApplications, error) {
	var networkApplications []NetworkApplications
	err := common.ReadAllPages(service.Client, networkApplicationsEndpoint, &networkApplications)
	return networkApplications, err
}
