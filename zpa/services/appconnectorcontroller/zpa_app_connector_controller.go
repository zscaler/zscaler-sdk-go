package appconnectorcontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig           = "/mgmtconfig/v1/admin/customers/"
	appConnectorEndpoint = "/connector"
	scheduleEndpoint     = "/assistantSchedule"
)

type AppConnector struct {
	ApplicationStartTime             string                       `json:"applicationStartTime,omitempty"`
	AppConnectorGroupID              string                       `json:"appConnectorGroupId,omitempty"`
	AppConnectorGroupName            string                       `json:"appConnectorGroupName,omitempty"`
	AssistantVersion                 AssistantVersion             `json:"assistantVersion,omitempty"`
	ControlChannelStatus             string                       `json:"controlChannelStatus,omitempty"`
	CreationTime                     string                       `json:"creationTime,omitempty"`
	CtrlBrokerName                   string                       `json:"ctrlBrokerName,omitempty"`
	CurrentVersion                   string                       `json:"currentVersion,omitempty"`
	Description                      string                       `json:"description,omitempty"`
	Enabled                          bool                         `json:"enabled,omitempty"`
	ExpectedUpgradeTime              string                       `json:"expectedUpgradeTime,omitempty"`
	ExpectedVersion                  string                       `json:"expectedVersion,omitempty"`
	Fingerprint                      string                       `json:"fingerprint,omitempty"`
	ID                               string                       `json:"id,omitempty"`
	IPACL                            string                       `json:"ipAcl,omitempty"`
	IssuedCertID                     string                       `json:"issuedCertId,omitempty"`
	LastBrokerConnectTime            string                       `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerConnectTimeDuration    string                       `json:"lastBrokerConnectTimeDuration,omitempty"`
	LastBrokerDisconnectTime         string                       `json:"lastBrokerDisconnectTime,omitempty"`
	LastBrokerDisconnectTimeDuration string                       `json:"lastBrokerDisconnectTimeDuration,omitempty"`
	LastUpgradeTime                  string                       `json:"lastUpgradeTime,omitempty"`
	Latitude                         string                       `json:"latitude,omitempty"`
	Location                         string                       `json:"location,omitempty"`
	Longitude                        string                       `json:"longitude,omitempty"`
	ModifiedBy                       string                       `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                       `json:"modifiedTime,omitempty"`
	Name                             string                       `json:"name,omitempty"`
	ProvisioningKeyID                string                       `json:"provisioningKeyId"`
	ProvisioningKeyName              string                       `json:"provisioningKeyName"`
	Platform                         string                       `json:"platform,omitempty"`
	PlatformDetail                   string                       `json:"platformDetail,omitempty"`
	PreviousVersion                  string                       `json:"previousVersion,omitempty"`
	PrivateIP                        string                       `json:"privateIp,omitempty"`
	PublicIP                         string                       `json:"publicIp,omitempty"`
	RuntimeOS                        string                       `json:"runtimeOS,omitempty"`
	SargeVersion                     string                       `json:"sargeVersion,omitempty"`
	EnrollmentCert                   map[string]interface{}       `json:"enrollmentCert,omitempty"`
	UpgradeAttempt                   string                       `json:"upgradeAttempt,omitempty"`
	UpgradeStatus                    string                       `json:"upgradeStatus,omitempty"`
	MicroTenantID                    string                       `json:"microtenantId,omitempty"`
	MicroTenantName                  string                       `json:"microtenantName,omitempty"`
	ZPNSubModuleUpgrade              []common.ZPNSubModuleUpgrade `json:"zpnSubModuleUpgradeList,omitempty"`
}

type AssistantVersion struct {
	ID                       string `json:"id,omitempty"`
	ApplicationStartTime     string `json:"applicationStartTime,omitempty"`
	AppConnectorGroupID      string `json:"appConnectorGroupId,omitempty"`
	BrokerId                 string `json:"brokerId,omitempty"`
	CreationTime             string `json:"creationTime,omitempty"`
	CtrlChannelStatus        string `json:"ctrlChannelStatus,omitempty"`
	CurrentVersion           string `json:"currentVersion,omitempty"`
	DisableAutoUpdate        bool   `json:"disableAutoUpdate,omitempty"`
	ExpectedVersion          string `json:"expectedVersion,omitempty"`
	LastBrokerConnectTime    string `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerDisconnectTime string `json:"lastBrokerDisconnectTime,omitempty"`
	LastUpgradedTime         string `json:"lastUpgradedTime,omitempty"`
	LoneWarrior              bool   `json:"loneWarrior,omitempty"`
	ModifiedBy               string `json:"modifiedBy,omitempty"`
	ModifiedTime             string `json:"modifiedTime,omitempty"`
	Latitude                 string `json:"latitude,omitempty"`
	Longitude                string `json:"longitude,omitempty"`
	MtunnelID                string `json:"mtunnelId,omitempty"`
	Platform                 string `json:"platform,omitempty"`
	PlatformDetail           string `json:"platformDetail,omitempty"`
	PreviousVersion          string `json:"previousVersion,omitempty"`
	PrivateIP                string `json:"privateIp,omitempty"`
	PublicIP                 string `json:"publicIp,omitempty"`
	RestartTimeInSec         string `json:"restartTimeInSec,omitempty"`
	RuntimeOS                string `json:"runtimeOS,omitempty"`
	SargeVersion             string `json:"sargeVersion,omitempty"`
	SystemStartTime          string `json:"systemStartTime,omitempty"`
	UpgradeAttempt           string `json:"upgradeAttempt,omitempty"`
	UpgradeStatus            string `json:"upgradeStatus,omitempty"`
	UpgradeNowOnce           bool   `json:"upgradeNowOnce,omitempty"`
}

type BulkDeleteRequest struct {
	IDs []string `json:"ids"`
}

// This function search the App Connector by ID
func Get(service *services.Service, appConnectorID string) (*AppConnector, *http.Response, error) {
	v := new(AppConnector)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+appConnectorEndpoint, appConnectorID)
	resp, err := service.Client.NewRequestDo("GET", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// This function search the App Connector by Name
func GetByName(service *services.Service, appConnectorName string) (*AppConnector, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppConnector](service.Client, relativeURL, common.Filter{Search: appConnectorName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, appConnectorName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no app connector named '%s' was found", appConnectorName)
}

func GetAll(service *services.Service) ([]AppConnector, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppConnector](service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

// Update Updates the App Connector details for the specified ID.
func Update(service *services.Service, appConnectorID string, appConnector AppConnector) (*AppConnector, *http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+appConnectorEndpoint, appConnectorID)
	_, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, appConnector, nil)
	if err != nil {
		return nil, nil, err
	}
	resource, resp, err := Get(service, appConnectorID)
	if err != nil {
		return nil, nil, err
	}
	return resource, resp, nil
}

// Delete Deletes the App Connector for the specified ID.
func Delete(service *services.Service, appConnectorID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appConnectorEndpoint, appConnectorID)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// BulkDelete Bulk deletes the App Connectors.
func BulkDelete(service *services.Service, appConnectorIDs []string) (*http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorEndpoint + "/bulkDelete"

	// Check if a microtenant ID is provided, else use the one from the service
	microTenantID := service.MicroTenantID()

	// Construct the filter with the microtenant ID if available
	filter := common.Filter{
		MicroTenantID: microTenantID,
	}

	resp, err := service.Client.NewRequestDo("POST", relativeURL, filter, BulkDeleteRequest{IDs: appConnectorIDs}, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
