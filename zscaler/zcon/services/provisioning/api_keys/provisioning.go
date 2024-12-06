package api_keys

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcon/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcon/services/common"
)

const (
	apiKeysEndpoint           = "/apiKeys"
	regenerateApiKeysEndpoint = "/regenerate"
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

func Get(ctx context.Context, service *services.Service, apiKeyID int) (*ProvisioningAPIKeys, error) {
	var apiKey ProvisioningAPIKeys
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", apiKeysEndpoint, apiKeyID), &apiKey)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Available API Key from Get: %d", apiKey.ID)
	return &apiKey, nil
}

func GetPartnerAPIKey(ctx context.Context, service *services.Service, apiKeyValue string, includePartnerKey bool) (*ProvisioningAPIKeys, error) {
	// Constructing the API endpoint URL
	url := fmt.Sprintf("%s?includePartnerKey=%t", apiKeysEndpoint, includePartnerKey)

	var apiKeys []ProvisioningAPIKeys
	err := service.Client.Read(ctx, url, &apiKeys)
	if err != nil {
		return nil, err
	}

	// Iterating through the API keys to find a match
	for _, key := range apiKeys {
		if key.KeyValue == apiKeyValue {
			return &key, nil
		}
	}

	return nil, fmt.Errorf("no partner api key found with key value: %s", apiKeyValue)
}

func GetAll(ctx context.Context, service *services.Service) ([]ProvisioningAPIKeys, error) {
	var apiKeys []ProvisioningAPIKeys
	err := common.ReadAllPages(ctx, service.Client, apiKeysEndpoint, &apiKeys)
	return apiKeys, err
}

func Create(ctx context.Context, service *services.Service, apiKeyValue *ProvisioningAPIKeys, includePartnerKey bool, keyId *int) (*ProvisioningAPIKeys, error) {
	// Handle nil apiKeyValue appropriately
	if apiKeyValue == nil {
		apiKeyValue = &ProvisioningAPIKeys{}
	}
	// Determine the API endpoint URL based on whether keyId is provided

	var url string
	if keyId != nil {
		// Regenerate API key
		url = fmt.Sprintf("%s/%d%s?includePartnerKey=%t", apiKeysEndpoint, *keyId, regenerateApiKeysEndpoint, includePartnerKey)
	} else {
		// Create new API key
		url = fmt.Sprintf("%s?includePartnerKey=%t", apiKeysEndpoint, includePartnerKey)
	}

	resp, err := service.Client.Create(ctx, url, *apiKeyValue)
	if err != nil {
		return nil, err
	}

	createdApiKeys, ok := resp.(*ProvisioningAPIKeys)
	if !ok {
		return nil, errors.New("object returned from API was not an API key pointer")
	}

	log.Printf("returning API key from create: %d", createdApiKeys.ID)
	return createdApiKeys, nil
}
