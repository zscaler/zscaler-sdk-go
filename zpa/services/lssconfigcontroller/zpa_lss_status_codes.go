package lssconfigcontroller

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
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

func (service *Service) GetStatusCodes() (*LSSStatusCodes, *http.Response, error) {
	v := new(LSSStatusCodes)
	relativeURL := fmt.Sprintf(mgmtConfigTypesAndFormats + lssStatusCodesEndpoint)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	service.Client.Config.Logger.Printf("[INFO] got LSSStatusCodes:%#v", v)
	return v, resp, nil
}
