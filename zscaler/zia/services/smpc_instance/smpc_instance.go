package smpc_instance

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	smpcInstanceEndpoint   = "/zia/api/v1/smpcInstance"
	smpcInstanceDCEndpoint = smpcInstanceEndpoint + "/dc"
)

type SmpcInstance struct {
	Data []Data `json:"data,omitempty"`
}

type Data struct {
	ID                          int    `json:"id,omitempty"`
	DatacenterName              string `json:"datacenterName,omitempty"`
	City                        string `json:"city,omitempty"`
	CountryCode                 string `json:"countryCode,omitempty"`
	ZscmIsolationEnabled        int    `json:"zscmIsolationEnabled,omitempty"`
	ZscmIsolationEnabledApiFlag bool   `json:"zscmIsolationEnabledApiFlag,omitempty"`
	ConnectionStatus            string `json:"connectionStatus,omitempty"`
	DCIsolationStatus           string `json:"dcIsolationStatus,omitempty"`
	Nodes                       []Node `json:"nodes,omitempty"`
}

type Node struct {
	ID                    int    `json:"id,omitempty"`
	NodeName              string `json:"nodeName,omitempty"`
	InstanceId            int    `json:"instanceId,omitempty"`
	AdminIsolationEnabled bool   `json:"adminIsolationEnabled,omitempty"`
	ConnectionState       int    `json:"connectionState,omitempty"`
	LastTimestamp         int    `json:"lastTimestamp,omitempty"`
}

// Update updates the data center and instance details based on the specified
// data center IDs.
//
// Both parameters are required by the API:
//   - isEnabled: set to true to enable DC isolation or false to disable it
//     (sent as the isEnabled query parameter).
//   - dcIDs: the list of data center IDs to update (sent as the request body,
//     a JSON array of integers).
func Update(ctx context.Context, service *zscaler.Service, isEnabled bool, dcIDs []int) (*SmpcInstance, *http.Response, error) {
	queryParams := url.Values{}
	queryParams.Set("isEnabled", strconv.FormatBool(isEnabled))
	endpoint := smpcInstanceDCEndpoint + "?" + queryParams.Encode()

	respBody, err := service.Client.UpdateWithSlicePayload(ctx, endpoint, dcIDs)
	if err != nil {
		return nil, nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] updated smpc instance dc isolation (isEnabled=%t) for ids: %v", isEnabled, dcIDs)

	if len(respBody) == 0 {
		return nil, nil, nil
	}

	var updated SmpcInstance
	if err := json.Unmarshal(respBody, &updated); err != nil {
		return nil, nil, err
	}
	return &updated, nil, nil
}

// GetAll returns all datacenters and their instance details. The API returns a
// single object with a flat "data" list; it does not support pagination.
func GetAll(ctx context.Context, service *zscaler.Service) ([]Data, error) {
	var instance SmpcInstance
	err := service.Client.Read(ctx, smpcInstanceEndpoint, &instance)
	if err != nil {
		return nil, err
	}
	return instance.Data, nil
}
