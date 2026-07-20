package outbound_email_dlp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
)

const (
	emailDlpRulesEndpoint = "/zia/api/v1/emailDlpRules"
)

type OutboundEmailDlp struct {
	ID                       int                       `json:"id,omitempty"`
	Order                    int                       `json:"order,omitempty"`
	Name                     string                    `json:"name,omitempty"`
	Description              string                    `json:"description,omitempty"`
	State                    string                    `json:"state,omitempty"`
	Action                   string                    `json:"action,omitempty"`
	Groups                   []common.IDNameExtensions `json:"groups,omitempty"`
	Departments              []common.IDNameExtensions `json:"departments,omitempty"`
	Users                    []common.IDNameExtensions `json:"users,omitempty"`
	ExcludedGroups           []common.IDNameExtensions `json:"excludedGroups,omitempty"`
	ExcludedDepartments      []common.IDNameExtensions `json:"excludedDepartments,omitempty"`
	ExcludedUsers            []common.IDNameExtensions `json:"excludedUsers,omitempty"`
	TimeWindows              []common.IDNameExtensions `json:"timeWindows,omitempty"`
	DlpEngines               []common.IDNameExtensions `json:"dlpEngines,omitempty"`
	FileTypes                []string                  `json:"fileTypes,omitempty"`
	MinSize                  int                       `json:"minSize,omitempty"`
	WithoutContentInspection bool                      `json:"withoutContentInspection,omitempty"`
	Auditor                  *common.IDNameExtensions  `json:"auditor,omitempty"`
	ExternalAuditorEmail     string                    `json:"externalAuditorEmail,omitempty"`
	NotificationTemplate     *common.IDNameExtensions  `json:"notificationTemplate,omitempty"`
	LastModifiedTime         int                       `json:"lastModifiedTime,omitempty"`
	LastModifiedBy           *common.IDNameExtensions  `json:"lastModifiedBy,omitempty"`
	Receiver                 *dlp_web_rules.Receiver   `json:"receiver,omitempty"`
	Labels                   []common.IDNameExtensions `json:"labels,omitempty"`
	IncludedDomainProfiles   []common.IDNameExtensions `json:"includedDomainProfiles,omitempty"`
	Severity                 string                    `json:"severity,omitempty"`
	UserRiskScoreLevels      []string                  `json:"userRiskScoreLevels,omitempty"`
	EmailTenants             []common.IDNameExtensions `json:"emailTenants,omitempty"`
	ContentLocations         []string                  `json:"contentLocations,omitempty"`
	ParentRule               int                       `json:"parentRule,omitempty"`
	SubRules                 []OutboundEmailDlp        `json:"subRules,omitempty"`
	CustomHeader             string                    `json:"customHeader,omitempty"`
	EmailRecipientProfiles   []common.IDNameExtensions `json:"emailRecipientProfiles,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*OutboundEmailDlp, error) {
	var emailDlpPolicy OutboundEmailDlp
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", emailDlpRulesEndpoint, ruleID), &emailDlpPolicy)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning outbound email dlp policy from Get: %d", emailDlpPolicy.ID)
	return &emailDlpPolicy, nil
}

func GetLite(ctx context.Context, service *zscaler.Service) ([]OutboundEmailDlp, error) {
	var emailDlpPolicies []OutboundEmailDlp
	err := service.Client.Read(ctx, emailDlpRulesEndpoint+"/lite", &emailDlpPolicies)
	return emailDlpPolicies, err
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*OutboundEmailDlp, error) {
	emailDlpPolicies, err := GetAll(ctx, service, nil)
	if err != nil {
		return nil, err
	}
	for _, dlpPolicy := range emailDlpPolicies {
		if strings.EqualFold(dlpPolicy.Name, ruleName) {
			return &dlpPolicy, nil
		}
	}
	return nil, fmt.Errorf("no outbound email dlp policy found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, dlpRule *OutboundEmailDlp) (*OutboundEmailDlp, *http.Response, error) {
	resp, err := service.Client.Create(ctx, emailDlpRulesEndpoint, *dlpRule)
	if err != nil {
		return nil, nil, err
	}

	createdOutboundEmail, ok := resp.(*OutboundEmailDlp)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a outbound email dlp policypointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new outbound email dlp policyfrom create: %d", createdOutboundEmail.ID)
	return createdOutboundEmail, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, dlpRule *OutboundEmailDlp) (*OutboundEmailDlp, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", emailDlpRulesEndpoint, ruleID), *dlpRule)
	if err != nil {
		return nil, nil, err
	}
	updatedOutboundEmail, _ := resp.(*OutboundEmailDlp)

	service.Client.GetLogger().Printf("[DEBUG]returning updates outbound email dlp policy from update: %d", updatedOutboundEmail.ID)
	return updatedOutboundEmail, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, profileID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", emailDlpRulesEndpoint, profileID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAll retrieves the list of all Outbound Email DLP rules.
//
// The orgId parameter is optional; when provided, it is sent as the orgId query
// parameter to filter the results for the specified organization. Pass nil to
// omit it.
func GetAll(ctx context.Context, service *zscaler.Service, orgID *int) ([]OutboundEmailDlp, error) {
	var emailDlpPolicies []OutboundEmailDlp
	endpoint := emailDlpRulesEndpoint

	if orgID != nil {
		queryParams := url.Values{}
		queryParams.Set("orgId", strconv.Itoa(*orgID))
		endpoint += "?" + queryParams.Encode()
	}

	err := service.Client.Read(ctx, endpoint, &emailDlpPolicies)
	return emailDlpPolicies, err
}

// GetActions retrieves a mapping of supported Outbound Email DLP rule actions
// for the specified email tenant applications as a CSV file.
//
// The tenantIds parameter is required and is sent as a repeated query parameter,
// one entry per email tenant application ID. The response body is CSV and is
// returned as raw bytes.
func GetActions(ctx context.Context, service *zscaler.Service, tenantIDs []int) ([]byte, error) {
	queryParams := url.Values{}
	for _, id := range tenantIDs {
		queryParams.Add("tenantIds", strconv.Itoa(id))
	}

	endpoint := emailDlpRulesEndpoint + "/actions"
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	// Pass "" so the request Content-Type defaults to application/json; using
	// text/csv on a body-less GET can yield a 415 from the OneAPI gateway.
	return service.Client.ReadRaw(ctx, endpoint, "")
}
