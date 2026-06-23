package webhooks

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	webhooksEndpoint = "/zia/api/v1/alertRuleConfiguration/webhooks"
)

type WebhookConfiguration struct {
	ID                 int                      `json:"id,omitempty"`
	Name               string                   `json:"name,omitempty"`
	UserName           string                   `json:"userName,omitempty"`
	Password           string                   `json:"password,omitempty"`
	AuthToken          string                   `json:"authToken,omitempty"`
	UrlText            string                   `json:"urlText,omitempty"`
	Status             bool                     `json:"status,omitempty"`
	LastTriggered      int                      `json:"lastTriggered,omitempty"`
	AuthenticationType string                   `json:"authenticationType,omitempty"`
	Deleted            bool                     `json:"deleted,omitempty"`
	LastModifiedTime   int                      `json:"lastModifiedTime,omitempty"`
	LastModifiedBy     *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`
}

func Create(ctx context.Context, service *zscaler.Service, webhookConfiguration *WebhookConfiguration) (*WebhookConfiguration, *http.Response, error) {
	resp, err := service.Client.Create(ctx, webhooksEndpoint, *webhookConfiguration)
	if err != nil {
		return nil, nil, err
	}

	createdWebhookConfiguration, ok := resp.(*WebhookConfiguration)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a webhook configuration pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new webhook configuration from create: %d", createdWebhookConfiguration.ID)
	return createdWebhookConfiguration, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, webhookID int, webhook *WebhookConfiguration) (*WebhookConfiguration, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", webhooksEndpoint, webhookID), *webhook)
	if err != nil {
		return nil, nil, err
	}
	updatedWebhookConfiguration, _ := resp.(*WebhookConfiguration)

	service.Client.GetLogger().Printf("[DEBUG]returning updates webhook configuration from update: %d", updatedWebhookConfiguration.ID)
	return updatedWebhookConfiguration, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, webhookID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", webhooksEndpoint, webhookID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]WebhookConfiguration, error) {
	var webHooks []WebhookConfiguration
	err := service.Client.Read(ctx, webhooksEndpoint, &webHooks)
	return webHooks, err
}

// TestWebhook tests a webhook configuration by sending a sample notification.
//
// The endpoint takes no path or query parameters; the webhook configuration is
// supplied in the request body. Any non-2xx HTTP response from the webhook
// endpoint is treated as a failure (redirects are not followed), which is
// surfaced as an error.
func TestWebhook(ctx context.Context, service *zscaler.Service, webhook *WebhookConfiguration) (*http.Response, error) {
	resp, err := service.Client.CreateWithNoContent(ctx, webhooksEndpoint+"/test", *webhook)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] webhook test notification sent successfully")
	return resp, nil
}
