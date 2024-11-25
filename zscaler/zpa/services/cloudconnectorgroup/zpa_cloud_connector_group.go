package cloudconnectorgroup

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                  = "/zpa/mgmtconfig/v1/admin/customers/"
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

func Get(ctx context.Context, service *zscaler.Service, cloudConnectorGroupID string) (*CloudConnectorGroup, *http.Response, error) {
	v := new(CloudConnectorGroup)
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + cloudConnectorGroupEndpoint + "/" + cloudConnectorGroupID
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, cloudConnectorGroupName string) (*CloudConnectorGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + cloudConnectorGroupEndpoint
	list, resp, err := common.GetAllPagesGeneric[CloudConnectorGroup](ctx, service.Client, relativeURL, "")
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

func GetAll(ctx context.Context, service *zscaler.Service) ([]CloudConnectorGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + cloudConnectorGroupEndpoint
	list, resp, err := common.GetAllPagesGeneric[CloudConnectorGroup](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
