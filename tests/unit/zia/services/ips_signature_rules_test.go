package services

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/ips_control_policies/ips_signature_rules"
)

const (
	ipsSignaturesPath         = "/zia/api/v1/ipsSignatureRules"
	ipsSignaturesExportPath   = "/zia/api/v1/ipsSignatureRules/export"
	ipsSignaturesValidatePath = "/zia/api/v1/ipsSignatureRules/validateRuleText"
)

const advancedSecurityCategoryID = 64

func sampleIPSSignatureRule(sid int) ips_signature_rules.IPSSignatureRules {
	return ips_signature_rules.IPSSignatureRules{
		Name:        "tests-ips-signature-rule",
		Description: "tests-ips-signature-description",
		RuleText: fmt.Sprintf(
			`alert http any any -> any any (msg:"Test HTTP rule sid %d"; content:"/admin"; http_uri; nocase; sid:%d; rev:1;)`,
			sid, sid,
		),
		Enabled: false,
		Category: &ips_signature_rules.IPSSignatureCategory{
			ID: advancedSecurityCategoryID,
		},
	}
}

func TestIPSSignatureRules_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 100
	rule := sampleIPSSignatureRule(1000001)
	rule.ID = ruleID

	server.On("GET", ipsSignaturesPath+"/100", common.SuccessResponse(rule))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ips_signature_rules.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
	assert.False(t, result.Enabled)
	assert.Equal(t, advancedSecurityCategoryID, result.Category.ID)
}

func TestIPSSignatureRules_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleName := "tests-ips-signature-rule"
	server.On("GET", ipsSignaturesPath, common.SuccessResponse([]ips_signature_rules.IPSSignatureRules{
		func() ips_signature_rules.IPSSignatureRules {
			r := sampleIPSSignatureRule(1000001)
			r.ID = 100
			return r
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ips_signature_rules.GetByName(context.Background(), service, ruleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleName, result.Name)
}

func TestIPSSignatureRules_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	created := sampleIPSSignatureRule(1000001)
	created.ID = 99999

	server.On("POST", ipsSignaturesPath, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := sampleIPSSignatureRule(1000001)
	result, _, err := ips_signature_rules.Create(context.Background(), service, &input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestIPSSignatureRules_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 100
	updated := sampleIPSSignatureRule(1000001)
	updated.ID = ruleID
	updated.Description = "updated-description"

	server.On("PUT", ipsSignaturesPath+"/100", common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := ips_signature_rules.Update(context.Background(), service, ruleID, &updated)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "updated-description", result.Description)
}

func TestIPSSignatureRules_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", ipsSignaturesPath+"/100", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = ips_signature_rules.Delete(context.Background(), service, 100)

	require.NoError(t, err)
}

func TestIPSSignatureRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", ipsSignaturesPath, common.SuccessResponse([]ips_signature_rules.IPSSignatureRules{
		sampleIPSSignatureRule(1000001),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ips_signature_rules.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestIPSSignatureRules_Export_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	csvBody := `+, "rule-name", "valid suricata rule text", "ips-category-name", "description", "Enable"`
	server.On("GET", ipsSignaturesExportPath, common.SuccessResponse(csvBody))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ips_signature_rules.ExportIPSSignatureRules(context.Background(), service)

	require.NoError(t, err)
	assert.Contains(t, string(result), "rule-name")
}

func TestIPSSignatureRules_ValidateRuleText_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", ipsSignaturesValidatePath, common.SuccessResponse(ips_signature_rules.IPSSignatureRulesValidation{
		Status: 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	ruleText := sampleIPSSignatureRule(1000001).RuleText
	result, err := ips_signature_rules.ValidateIPSSignatureRuleText(context.Background(), service, ruleText)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, result.Status)
}

func TestIPSSignatureRules_ValidateRuleText_EmptyInput_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ips_signature_rules.ValidateIPSSignatureRuleText(context.Background(), service, "")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestIPSSignatureRules_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		rule := sampleIPSSignatureRule(1000001)
		rule.ID = 100

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"enabled":false`)
		assert.Contains(t, string(data), `"ruleText"`)
	})
}
