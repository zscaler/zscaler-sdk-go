package devices

import (
	"context"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	getDevicesEndpoint = "/zcc/papi/public/v1/getDevices"
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

func GetAll(ctx context.Context, service *zscaler.Service, username, osType string) ([]GetDevices, error) {
	queryParams := GetDevicesQueryParams{
		Username: username,
		OsType:   osType,
	}
	return common.ReadAllPages[GetDevices](ctx, service.Client, getDevicesEndpoint, queryParams, 1000)
}
