package applicationsegmentbytype

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig              = "/mgmtconfig/v1/admin/customers/"
	applicationTypeEndpoint = "/application/getAppsByType"
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

func GetByApplicationType(service *services.Service, appName, applicationType string, expandAll bool) ([]AppSegmentBaseAppDto, *http.Response, error) {
	validApplicationTypes := map[string]bool{
		"BROWSER_ACCESS":       true,
		"INSPECT":              true,
		"SECURE_REMOTE_ACCESS": true,
	}

	if !validApplicationTypes[applicationType] {
		return nil, nil, fmt.Errorf("invalid applicationType '%s'. Valid types are 'BROWSER_ACCESS', 'INSPECT', 'SECURE_REMOTE_ACCESS'", applicationType)
	}

	relativeURL := mgmtConfig + service.Client.Config.CustomerID + applicationTypeEndpoint

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

	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AppSegmentBaseAppDto](service.Client, constructedURL, filter)
	if err != nil {
		return nil, nil, err
	}

	return list, resp, nil
}
