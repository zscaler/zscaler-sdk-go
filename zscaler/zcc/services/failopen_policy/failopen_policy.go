package failopen_policy

import (
	"context"
	"errors"
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	baseFailOpenPolicy = "/zcc/papi/public/v1/webFailOpenPolicy"
)

type WebFailOpenPolicy struct {
	Active                            string `json:"active"`
	CaptivePortalWebSecDisableMinutes int    `json:"captivePortalWebSecDisableMinutes"`
	CompanyID                         string `json:"companyId"`
	CreatedBy                         string `json:"createdBy"`
	EditedBy                          string `json:"editedBy"`
	EnableCaptivePortalDetection      int    `json:"enableCaptivePortalDetection"`
	EnableFailOpen                    int    `json:"enableFailOpen"`
	EnableStrictEnforcementPrompt     int    `json:"enableStrictEnforcementPrompt"`
	EnableWebSecOnProxyUnreachable    string `json:"enableWebSecOnProxyUnreachable"`
	EnableWebSecOnTunnelFailure       string `json:"enableWebSecOnTunnelFailure"`
	ID                                string `json:"id"`
	StrictEnforcementPromptDelayMins  int    `json:"strictEnforcementPromptDelayMinutes"`
	StrictEnforcementPromptMessage    string `json:"strictEnforcementPromptMessage"`
	TunnelFailureRetryCount           int    `json:"tunnelFailureRetryCount"`
}

func GetFailOpenPolicy(ctx context.Context, service *zscaler.Service, pageSize int) ([]WebFailOpenPolicy, error) {
	endpoint := fmt.Sprintf("%s/listByCompany", baseFailOpenPolicy)
	policies, err := common.ReadAllPages[WebFailOpenPolicy](ctx, service.Client, endpoint, common.QueryParams{}, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fail open policies: %w", err)
	}
	return policies, nil
}

func GetFailOpenPolicyByID(ctx context.Context, service *zscaler.Service, id string) (*WebFailOpenPolicy, error) {
	policies, err := GetFailOpenPolicy(ctx, service, 1000)
	if err != nil {
		return nil, err
	}
	for i := range policies {
		if policies[i].ID == id {
			return &policies[i], nil
		}
	}
	return nil, fmt.Errorf("fail open policy with ID %q not found", id)
}

func UpdateFailOpenPolicy(ctx context.Context, service *zscaler.Service, openPolicy *WebFailOpenPolicy) (*WebFailOpenPolicy, error) {
	if openPolicy == nil {
		return nil, errors.New("open policy is required")
	}

	// Construct the URL for the update endpoint
	url := fmt.Sprintf("%s/edit", baseFailOpenPolicy)

	// Initialize a variable to hold the response
	var updatedPolicy WebFailOpenPolicy

	// Make the PUT request to update the web policy
	_, err := service.Client.NewZccRequestDo(ctx, "PUT", url, nil, openPolicy, &updatedPolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to update web policy: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] returning web policy from update: %+v", updatedPolicy)
	return &updatedPolicy, nil
}
