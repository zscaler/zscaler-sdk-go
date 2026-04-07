package web_privacy

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	setWebPrivacyInfoEndpoint = "/zcc/papi/public/v1/setWebPrivacyInfo"
	getWebPrivacyInfoEndpoint = "/zcc/papi/public/v1/getWebPrivacyInfo"
)

type WebPrivacyInfo struct {
	ID                            string `json:"id"`
	Active                        string `json:"active"`
	CollectUserInfo               string `json:"collectUserInfo"`
	CollectMachineHostname        string `json:"collectMachineHostname"`
	CollectZdxLocation            string `json:"collectZdxLocation"`
	EnablePacketCapture           string `json:"enablePacketCapture"`
	DisableCrashlytics            string `json:"disableCrashlytics"`
	OverrideT2ProtocolSetting     string `json:"overrideT2ProtocolSetting"`
	RestrictRemotePacketCapture   string `json:"restrictRemotePacketCapture"`
	GrantAccessToZscalerLogFolder string `json:"grantAccessToZscalerLogFolder"`
	ExportLogsForNonAdmin         string `json:"exportLogsForNonAdmin"`
	EnableAutoLogSnippet          string `json:"enableAutoLogSnippet"`
	EnforceSecurePacUrls          string `json:"enforceSecurePacUrls"`
	EnableFQDNMatchForVpnBypasses string `json:"enableFQDNMatchForVpnBypasses"`
}

func GetWebPrivacyInfo(ctx context.Context, service *zscaler.Service) (*WebPrivacyInfo, error) {
	var privacyInfo WebPrivacyInfo

	resp, err := service.Client.NewZccRequestDo(ctx, "GET", getWebPrivacyInfoEndpoint, nil, nil, &privacyInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve web privacy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve web privacy info: received status code %d", resp.StatusCode)
	}

	return &privacyInfo, nil
}

func UpdateWebPrivacyInfo(ctx context.Context, service *zscaler.Service, info *WebPrivacyInfo) (*WebPrivacyInfo, error) {
	if info == nil {
		return nil, errors.New("web privacy info is required")
	}

	_, err := service.Client.NewZccRequestDo(ctx, "PUT", setWebPrivacyInfoEndpoint, nil, info, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update web privacy info: %w", err)
	}

	return GetWebPrivacyInfo(ctx, service)
}
