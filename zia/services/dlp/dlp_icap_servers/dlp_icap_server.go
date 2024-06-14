package dlp_icap_servers

import (
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	dlpIcapServersEndpoint = "/icapServers"
)

type DLPICAPServers struct {
	// The unique identifier for a DLP server.
	ID int `json:"id"`

	// The DLP server name.
	Name string `json:"name,omitempty"`

	// The DLP server URL.
	URL string `json:"url,omitempty"`

	// The DLP server status
	Status string `json:"status,omitempty"`
}

func Get(service *services.Service, icapServerID int) (*DLPICAPServers, error) {
	var icapServers DLPICAPServers
	err := service.Client.Read(fmt.Sprintf("%s/%d", dlpIcapServersEndpoint, icapServerID), &icapServers)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning dlp icap server from Get: %d", icapServers.ID)
	return &icapServers, nil
}

func GetByName(service *services.Service, icapServerName string) (*DLPICAPServers, error) {
	var icapServers []DLPICAPServers
	err := common.ReadAllPages(service.Client, dlpIcapServersEndpoint, &icapServers)
	if err != nil {
		return nil, err
	}
	for _, icap := range icapServers {
		if strings.EqualFold(icap.Name, icapServerName) {
			return &icap, nil
		}
	}
	return nil, fmt.Errorf("no dlp icap server found with name: %s", icapServerName)
}

func GetAll(service *services.Service) ([]DLPICAPServers, error) {
	var icapServers []DLPICAPServers
	err := common.ReadAllPages(service.Client, dlpIcapServersEndpoint, &icapServers)
	return icapServers, err
}
