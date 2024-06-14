package serviceedgecontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig                    = "/mgmtconfig/v1/admin/customers/"
	serviceEdgeControllerEndpoint = "/serviceEdge"
	scheduleEndpoint              = "/serviceEdgeSchedule"
)

type ServiceEdgeController struct {
	ApplicationStartTime             string                    `json:"applicationStartTime,omitempty"`
	ServiceEdgeGroupID               string                    `json:"serviceEdgeGroupId,omitempty"`
	ServiceEdgeGroupName             string                    `json:"serviceEdgeGroupName,omitempty"`
	ControlChannelStatus             string                    `json:"controlChannelStatus,omitempty"`
	CreationTime                     string                    `json:"creationTime,omitempty"`
	CtrlBrokerName                   string                    `json:"ctrlBrokerName,omitempty"`
	CurrentVersion                   string                    `json:"currentVersion,omitempty"`
	Description                      string                    `json:"description,omitempty"`
	Enabled                          bool                      `json:"enabled,omitempty"`
	ExpectedUpgradeTime              string                    `json:"expectedUpgradeTime,omitempty"`
	ExpectedVersion                  string                    `json:"expectedVersion,omitempty"`
	Fingerprint                      string                    `json:"fingerprint,omitempty"`
	ID                               string                    `json:"id,omitempty"`
	IPACL                            string                    `json:"ipAcl,omitempty"`
	IssuedCertID                     string                    `json:"issuedCertId,omitempty"`
	LastBrokerConnectTime            string                    `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerConnectTimeDuration    string                    `json:"lastBrokerConnectTimeDuration,omitempty"`
	LastBrokerDisconnectTime         string                    `json:"lastBrokerDisconnectTime,omitempty"`
	LastBrokerDisconnectTimeDuration string                    `json:"lastBrokerDisconnectTimeDuration,omitempty"`
	LastUpgradeTime                  string                    `json:"lastUpgradeTime,omitempty"`
	Latitude                         string                    `json:"latitude,omitempty"`
	Location                         string                    `json:"location,omitempty"`
	Longitude                        string                    `json:"longitude,omitempty"`
	ListenIPs                        string                    `json:"listenIps,omitempty"`
	ModifiedBy                       string                    `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                    `json:"modifiedTime,omitempty"`
	Name                             string                    `json:"name,omitempty"`
	ProvisioningKeyID                string                    `json:"provisioningKeyId"`
	ProvisioningKeyName              string                    `json:"provisioningKeyName"`
	Platform                         string                    `json:"platform,omitempty"`
	PreviousVersion                  string                    `json:"previousVersion,omitempty"`
	PrivateIP                        string                    `json:"privateIp,omitempty"`
	PublicIP                         string                    `json:"publicIp,omitempty"`
	PublishIPs                       string                    `json:"publishIps,omitempty"`
	SargeVersion                     string                    `json:"sargeVersion,omitempty"`
	EnrollmentCert                   map[string]interface{}    `json:"enrollmentCert,omitempty"`
	UpgradeAttempt                   string                    `json:"upgradeAttempt,omitempty"`
	UpgradeStatus                    string                    `json:"upgradeStatus,omitempty"`
	MicroTenantID                    string                    `json:"microtenantId,omitempty"`
	MicroTenantName                  string                    `json:"microtenantName,omitempty"`
	ZPNSubModuleUpgradeList          []ZPNSubModuleUpgradeList `json:"zpnSubModuleUpgradeList,omitempty"`
}

type ZPNSubModuleUpgradeList struct {
	ID              string `json:"id,omitempty"`
	CreationTime    string `json:"creationTime,omitempty"`
	CurrentVersion  string `json:"currentVersion,omitempty"`
	EntityGid       string `json:"entityGid,omitempty"`
	EntityType      string `json:"entityType,omitempty"`
	ExpectedVersion string `json:"expectedVersion,omitempty"`
	ModifiedBy      string `json:"modifiedBy,omitempty"`
	ModifiedTime    string `json:"modifiedTime,omitempty"`
	PreviousVersion string `json:"previousVersion,omitempty"`
	Role            string `json:"role,omitempty"`
	UpgradeStatus   string `json:"upgradeStatus,omitempty"`
	UpgradeTime     string `json:"upgradeTime,omitempty"`
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

func Get(service *services.Service, serviceEdgeID string) (*ServiceEdgeController, *http.Response, error) {
	v := new(ServiceEdgeController)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeControllerEndpoint, serviceEdgeID)
	resp, err := service.Client.NewRequestDo("GET", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(service *services.Service, serviceEdgeName string) (*ServiceEdgeController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + serviceEdgeControllerEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ServiceEdgeController](service.Client, relativeURL, common.Filter{Search: serviceEdgeName, MicroTenantID: service.MicroTenantID()})
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

func GetAll(service *services.Service) ([]ServiceEdgeController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + serviceEdgeControllerEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ServiceEdgeController](service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

type BulkDeleteRequest struct {
	IDs []string `json:"ids"`
}

// Update Updates the Service Edge details for the specified ID.
func Update(service *services.Service, serviceEdgeID string, serviceEdge ServiceEdgeController) (*ServiceEdgeController, *http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeControllerEndpoint, serviceEdgeID)
	_, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, serviceEdge, nil)
	if err != nil {
		return nil, nil, err
	}
	resource, resp, err := Get(service, serviceEdgeID)
	if err != nil {
		return nil, nil, err
	}
	return resource, resp, nil
}

// Delete Deletes the Service Edge for the specified ID.
func Delete(service *services.Service, serviceEdgeID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+serviceEdgeControllerEndpoint, serviceEdgeID)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// BulkDelete Bulk deletes the Service Edge.
func BulkDelete(service *services.Service, serviceEdgeIDs []string) (*http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + serviceEdgeControllerEndpoint + "/bulkDelete"
	resp, err := service.Client.NewRequestDo("POST", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, BulkDeleteRequest{IDs: serviceEdgeIDs}, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
