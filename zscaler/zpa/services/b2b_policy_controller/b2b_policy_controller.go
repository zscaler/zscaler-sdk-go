package b2b_policy_controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

const (
	mgmtConfig        = "/zpa/mgmtconfig/v1/admin/customers/"
	b2bPolicyEndpoint = "/policySet/rules/policyType/GLOBAL_POLICY/guest"
)

// PolicyRule is the B2B (guest) policy rule. It uses the exact same attributes
// as the shared policy rule defined in policysetcontrollerv2, so it is aliased
// rather than redefined to keep a single source of truth for the struct and
// all of its nested types.
type PolicyRule = policysetcontrollerv2.PolicyRule

// GetAll returns the B2B (guest) policy rules created by a partner for the
// given guestID. The endpoint does not support the microtenant filter, so an
// empty common.Filter is passed (the nil MicroTenantID is omitted from the
// query string). Pagination is handled by the shared ReadAll paginator.
func GetAll(ctx context.Context, service *zscaler.Service, guestID string) ([]PolicyRule, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+b2bPolicyEndpoint, guestID)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyRule](ctx, service.Client, relativeURL, common.Filter{})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
