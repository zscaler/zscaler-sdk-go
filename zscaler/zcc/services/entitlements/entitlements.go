package entitlements

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	getZdxGroupEndpoint    = "/zcc/papi/public/v1/getZdxGroupEntitlements"
	updateZdxGroupEndpoint = "/zcc/papi/public/v1/updateZdxGroupEntitlement"
	getZpaGroupEndpoint    = "/zcc/papi/public/v1/getZpaGroupEntitlements"
	updateZpaGroupEndpoint = "/zcc/papi/public/v1/updateZpaGroupEntitlement"
)

type ZdxGroupEntitlements struct {
	CollectZdxLocation        int           `json:"collectZdxLocation"`
	ComputeDeviceGroupsForZDX int           `json:"computeDeviceGroupsForZDX"`
	LogoutZCCForZDXService    int           `json:"logoutZCCForZDXService"`
	TotalCount                int           `json:"totalCount"`
	UpmDeviceGroupList        []DeviceGroup `json:"upmDeviceGroupList"`
	UpmEnableForAll           int           `json:"upmEnableForAll"`
	UpmGroupList              []DeviceGroup `json:"upmGroupList"`
}

type DeviceGroup struct {
	Active     int    `json:"active"`
	AuthType   string `json:"authType"`
	GroupID    int    `json:"groupId"`
	GroupName  string `json:"groupName"`
	UpmEnabled int    `json:"upmEnabled"`
}

type ZpaGroupEntitlements struct {
	ComputeDeviceGroupsForZPA int               `json:"computeDeviceGroupsForZPA"`
	DeviceGroupList           []DeviceGroupItem `json:"deviceGroupList"`
	GroupList                 []GroupListItem   `json:"groupList"`
	MachineTunEnabledForAll   int               `json:"machineTunEnabledForAll"`
	TotalCount                int               `json:"totalCount"`
	ZpaEnableForAll           int               `json:"zpaEnableForAll"`
}

type DeviceGroupItem struct {
	Active     int    `json:"active"`
	AuthType   string `json:"authType"`
	GroupID    int    `json:"groupId"`
	GroupName  string `json:"groupName"`
	ZpaEnabled int    `json:"zpaEnabled"`
}

type GroupListItem struct {
	Active     int    `json:"active"`
	AuthType   string `json:"authType"`
	GroupID    int    `json:"groupId"`
	GroupName  string `json:"groupName"`
	ZpaEnabled int    `json:"zpaEnabled"`
}

func GetZdxGroupEntitlements(ctx context.Context, service *zscaler.Service, search string, pageSize int) ([]ZdxGroupEntitlements, error) {
	// Construct query parameters dynamically
	queryParams := map[string]interface{}{}
	if search != "" {
		queryParams["search"] = search
	}

	// Leverage ReadAllPages to handle pagination
	return common.ReadAllPages[ZdxGroupEntitlements](ctx, service.Client, getZdxGroupEndpoint, queryParams, pageSize)
}

func UpdateZdxGroupEntitlements(ctx context.Context, service *zscaler.Service, updateZdxGroup *ZdxGroupEntitlements) (*ZdxGroupEntitlements, error) {
	if updateZdxGroup == nil {
		return nil, errors.New("updateZdxGroup is required")
	}

	// Marshal the DeviceCleanupInfo struct into JSON
	body, err := json.Marshal(updateZdxGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update Zdx Group request: %w", err)
	}

	// Make the PUT request
	resp, err := service.Client.NewZccRequestDo(ctx, "PUT", updateZdxGroupEndpoint, nil, bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update Zdx Group: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP response codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update Zdx Group: received status code %d", resp.StatusCode)
	}

	// Decode the response body into a DeviceCleanupInfo struct
	var response ZdxGroupEntitlements
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func GetZpaGroupEntitlements(ctx context.Context, service *zscaler.Service, search string, pageSize int) ([]ZpaGroupEntitlements, error) {
	// Construct query parameters dynamically
	queryParams := map[string]interface{}{}
	if search != "" {
		queryParams["search"] = search
	}

	// Leverage ReadAllPages to handle pagination
	return common.ReadAllPages[ZpaGroupEntitlements](ctx, service.Client, getZpaGroupEndpoint, queryParams, pageSize)
}

func UpdateZpaGroupEntitlements(ctx context.Context, service *zscaler.Service, updateZpaGroup *ZpaGroupEntitlements) (*ZpaGroupEntitlements, error) {
	if updateZpaGroup == nil {
		return nil, errors.New("updateZpaGroup is required")
	}

	body, err := json.Marshal(updateZpaGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update Zpa Group request: %w", err)
	}

	resp, err := service.Client.NewZccRequestDo(ctx, "PUT", updateZpaGroupEndpoint, nil, bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update Zpa Group: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update Zpa Group: received status code %d", resp.StatusCode)
	}

	var response ZpaGroupEntitlements
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
