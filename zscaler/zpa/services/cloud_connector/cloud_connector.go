package cloud_connector

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig             = "/zpa/mgmtconfig/v1/admin/customers/"
	cloudConnectorEndpoint = "/cloudConnector"
)

type CloudConnector struct {
	CreationTime           string                 `json:"creationTime,omitempty"`
	Description            string                 `json:"description,omitempty"`
	EdgeConnectorGroupID   string                 `json:"edgeConnectorGroupId,omitempty"`
	EdgeConnectorGroupName string                 `json:"edgeConnectorGroupName,omitempty"`
	Enabled                bool                   `json:"enabled,omitempty"`
	Fingerprint            string                 `json:"fingerprint,omitempty"`
	ID                     string                 `json:"id,omitempty"`
	IpAcl                  []string               `json:"ipAcl,omitempty"`
	IssuedCertID           string                 `json:"issuedCertId,omitempty"`
	ModifiedBy             string                 `json:"modifiedBy,omitempty"`
	ModifiedTime           int                    `json:"modifiedTime,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	EnrollmentCert         map[string]interface{} `json:"enrollmentCert,omitempty"`
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]CloudConnector, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + cloudConnectorEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[CloudConnector](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, cloudConnectorName string) (*CloudConnector, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + cloudConnectorEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[CloudConnector](ctx, service.Client, relativeURL, common.Filter{Search: cloudConnectorName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, cloudConnectorName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no cloud connector group named '%s' was found", cloudConnectorName)
}
