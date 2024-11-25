package dlp_notification_templates

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	dlpNotificationTemplatesEndpoint = "/zia/api/v1/dlpNotificationTemplates"
)

type DlpNotificationTemplates struct {
	// The unique identifier for a DLP notification template.
	ID int `json:"id"`

	// The DLP notification template name.
	Name string `json:"name,omitempty"`

	// The Subject line that is displayed within the DLP notification email.
	Subject string `json:"subject,omitempty"`

	// If set to true, the content that is violation is attached to the DLP notification email.
	AttachContent bool `json:"attachContent,omitempty"`

	// The template for the plain text UTF-8 message body that must be displayed in the DLP notification email.
	PlainTextMessage string `json:"plainTextMessage,omitempty"`

	// The template for the HTML message body that must be displayed in the DLP notification email.
	HtmlMessage string `json:"htmlMessage,omitempty"`

	// The template for the HTML message body that must be displayed in the DLP notification email.
	TLSEnabled bool `json:"tlsEnabled,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, dlpTemplateID int) (*DlpNotificationTemplates, error) {
	var dlpTemplates DlpNotificationTemplates
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", dlpNotificationTemplatesEndpoint, dlpTemplateID), &dlpTemplates)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning dlp notification template from Get: %d", dlpTemplates.ID)
	return &dlpTemplates, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, templateName string) (*DlpNotificationTemplates, error) {
	var dlpTemplates []DlpNotificationTemplates
	err := common.ReadAllPages(ctx, service.Client, dlpNotificationTemplatesEndpoint, &dlpTemplates)
	if err != nil {
		return nil, err
	}
	for _, template := range dlpTemplates {
		if strings.EqualFold(template.Name, templateName) {
			return &template, nil
		}
	}
	return nil, fmt.Errorf("no dictionary found with name: %s", templateName)
}

func Create(ctx context.Context, service *zscaler.Service, dlpTemplateID *DlpNotificationTemplates) (*DlpNotificationTemplates, *http.Response, error) {
	resp, err := service.Client.Create(ctx, dlpNotificationTemplatesEndpoint, *dlpTemplateID)
	if err != nil {
		return nil, nil, err
	}

	createdDlpTemplate, ok := resp.(*DlpNotificationTemplates)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a dlp dictionary pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new dlp notification template from create: %d", createdDlpTemplate.ID)
	return createdDlpTemplate, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, dlpTemplateID int, dlpTemplates *DlpNotificationTemplates) (*DlpNotificationTemplates, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", dlpNotificationTemplatesEndpoint, dlpTemplateID), *dlpTemplates)
	if err != nil {
		return nil, nil, err
	}
	updatedDlpTemplate, _ := resp.(*DlpNotificationTemplates)

	service.Client.GetLogger().Printf("[DEBUG]returning updates from dlp notification template from update: %d", updatedDlpTemplate.ID)
	return updatedDlpTemplate, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, dlpTemplateID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", dlpNotificationTemplatesEndpoint, dlpTemplateID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]DlpNotificationTemplates, error) {
	var dlpTemplates []DlpNotificationTemplates
	err := common.ReadAllPages(ctx, service.Client, dlpNotificationTemplatesEndpoint, &dlpTemplates)
	return dlpTemplates, err
}
