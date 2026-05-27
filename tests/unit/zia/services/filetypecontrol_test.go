package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/filetypecontrol"
)

const fileTypeRulesPath = "/zia/api/v1/fileTypeRules"
const fileTypeCategoriesPath = "/zia/api/v1/fileTypeCategories"
const customFileTypesListPath = "/zia/api/v1/customFileTypes"

func sampleFileTypeRule() filetypecontrol.FileTypeRules {
	return filetypecontrol.FileTypeRules{
		Name:              "tests-filetype-rule",
		Description:       "tests-filetype-rule",
		Order:             1,
		Rank:              7,
		State:             "ENABLED",
		FilteringAction:   "ALLOW",
		Operation:         "DOWNLOAD",
		Protocols:         []string{"FOHTTP_RULE", "FTP_RULE", "HTTPS_RULE", "HTTP_RULE"},
		DeviceTrustLevels: []string{"UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"},
		FileTypes:         []string{"FTCATEGORY_ALZ", "FTCATEGORY_P7Z", "FTCATEGORY_B64"},
	}
}

func TestFileTypeControl_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	rule := sampleFileTypeRule()
	rule.ID = ruleID

	server.On("GET", fileTypeRulesPath+"/12345", common.SuccessResponse(rule))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filetypecontrol.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
	assert.Equal(t, "ALLOW", result.FilteringAction)
	assert.Equal(t, "DOWNLOAD", result.Operation)
}

func TestFileTypeControl_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleName := "tests-filetype-rule"
	server.On("GET", fileTypeRulesPath, common.SuccessResponse([]filetypecontrol.FileTypeRules{
		{ID: 1, Name: "Other Rule"},
		func() filetypecontrol.FileTypeRules {
			r := sampleFileTypeRule()
			r.ID = 2
			return r
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filetypecontrol.GetByName(context.Background(), service, ruleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleName, result.Name)
}

func TestFileTypeControl_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	created := sampleFileTypeRule()
	created.ID = 99999

	server.On("POST", fileTypeRulesPath, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := sampleFileTypeRule()
	result, err := filetypecontrol.Create(context.Background(), service, &input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestFileTypeControl_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	updated := sampleFileTypeRule()
	updated.ID = ruleID
	updated.Name = "updated-filetype-rule"

	server.On("PUT", fileTypeRulesPath+"/12345", common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filetypecontrol.Update(context.Background(), service, ruleID, &updated)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "updated-filetype-rule", result.Name)
}

func TestFileTypeControl_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", fileTypeRulesPath+"/12345", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = filetypecontrol.Delete(context.Background(), service, 12345)

	require.NoError(t, err)
}

func TestFileTypeControl_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", fileTypeRulesPath, common.SuccessResponse([]filetypecontrol.FileTypeRules{
		sampleFileTypeRule(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filetypecontrol.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestFileTypeControl_GetFileTypeCategories_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", fileTypeCategoriesPath, common.SuccessResponse([]filetypecontrol.FileTypeCategory{
		{ID: 1, Name: "ALZ", Parent: "FTCATEGORY_ALZ"},
		{ID: 2, Name: "P7Z", Parent: "FTCATEGORY_P7Z"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filetypecontrol.GetFileTypeCategories(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestFileTypeControl_GetCustomFileTypes_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", customFileTypesListPath, common.SuccessResponse([]filetypecontrol.CustomFileTypes{
		{ID: 1, Name: "Custom Type 1", Extension: "xyz"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filetypecontrol.GetCustomFileTypes(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestFileTypeControl_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		rule := sampleFileTypeRule()
		rule.ID = 12345

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"filteringAction":"ALLOW"`)
		assert.Contains(t, string(data), `"FTCATEGORY_ALZ"`)
	})
}
