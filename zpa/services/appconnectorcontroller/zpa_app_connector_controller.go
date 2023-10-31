package appconnectorcontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig                = "/mgmtconfig/v1/admin/customers/"
	appConnectorEndpoint      = "/connector"
	assistantScheduleEndpoint = "/assistantSchedule"
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
	ID                string `json:"id,omitempty"`
	CustomerId        string `json:"customerId,omitempty"`
	DeleteDisabled    bool   `json:"deleteDisabled,omitempty"`
	Enabled           bool   `json:"enabled,omitempty"`
	Frequency         string `json:"frequency,omitempty"`
	FrequencyInterval string `json:"frequencyInterval,omitempty"`
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

func (service *Service) GetAllAssistantSchedule() (*AssistantSchedule, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + assistantScheduleEndpoint

	// Make the GET request
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	// Read and unmarshal the response
	var schedule AssistantSchedule
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&schedule)
	if err != nil {
		return nil, nil, err
	}

	return &schedule, resp, nil
}

func (service *Service) CreateScheduleAssistant(scheduleAssistant *AssistantSchedule) (*AssistantSchedule, *http.Response, error) {
	v := new(AssistantSchedule)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+assistantScheduleEndpoint, common.Filter{MicroTenantID: service.microTenantID}, scheduleAssistant, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// PUT - /mgmtconfig/v1/admin/customers/{customerId}/assistantSchedule/{id}
func (service *Service) UpdateScheduleAssistant(scheduleAssistant *AssistantSchedule) (*AssistantSchedule, *http.Response, error) {
	// Construct the URL
	if scheduleAssistant.ID == "" {
		return nil, nil, errors.New("the AssistantSchedule provided must have a valid ID")
	}
	relativeURL := fmt.Sprintf("%s%s%s/%s", mgmtConfig, service.Client.Config.CustomerID, assistantScheduleEndpoint, scheduleAssistant.ID)

	// Define the response variable
	v := new(AssistantSchedule)

	// Send the PUT request
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, scheduleAssistant, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
