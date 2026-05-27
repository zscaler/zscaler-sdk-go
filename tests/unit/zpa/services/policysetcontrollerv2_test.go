package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

func TestPolicySetControllerV2_GetByPolicyType_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policyType := "ACCESS_POLICY"
	path := common.ZPAPath(api.CustomerID, "policySet", "policyType", policyType)

	api.On("GET", path, common.SuccessResponse(policysetcontrollerv2.PolicySet{
		ID:         "policy-123",
		Name:       "Access Policy",
		PolicyType: policyType,
	}))

	result, _, err := policysetcontrollerv2.GetByPolicyType(context.Background(), api.Service, policyType)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, policyType, result.PolicyType)
}

func TestPolicySetControllerV2_GetAllByType_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policyType := "ACCESS_POLICY"
	path := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", policyType)

	api.On("GET", path, common.SuccessResponse(common.ZPAList([]policysetcontrollerv2.PolicyRuleResource{
		{ID: "rule-001", Name: "Rule 1"},
		{ID: "rule-002", Name: "Rule 2"},
	})))

	result, _, err := policysetcontrollerv2.GetAllByType(context.Background(), api.Service, policyType)
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestPolicySetControllerV2_GetPolicyRule_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policySetID := "policy-set-123"
	ruleID := "rule-456"
	path := common.ZPAPath(api.CustomerID, "policySet", policySetID, "rule", ruleID)

	api.On("GET", path, common.SuccessResponse(policysetcontrollerv2.PolicyRuleResource{
		ID:     ruleID,
		Name:   "Test Rule",
		Action: "ALLOW",
	}))

	result, _, err := policysetcontrollerv2.GetPolicyRule(context.Background(), api.Service, policySetID, ruleID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
}

func TestPolicySetControllerV2_CreateRule_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policySetID := "policy-set-123"
	path := common.ZPAv2Path(api.CustomerID, "policySet", policySetID, "rule")

	api.On("POST", path, common.SuccessResponse(policysetcontrollerv2.PolicyRule{
		ID:   "new-rule-789",
		Name: "New Rule",
	}))

	newRule := &policysetcontrollerv2.PolicyRule{
		Name:        "New Rule",
		Action:      "ALLOW",
		PolicySetID: policySetID,
	}
	result, _, err := policysetcontrollerv2.CreateRule(context.Background(), api.Service, newRule)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-rule-789", result.ID)
}

func TestPolicySetControllerV2_UpdateRule_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policySetID := "policy-set-123"
	ruleID := "rule-456"
	path := common.ZPAv2Path(api.CustomerID, "policySet", policySetID, "rule", ruleID)

	api.On("PUT", path, common.NoContentResponse())

	resp, err := policysetcontrollerv2.UpdateRule(context.Background(), api.Service, policySetID, ruleID, &policysetcontrollerv2.PolicyRule{
		ID:   ruleID,
		Name: "Updated Rule",
	})
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPolicySetControllerV2_Delete_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policySetID := "policy-set-123"
	ruleID := "rule-456"
	path := common.ZPAPath(api.CustomerID, "policySet", policySetID, "rule", ruleID)

	api.On("DELETE", path, common.NoContentResponse())

	resp, err := policysetcontrollerv2.Delete(context.Background(), api.Service, policySetID, ruleID)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPolicySetControllerV2_GetByNameAndType_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policyType := "ACCESS_POLICY"
	ruleName := "Production Rule"
	path := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", policyType)

	api.On("GET", path, common.SuccessResponse(common.ZPAList([]policysetcontrollerv2.PolicyRuleResource{
		{ID: "rule-001", Name: "Other Rule"},
		{ID: "rule-002", Name: ruleName},
	})))

	result, _, err := policysetcontrollerv2.GetByNameAndType(context.Background(), api.Service, policyType, ruleName)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "rule-002", result.ID)
}

func TestPolicySetControllerV2_GetByNameAndTypes_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	ruleName := "Shared Rule"
	path := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", "ACCESS_POLICY")

	api.On("GET", path, common.SuccessResponse(common.ZPAList([]policysetcontrollerv2.PolicyRuleResource{
		{ID: "rule-002", Name: ruleName},
	})))

	result, _, err := policysetcontrollerv2.GetByNameAndTypes(context.Background(), api.Service, []string{"ACCESS_POLICY", "TIMEOUT_POLICY"}, ruleName)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleName, result.Name)
}

func TestPolicySetControllerV2_Reorder_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policySetID := "policy-set-123"
	ruleID := "rule-456"
	path := common.ZPAPath(api.CustomerID, "policySet", policySetID, "rule", ruleID, "reorder", "1")

	api.On("PUT", path, common.NoContentResponse())

	resp, err := policysetcontrollerv2.Reorder(context.Background(), api.Service, policySetID, ruleID, 1)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPolicySetControllerV2_BulkReorder_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policyType := "ACCESS_POLICY"
	policySetPath := common.ZPAPath(api.CustomerID, "policySet", "policyType", policyType)
	rulesPath := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", policyType)
	reorderPath := common.ZPAPath(api.CustomerID, "policySet", "policy-123", "reorder")

	api.On("GET", policySetPath, common.SuccessResponse(policysetcontrollerv2.PolicySet{
		ID:         "policy-123",
		PolicyType: policyType,
	}))
	api.On("GET", rulesPath, common.SuccessResponse(common.ZPAList([]policysetcontrollerv2.PolicyRuleResource{
		{ID: "rule-001", Name: "Rule 1", RuleOrder: "1"},
		{ID: "rule-002", Name: "Rule 2", RuleOrder: "2"},
		{ID: "rule-default", Name: "Default_Rule", RuleOrder: "3"},
	})))
	api.On("PUT", reorderPath, common.NoContentResponse())

	resp, err := policysetcontrollerv2.BulkReorder(context.Background(), api.Service, policyType, map[string]int{
		"rule-002": 1,
		"rule-001": 2,
	})
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPolicySetControllerV2_GetPolicyCount_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policyType := "ACCESS_POLICY"
	path := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", policyType, "count")

	api.On("GET", path, common.SuccessResponse(common.ZPAList([]policysetcontrollerv2.PolicyRuleResource{
		{ID: "rule-001"},
	})))

	result, _, err := policysetcontrollerv2.GetPolicyCount(context.Background(), api.Service, policyType)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestPolicySetControllerV2_GetPolicyByApplication_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policyType := "ACCESS_POLICY"
	appID := "app-123"
	path := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", policyType, "application", appID)

	api.On("GET", path, common.SuccessResponse(common.ZPAList([]policysetcontrollerv2.PolicyRuleResource{
		{ID: "rule-001", Name: "App Rule"},
	})))

	result, _, err := policysetcontrollerv2.GetPolicyByApplication(context.Background(), api.Service, policyType, appID)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestPolicySetControllerV2_GetRiskScoreValues_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	path := common.ZPAv2Path(api.CustomerID, "riskScoreValues")

	excludeUnknown := true
	api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Equal(t, "true", r.URL.Query().Get("excludeUnknown"))
		return common.SuccessResponse([]string{"LOW", "MEDIUM", "HIGH"})
	})

	result, _, err := policysetcontrollerv2.GetRiskScoreValues(context.Background(), api.Service, &excludeUnknown)
	require.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestPolicySetControllerV2_GetPolicyRule_NotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "policySet", "policy-set-123", "rule", "missing")

	api.On("GET", path, common.NotFoundResponse())

	result, _, err := policysetcontrollerv2.GetPolicyRule(context.Background(), api.Service, "policy-set-123", "missing")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestPolicySetControllerV2_GetByNameAndType_NotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	policyType := "ACCESS_POLICY"
	path := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", policyType)

	api.On("GET", path, common.SuccessResponse(common.ZPAList([]policysetcontrollerv2.PolicyRuleResource{
		{ID: "rule-001", Name: "Other Rule"},
	})))

	result, _, err := policysetcontrollerv2.GetByNameAndType(context.Background(), api.Service, policyType, "Missing Rule")
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no policy rule named 'Missing Rule' found")
}

func TestPolicySetControllerV2_GetByNameAndTypes_NoMatchAcrossTypes_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	ruleName := "Ghost Rule"
	pathAccess := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", "ACCESS_POLICY")
	pathTimeout := common.ZPAPath(api.CustomerID, "policySet", "rules", "policyType", "TIMEOUT_POLICY")

	api.On("GET", pathAccess, common.SuccessResponse(common.ZPAList([]policysetcontrollerv2.PolicyRuleResource{
		{ID: "r1", Name: "Allow All"},
	})))
	api.On("GET", pathTimeout, common.SuccessResponse(common.ZPAList([]policysetcontrollerv2.PolicyRuleResource{
		{ID: "r2", Name: "Timeout Default"},
	})))

	result, _, err := policysetcontrollerv2.GetByNameAndTypes(context.Background(), api.Service, []string{"ACCESS_POLICY", "TIMEOUT_POLICY"}, ruleName)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no policy rule named 'Ghost Rule' found in any policy type")
}
