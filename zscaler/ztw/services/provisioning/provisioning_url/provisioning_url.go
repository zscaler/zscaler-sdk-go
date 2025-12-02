package provisioning_url

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/locationmanagement/locationtemplate"
)

const (
	provisioningUrlEndpoint = "/ztw/api/v1/provUrl"
)

type ProvisioningURL struct {
	ID             int                       `json:"id,omitempty"`
	Name           string                    `json:"name,omitempty"`
	Desc           string                    `json:"desc,omitempty"`
	ProvUrl        string                    `json:"provUrl,omitempty"`
	ProvUrlType    string                    `json:"provUrlType,omitempty"`
	LastModTime    int                       `json:"lastModTime,omitempty"`
	Status         string                    `json:"status,omitempty"`
	ProvUrlData    ProvUrlData               `json:"provUrlData,omitempty"`
	LastModUid     *common.IDNameExtensions  `json:"lastModUid,omitempty"`
	UsedInEcGroups []common.IDNameExtensions `json:"usedInEcGroups,omitempty"`
}

type ProvUrlData struct {
	ZsCloudDomain      string                            `json:"zsCloudDomain,omitempty"`
	OrgID              int                               `json:"orgId,omitempty"`
	ConfigServer       string                            `json:"configServer,omitempty"`
	RegistrationServer string                            `json:"registrationServer,omitempty"`
	ApiServer          string                            `json:"apiServer,omitempty"`
	PacServer          string                            `json:"pacServer,omitempty"`
	CloudProviderType  string                            `json:"cloudProviderType,omitempty"`
	FormFactor         string                            `json:"formFactor,omitempty"`
	LocationTemplate   locationtemplate.LocationTemplate `json:"locationTemplate,omitempty"`
	AutoScaleDetails   AutoScaleDetails                  `json:"autoScaleDetails,omitempty"`
	CellEdgeDeploy     bool                              `json:"cellEdgeDeploy,omitempty"`
	ReleaseChannel     string                            `json:"releaseChannel,omitempty"`
}

type AutoScaleDetails struct {
	AutoScale bool `json:"autoScale,omitempty"`
}

type BcGroup struct {
	ID                    int                            `json:"id,omitempty"`
	Name                  string                         `json:"name,omitempty"`
	Desc                  string                         `json:"desc,omitempty"`
	DeployType            string                         `json:"deployType,omitempty"`
	Platform              string                         `json:"platform,omitempty"`
	AwsAvailabilityZone   string                         `json:"awsAvailabilityZone,omitempty"`
	AzureAvailabilityZone string                         `json:"azureAvailabilityZone,omitempty"`
	MaxEcCount            int                            `json:"maxEcCount,omitempty"`
	TunnelMode            string                         `json:"tunnelMode,omitempty"`
	Status                []string                       `json:"status,omitempty"`
	EcVMs                 []EcVM                         `json:"ecVMs,omitempty"`
	ProvTemplate          *common.CommonIDNameExternalID `json:"provTemplate,omitempty"`
	Location              *common.CommonIDNameExternalID `json:"location,omitempty"`
}

type EcVM struct {
	ID                int          `json:"id,omitempty"`
	Name              string       `json:"name,omitempty"`
	Status            []string     `json:"status,omitempty"`
	OperationalStatus string       `json:"operationalStatus,omitempty"`
	FormFactor        string       `json:"formFactor,omitempty"`
	ManagementNw      ManagementNw `json:"managementNw,omitempty"`
	EcInstances       []EcInstance `json:"ecInstances,omitempty"`
	CityGeoId         int          `json:"cityGeoId,omitempty"`
	NatIp             string       `json:"natIp,omitempty"`
	ZiaGateway        string       `json:"ziaGateway,omitempty"`
	ZpaBroker         string       `json:"zpaBroker,omitempty"`
	BuildVersion      string       `json:"buildVersion,omitempty"`
	LastUpgradeTime   int          `json:"lastUpgradeTime,omitempty"`
	UpgradeStatus     int          `json:"upgradeStatus,omitempty"`
	UpgradeStartTime  int          `json:"upgradeStartTime,omitempty"`
	UpgradeEndTime    int          `json:"upgradeEndTime,omitempty"`
	UpgradeDayOfWeek  int          `json:"upgradeDayOfWeek,omitempty"`
}

type ManagementNw struct {
	ID             int    `json:"id,omitempty"`
	IpStart        string `json:"ipStart,omitempty"`
	IpEnd          string `json:"ipEnd,omitempty"`
	Netmask        string `json:"netmask,omitempty"`
	DefaultGateway string `json:"defaultGateway,omitempty"`
	NwType         string `json:"nwType,omitempty"`
	DNS            DNS    `json:"dns,omitempty"`
}

type DNS struct {
	ID      int      `json:"id,omitempty"`
	Ips     []string `json:"ips,omitempty"`
	DNSType string   `json:"dnsType,omitempty"`
}

type EcInstance struct {
	ID           int        `json:"id,omitempty"`
	InstanceType string     `json:"instanceType,omitempty"`
	ServiceIps   ServiceIps `json:"serviceIps,omitempty"`
	LbIpAddr     ServiceIps `json:"lbIpAddr,omitempty"`
	OutGwIp      string     `json:"outGwIp,omitempty"`
	NatIp        string     `json:"natIp,omitempty"`
	DnsIp        []string   `json:"dnsIp,omitempty"`
}

type ServiceIps struct {
	IpStart string `json:"ipStart,omitempty"`
	IpEnd   string `json:"ipEnd,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, provUrlID int) (*ProvisioningURL, error) {
	var provURLs ProvisioningURL
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", provisioningUrlEndpoint, provUrlID), &provURLs)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning provisioning URL from Get: %d", provURLs.ID)
	return &provURLs, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, provUrlName string) (*ProvisioningURL, error) {
	var provURLs []ProvisioningURL
	err := common.ReadAllPages(ctx, service.Client, provisioningUrlEndpoint, &provURLs)
	if err != nil {
		return nil, err
	}
	for _, provURL := range provURLs {
		if strings.EqualFold(provURL.Name, provUrlName) {
			return &provURL, nil
		}
	}
	return nil, fmt.Errorf("no provisioning URL found with name: %s", provUrlName)
}

func Create(ctx context.Context, service *zscaler.Service, ProvURL *ProvisioningURL) (*ProvisioningURL, *http.Response, error) {
	resp, err := service.Client.Create(ctx, provisioningUrlEndpoint, *ProvURL)
	if err != nil {
		return nil, nil, err
	}

	createdProvisioningURL, ok := resp.(*ProvisioningURL)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a provisioning URL pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new provisioning URL from create: %d", createdProvisioningURL.ID)
	return createdProvisioningURL, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, provisioningUrlID int, provisioningUrl *ProvisioningURL) (*ProvisioningURL, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", provisioningUrlEndpoint, provisioningUrlID), *provisioningUrl)
	if err != nil {
		return nil, nil, err
	}
	updatedProvisioningURL, _ := resp.(*ProvisioningURL)

	service.Client.GetLogger().Printf("[DEBUG]returning updates provisioning URL from update: %d", updatedProvisioningURL.ID)
	return updatedProvisioningURL, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, provisioningUrlID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", provisioningUrlEndpoint, provisioningUrlID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ProvisioningURL, error) {
	var provURLs []ProvisioningURL
	err := common.ReadAllPages(ctx, service.Client, provisioningUrlEndpoint, &provURLs)
	return provURLs, err
}
