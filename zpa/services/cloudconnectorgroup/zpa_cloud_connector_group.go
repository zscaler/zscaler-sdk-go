package cloudconnectorgroup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig                  = "/mgmtconfig/v1/admin/customers/"
	cloudConnectorGroupEndpoint = "/cloudConnectorGroup"
)

type CloudConnectorGroup struct {
	CreationTime    string            `json:"creationTime,omitempty"`
	Description     string            `json:"description,omitempty"`
	CloudConnectors []CloudConnectors `json:"cloudConnectors,omitempty"`
	Enabled         bool              `json:"enabled,omitempty"`
	GeolocationID   string            `json:"geoLocationId,omitempty"`
	ID              string            `json:"id,omitempty"`
	ModifiedBy      string            `json:"modifiedBy,omitempty"`
	ModifiedTime    string            `json:"modifiedTime,omitempty"`
	Name            string            `json:"name,omitempty"`
	ZiaCloud        string            `json:"ziaCloud,omitempty"`
	ZiaOrgid        string            `json:"ziaOrgId,omitempty"`
}

type CloudConnectors struct {
	ID              string                 `json:"id,omitempty"`
	Name            string                 `json:"name,omitempty"`
	CreationTime    string                 `json:"creationTime,omitempty"`
	Description     string                 `json:"description,omitempty"`
	Enabled         bool                   `json:"enabled,omitempty"`
	Fingerprint     string                 `json:"fingerprint,omitempty"`
	IPACL           []string               `json:"ipAcl,omitempty"`
	IssuedCertID    string                 `json:"issuedCertId,omitempty"`
	ModifiedBy      string                 `json:"modifiedBy,omitempty"`
	ModifiedTime    string                 `json:"modifiedTime,omitempty"`
	SigningCert     map[string]interface{} `json:"signingCert,omitempty"`
	MicroTenantID   string                 `json:"microtenantId,omitempty"`
	MicroTenantName string                 `json:"microtenantName,omitempty"`
}

func Get(service *services.Service, cloudConnectorGroupID string) (*CloudConnectorGroup, *http.Response, error) {
	v := new(CloudConnectorGroup)
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + cloudConnectorGroupEndpoint + "/" + cloudConnectorGroupID
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(service *services.Service, cloudConnectorGroupName string) (*CloudConnectorGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + cloudConnectorGroupEndpoint
	list, resp, err := common.GetAllPagesGeneric[CloudConnectorGroup](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, cloudConnectorGroupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application named '%s' was found", cloudConnectorGroupName)
}

func GetAll(service *services.Service) ([]CloudConnectorGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + cloudConnectorGroupEndpoint
	list, resp, err := common.GetAllPagesGeneric[CloudConnectorGroup](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
