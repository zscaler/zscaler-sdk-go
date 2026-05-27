// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"encoding/json"
	"net/http"
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

func TestPolicySetController_GetPolicyRule_NotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "policySet", "ps-1", "rule", "missing")
	api.On("GET", path, common.NotFoundResponse())

	got, _, err := policysetcontroller.GetPolicyRule(context.Background(), api.Service, "ps-1", "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestPolicySetController_GetByPolicyType_NotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policyType := "ACCESS_POLICY"
	path := common.ZPAPath(api.CustomerID, "policySet", "policyType", policyType)
	api.On("GET", path, common.NotFoundResponse())

	got, _, err := policysetcontroller.GetByPolicyType(context.Background(), api.Service, policyType)
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestPolicySetController_GetByNameAndType_NotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policyType := "ACCESS_POLICY"
	path := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", policyType)
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]policysetcontroller.PolicyRule{
		{ID: "rule-001", Name: "Other Rule"},
	})))

	got, _, err := policysetcontroller.GetByNameAndType(context.Background(), api.Service, policyType, "missing-rule")
	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "no policy rule named :missing-rule found")
}

func TestPolicySetController_GetByNameAndTypes_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	wantName := "Target Rule"

	pathIsolation := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", "ISOLATION_POLICY")
	api.On("GET", pathIsolation, common.SuccessResponse(common.ZPAList([]policysetcontroller.PolicyRule{})))

	pathAccess := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", "ACCESS_POLICY")
	api.On("GET", pathAccess, common.SuccessResponse(common.ZPAList([]policysetcontroller.PolicyRule{
		{ID: "rule-001", Name: "Noise"},
		{ID: "rule-002", Name: wantName},
	})))

	got, _, err := policysetcontroller.GetByNameAndTypes(context.Background(), api.Service, []string{"ISOLATION_POLICY", "ACCESS_POLICY"}, wantName)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "rule-002", got.ID)
	assert.Equal(t, wantName, got.Name)
}

func TestPolicySetController_BulkReorder_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policySetType := "ACCESS_POLICY"

	policyTypePath := common.ZPAPath(api.CustomerID, "policySet", "policyType", policySetType)
	api.On("GET", policyTypePath, common.SuccessResponse(policysetcontroller.PolicySet{
		ID:   "policy-set-1",
		Name: "Access Policy",
	}))

	rulesPath := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", policySetType)
	api.On("GET", rulesPath, common.SuccessResponse(common.ZPAList([]policysetcontroller.PolicyRule{
		{ID: "rule-b", Name: "B"},
		{ID: "rule-a", Name: "A"},
		{ID: "rule-default", Name: "Default_Rule"},
	})))

	reorderPath := common.ZPAPath(api.CustomerID, "policySet", "policy-set-1", "reorder")
	api.OnFunc("PUT", reorderPath, func(req *http.Request, body []byte) common.MockResponse {
		var got []string
		require.NoError(t, json.Unmarshal(body, &got))
		require.Equal(t, []string{"rule-a", "rule-b", "rule-default"}, got)
		return common.NoContentResponse()
	})

	ruleOrders := map[string]int{"rule-a": 1, "rule-b": 2}
	resp, err := policysetcontroller.BulkReorder(context.Background(), api.Service, policySetType, ruleOrders)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestPolicySetController_BulkReorder_Error_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policySetType := "ACCESS_POLICY"

	policyTypePath := common.ZPAPath(api.CustomerID, "policySet", "policyType", policySetType)
	api.On("GET", policyTypePath, common.SuccessResponse(policysetcontroller.PolicySet{
		ID: "policy-set-1",
	}))

	rulesPath := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", policySetType)
	api.On("GET", rulesPath, common.SuccessResponse(common.ZPAList([]policysetcontroller.PolicyRule{
		{ID: "rule-1", Name: "R1"},
		{ID: "rule-def", Name: "Default_Rule"},
	})))

	reorderPath := common.ZPAPath(api.CustomerID, "policySet", "policy-set-1", "reorder")
	api.On("PUT", reorderPath, common.MockResponse{
		StatusCode: 400,
		Body:       `reorder rejected`,
	})

	resp, err := policysetcontroller.BulkReorder(context.Background(), api.Service, policySetType, map[string]int{"rule-1": 1})
	require.Error(t, err)
	assert.Nil(t, resp, "failed PUT surfaces as error before BulkReorder parses body")
}

func TestPolicySetController_UpdateRule_ConditionsNormalization_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policySetID := "ps-1"
	ruleID := "rule-456"
	putPath := common.ZPAPath(api.CustomerID, "policySet", policySetID, "rule", ruleID)

	api.OnFunc("PUT", putPath, func(req *http.Request, body []byte) common.MockResponse {
		var got policysetcontroller.PolicyRule
		require.NoError(t, json.Unmarshal(body, &got))
		require.NotNil(t, got.Conditions)
		assert.Len(t, got.Conditions, 2)
		assert.Empty(t, got.Conditions[0].Operands)
		assert.Len(t, got.Conditions[1].Operands, 1)
		assert.Equal(t, "", got.Conditions[1].Operands[0].Name, "operand Name should be cleared when non-empty")
		assert.Equal(t, "rhs-val", got.Conditions[1].Operands[0].RHS)
		return common.NoContentResponse()
	})

	rule := &policysetcontroller.PolicyRule{
		PolicySetID: policySetID,
		Conditions: []policysetcontroller.Conditions{
			{Operands: []policysetcontroller.Operands{}},
			{Operands: []policysetcontroller.Operands{
				{Name: "stripMe", RHS: "rhs-val", ObjectType: "OT"},
			}},
		},
	}

	resp, err := policysetcontroller.UpdateRule(context.Background(), api.Service, policySetID, ruleID, rule)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestPolicySetController_UpdateRule_NilConditionsBecomesEmptySlice_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policySetID := "ps-1"
	ruleID := "rule-456"
	putPath := common.ZPAPath(api.CustomerID, "policySet", policySetID, "rule", ruleID)

	api.OnFunc("PUT", putPath, func(req *http.Request, body []byte) common.MockResponse {
		var raw map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(body, &raw))
		condRaw, ok := raw["conditions"]
		require.True(t, ok, "conditions key should be present")
		var conds []policysetcontroller.Conditions
		require.NoError(t, json.Unmarshal(condRaw, &conds))
		assert.Len(t, conds, 0)
		return common.NoContentResponse()
	})

	rule := &policysetcontroller.PolicyRule{PolicySetID: policySetID, Conditions: nil}
	resp, err := policysetcontroller.UpdateRule(context.Background(), api.Service, policySetID, ruleID, rule)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
