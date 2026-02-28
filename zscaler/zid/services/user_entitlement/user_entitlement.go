package user_entitlement

import (
	"context"
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zid/services/common"
)

const (
	entitlementEndpoint = "/admin/api/v1/users"
)

type Entitlements struct {
	Roles   []common.IDNameDisplayName `json:"roles,omitempty"`
	Scope   common.IDNameDisplayName   `json:"scope,omitempty"`
	Service Service                    `json:"service,omitempty"`
}

type Scope struct {
	Scope []common.IDNameDisplayName `json:"scope,omitempty"`
}

type Service struct {
	ID              string `json:"id,omitempty"`
	ServiceName     string `json:"serviceName,omitempty"`
	CloudName       string `json:"cloudName,omitempty"`
	CloudDomainName string `json:"cloudDomainName,omitempty"`
	OrgName         string `json:"orgName,omitempty"`
	OrgID           string `json:"orgId,omitempty"`
}

func GetAdminEntitlement(ctx context.Context, service *zscaler.Service, userID string) ([]Entitlements, error) {
	var adminEntitlements []Entitlements
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s/admin-entitlements", entitlementEndpoint, userID), &adminEntitlements)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning admin entitlements for user: %s", userID)
	return adminEntitlements, nil
}

func GetServiceEntitlement(ctx context.Context, service *zscaler.Service, userID string) ([]Service, error) {
	var serviceEntitlements []Service
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s/service-entitlements", entitlementEndpoint, userID), &serviceEntitlements)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning service entitlements for user: %s", userID)
	return serviceEntitlements, nil
}
