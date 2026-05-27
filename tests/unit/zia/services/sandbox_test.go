// Package services provides unit tests for ZIA services
package services

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_report"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_settings"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_submission"
)

// =====================================================
// sandbox_rules
// =====================================================

func TestSandboxRules_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/sandboxRules/12345"
	server.On("GET", path, common.SuccessResponse(sandbox_rules.SandboxRules{
		ID: ruleID, Name: "tests-sandbox-rule", Order: 1, Rank: 7, State: "ENABLED",
		BaRuleAction: "BLOCK", FirstTimeEnable: true, MLActionEnabled: true,
		FirstTimeOperation: "ALLOW_SCAN",
		Protocols:          []string{"FOHTTP_RULE", "FTP_RULE", "HTTPS_RULE", "HTTP_RULE"},
		BaPolicyCategories: []string{"ADWARE_BLOCK", "BOTMAL_BLOCK", "RANSOMWARE_BLOCK"},
		FileTypes:          []string{"FTCATEGORY_P7Z", "FTCATEGORY_BZIP2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sandbox_rules.Get(context.Background(), service, ruleID)
	require.NoError(t, err)
	assert.Equal(t, "BLOCK", result.BaRuleAction)
	assert.Equal(t, "ALLOW_SCAN", result.FirstTimeOperation)
}

func TestSandboxRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sandboxRules"
	server.On("GET", path, common.SuccessResponse([]sandbox_rules.SandboxRules{
		{ID: 1, Name: "Rule 1", State: "ENABLED", BaRuleAction: "BLOCK"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sandbox_rules.GetAll(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestSandboxRules_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sandboxRules"
	server.On("POST", path, common.SuccessResponse(sandbox_rules.SandboxRules{
		ID: 99999, Name: "tests-sandbox-rule", Order: 1, Rank: 7, State: "ENABLED", BaRuleAction: "BLOCK",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &sandbox_rules.SandboxRules{
		Name: "tests-sandbox-rule", Order: 1, Rank: 7, State: "ENABLED",
		BaRuleAction: "BLOCK", FirstTimeEnable: true, MLActionEnabled: true,
		FirstTimeOperation: "ALLOW_SCAN",
		Protocols:          []string{"HTTPS_RULE", "HTTP_RULE"},
		BaPolicyCategories: []string{"SUSPICIOUS_BLOCK"},
		FileTypes:          []string{"FTCATEGORY_P7Z"},
	}

	result, err := sandbox_rules.Create(context.Background(), service, newRule)
	require.NoError(t, err)
	assert.Equal(t, 99999, result.ID)
}

func TestSandboxRules_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/sandboxRules/12345"
	server.On("PUT", path, common.SuccessResponse(sandbox_rules.SandboxRules{
		ID: ruleID, Name: "tests-sandbox-updated",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	update := &sandbox_rules.SandboxRules{ID: ruleID, Name: "tests-sandbox-updated"}
	result, err := sandbox_rules.Update(context.Background(), service, ruleID, update)
	require.NoError(t, err)
	assert.Equal(t, "tests-sandbox-updated", result.Name)
}

func TestSandboxRules_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sandboxRules/12345"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = sandbox_rules.Delete(context.Background(), service, 12345)
	require.NoError(t, err)
}

func TestSandboxRules_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	name := "tests-sandbox-rule"
	path := "/zia/api/v1/sandboxRules"
	server.On("GET", path, common.SuccessResponse([]sandbox_rules.SandboxRules{
		{ID: 1, Name: "other"},
		{ID: 2, Name: name, BaRuleAction: "BLOCK"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sandbox_rules.GetByName(context.Background(), service, name)
	require.NoError(t, err)
	assert.Equal(t, name, result.Name)
}

func TestSandboxRules_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sandboxRules/99999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sandbox_rules.Get(context.Background(), service, 99999)
	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// sandbox_settings
// =====================================================

func TestSandboxSettings_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/behavioralAnalysisAdvancedSettings"
	server.On("GET", path, common.SuccessResponse(sandbox_settings.BaAdvancedSettings{
		FileHashesToBeBlocked: []string{
			"42914d6d213a20a2684064be5c80ffa9",
			"c0202cf6aeab8437c638533d14563d35",
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sandbox_settings.Get(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result.FileHashesToBeBlocked, 2)
}

func TestSandboxSettings_Getv2_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/behavioralAnalysisAdvancedSettings"
	server.On("GET", path, common.SuccessResponse(sandbox_settings.Md5HashValueListPayload{
		Md5HashValueList: []sandbox_settings.Md5HashValue{
			{URL: "42914d6d213a20a2684064be5c80ffa9", Type: "MALWARE"},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sandbox_settings.Getv2(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result.Md5HashValueList, 1)
	assert.Equal(t, "MALWARE", result.Md5HashValueList[0].Type)
}

func TestSandboxSettings_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/behavioralAnalysisAdvancedSettings"
	server.On("PUT", path, common.SuccessResponse(sandbox_settings.BaAdvancedSettings{
		FileHashesToBeBlocked: []string{"42914d6d213a20a2684064be5c80ffa9"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	settings := sandbox_settings.BaAdvancedSettings{
		FileHashesToBeBlocked: []string{"42914d6d213a20a2684064be5c80ffa9"},
	}
	result, err := sandbox_settings.Update(context.Background(), service, settings)
	require.NoError(t, err)
	assert.Len(t, result.FileHashesToBeBlocked, 1)
}

func TestSandboxSettings_Updatev2_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/behavioralAnalysisAdvancedSettings"
	server.On("PUT", path, common.SuccessResponse(sandbox_settings.Md5HashValueListPayload{
		Md5HashValueList: []sandbox_settings.Md5HashValue{},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	payload := sandbox_settings.Md5HashValueListPayload{Md5HashValueList: []sandbox_settings.Md5HashValue{}}
	result, err := sandbox_settings.Updatev2(context.Background(), service, payload)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSandboxSettings_GetFileHashCount_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/behavioralAnalysisAdvancedSettings/fileHashCount"
	server.On("GET", path, common.SuccessResponse(sandbox_settings.FileHashCount{
		BlockedFileHashesCount: 5,
		RemainingFileHashes:    995,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sandbox_settings.GetFileHashCount(context.Background(), service)
	require.NoError(t, err)
	assert.Equal(t, 5, result.BlockedFileHashesCount)
	assert.Equal(t, 995, result.RemainingFileHashes)
}

// =====================================================
// sandbox_report
// =====================================================

func TestSandboxReport_GetRatingQuota_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sandbox/report/quota"
	server.On("GET", path, common.SuccessResponse([]sandbox_report.RatingQuota{
		{Allowed: 100, Used: 10},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sandbox_report.GetRatingQuota(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestSandboxReport_GetReportMD5Hash_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	md5Hash := "42914d6d213a20a2684064be5c80ffa9"
	path := "/zia/api/v1/sandbox/report/" + md5Hash
	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"Summary": map[string]interface{}{
			"Summary": map[string]interface{}{
				"Status":   "COMPLETED",
				"Category": "MALICIOUS",
				"FileType": "PDF",
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sandbox_report.GetReportMD5Hash(context.Background(), service, md5Hash, "summary")
	require.NoError(t, err)
	require.NotNil(t, result.Details)
	assert.Equal(t, "MALICIOUS", result.Details.Summary.Category)
}

// =====================================================
// sandbox_submission
// =====================================================

func TestSandboxSubmission_SubmitFile_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zscsb/submit"
	server.On("POST", path, common.SuccessResponse(sandbox_submission.ScanResult{
		Code: 0, Message: "Success", Md5: "42914d6d213a20a2684064be5c80ffa9",
		FileType: "PDF", SandboxSubmission: "SUBMITTED",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	file := strings.NewReader("test file content")
	result, err := sandbox_submission.SubmitFile(context.Background(), service, "test.pdf", file, "")
	require.NoError(t, err)
	assert.Equal(t, "42914d6d213a20a2684064be5c80ffa9", result.Md5)
}

func TestSandboxSubmission_Discan_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zscsb/discan"
	server.On("POST", path, common.SuccessResponse(sandbox_submission.ScanResult{
		Code: 0, Message: "Success", Md5: "c0202cf6aeab8437c638533d14563d35",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	file := bytes.NewReader([]byte("discan payload"))
	result, err := sandbox_submission.Discan(context.Background(), service, "sample.exe", file)
	require.NoError(t, err)
	assert.Equal(t, "c0202cf6aeab8437c638533d14563d35", result.Md5)
}
