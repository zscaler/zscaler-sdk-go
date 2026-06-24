package private_cloud

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig           = "/zpa/mgmtconfig/v1/admin/customers/"
	privateCloudEndpoint = "/privateCloud"
)

type PrivateCloudController struct {
	ModifiedTime            string                 `json:"modifiedTime,omitempty"`
	CreationTime            string                 `json:"creationTime,omitempty"`
	ModifiedBy              string                 `json:"modifiedBy,omitempty"`
	ID                      string                 `json:"id,omitempty"`
	Name                    string                 `json:"name,omitempty"`
	Description             string                 `json:"description,omitempty"`
	ReEnrollPeriod          string                 `json:"reEnrollPeriod,omitempty"`
	FireDrillEnabled        bool                   `json:"fireDrillEnabled,omitempty"`
	SitecPreferred          bool                   `json:"sitecPreferred,omitempty"`
	Enabled                 bool                   `json:"enabled,omitempty"`
	RemoteLss               bool                   `json:"remoteLss,omitempty"`
	ReadOnly                bool                   `json:"readOnly,omitempty"`
	ZscalerManaged          bool                   `json:"zscalerManaged,omitempty"`
	MicrotenantName         string                 `json:"microtenantName,omitempty"`
	AssistantGroupsIDs      []common.CommonSummary `json:"assistantGroupsIds,omitempty"`
	SiteControllerGroupIDs  []common.CommonSummary `json:"siteControllerGroupIds,omitempty"`
	SiemIDs                 []common.CommonSummary `json:"siemIds,omitempty"`
	PrivateExporterGroupIDs []common.CommonSummary `json:"privateExporterGroupIds,omitempty"`
	PrivateBrokerGroupIDs   []common.CommonSummary `json:"privateBrokerGroupIds,omitempty"`
	ZPNFireDrillSite        *ZPNFireDrillSite      `json:"zpnFireDrillSite,omitempty"`
}

type ZPNFireDrillSite struct {
	ID                        string `json:"id,omitempty"`
	ModifiedTime              string `json:"modifiedTime,omitempty"`
	CreationTime              string `json:"creationTime,omitempty"`
	ModifiedBy                string `json:"modifiedBy,omitempty"`
	MicrotenantID             string `json:"microtenantId,omitempty"`
	MicrotenantName           string `json:"microtenantName,omitempty"`
	FireDrillInterval         string `json:"fireDrillInterval,omitempty"`
	FireDrillIntervalTimeUnit string `json:"fireDrillIntervalTimeUnit,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, privatecloudID string) (*PrivateCloudController, *http.Response, error) {
	v := new(PrivateCloudController)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+privateCloudEndpoint, privatecloudID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, privateCloudName string) (*PrivateCloudController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + privateCloudEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PrivateCloudController](ctx, service.Client, relativeURL, common.Filter{Search: privateCloudName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, privateCloudName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no app connector group named '%s' was found", privateCloudName)
}

func Create(ctx context.Context, service *zscaler.Service, privateCloud PrivateCloudController) (*PrivateCloudController, *http.Response, error) {
	v := new(PrivateCloudController)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+privateCloudEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, privateCloud, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, privateCloudID string, privateCloud *PrivateCloudController) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+privateCloudEndpoint, privateCloudID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, privateCloud, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, appConnectorGroupID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+privateCloudEndpoint, appConnectorGroupID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]PrivateCloudController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + privateCloudEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PrivateCloudController](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
