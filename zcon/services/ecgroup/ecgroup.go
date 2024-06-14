package ecgroup

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zcon/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zcon/services/common"
)

const (
	ecGroupEndpoint     = "/ecgroup"
	ecGroupLiteEndpoint = "/ecgroup/lite"
)

type EcGroup struct {
	ID                    int                    `json:"id,omitempty"`
	Name                  string                 `json:"name,omitempty"`
	Description           string                 `json:"desc,omitempty"`
	DeployType            string                 `json:"deployType,omitempty"`
	Status                []string               `json:"status,omitempty"`
	Platform              string                 `json:"platform,omitempty"`
	AWSAvailabilityZone   string                 `json:"awsAvailabilityZone,omitempty"`
	AzureAvailabilityZone string                 `json:"azureAvailabilityZone,omitempty"`
	MaxEcCount            int                    `json:"maxEcCount,omitempty"`
	TunnelMode            string                 `json:"tunnelMode,omitempty"`
	Location              *common.GeneralPurpose `json:"location,omitempty"`
	ProvTemplate          *common.GeneralPurpose `json:"provTemplate,omitempty"`
	ECVMs                 []common.ECVMs         `json:"ecVMs,omitempty"`
}

type ManagementNw struct {
	ID             int    `json:"id,omitempty"`
	IPStart        string `json:"ipStart,omitempty"`
	IPEnd          string `json:"ipEnd,omitempty"`
	Netmask        string `json:"netmask,omitempty"`
	DefaultGateway string `json:"defaultGateway,omitempty"`
	NWType         string `json:"nwType,omitempty"`
	DNS            *DNS   `json:"dns,omitempty"`
}

type DNS struct {
	ID      int      `json:"id,omitempty"`
	IPs     []string `json:"ips,omitempty"`
	DNSType string   `json:"dnsType,omitempty"`
}

type ECInstances struct {
	ServiceNw      *ManagementNw `json:"serviceNw,omitempty"`
	VirtualNw      *ManagementNw `json:"virtualNw,omitempty"`
	ECInstanceType string        `json:"ecInstanceType,omitempty"`
	OutGwIp        string        `json:"outGwIp,omitempty"`
	NatIP          string        `json:"natIp,omitempty"`
	DNSIp          string        `json:"dnsIp,omitempty"`
}

type LBIPAddr struct {
	IPStart string `json:"ipStart,omitempty"`
	IPEnd   string `json:"ipEnd,omitempty"`
}

func Get(service *services.Service, ecGroupID int) (*EcGroup, error) {
	var ecGroup EcGroup
	err := service.Client.Read(fmt.Sprintf("%s/%d", ecGroupEndpoint, ecGroupID), &ecGroup)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Cloud & Branch Connector Group from Get: %d", ecGroup.ID)
	return &ecGroup, nil
}

func GetByName(service *services.Service, ecGroupName string) (*EcGroup, error) {
	var ecGroup []EcGroup
	// We are assuming this provisioning url name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(service.Client, ecGroupEndpoint, &ecGroup)
	if err != nil {
		return nil, err
	}
	for _, ec := range ecGroup {
		if strings.EqualFold(ec.Name, ecGroupName) {
			return &ec, nil
		}
	}
	return nil, fmt.Errorf("no Cloud & Branch Connector Group found with name: %s", ecGroupName)
}

func Delete(service *services.Service, ecGroupID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", ecGroupEndpoint, ecGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(service *services.Service) ([]EcGroup, error) {
	var ecgroups []EcGroup
	err := common.ReadAllPages(service.Client, ecGroupEndpoint, &ecgroups)
	return ecgroups, err
}

func GetEcGroupLiteID(service *services.Service, ecGroupID int) (*EcGroup, error) {
	var ecgroupLite EcGroup
	err := service.Client.Read(fmt.Sprintf("%s/%d", ecGroupLiteEndpoint, ecGroupID), &ecgroupLite)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]returning Cloud & Branch Connector Group from Get: %d", ecgroupLite.ID)
	return &ecgroupLite, nil
}

func GetEcGroupLiteByName(service *services.Service, ecGroupLiteName string) (*EcGroup, error) {
	var ecgroupLite []EcGroup
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?name=%s", ecGroupLiteEndpoint, url.QueryEscape(ecGroupLiteName)), &ecgroupLite)
	if err != nil {
		return nil, err
	}
	for _, ecgroupLite := range ecgroupLite {
		if strings.EqualFold(ecgroupLite.Name, ecGroupLiteName) {
			return &ecgroupLite, nil
		}
	}
	return nil, fmt.Errorf("no Cloud & Branch Connector Group found with name: %s", ecGroupLiteName)
}
