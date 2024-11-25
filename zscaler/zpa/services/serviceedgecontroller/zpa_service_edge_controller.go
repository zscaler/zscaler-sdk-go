package serviceedgecontroller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                    = "/zpa/mgmtconfig/v1/admin/customers/"
	serviceEdgeControllerEndpoint = "/serviceEdge"
	scheduleEndpoint              = "/serviceEdgeSchedule"
)

type ServiceEdgeController struct {
	ApplicationStartTime             string                 `json:"applicationStartTime,omitempty"`
	ServiceEdgeGroupID               string                 `json:"serviceEdgeGroupId,omitempty"`
	ServiceEdgeGroupName             string                 `json:"serviceEdgeGroupName,omitempty"`
	ControlChannelStatus             string                 `json:"controlChannelStatus,omitempty"`
	CreationTime                     string                 `json:"creationTime,omitempty"`
	CtrlBrokerName                   string                 `json:"ctrlBrokerName,omitempty"`
	CurrentVersion                   string                 `json:"currentVersion,omitempty"`
	Description                      string                 `json:"description,omitempty"`
	Enabled                          bool                   `json:"enabled,omitempty"`
	ExpectedUpgradeTime              string                 `json:"expectedUpgradeTime,omitempty"`
	ExpectedVersion                  string                 `json:"expectedVersion,omitempty"`
	Fingerprint                      string                 `json:"fingerprint,omitempty"`
	ID                               string                 `json:"id,omitempty"`
	IPACL                            string                 `json:"ipAcl,omitempty"`
	IssuedCertID                     string                 `json:"issuedCertId,omitempty"`
	LastBrokerConnectTime            string                 `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerConnectTimeDuration    string                 `json:"lastBrokerConnectTimeDuration,omitempty"`
	LastBrokerDisconnectTime         string                 `json:"lastBrokerDisconnectTime,omitempty"`
	LastBrokerDisconnectTimeDuration string                 `json:"lastBrokerDisconnectTimeDuration,omitempty"`
	LastUpgradeTime                  string                 `json:"lastUpgradeTime,omitempty"`
	Latitude                         string                 `json:"latitude,omitempty"`
	Location                         string                 `json:"location,omitempty"`
	Longitude                        string                 `json:"longitude,omitempty"`
	ListenIPs                        string                 `json:"listenIps,omitempty"`
	ModifiedBy                       string                 `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                 `json:"modifiedTime,omitempty"`
	Name                             string                 `json:"name,omitempty"`
	ProvisioningKeyID                string                 `json:"provisioningKeyId"`
	ProvisioningKeyName              string                 `json:"provisioningKeyName"`
	Platform                         string                 `json:"platform,omitempty"`
	PlatformDetail                   string                 `json:"platformDetail,omitempty"`
	PreviousVersion                  string                 `json:"previousVersion,omitempty"`
	PrivateIP                        string                 `json:"privateIp,omitempty"`
	PublicIP                         string                 `json:"publicIp,omitempty"`
	PublishIPs                       []string               `json:"publishIps,omitempty"`
	PublishIPv6                      bool                   `json:"publishIpv6,omitempty"`
	RuntimeOS                        string                 `json:"runtimeOS,omitempty"`
	SargeVersion                     string                 `json:"sargeVersion,omitempty"`
	EnrollmentCert                   map[string]interface{} `json:"enrollmentCert,omitempty"`
	UpgradeAttempt                   string                 `json:"upgradeAttempt,omitempty"`
	UpgradeStatus                    string                 `json:"upgradeStatus,omitempty"`
	MicroTenantID                    string                 `json:"microtenantId,omitempty"`
	MicroTenantName                  string                 `json:"microtenantName,omitempty"`
	PrivateBrokerVersion             PrivateBrokerVersion   `json:"privateBrokerVersion,omitempty"`
}

type PrivateBrokerVersion struct {
	ID                      string                       `json:"id,omitempty"`
	ApplicationStartTime    string                       `json:"applicationStartTime,omitempty"`
	BrokerId                string                       `json:"brokerId,omitempty"`
	CreationTime            string                       `json:"creationTime,omitempty"`
	CtrlChannelStatus       string                       `json:"ctrlChannelStatus,omitempty"`
	CurrentVersion          string                       `json:"currentVersion,omitempty"`
	DisableAutoUpdate       bool                         `json:"disableAutoUpdate,omitempty"`
	LastConnectTime         string                       `json:"lastConnectTime,omitempty"`
	LastDisconnectTime      string                       `json:"lastDisconnectTime,omitempty"`
	LastUpgradedTime        string                       `json:"lastUpgradedTime,omitempty"`
	LoneWarrior             bool                         `json:"loneWarrior,omitempty"`
	ModifiedBy              string                       `json:"modifiedBy,omitempty"`
	ModifiedTime            string                       `json:"modifiedTime,omitempty"`
	Platform                string                       `json:"platform,omitempty"`
	PlatformDetail          string                       `json:"platformDetail,omitempty"`
	PreviousVersion         string                       `json:"previousVersion,omitempty"`
	ServiceEdgeGroupID      string                       `json:"serviceEdgeGroupId,omitempty"`
	PrivateIP               string                       `json:"privateIp,omitempty"`
	PublicIP                string                       `json:"publicIp,omitempty"`
	RestartInstructions     string                       `json:"restartInstructions,omitempty"`
	RestartTimeInSec        string                       `json:"restartTimeInSec,omitempty"`
	RuntimeOS               string                       `json:"runtimeOS,omitempty"`
	SargeVersion            string                       `json:"sargeVersion,omitempty"`
	SystemStartTime         string                       `json:"systemStartTime,omitempty"`
	TunnelId                string                       `json:"tunnelId,omitempty"`
	UpgradeAttempt          string                       `json:"upgradeAttempt,omitempty"`
	UpgradeStatus           string                       `json:"upgradeStatus,omitempty"`
	UpgradeNowOnce          bool                         `json:"upgradeNowOnce,omitempty"`
	ZPNSubModuleUpgradeList []common.ZPNSubModuleUpgrade `json:"zpnSubModuleUpgradeList,omitempty"`
}

type AssistantSchedule struct {
	// The unique identifier for the Service Edge Controller auto deletion configuration for a customer. This field is only required for the PUT request to update the frequency of the App Connector Settings.
	ID string `json:"id,omitempty"`

	// The unique identifier of the ZPA tenant.
	CustomerID string `json:"customerId"`

	// Indicates if the Service Edge Controller are included for deletion if they are in a disconnected state based on frequencyInterval and frequency values.
	DeleteDisabled bool `json:"deleteDisabled"`

	// Indicates if the setting for deleting Service Edge Controller is enabled or disabled.
	Enabled bool `json:"enabled"`

	// The scheduled frequency at which the disconnected Service Edge Controller are deleted.
	Frequency string `json:"frequency"`

	// The interval for the configured frequency value. The minimum supported value is 5.
	FrequencyInterval string `json:"frequencyInterval"`
}

func Get(ctx context.Context, service *zscaler.Service, serviceEdgeID string) (*ServiceEdgeController, *http.Response, error) {
	v := new(ServiceEdgeController)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+serviceEdgeControllerEndpoint, serviceEdgeID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, serviceEdgeName string) (*ServiceEdgeController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + serviceEdgeControllerEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ServiceEdgeController](ctx, service.Client, relativeURL, common.Filter{Search: serviceEdgeName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, service := range list {
		if strings.EqualFold(service.Name, serviceEdgeName) {
			return &service, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no service edge named '%s' was found", serviceEdgeName)
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ServiceEdgeController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + serviceEdgeControllerEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ServiceEdgeController](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

type BulkDeleteRequest struct {
	IDs []string `json:"ids"`
}

// Update Updates the Service Edge details for the specified ID.
func Update(ctx context.Context, service *zscaler.Service, serviceEdgeID string, serviceEdge ServiceEdgeController) (*ServiceEdgeController, *http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+serviceEdgeControllerEndpoint, serviceEdgeID)
	_, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, serviceEdge, nil)
	if err != nil {
		return nil, nil, err
	}
	resource, resp, err := Get(ctx, service, serviceEdgeID)
	if err != nil {
		return nil, nil, err
	}
	return resource, resp, nil
}

// Delete Deletes the Service Edge for the specified ID.
func Delete(ctx context.Context, service *zscaler.Service, serviceEdgeID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+serviceEdgeControllerEndpoint, serviceEdgeID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// BulkDelete Bulk deletes the Service Edge.
func BulkDelete(ctx context.Context, service *zscaler.Service, serviceEdgeIDs []string) (*http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + serviceEdgeControllerEndpoint + "/bulkDelete"
	resp, err := service.Client.NewRequestDo(ctx, "POST", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, BulkDeleteRequest{IDs: serviceEdgeIDs}, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
