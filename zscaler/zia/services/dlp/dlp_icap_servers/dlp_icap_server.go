package dlp_icap_servers

import (
	"context"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	dlpIcapServersEndpoint = "/zia/api/v1/icapServers"
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

func Get(ctx context.Context, service *zscaler.Service, icapServerID int) (*DLPICAPServers, error) {
	var icapServers DLPICAPServers
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", dlpIcapServersEndpoint, icapServerID), &icapServers)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning dlp icap server from Get: %d", icapServers.ID)
	return &icapServers, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, icapServerName string) (*DLPICAPServers, error) {
	var icapServers []DLPICAPServers
	err := common.ReadAllPages(ctx, service.Client, dlpIcapServersEndpoint, &icapServers)
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

func GetAll(ctx context.Context, service *zscaler.Service) ([]DLPICAPServers, error) {
	var icapServers []DLPICAPServers
	err := common.ReadAllPages(ctx, service.Client, dlpIcapServersEndpoint, &icapServers)
	return icapServers, err
}
