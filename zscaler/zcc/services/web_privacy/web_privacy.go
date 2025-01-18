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
	CollectMachineHostname        string `json:"collectMachineHostname"`
	CollectUserInfo               string `json:"collectUserInfo"`
	CollectZdxLocation            string `json:"collectZdxLocation"`
	DisableCrashlytics            string `json:"disableCrashlytics"`
	EnablePacketCapture           string `json:"enablePacketCapture"`
	ExportLogsForNonAdmin         string `json:"exportLogsForNonAdmin"`
	GrantAccessToZscalerLogFolder string `json:"grantAccessToZscalerLogFolder"`
	OverrideT2ProtocolSetting     string `json:"overrideT2ProtocolSetting"`
	RestrictRemotePacketCapture   string `json:"restrictRemotePacketCapture"`
}

func GetWebPrivacyInfo(ctx context.Context, service *zscaler.Service) (*WebPrivacyInfo, error) {
	// Initialize a variable to hold the response
	var privacyInfo WebPrivacyInfo

	// Make the GET request
	resp, err := service.Client.NewZccRequestDo(ctx, "GET", getWebPrivacyInfoEndpoint, nil, nil, &privacyInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve web privacy: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP response codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve web privacy info: received status code %d", resp.StatusCode)
	}

	return &privacyInfo, nil
}

func UpdatePrivacyInfo(ctx context.Context, service *zscaler.Service, info *WebPrivacyInfo) (*WebPrivacyInfo, error) {
	if info == nil {
		return nil, errors.New("web policy is required")
	}

	// Construct the URL for the update endpoint
	url := setWebPrivacyInfoEndpoint

	// Initialize a variable to hold the response
	var updatedPrivacyInfo WebPrivacyInfo

	// Make the PUT request to update the web privacy info
	_, err := service.Client.NewRequestDo(ctx, "PUT", url, nil, info, &updatedPrivacyInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to update web privacy info: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning web privacy info from update: %s", updatedPrivacyInfo)
	return &updatedPrivacyInfo, nil
}
