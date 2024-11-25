package clienttypes

import (
	"context"
	"net/http"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	mgmtConfig          = "/zpa/mgmtconfig/v1/admin/customers/"
	clientTypesEndpoint = "/clientTypes"
)

type ClientTypes struct {
	ZPNClientTypeExplorer         string `json:"zpn_client_type_exporter"`
	ZPNClientTypeNoAuth           string `json:"zpn_client_type_exporter_noauth"`
	ZPNClientTypeBrowserIsolation string `json:"zpn_client_type_browser_isolation"`
	ZPNClientTypeMachineTunnel    string `json:"zpn_client_type_machine_tunnel"`
	ZPNClientTypeIPAnchoring      string `json:"zpn_client_type_ip_anchoring"`
	ZPNClientTypeEdgeConnector    string `json:"zpn_client_type_edge_connector"`
	ZPNClientTypeZAPP             string `json:"zpn_client_type_zapp"`
	ZPNClientTypeSlogger          string `json:"zpn_client_type_slogger"`
	ZPNClientTypeBranchConnector  string `json:"zpn_client_type_branch_connector"`
	ZPNClientTypePartner          string `json:"zpn_client_type_zapp_partner"`
}

func GetAllClientTypes(ctx context.Context, service *zscaler.Service) (*ClientTypes, *http.Response, error) {
	v := new(ClientTypes)
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + clientTypesEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
