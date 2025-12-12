package private_cloud_controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                     = "/zpa/mgmtconfig/v1/admin/customers/"
	privateCloudControllerEndpoint = "/privateCloudController"
)

type PrivateCloudController struct {
	ApplicationStartTime             string                 `json:"applicationStartTime,omitempty"`
	ControlChannelStatus             string                 `json:"controlChannelStatus,omitempty"`
	CreationTime                     string                 `json:"creationTime,omitempty"`
	CtrlBrokerName                   string                 `json:"ctrlBrokerName,omitempty"`
	CurrentVersion                   string                 `json:"currentVersion,omitempty"`
	Description                      string                 `json:"description,omitempty"`
	Enabled                          bool                   `json:"enabled,omitempty"`
	ExpectedSargeVersion             string                 `json:"expectedSargeVersion,omitempty"`
	ExpectedUpgradeTime              string                 `json:"expectedUpgradeTime,omitempty"`
	ExpectedVersion                  string                 `json:"expectedVersion,omitempty"`
	Fingerprint                      string                 `json:"fingerprint,omitempty"`
	ID                               string                 `json:"id,omitempty"`
	IpAcl                            []string               `json:"ipAcl,omitempty"`
	IssuedCertId                     string                 `json:"issuedCertId,omitempty"`
	LastBrokerConnectTime            string                 `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerConnectTimeDuration    string                 `json:"lastBrokerConnectTimeDuration,omitempty"`
	LastBrokerDisconnectTime         string                 `json:"lastBrokerDisconnectTime,omitempty"`
	LastBrokerDisconnectTimeDuration string                 `json:"lastBrokerDisconnectTimeDuration,omitempty"`
	LastOsUpgradeTime                string                 `json:"lastOSUpgradeTime,omitempty"`
	LastSargeUpgradeTime             string                 `json:"lastSargeUpgradeTime,omitempty"`
	LastUpgradeTime                  string                 `json:"lastUpgradeTime,omitempty"`
	Latitude                         string                 `json:"latitude,omitempty"`
	ListenIps                        []string               `json:"listenIps,omitempty"`
	Location                         string                 `json:"location,omitempty"`
	Longitude                        string                 `json:"longitude,omitempty"`
	MasterLastSyncTime               string                 `json:"masterLastSyncTime,omitempty"`
	ModifiedBy                       string                 `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                 `json:"modifiedTime,omitempty"`
	Name                             string                 `json:"name,omitempty"`
	ProvisioningKeyId                string                 `json:"provisioningKeyId,omitempty"`
	ProvisioningKeyName              string                 `json:"provisioningKeyName,omitempty"`
	OsUpgradeEnabled                 bool                   `json:"osUpgradeEnabled,omitempty"`
	OsUpgradeStatus                  string                 `json:"osUpgradeStatus,omitempty"`
	Platform                         string                 `json:"platform,omitempty"`
	PlatformDetail                   string                 `json:"platformDetail,omitempty"`
	PlatformVersion                  string                 `json:"platformVersion,omitempty"`
	PreviousVersion                  string                 `json:"previousVersion,omitempty"`
	PrivateIp                        string                 `json:"privateIp,omitempty"`
	PublicIp                         string                 `json:"publicIp,omitempty"`
	PublishIps                       []string               `json:"publishIps,omitempty"`
	ReadOnly                         bool                   `json:"readOnly,omitempty"`
	RestrictionType                  string                 `json:"restrictionType,omitempty"`
	Runtime                          string                 `json:"runtimeOS,omitempty"`
	SargeUpgradeAttempt              string                 `json:"sargeUpgradeAttempt,omitempty"`
	SargeUpgradeStatus               string                 `json:"sargeUpgradeStatus,omitempty"`
	SargeVersion                     string                 `json:"sargeVersion,omitempty"`
	MicrotenantId                    string                 `json:"microtenantId,omitempty"`
	MicrotenantName                  string                 `json:"microtenantName,omitempty"`
	ShardLastSyncTime                string                 `json:"shardLastSyncTime,omitempty"`
	EnrollmentCert                   map[string]interface{} `json:"enrollmentCert,omitempty"`
	PrivateCloudControllerGroupId    string                 `json:"privateCloudControllerGroupId,omitempty"`
	PrivateCloudControllerGroupName  string                 `json:"privateCloudControllerGroupName,omitempty"`
	PrivateCloudControllerVersion    map[string]interface{} `json:"privateCloudControllerVersion,omitempty"`
	SiteSpDnsName                    string                 `json:"siteSpDnsName,omitempty"`
	UpgradeAttempt                   string                 `json:"upgradeAttempt,omitempty"`
	UpgradeStatus                    string                 `json:"upgradeStatus,omitempty"`
	UserdbLastSyncTime               string                 `json:"userdbLastSyncTime,omitempty"`
	ZpnSubModuleUpgradeList          []interface{}          `json:"zpnSubModuleUpgradeList,omitempty"`
	ZscalerManaged                   bool                   `json:"zscalerManaged,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, controllerID string) (*PrivateCloudController, *http.Response, error) {
	v := new(PrivateCloudController)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+privateCloudControllerEndpoint, controllerID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// This function search the private cloud controller by Name
func GetByName(ctx context.Context, service *zscaler.Service, controllerName string) (*PrivateCloudController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + privateCloudControllerEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PrivateCloudController](ctx, service.Client, relativeURL, common.Filter{Search: controllerName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, controllerName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no private cloud controller named '%s' was found", controllerName)
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]PrivateCloudController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + privateCloudControllerEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PrivateCloudController](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

// Update Updates the private cloud controller details for the specified ID.
func Update(ctx context.Context, service *zscaler.Service, controllerID string, pcController PrivateCloudController) (*PrivateCloudController, *http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+privateCloudControllerEndpoint, controllerID)
	_, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, pcController, nil)
	if err != nil {
		return nil, nil, err
	}
	resource, resp, err := Get(ctx, service, controllerID)
	if err != nil {
		return nil, nil, err
	}
	return resource, resp, nil
}

func ControllerRestart(ctx context.Context, service *zscaler.Service, controllerID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+privateCloudControllerEndpoint+"/restart", controllerID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Delete Deletes the private cloud controller for the specified ID.
func Delete(ctx context.Context, service *zscaler.Service, controllerID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+privateCloudControllerEndpoint, controllerID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
