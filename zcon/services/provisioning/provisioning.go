package provisioning

import (
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v2/zcon/services/common"
)

const (
	apiKeysEndpoint = "/apiKeys"
)

type ProvisioningAPIKeys struct {
	// ID of the API key. This is used to regenerate the API key
	ID int `json:"id,omitempty"`

	// API key value (12 alphanumeric characters in length)
	KeyValue string `json:"keyValue,omitempty"`

	// List of functional areas to which this API key applies. This attribute is subject to change
	Permissions []string `json:"permissions,omitempty"`

	// Indicates whether API key is enabled
	Enabled bool `json:"enabled,omitempty"`

	// Last time API key was modified
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// Last user to modify API key
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// Not applicable to Cloud & Branch Connector
	PartnerUrl string `json:"partnerUrl,omitempty"`
}

func (service *Service) Get(apiKeyID int) (*ProvisioningAPIKeys, error) {
	var apiKey ProvisioningAPIKeys
	err := service.Client.Read(fmt.Sprintf("%s/%d", apiKeysEndpoint, apiKeyID), &apiKey)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Available API Key from Get: %d", apiKey.ID)
	return &apiKey, nil
}

func (service *Service) GetAll() ([]ProvisioningAPIKeys, error) {
	var apiKeys []ProvisioningAPIKeys
	err := common.ReadAllPages(service.Client, apiKeysEndpoint, &apiKeys)
	return apiKeys, err
}
