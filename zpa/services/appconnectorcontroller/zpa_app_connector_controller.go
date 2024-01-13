package appconnectorcontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig           = "/mgmtconfig/v1/admin/customers/"
	appConnectorEndpoint = "/connector"
	scheduleEndpoint     = "/assistantSchedule"
)

type AppConnector struct {
	ApplicationStartTime             string                 `json:"applicationStartTime,omitempty"`
	AppConnectorGroupID              string                 `json:"appConnectorGroupId,omitempty"`
	AppConnectorGroupName            string                 `json:"appConnectorGroupName,omitempty"`
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
	ModifiedBy                       string                 `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                 `json:"modifiedTime,omitempty"`
	Name                             string                 `json:"name,omitempty"`
	ProvisioningKeyID                string                 `json:"provisioningKeyId"`
	ProvisioningKeyName              string                 `json:"provisioningKeyName"`
	Platform                         string                 `json:"platform,omitempty"`
	PreviousVersion                  string                 `json:"previousVersion,omitempty"`
	PrivateIP                        string                 `json:"privateIp,omitempty"`
	PublicIP                         string                 `json:"publicIp,omitempty"`
	SargeVersion                     string                 `json:"sargeVersion,omitempty"`
	EnrollmentCert                   map[string]interface{} `json:"enrollmentCert,omitempty"`
	UpgradeAttempt                   string                 `json:"upgradeAttempt,omitempty"`
	UpgradeStatus                    string                 `json:"upgradeStatus,omitempty"`
	MicroTenantID                    string                 `json:"microtenantId,omitempty"`
	MicroTenantName                  string                 `json:"microtenantName,omitempty"`
}

type AssistantSchedule struct {
	// The unique identifier for the App Connector auto deletion configuration for a customer. This field is only required for the PUT request to update the frequency of the App Connector Settings.
	ID string `json:"id,omitempty"`

	// The unique identifier of the ZPA tenant.
	CustomerID string `json:"customerId"`

	// Indicates if the App Connectors are included for deletion if they are in a disconnected state based on frequencyInterval and frequency values.
	DeleteDisabled bool `json:"deleteDisabled"`

	// Indicates if the setting for deleting App Connectors is enabled or disabled.
	Enabled bool `json:"enabled"`

	// The scheduled frequency at which the disconnected App Connectors are deleted.
	Frequency string `json:"frequency"`

	// The interval for the configured frequency value. The minimum supported value is 5.
	FrequencyInterval string `json:"frequencyInterval"`
}

// This function search the App Connector by ID
func (service *Service) Get(appConnectorID string) (*AppConnector, *http.Response, error) {
	v := new(AppConnector)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+appConnectorEndpoint, appConnectorID)
	resp, err := service.Client.NewRequestDo("GET", path, common.Filter{MicroTenantID: service.microTenantID}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// This function search the App Connector by Name
func (service *Service) GetByName(appConnectorName string) (*AppConnector, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppConnector](service.Client, relativeURL, common.Filter{Search: appConnectorName, MicroTenantID: service.microTenantID})
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

func (service *Service) GetAll() ([]AppConnector, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppConnector](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

type BulkDeleteRequest struct {
	IDs []string `json:"ids"`
}

// Update Updates the App Connector details for the specified ID.
func (service *Service) Update(appConnectorID string, appConnector AppConnector) (*AppConnector, *http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+appConnectorEndpoint, appConnectorID)
	_, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, appConnector, nil)
	if err != nil {
		return nil, nil, err
	}
	resource, resp, err := service.Get(appConnectorID)
	if err != nil {
		return nil, nil, err
	}
	return resource, resp, nil
}

// Delete Deletes the App Connector for the specified ID.
func (service *Service) Delete(appConnectorID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appConnectorEndpoint, appConnectorID)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// BulkDelete Bulk deletes the App Connectors.
func (service *Service) BulkDelete(appConnectorIDs []string) (*http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorEndpoint + "/bulkDelete"
	resp, err := service.Client.NewRequestDo("POST", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, BulkDeleteRequest{IDs: appConnectorIDs}, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Get a Configured App Connector schedule frequency.
func (service *Service) GetSchedule() (*AssistantSchedule, *http.Response, error) {
	v := new(AssistantSchedule)
	path := fmt.Sprintf("%v", mgmtConfig+service.Client.Config.CustomerID+scheduleEndpoint)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Configure a App Connector schedule frequency to delete the in active connectors with configured frequency.
func (service *Service) CreateSchedule(assistantSchedule AssistantSchedule) (*AssistantSchedule, *http.Response, error) {
	v := new(AssistantSchedule)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+scheduleEndpoint, nil, assistantSchedule, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) UpdateSchedule(schedulerID string, assistantSchedule *AssistantSchedule) (*http.Response, error) {
	// Validate FrequencyInterval
	validIntervals := map[string]bool{"5": true, "7": true, "14": true, "30": true, "60": true, "90": true}
	if _, valid := validIntervals[assistantSchedule.FrequencyInterval]; !valid {
		return nil, fmt.Errorf("invalid FrequencyInterval: %s", assistantSchedule.FrequencyInterval)
	}

	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+scheduleEndpoint, schedulerID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, nil, assistantSchedule, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
