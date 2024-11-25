package lssconfigcontroller

import (
	"context"
	"net/http"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	lssClientTypesEndpoint = "lssConfig/clientTypes"
)

type LSSClientTypes struct {
	ZPNClientTypeExporter      string `json:"zpn_client_type_exporter"`
	ZPNClientTypeMachineTunnel string `json:"zpn_client_type_machine_tunnel"`
	ZPNClientTypeIPAnchoring   string `json:"zpn_client_type_ip_anchoring"`
	ZPNClientTypeEdgeConnector string `json:"zpn_client_type_edge_connector"`
	ZPNClientTypeZAPP          string `json:"zpn_client_type_zapp"`
	ZPNClientTypeSlogger       string `json:"zpn_client_type_slogger,omitempty"`
}

func GetClientTypes(ctx context.Context, service *zscaler.Service) (*LSSClientTypes, *http.Response, error) {
	v := new(LSSClientTypes)
	relativeURL := mgmtConfigTypesAndFormats + lssClientTypesEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
