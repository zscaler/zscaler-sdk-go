package user_entitlement

import (
	"context"
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ziam/services/common"
)

const (
	entitlementEndpoint = "/ziam/admin/api/v1/users"
)

// Entitlements is one administrative entitlement record: a set of roles, the
// scope they apply within, and the service they apply to.
//
// Scope and Service are pointers because they are genuinely optional in the
// response. A ZIAM-service entitlement carries no `scope` key at all, and as
// value types the absence was indistinguishable from `{"id": ""}` — a
// distinction a caller mapping this into Terraform state has to be able to
// make. Note also that `omitempty` has no effect on a struct value, so the tags
// these fields carried were inert.
type Entitlements struct {
	Roles   []common.IDNameDisplayName `json:"roles,omitempty"`
	Scope   *common.IDNameDisplayName  `json:"scope,omitempty"`
	Service *Service                   `json:"service,omitempty"`
}

// ServiceEntitlement is one element of the service entitlements response.
//
// The element is an object wrapping a `service` key rather than a bare service
// object, which mirrors Entitlements above minus the roles and scope. Reading
// the array straight into []Service — as this package used to — matched no keys
// and produced one zero-valued Service per entitlement, losing the entire
// response without an error.
type ServiceEntitlement struct {
	Service *Service `json:"service,omitempty"`
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

// GetServiceEntitlement retrieves the service entitlements of a user.
//
// Each element wraps a `service` object; see ServiceEntitlement for why that
// matters.
func GetServiceEntitlement(ctx context.Context, service *zscaler.Service, userID string) ([]ServiceEntitlement, error) {
	var serviceEntitlements []ServiceEntitlement
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s/service-entitlements", entitlementEndpoint, userID), &serviceEntitlements)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning service entitlements for user: %s", userID)
	return serviceEntitlements, nil
}
