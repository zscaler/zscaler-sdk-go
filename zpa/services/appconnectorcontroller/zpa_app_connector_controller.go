package appconnectorcontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	mgmtConfig           = "/mgmtconfig/v1/admin/customers/"
	appConnectorEndpoint = "/connector"
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
}

// This function search the App Connector by ID
func (service *Service) Get(appConnectorID string) (*AppConnector, *http.Response, error) {
	v := new(AppConnector)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+appConnectorEndpoint, appConnectorID)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// This function search the App Connector by Name
func (service *Service) GetByName(appConnectorName string) (*AppConnector, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorEndpoint
	list, resp, err := common.GetAllPagesGeneric[AppConnector](service.Client, relativeURL, "")
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

func (service *Service) GetAll() ([]AppConnector, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorEndpoint
	list, _, err := common.GetAllPagesGeneric[AppConnector](service.Client, relativeURL, "")
	if err != nil {
		return nil, err
	}
	return list, nil
}

type BulkDeleteRequest struct {
	IDs []string `json:"ids"`
}

// Update Updates the App Connector details for the specified ID.
func (service *Service) Update(appConnectorID string, appConnector AppConnector) (*AppConnector, *http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+appConnectorEndpoint, appConnectorID)
	_, err := service.Client.NewRequestDo("PUT", path, nil, appConnector, nil)
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
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// BulkDelete Bulk deletes the App Connectors.
func (service *Service) BulkDelete(appConnectorIDs []string) (*http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appConnectorEndpoint + "/bulkDelete"
	resp, err := service.Client.NewRequestDo("POST", relativeURL, nil, BulkDeleteRequest{IDs: appConnectorIDs}, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
