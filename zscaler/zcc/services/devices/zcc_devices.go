package devices

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	getDevicesEndpoint       = "/zcc/papi/public/v1/getDevices"
	getDeviceDetailsEndpoint = "/zcc/papi/public/v1/getDeviceDetails"
	getDeviceCleanupEndpoint = "/zcc/papi/public/v1/getDeviceCleanupInfo"
	setDeviceCleanupEndpoint = "/zcc/papi/public/v1/setDeviceCleanupInfo"
)

type GetDevices struct {
	AgentVersion            string `json:"agentVersion"`
	CompanyName             string `json:"companyName"`
	ConfigDownloadTime      string `json:"config_download_time"`
	DeregistrationTimestamp string `json:"deregistrationTimestamp"`
	Detail                  string `json:"detail"`
	DownloadCount           int    `json:"download_count"`
	HardwareFingerprint     string `json:"hardwareFingerprint"`
	KeepAliveTime           string `json:"keepAliveTime"`
	LastSeenTime            string `json:"last_seen_time"`
	MacAddress              string `json:"macAddress"`
	MachineHostname         string `json:"machineHostname"`
	Manufacturer            string `json:"manufacturer"`
	OsVersion               string `json:"osVersion"`
	Owner                   string `json:"owner"`
	PolicyName              string `json:"policyName"`
	RegistrationState       string `json:"registrationState"`
	RegistrationTime        string `json:"registration_time"`
	State                   int    `json:"state"`
	TunnelVersion           string `json:"tunnelVersion"`
	Type                    int    `json:"type"`
	Udid                    string `json:"udid"`
	UpmVersion              string `json:"upmVersion"`
	User                    string `json:"user"`
	VpnState                int    `json:"vpnState"`
	ZappArch                string `json:"zappArch"`
}

type GetDevicesQueryParams struct {
	Username string `url:"username,omitempty"`
	OsType   string `url:"osType,omitempty"`
}

type DeviceCleanupInfo struct {
	ID                    string `json:"id"`
	Active                string `json:"active"`
	AutoPurgeDays         string `json:"autoPurgeDays"`
	AutoRemovalDays       string `json:"autoRemovalDays"`
	CompanyID             string `json:"companyId"`
	CreatedBy             string `json:"createdBy"`
	DeviceExceedLimit     string `json:"deviceExceedLimit"`
	EditedBy              string `json:"editedBy"`
	ForceRemoveType       string `json:"forceRemoveType"`
	ForceRemoveTypeString string `json:"forceRemoveTypeString"`
}

type DeviceDetails struct {
	AgentVersion        string `json:"agent_version"`
	Carrier             string `json:"carrier"`
	ConfigDownloadTime  string `json:"config_download_time"`
	DeregistrationTime  string `json:"deregistration_time"`
	DevicePolicyName    string `json:"devicePolicyName"`
	DeviceLocale        string `json:"device_locale"`
	DownloadCount       int    `json:"download_count"`
	ExternalModel       string `json:"external_model"`
	HardwareFingerprint string `json:"hardwareFingerprint"`
	KeepAliveTime       string `json:"keep_alive_time"`
	LastSeenTime        string `json:"last_seen_time"`
	MacAddress          string `json:"mac_address"`
	MachineHostname     string `json:"machineHostname"`
	Manufacturer        string `json:"manufacturer"`
	OSVersion           string `json:"os_version"`
	Owner               string `json:"owner"`
	RegistrationTime    string `json:"registration_time"`
	Rooted              int    `json:"rooted"`
	State               string `json:"state"`
	TunnelVersion       string `json:"tunnelVersion"`
	Type                string `json:"type"`
	UniqueID            string `json:"unique_id"`
	UpmVersion          string `json:"upmVersion"`
	UserName            string `json:"user_name"`
	ZadVersion          string `json:"zadVersion"`
	ZappArch            string `json:"zappArch"`
}

func GetAll(ctx context.Context, service *zscaler.Service, username, osType string) ([]GetDevices, error) {
	queryParams := GetDevicesQueryParams{
		Username: username,
		OsType:   osType,
	}
	return common.ReadAllPages[GetDevices](ctx, service.Client, getDevicesEndpoint, queryParams, 1000)
}

func GetDeviceCleanupInfo(ctx context.Context, service *zscaler.Service) (*DeviceCleanupInfo, error) {
	// Make the GET request
	resp, err := service.Client.NewZccRequestDo(ctx, "GET", getDeviceCleanupEndpoint, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve device cleanup info: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP response codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve device cleanup info: received status code %d", resp.StatusCode)
	}

	// Parse the response body into DeviceCleanupInfo struct
	var cleanupInfo DeviceCleanupInfo
	err = json.NewDecoder(resp.Body).Decode(&cleanupInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to decode device cleanup info response: %w", err)
	}

	return &cleanupInfo, nil
}

func GetDeviceDetails(ctx context.Context, service *zscaler.Service, username, udid string) ([]DeviceDetails, error) {
	// Construct query parameters with optional username and udid
	queryParams := url.Values{}
	if username != "" {
		queryParams.Set("username", username)
	}
	if udid != "" {
		queryParams.Set("udid", udid)
	}

	// Construct the full endpoint with query parameters
	fullEndpoint := getDeviceDetailsEndpoint
	if len(queryParams) > 0 {
		fullEndpoint = fmt.Sprintf("%s?%s", fullEndpoint, queryParams.Encode())
	}

	// Fetch device details
	var deviceDetails []DeviceDetails
	err := service.Client.Read(ctx, fullEndpoint, &deviceDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch device details: %w", err)
	}

	return deviceDetails, nil
}

func SetDeviceCleanupInfo(ctx context.Context, service *zscaler.Service, cleanupInfo *DeviceCleanupInfo) (*DeviceCleanupInfo, error) {
	if cleanupInfo == nil {
		return nil, errors.New("cleanupInfo is required")
	}

	// Marshal the DeviceCleanupInfo struct into JSON
	body, err := json.Marshal(cleanupInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal device cleanup info request: %w", err)
	}

	// Make the PUT request
	resp, err := service.Client.NewZccRequestDo(ctx, "PUT", setDeviceCleanupEndpoint, nil, bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to set device cleanup info: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP response codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to set device cleanup info: received status code %d", resp.StatusCode)
	}

	// Decode the response body into a DeviceCleanupInfo struct
	var response DeviceCleanupInfo
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
