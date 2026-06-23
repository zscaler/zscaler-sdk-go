package dlp_global_options

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	webDlpGlobalOptionsEndpoint = "/zia/api/v1/webDlpGlobalOptions"
)

type WebDlpGlobal struct {
	Applications                []string                  `json:"applications,omitempty"`
	Urls                        []string                  `json:"urls,omitempty"`
	URLCategories               []common.IDNameExtensions `json:"urlCategories,omitempty"`
	ExemptUrlEncodedData        bool                      `json:"exemptUrlEncodedData,omitempty"`
	EnableNpkEdmTemplates       bool                      `json:"enableNpkEdmTemplates,omitempty"`
	EnableNpkEdmTemplatesForOrg bool                      `json:"enableNpkEdmTemplatesForOrg,omitempty"`
	EnableInlineDlpOcr          bool                      `json:"enableInlineDlpOcr,omitempty"`
	EnableCasbOcr               bool                      `json:"enableCasbOcr,omitempty"`
	EnableEmailDlpOcr           bool                      `json:"enableEmailDlpOcr,omitempty"`
	EnableEvaluateAllDlpRules   bool                      `json:"enableEvaluateAllDlpRules,omitempty"`
	EnableEdmPopularFormat      bool                      `json:"enableEdmPopularFormat,omitempty"`
	HttpGetCustomUrlCategories  []string                  `json:"httpGetCustomUrlCategories,omitempty"`
}

func GetDLPGlobalOptions(ctx context.Context, service *zscaler.Service) (*WebDlpGlobal, error) {
	var dlpGlobal WebDlpGlobal
	err := service.Client.Read(ctx, webDlpGlobalOptionsEndpoint, &dlpGlobal)
	if err != nil {
		return nil, err
	}
	return &dlpGlobal, nil
}

func UpdateDLPGlobalOptions(ctx context.Context, service *zscaler.Service, dlpGlobal WebDlpGlobal) (*WebDlpGlobal, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, webDlpGlobalOptionsEndpoint, dlpGlobal)
	if err != nil {
		return nil, nil, err
	}

	browserSettings, ok := resp.(*WebDlpGlobal)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected response type")
	}
	service.Client.GetLogger().Printf("[DEBUG] Updated DLP Global Options : %+v", dlpGlobal)
	return browserSettings, nil, nil
}
