package lssconfigcontroller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	lssStatusCodesEndpoint = "lssConfig/statusCodes"
)

type LSSStatusCodes struct {
	ZPNAuthLog    map[string]interface{} `json:"zpn_auth_log"`
	ZPNAstAuthLog map[string]interface{} `json:"zpn_ast_auth_log"`
	ZPNTransLog   map[string]interface{} `json:"zpn_trans_log"`
	ZPNSysAuthLog map[string]interface{} `json:"zpn_sys_auth_log"`
}

func GetStatusCodes(ctx context.Context, service *zscaler.Service) (*LSSStatusCodes, *http.Response, error) {
	v := new(LSSStatusCodes)
	relativeURL := fmt.Sprintf(mgmtConfigTypesAndFormats + lssStatusCodesEndpoint)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	service.Client.GetLogger().Printf("[INFO] got LSSStatusCodes:%#v", v)
	return v, resp, nil
}
