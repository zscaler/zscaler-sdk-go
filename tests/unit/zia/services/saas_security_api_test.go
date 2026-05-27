// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/saas_security_api"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/saas_security_api/casb_dlp_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/saas_security_api/casb_malware_rules"
	ziacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

// =====================================================
// saas_security_api
// =====================================================

func TestSaaSSecurityAPI_GetDomainProfiles_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/domainProfiles"
	server.On("GET", path, common.SuccessResponse([]saas_security_api.DomainProfiles{
		{ProfileID: 1, ProfileName: "Corp Domains", IncludeCompanyDomains: true, CustomDomains: []string{"example.com"}},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := saas_security_api.GetDomainProfiles(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Corp Domains", result[0].ProfileName)
}

func TestSaaSSecurityAPI_GetQuarantineTombstoneLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/quarantineTombstoneTemplate/lite"
	server.On("GET", path, common.SuccessResponse([]saas_security_api.QuarantineTombstoneLite{
		{ID: 1, Name: "Default Tombstone"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := saas_security_api.GetQuarantineTombstoneLite(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestSaaSSecurityAPI_GetCasbEmailLabelLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbEmailLabel/lite"
	server.On("GET", path, common.SuccessResponse([]saas_security_api.CasbEmailLabel{
		{ID: 1, Name: "Confidential"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := saas_security_api.GetCasbEmailLabelLite(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestSaaSSecurityAPI_GetCasbTenantTagPolicy_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tenantID := 15881081
	path := "/zia/api/v1/casbTenant/15881081/tags/policy"
	server.On("GET", path, common.SuccessResponse([]saas_security_api.CasbTenantTags{
		{TagID: 100, TagName: "Finance"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := saas_security_api.GetCasbTenantTagPolicy(context.Background(), service, tenantID)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestSaaSSecurityAPI_GetCasbTenantLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbTenant/lite"
	server.On("GET", path, common.SuccessResponse([]saas_security_api.CasbTenants{
		{TenantID: 15881081, TenantName: "O365 Tenant"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := saas_security_api.GetCasbTenantLite(context.Background(), service, nil)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 15881081, result[0].TenantID)
}

func TestSaaSSecurityAPI_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbTenant/scanInfo"
	server.On("GET", path, common.SuccessResponse([]saas_security_api.CasbTenantScanInfo{
		{TenantID: 15881081},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := saas_security_api.GetAll(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

// =====================================================
// casb_dlp_rules
// =====================================================

func TestCasbDLPRules_GetByRuleID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	ruleType := "OFLCASB_DLP_ITSM"
	path := "/zia/api/v1/casbDlpRules/12345"

	server.On("GET", path, common.SuccessResponse(casb_dlp_rules.CasbDLPRules{
		ID: ruleID, Name: "tests-casb-dlp", Type: ruleType, Order: 1, Rank: 7,
		State: "ENABLED", Severity: "RULE_SEVERITY_HIGH", Action: "OFLCASB_DLP_REPORT_INCIDENT",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := casb_dlp_rules.GetByRuleID(context.Background(), service, ruleType, ruleID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
	assert.Equal(t, "OFLCASB_DLP_REPORT_INCIDENT", result.Action)
}

func TestCasbDLPRules_GetByRuleType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbDlpRules"
	server.On("GET", path, common.SuccessResponse([]casb_dlp_rules.CasbDLPRules{
		{ID: 1, Name: "Rule 1", Type: "OFLCASB_DLP_ITSM", State: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := casb_dlp_rules.GetByRuleType(context.Background(), service, "OFLCASB_DLP_ITSM")
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestCasbDLPRules_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbDlpRules"
	server.On("POST", path, common.SuccessResponse(casb_dlp_rules.CasbDLPRules{
		ID: 99999, Name: "tests-casb-dlp", Type: "OFLCASB_DLP_ITSM", Order: 1, Rank: 7, State: "ENABLED",
		Severity: "RULE_SEVERITY_HIGH", Action: "OFLCASB_DLP_REPORT_INCIDENT",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &casb_dlp_rules.CasbDLPRules{
		Name: "tests-casb-dlp", Type: "OFLCASB_DLP_ITSM", Order: 1, Rank: 7, State: "ENABLED",
		Severity: "RULE_SEVERITY_HIGH", Action: "OFLCASB_DLP_REPORT_INCIDENT",
		DeviceTrustLevels: []string{"UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"},
		Components:        []string{"COMPONENT_ITSM_ATTACHMENTS", "COMPONENT_ITSM_OBJECTS"},
		CloudAppTenants:   []ziacommon.IDNameExtensions{{ID: 15881081}},
	}

	result, err := casb_dlp_rules.Create(context.Background(), service, newRule)
	require.NoError(t, err)
	assert.Equal(t, 99999, result.ID)
}

func TestCasbDLPRules_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/casbDlpRules/12345"
	server.On("PUT", path, common.SuccessResponse(casb_dlp_rules.CasbDLPRules{
		ID: ruleID, Name: "tests-casb-dlp-updated",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	update := &casb_dlp_rules.CasbDLPRules{ID: ruleID, Name: "tests-casb-dlp-updated"}
	result, err := casb_dlp_rules.Update(context.Background(), service, ruleID, update)
	require.NoError(t, err)
	assert.Equal(t, "tests-casb-dlp-updated", result.Name)
}

func TestCasbDLPRules_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbDlpRules/12345"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = casb_dlp_rules.Delete(context.Background(), service, "OFLCASB_DLP_ITSM", 12345)
	require.NoError(t, err)
}

func TestCasbDLPRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbDlpRules/all"
	server.On("GET", path, common.SuccessResponse([]casb_dlp_rules.CasbDLPRules{
		{ID: 1, Name: "Rule 1", Type: "OFLCASB_DLP_ITSM"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := casb_dlp_rules.GetAll(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestCasbDLPRules_GetByRuleID_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbDlpRules/99999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := casb_dlp_rules.GetByRuleID(context.Background(), service, "OFLCASB_DLP_ITSM", 99999)
	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// casb_malware_rules
// =====================================================

func TestCasbMalwareRules_GetByRuleID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/casbMalwareRules/12345"
	server.On("GET", path, common.SuccessResponse(casb_malware_rules.CasbMalwareRules{
		ID: ruleID, Name: "tests-casb-avp", Type: "OFLCASB_AVP_ITSM", Order: 1, State: "ENABLED",
		Action: "OFLCASB_AVP_REPORT_MALWARE",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := casb_malware_rules.GetByRuleID(context.Background(), service, "OFLCASB_AVP_ITSM", ruleID)
	require.NoError(t, err)
	assert.Equal(t, "OFLCASB_AVP_REPORT_MALWARE", result.Action)
}

func TestCasbMalwareRules_GetByRuleType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbMalwareRules"
	server.On("GET", path, common.SuccessResponse([]casb_malware_rules.CasbMalwareRules{
		{ID: 1, Type: "OFLCASB_AVP_ITSM", State: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := casb_malware_rules.GetByRuleType(context.Background(), service, "OFLCASB_AVP_ITSM")
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestCasbMalwareRules_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbMalwareRules"
	server.On("POST", path, common.SuccessResponse(casb_malware_rules.CasbMalwareRules{
		ID: 99999, Name: "tests-casb-avp", Type: "OFLCASB_AVP_ITSM", Order: 1, State: "ENABLED",
		Action: "OFLCASB_AVP_REPORT_MALWARE",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &casb_malware_rules.CasbMalwareRules{
		Name: "tests-casb-avp", Type: "OFLCASB_AVP_ITSM", Order: 1, State: "ENABLED",
		Action: "OFLCASB_AVP_REPORT_MALWARE",
		CloudAppTenants: []ziacommon.IDNameExtensions{{ID: 15881081}},
	}

	result, err := casb_malware_rules.Create(context.Background(), service, newRule)
	require.NoError(t, err)
	assert.Equal(t, 99999, result.ID)
}

func TestCasbMalwareRules_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbMalwareRules/12345"
	server.On("PUT", path, common.SuccessResponse(casb_malware_rules.CasbMalwareRules{
		ID: 12345, Name: "tests-casb-avp-updated",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	update := &casb_malware_rules.CasbMalwareRules{ID: 12345, Name: "tests-casb-avp-updated"}
	result, err := casb_malware_rules.Update(context.Background(), service, 12345, update)
	require.NoError(t, err)
	assert.Equal(t, "tests-casb-avp-updated", result.Name)
}

func TestCasbMalwareRules_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbMalwareRules/12345"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = casb_malware_rules.Delete(context.Background(), service, "OFLCASB_AVP_ITSM", 12345)
	require.NoError(t, err)
}

func TestCasbMalwareRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/casbMalwareRules/all"
	server.On("GET", path, common.SuccessResponse([]casb_malware_rules.CasbMalwareRules{
		{ID: 1, Type: "OFLCASB_AVP_ITSM"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := casb_malware_rules.GetAll(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}
