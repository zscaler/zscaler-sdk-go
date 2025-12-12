// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func TestPolicySetController_GetByPolicyType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	policyType := "ACCESS_POLICY"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/policySet/policyType/" + policyType

	server.On("GET", path, common.SuccessResponse(policysetcontroller.PolicySet{
		ID:         "policy-123",
		Name:       "Access Policy",
		PolicyType: policyType,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := policysetcontroller.GetByPolicyType(context.Background(), service, policyType)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, policyType, result.PolicyType)
}

func TestPolicySetController_GetAllByType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	policyType := "ACCESS_POLICY"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/policySet/rules/policyType/" + policyType

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []policysetcontroller.PolicyRule{{ID: "rule-001"}, {ID: "rule-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := policysetcontroller.GetAllByType(context.Background(), service, policyType)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
