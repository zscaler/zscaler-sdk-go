package cloudapplications

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
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

func GetCloudApplicationPolicy(ctx context.Context, service *zscaler.Service, params map[string]interface{}) ([]CloudApplications, error) {
	queryParams := url.Values{}
	// log.Printf("[DEBUG] Received params: %v", params)

	if appClasses, ok := params["appClass"].([]interface{}); ok {
		for _, appClass := range appClasses {
			queryParams.Add("appClass", appClass.(string))
		}
	}
	if groupResults, ok := params["groupResults"].(bool); ok {
		queryParams.Set("groupResults", strconv.FormatBool(groupResults))
	}

	endpoint := fmt.Sprintf("%s?%s", cloudAppPolicyEndpoint, queryParams.Encode())
	// log.Printf("[DEBUG] Constructed endpoint: %s", endpoint)

	var results []CloudApplications
	err := common.ReadAllPages(ctx, service.Client, endpoint, &results)
	if err != nil {
		return nil, fmt.Errorf("error fetching cloud application policies: %w", err)
	}
	return results, nil
}

func GetCloudApplicationSSLPolicy(ctx context.Context, service *zscaler.Service, params map[string]interface{}) ([]CloudApplications, error) {
	queryParams := url.Values{}
	// log.Printf("[DEBUG] Received params: %v", params)

	if appClasses, ok := params["appClass"].([]interface{}); ok {
		for _, appClass := range appClasses {
			queryParams.Add("appClass", appClass.(string))
		}
	}
	if groupResults, ok := params["groupResults"].(bool); ok {
		queryParams.Set("groupResults", strconv.FormatBool(groupResults))
	}

	endpoint := fmt.Sprintf("%s?%s", cloudAppSSLPolicyEndpoint, queryParams.Encode())
	//log.Printf("[DEBUG] Constructed endpoint: %s", endpoint)

	var results []CloudApplications
	err := common.ReadAllPages(ctx, service.Client, endpoint, &results)
	if err != nil {
		return nil, fmt.Errorf("error fetching cloud application SSL policies: %w", err)
	}
	return results, nil
}
