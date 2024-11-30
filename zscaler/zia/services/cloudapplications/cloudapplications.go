package cloudapplications

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	cloudAppPolicyEndpoint    = "/zia/api/v1/cloudApplications/policy"
	cloudAppSSLPolicyEndpoint = "/zia/api/v1/cloudApplications/sslPolicy"
)

type CloudApplications struct {
	// Application enum constant
	App string `json:"app,omitempty"`

	// Cloud application name
	AppName string `json:"appName,omitempty"`

	// Application category enum constant
	Parent string `json:"parent,omitempty"`

	// Name of the cloud application category
	ParentName string `json:"parentName,omitempty"`
}

func GetCloudApplicationPolicy(ctx context.Context, service *zscaler.Service, params map[string]interface{}) (interface{}, error) {
	// Build query parameters
	queryParams := ""
	if len(params) > 0 {
		q := url.Values{}
		if appClasses, ok := params["appClass"].([]string); ok {
			for _, appClass := range appClasses {
				q.Add("appClass", appClass)
			}
		}
		if groupResults, ok := params["groupResults"].(bool); ok {
			q.Set("groupResults", strconv.FormatBool(groupResults))
		}
		queryParams = "?" + q.Encode()
	}

	// Construct the full endpoint with optional query parameters
	endpoint := cloudAppPolicyEndpoint + queryParams

	// Fetch raw response from the API
	var rawResponse json.RawMessage
	err := service.Client.Read(ctx, endpoint, &rawResponse)
	if err != nil {
		return nil, err
	}

	// Attempt to unmarshal into filtered format (map)
	var filteredResponse map[string]int
	if err := json.Unmarshal(rawResponse, &filteredResponse); err == nil {
		return filteredResponse, nil
	}

	// Attempt to unmarshal into unfiltered format (list)
	var unfilteredResponse []CloudApplications
	if err := json.Unmarshal(rawResponse, &unfilteredResponse); err == nil {
		return unfilteredResponse, nil
	}

	// If both attempts fail, return an error
	return nil, fmt.Errorf("unexpected response format: %s", string(rawResponse))
}

func GetCloudApplicationSSLPolicy(ctx context.Context, service *zscaler.Service, params map[string]interface{}) (interface{}, error) {
	// Build query parameters
	queryParams := ""
	if len(params) > 0 {
		q := url.Values{}
		if appClasses, ok := params["appClass"].([]string); ok {
			for _, appClass := range appClasses {
				q.Add("appClass", appClass)
			}
		}
		if groupResults, ok := params["groupResults"].(bool); ok {
			q.Set("groupResults", strconv.FormatBool(groupResults))
		}
		queryParams = "?" + q.Encode()
	}

	// Construct the full endpoint with optional query parameters
	endpoint := cloudAppSSLPolicyEndpoint + queryParams

	// Fetch raw response from the API
	var rawResponse json.RawMessage
	err := service.Client.Read(ctx, endpoint, &rawResponse)
	if err != nil {
		return nil, err
	}

	// Attempt to unmarshal into filtered format (map)
	var filteredResponse map[string]int
	if err := json.Unmarshal(rawResponse, &filteredResponse); err == nil {
		return filteredResponse, nil
	}

	// Attempt to unmarshal into unfiltered format (list)
	var unfilteredResponse []CloudApplications
	if err := json.Unmarshal(rawResponse, &unfilteredResponse); err == nil {
		return unfilteredResponse, nil
	}

	// If both attempts fail, return an error
	return nil, fmt.Errorf("unexpected response format: %s", string(rawResponse))
}
