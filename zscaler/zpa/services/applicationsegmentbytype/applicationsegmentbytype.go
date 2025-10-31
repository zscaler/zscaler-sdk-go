package applicationsegmentbytype

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig          = "/zpa/mgmtconfig/v1/admin/customers/"
	applicationEndpoint = "/application"
)

type AppSegmentBaseAppDto struct {
	ID                  string `json:"id,omitempty"`
	AppID               string `json:"appId,omitempty"`
	Name                string `json:"name,omitempty"`
	Enabled             bool   `json:"enabled"`
	Domain              string `json:"domain,omitempty"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	CertificateID       string `json:"certificateId,omitempty"`
	CertificateName     string `json:"certificateName,omitempty"`
	MicroTenantID       string `json:"microtenantId,omitempty"`
	MicroTenantName     string `json:"microtenantName,omitempty"`
}

func GetByApplicationType(ctx context.Context, service *zscaler.Service, appName, applicationType string, expandAll bool) ([]AppSegmentBaseAppDto, *http.Response, error) {
	validApplicationTypes := map[string]bool{
		"BROWSER_ACCESS":       true,
		"INSPECT":              true,
		"SECURE_REMOTE_ACCESS": true,
	}

	if !validApplicationTypes[applicationType] {
		return nil, nil, fmt.Errorf("invalid applicationType '%s'. Valid types are 'BROWSER_ACCESS', 'INSPECT', 'SECURE_REMOTE_ACCESS'", applicationType)
	}

	relativeURL := mgmtConfig + service.Client.GetCustomerID() + applicationEndpoint + "/getAppsByType"

	// Construct the URL with expandAll and applicationType parameters
	query := url.Values{}
	query.Set("expandAll", fmt.Sprintf("%t", expandAll))
	query.Set("applicationType", applicationType)
	if appName != "" {
		query.Set("search", appName)
	}

	constructedURL := relativeURL + "?" + query.Encode()
	log.Printf("Constructed URL: %s\n", constructedURL)

	// Construct the filter
	filter := common.Filter{
		MicroTenantID: service.MicroTenantID(),
	}
	if appName != "" {
		filter.Search = appName
	}

	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppSegmentBaseAppDto](ctx, service.Client, constructedURL, filter)
	if err != nil {
		return nil, nil, err
	}

	return list, resp, nil
}

func DeleteByApplicationType(ctx context.Context, service *zscaler.Service, applicationID, applicationType string) (*http.Response, error) {
	validApplicationTypes := map[string]bool{
		"BROWSER_ACCESS":       true,
		"INSPECT":              true,
		"SECURE_REMOTE_ACCESS": true,
	}

	if !validApplicationTypes[applicationType] {
		return nil, fmt.Errorf("invalid applicationType '%s'. Valid types are 'BROWSER_ACCESS', 'INSPECT', 'SECURE_REMOTE_ACCESS'", applicationType)
	}

	// Construct the URL with applicationID and applicationType parameter
	relativeURL := fmt.Sprintf("%s%s%s/%s/deleteAppByType", mgmtConfig, service.Client.GetCustomerID(), applicationEndpoint, applicationID)
	query := url.Values{}
	query.Set("applicationType", applicationType)

	constructedURL := relativeURL + "?" + query.Encode()
	log.Printf("Constructed DELETE URL: %s\n", constructedURL)

	// Execute DELETE request
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", constructedURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
