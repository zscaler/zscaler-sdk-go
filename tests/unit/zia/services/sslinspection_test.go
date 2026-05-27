// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sslinspection"
)

func TestSSLInspection_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/sslInspectionRules/12345"
	server.On("GET", path, common.SuccessResponse(sslinspection.SSLInspectionRules{
		ID: ruleID, Name: "tests-ssl-rule", Order: 1, Rank: 7, State: "ENABLED",
		Action: sslinspection.Action{
			Type:                       "DECRYPT",
			OverrideDefaultCertificate: true,
			DecryptSubActions: &sslinspection.DecryptSubActions{
				ServerCertificates:  "BLOCK",
				OcspCheck:           true,
				MinClientTLSVersion: "CLIENT_TLS_1_0",
				MinServerTLSVersion: "SERVER_TLS_1_0",
				HTTP2Enabled:        true,
			},
		},
		CloudApplications: []string{"CHATGPT_AI", "ANDI"},
		Platforms:         []string{"SCAN_IOS", "SCAN_ANDROID", "SCAN_MACOS", "SCAN_WINDOWS"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sslinspection.Get(context.Background(), service, ruleID)
	require.NoError(t, err)
	assert.Equal(t, "DECRYPT", result.Action.Type)
	assert.Equal(t, "CLIENT_TLS_1_0", result.Action.DecryptSubActions.MinClientTLSVersion)
}

func TestSSLInspection_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sslInspectionRules"
	server.On("GET", path, common.SuccessResponse([]sslinspection.SSLInspectionRules{
		{ID: 1, Name: "Rule 1", State: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sslinspection.GetAll(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestSSLInspection_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sslInspectionRules"
	server.On("POST", path, common.SuccessResponse(sslinspection.SSLInspectionRules{
		ID: 99999, Name: "tests-ssl-rule", Order: 1, Rank: 7, State: "ENABLED",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &sslinspection.SSLInspectionRules{
		Name: "tests-ssl-rule", Order: 1, Rank: 7, State: "ENABLED",
		Action: sslinspection.Action{
			Type:                       "DECRYPT",
			OverrideDefaultCertificate: true,
			DecryptSubActions: &sslinspection.DecryptSubActions{
				ServerCertificates:  "BLOCK",
				OcspCheck:           true,
				MinClientTLSVersion: "CLIENT_TLS_1_0",
				MinServerTLSVersion: "SERVER_TLS_1_0",
			},
		},
		CloudApplications: []string{"CHATGPT_AI"},
		Platforms:         []string{"SCAN_WINDOWS"},
	}

	result, err := sslinspection.Create(context.Background(), service, newRule)
	require.NoError(t, err)
	assert.Equal(t, 99999, result.ID)
}

func TestSSLInspection_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/sslInspectionRules/12345"
	server.On("PUT", path, common.SuccessResponse(sslinspection.SSLInspectionRules{
		ID: ruleID, Name: "tests-ssl-updated",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	update := &sslinspection.SSLInspectionRules{
		ID: ruleID, Name: "tests-ssl-updated",
		Action: sslinspection.Action{
			Type:                       "DECRYPT",
			OverrideDefaultCertificate: true,
			DecryptSubActions: &sslinspection.DecryptSubActions{
				ServerCertificates:  "BLOCK",
				OcspCheck:           true,
				MinClientTLSVersion: "CLIENT_TLS_1_0",
				MinServerTLSVersion: "SERVER_TLS_1_0",
			},
		},
	}
	result, err := sslinspection.Update(context.Background(), service, ruleID, update)
	require.NoError(t, err)
	assert.Equal(t, "tests-ssl-updated", result.Name)
}

func TestSSLInspection_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sslInspectionRules/12345"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = sslinspection.Delete(context.Background(), service, 12345)
	require.NoError(t, err)
}

func TestSSLInspection_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	name := "tests-ssl-rule"
	path := "/zia/api/v1/sslInspectionRules"
	server.On("GET", path, common.SuccessResponse([]sslinspection.SSLInspectionRules{
		{ID: 1, Name: name, State: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sslinspection.GetByName(context.Background(), service, name)
	require.NoError(t, err)
	assert.Equal(t, name, result.Name)
}

func TestSSLInspection_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sslInspectionRules/99999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sslinspection.Get(context.Background(), service, 99999)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestSSLInspection_Create_ValidationError_SDK(t *testing.T) {
	service, err := common.CreateTestService(context.Background(), common.NewTestServer(), "123456")
	require.NoError(t, err)

	_, err = sslinspection.Create(context.Background(), service, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestSSLInspection_Create_DecryptValidationError_SDK(t *testing.T) {
	service, err := common.CreateTestService(context.Background(), common.NewTestServer(), "123456")
	require.NoError(t, err)

	_, err = sslinspection.Create(context.Background(), service, &sslinspection.SSLInspectionRules{
		Name: "bad-rule",
		Action: sslinspection.Action{Type: "DECRYPT"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decryptSubActions")
}

func TestSSLInspection_Create_DoNotDecryptValidationError_SDK(t *testing.T) {
	service, err := common.CreateTestService(context.Background(), common.NewTestServer(), "123456")
	require.NoError(t, err)

	_, err = sslinspection.Create(context.Background(), service, &sslinspection.SSLInspectionRules{
		Name: "bad-rule",
		Action: sslinspection.Action{Type: "DO_NOT_DECRYPT"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "doNotDecryptSubActions")
}

func TestSSLInspection_Create_BlockValidationError_SDK(t *testing.T) {
	service, err := common.CreateTestService(context.Background(), common.NewTestServer(), "123456")
	require.NoError(t, err)

	_, err = sslinspection.Create(context.Background(), service, &sslinspection.SSLInspectionRules{
		Name: "bad-rule",
		Action: sslinspection.Action{
			Type:              "BLOCK",
			DecryptSubActions: &sslinspection.DecryptSubActions{},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "neither decryptSubActions nor doNotDecryptSubActions")
}

func TestSSLInspection_Create_DoNotDecryptBypassValidation_SDK(t *testing.T) {
	service, err := common.CreateTestService(context.Background(), common.NewTestServer(), "123456")
	require.NoError(t, err)

	_, err = sslinspection.Create(context.Background(), service, &sslinspection.SSLInspectionRules{
		Name: "bad-rule",
		Action: sslinspection.Action{
			Type: "DO_NOT_DECRYPT",
			DoNotDecryptSubActions: &sslinspection.DoNotDecryptSubActions{
				BypassOtherPolicies: true,
				ServerCertificates:  "BLOCK",
			},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "bypassOtherPolicies is true")
}

func TestSSLInspection_Create_DoNotDecryptSuccess_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sslInspectionRules"
	server.On("POST", path, common.SuccessResponse(sslinspection.SSLInspectionRules{
		ID: 99999, Name: "tests-dnd-rule",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sslinspection.Create(context.Background(), service, &sslinspection.SSLInspectionRules{
		Name: "tests-dnd-rule",
		Action: sslinspection.Action{
			Type: "DO_NOT_DECRYPT",
			DoNotDecryptSubActions: &sslinspection.DoNotDecryptSubActions{
				ServerCertificates: "BLOCK",
				MinTLSVersion:      "CLIENT_TLS_1_0",
			},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 99999, result.ID)
}

func TestSSLInspection_Create_BlockSuccess_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sslInspectionRules"
	server.On("POST", path, common.SuccessResponse(sslinspection.SSLInspectionRules{
		ID: 88888, Name: "tests-block-rule",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sslinspection.Create(context.Background(), service, &sslinspection.SSLInspectionRules{
		Name: "tests-block-rule",
		Action: sslinspection.Action{
			Type:                       "BLOCK",
			OverrideDefaultCertificate: true,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 88888, result.ID)
}

func TestSSLInspection_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/sslInspectionRules"
	server.On("GET", path, common.SuccessResponse([]sslinspection.SSLInspectionRules{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sslinspection.GetByName(context.Background(), service, "missing-rule")
	require.Error(t, err)
	assert.Nil(t, result)
}
