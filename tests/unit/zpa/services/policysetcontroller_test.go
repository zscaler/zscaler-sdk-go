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

func TestPolicySetController_GetPolicyRule_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	policySetID := "policy-set-123"
	ruleID := "rule-456"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/policySet/" + policySetID + "/rule/" + ruleID

	server.On("GET", path, common.SuccessResponse(policysetcontroller.PolicyRule{
		ID:          ruleID,
		Name:        "Test Rule",
		Description: "Test description",
		Action:      "ALLOW",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := policysetcontroller.GetPolicyRule(context.Background(), service, policySetID, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
	assert.Equal(t, "Test Rule", result.Name)
}

func TestPolicySetController_CreateRule_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	policySetID := "policy-set-123"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/policySet/" + policySetID + "/rule"

	server.On("POST", path, common.SuccessResponse(policysetcontroller.PolicyRule{
		ID:   "new-rule-789",
		Name: "New Rule",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newRule := &policysetcontroller.PolicyRule{
		Name:        "New Rule",
		Description: "New rule description",
		Action:      "ALLOW",
		PolicySetID: policySetID,
	}

	result, _, err := policysetcontroller.CreateRule(context.Background(), service, newRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-rule-789", result.ID)
}

func TestPolicySetController_UpdateRule_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	policySetID := "policy-set-123"
	ruleID := "rule-456"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/policySet/" + policySetID + "/rule/" + ruleID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateRule := &policysetcontroller.PolicyRule{
		ID:          ruleID,
		Name:        "Updated Rule",
		Description: "Updated description",
	}

	resp, err := policysetcontroller.UpdateRule(context.Background(), service, policySetID, ruleID, updateRule)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPolicySetController_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	policySetID := "policy-set-123"
	ruleID := "rule-456"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/policySet/" + policySetID + "/rule/" + ruleID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := policysetcontroller.Delete(context.Background(), service, policySetID, ruleID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPolicySetController_GetByNameAndType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	policyType := "ACCESS_POLICY"
	ruleName := "Production Rule"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/policySet/rules/policyType/" + policyType

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []policysetcontroller.PolicyRule{
			{ID: "rule-001", Name: "Other Rule"},
			{ID: "rule-002", Name: ruleName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := policysetcontroller.GetByNameAndType(context.Background(), service, policyType, ruleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "rule-002", result.ID)
	assert.Equal(t, ruleName, result.Name)
}

func TestPolicySetController_Reorder_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	policySetID := "policy-set-123"
	ruleID := "rule-456"
	order := 1
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/policySet/" + policySetID + "/rule/" + ruleID + "/reorder/1"

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := policysetcontroller.Reorder(context.Background(), service, policySetID, ruleID, order)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
