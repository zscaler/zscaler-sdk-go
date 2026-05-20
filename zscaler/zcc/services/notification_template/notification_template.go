package notification_template

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	notificationTemplateEndpointV2 = "/zcc/papi/public/v2/notification-templates"
)

// GetAllFilterOptions models the documented optional query parameters for
// GET /zcc/papi/public/v2/notification-templates. Pagination (skip/perPage)
// is handled by the pagination helper; callers only supply filters here.
type GetAllFilterOptions struct {
	// Keyword filters records by name (substring match).
	Keyword string
}

type NotificationTemplate struct {
	ID                      int                     `json:"id,omitempty"`
	Name                    string                  `json:"name,omitempty"`
	IsDefaultTemplate       bool                    `json:"isDefaultTemplate"`
	EnableClient            bool                    `json:"enableClient"`
	EnableZia               bool                    `json:"enableZia"`
	EnableAppUpdates        bool                    `json:"enableAppUpdates"`
	EnableServiceStatus     bool                    `json:"enableServiceStatus"`
	DurationInSeconds       int                     `json:"durationInSeconds,omitempty"`
	EnablePersistent        bool                    `json:"enablePersistent"`
	EnableDoNotDisturb      bool                    `json:"enableDoNotDisturb"`
	CreatedBy               int                     `json:"createdBy,omitempty"`
	EditedBy                int                     `json:"editedBy,omitempty"`
	ZIANotificationTemplate ZIANotificationTemplate `json:"ziaNotificationTemplate"`
	ZPANotificationTemplate ZPANotificationTemplate `json:"zpaNotificationTemplate"`
}

type ZIANotificationTemplate struct {
	EnableZiaFirewall      bool `json:"enableZiaFirewall"`
	EnableZiaFirewallPopup bool `json:"enableZiaFirewallPopup"`
	EnableZiaDNS           bool `json:"enableZiaDNS"`
	EnableZiaDNSPopup      bool `json:"enableZiaDNSPopup"`
	EnableZiaIPS           bool `json:"enableZiaIPS"`
	EnableZiaIPSPopup      bool `json:"enableZiaIPSPopup"`
	EnableZiaPersistent    bool `json:"enableZiaPersistent"`
}

type ZPANotificationTemplate struct {
	EnableDevicePostureFailure bool `json:"enableDevicePostureFailure"`
	EnableZpaReauth            bool `json:"enableZpaReauth"`
	ZpaReauthIntervalInMinutes int  `json:"zpaReauthIntervalInMinutes,omitempty"`
	DelayPostureFailureSeconds int  `json:"delayPostureFailureSeconds"`
}

func Get(ctx context.Context, service *zscaler.Service, templateID int) (*NotificationTemplate, error) {
	var notificationTemplate NotificationTemplate
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", notificationTemplateEndpointV2, templateID), &notificationTemplate)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning notification template from Get: %d", notificationTemplate.ID)
	return &notificationTemplate, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, templateName string) (*NotificationTemplate, error) {
	// Narrow the server-side search via keyword; final exact match is
	// still done client-side because keyword is substring, not equality.
	templates, err := GetAll(ctx, service, &GetAllFilterOptions{Keyword: templateName})
	if err != nil {
		return nil, err
	}
	for _, template := range templates {
		if strings.EqualFold(template.Name, templateName) {
			return &template, nil
		}
	}
	return nil, fmt.Errorf("no notification template found with name: %s", templateName)
}

func Create(ctx context.Context, service *zscaler.Service, template *NotificationTemplate) (*NotificationTemplate, *http.Response, error) {
	if template == nil {
		return nil, nil, errors.New("notification template is required")
	}

	var created NotificationTemplate
	resp, err := service.Client.NewZccRequestDo(ctx, "POST", notificationTemplateEndpointV2, nil, template, &created)
	if err != nil {
		return nil, resp, err
	}

	service.Client.GetLogger().Printf("[DEBUG] returning new notification template from create: %d", created.ID)
	return &created, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, templateID int, template *NotificationTemplate) (*NotificationTemplate, *http.Response, error) {
	if template == nil {
		return nil, nil, errors.New("notification template is required")
	}

	endpoint := fmt.Sprintf("%s/%d", notificationTemplateEndpointV2, templateID)
	var updated NotificationTemplate
	resp, err := service.Client.NewZccRequestDo(ctx, "PUT", endpoint, nil, template, &updated)
	if err != nil {
		return nil, resp, err
	}

	// The ZCC notification-templates PUT endpoint echoes the full template
	// body but omits `id` (unlike GET and the trusted-networks PUT). Pin the
	// path id onto the returned object so callers can rely on it.
	if updated.ID == 0 {
		updated.ID = templateID
	}

	service.Client.GetLogger().Printf("[DEBUG] returning updated notification template from update: %d", updated.ID)
	return &updated, resp, nil
}

func PartialUpdate(ctx context.Context, service *zscaler.Service, templateID int, template *NotificationTemplate) (*NotificationTemplate, *http.Response, error) {
	if template == nil {
		return nil, nil, errors.New("notification template is required")
	}

	endpoint := fmt.Sprintf("%s/%d", notificationTemplateEndpointV2, templateID)
	var updated NotificationTemplate
	resp, err := service.Client.NewZccRequestDo(ctx, "PATCH", endpoint, nil, template, &updated)
	if err != nil {
		return nil, resp, err
	}

	// PATCH responds with the same shape as PUT (no `id`); preserve the path id.
	if updated.ID == 0 {
		updated.ID = templateID
	}

	service.Client.GetLogger().Printf("[DEBUG] returning notification template from partial update: %d", updated.ID)
	return &updated, resp, nil
}

func Delete(ctx context.Context, service *zscaler.Service, templateID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", notificationTemplateEndpointV2, templateID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) ([]NotificationTemplate, error) {
	params := common.QueryParamsV2{}
	if opts != nil {
		params.Keyword = opts.Keyword
	}
	return common.ReadAllPagesV2[NotificationTemplate](ctx, service.Client, notificationTemplateEndpointV2, params, common.DefaultPageSize)
}
