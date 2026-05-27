// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	ziacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/c2c_incident_receiver"
)

const (
	c2cIRBasePath   = "/zia/api/v1/cloudToCloudIR"
	c2cIRCountPath  = "/zia/api/v1/cloudToCloudIR/count"
	c2cIRLitePath   = "/zia/api/v1/cloudToCloudIR/lite"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestC2CIncidentReceiver_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	receiverID := 1001
	path := "/zia/api/v1/cloudToCloudIR/1001"

	server.On("GET", path, common.SuccessResponse(c2c_incident_receiver.C2CIncidentReceiver{
		ID:     receiverID,
		Name:   "Slack Receiver",
		Status: []string{"ACTIVE"},
		LastValidationMsg: &c2c_incident_receiver.LastValidationMsg{
			ErrorMsg:  "",
			ErrorCode: "OK",
		},
		OnboardableEntity: &c2c_incident_receiver.OnboardableEntity{
			ID:   10,
			Name: "Slack Tenant",
			Type: "SLACK",
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.Get(context.Background(), service, receiverID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, receiverID, result.ID)
	assert.Equal(t, "Slack Receiver", result.Name)
	require.NotNil(t, result.LastValidationMsg)
	assert.Equal(t, "OK", result.LastValidationMsg.ErrorCode)
}

func TestC2CIncidentReceiver_Get_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/cloudToCloudIR/9999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.Get(context.Background(), service, 9999)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestC2CIncidentReceiver_GetC2CIRName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	receiverName := "Teams Receiver"
	server.On("GET", c2cIRBasePath, common.SuccessResponse([]c2c_incident_receiver.C2CIncidentReceiver{
		{ID: 1, Name: "Other Receiver"},
		{ID: 2, Name: receiverName, Status: []string{"ACTIVE"}},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.GetC2CIRName(context.Background(), service, receiverName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, receiverName, result.Name)
	assert.Equal(t, 2, result.ID)
}

func TestC2CIncidentReceiver_GetC2CIRName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", c2cIRBasePath, common.SuccessResponse([]c2c_incident_receiver.C2CIncidentReceiver{
		{ID: 5, Name: "Salesforce IR"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.GetC2CIRName(context.Background(), service, "salesforce ir")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Salesforce IR", result.Name)
}

func TestC2CIncidentReceiver_GetC2CIRName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", c2cIRBasePath, common.SuccessResponse([]c2c_incident_receiver.C2CIncidentReceiver{
		{ID: 1, Name: "Existing"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.GetC2CIRName(context.Background(), service, "Missing")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no incident receiver found with name: Missing")
}

func TestC2CIncidentReceiver_GetC2CIRName_APIError_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", c2cIRBasePath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.GetC2CIRName(context.Background(), service, "Any")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestC2CIncidentReceiver_ValidateDelete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	receiverID := 2002
	path := "/zia/api/v1/cloudToCloudIR/config/2002/validateDelete"

	server.On("GET", path, common.SuccessResponse(c2c_incident_receiver.C2CIncidentReceiver{
		ID:     receiverID,
		Name:   "Deletable Receiver",
		Status: []string{"ACTIVE"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.ValidateDelete(context.Background(), service, receiverID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, receiverID, result.ID)
}

func TestC2CIncidentReceiver_ValidateDelete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/cloudToCloudIR/config/404/validateDelete"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.ValidateDelete(context.Background(), service, 404)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestC2CIncidentReceiver_C2CIRCount_NoSearch_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", c2cIRCountPath, common.SuccessResponse(7))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	count, err := c2c_incident_receiver.C2CIRCount(context.Background(), service, "")

	require.NoError(t, err)
	assert.Equal(t, 7, count)
	require.Len(t, server.Handler.Requests, 1)
	assert.Empty(t, server.Handler.Requests[0].Query)
}

func TestC2CIncidentReceiver_C2CIRCount_WithSearch_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", c2cIRCountPath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Equal(t, "search=slack+receiver", r.URL.RawQuery)
		return common.SuccessResponse(3)
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	count, err := c2c_incident_receiver.C2CIRCount(context.Background(), service, "slack receiver")

	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestC2CIncidentReceiver_C2CIRCount_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", c2cIRCountPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	count, err := c2c_incident_receiver.C2CIRCount(context.Background(), service, "")

	require.Error(t, err)
	assert.Equal(t, 0, count)
}

func TestC2CIncidentReceiver_GetAllLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", c2cIRLitePath, common.SuccessResponse([]c2c_incident_receiver.C2CIncidentReceiver{
		{ID: 1, Name: "Lite Receiver 1"},
		{ID: 2, Name: "Lite Receiver 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.GetAllLite(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestC2CIncidentReceiver_GetAllLite_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", c2cIRLitePath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.GetAllLite(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestC2CIncidentReceiver_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", c2cIRBasePath, common.SuccessResponse([]c2c_incident_receiver.C2CIncidentReceiver{
		{ID: 10, Name: "Receiver A", ModifiedTime: 1699000000},
		{ID: 11, Name: "Receiver B", ModifiedTime: 1699100000},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Receiver A", result[0].Name)
}

func TestC2CIncidentReceiver_GetAll_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", c2cIRBasePath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := c2c_incident_receiver.GetAll(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// Structure Tests
// =====================================================

func TestC2CIncidentReceiver_Structure(t *testing.T) {
	t.Parallel()

	t.Run("C2CIncidentReceiver JSON marshaling", func(t *testing.T) {
		receiver := c2c_incident_receiver.C2CIncidentReceiver{
			ID:                       100,
			Name:                     "AWS S3 Receiver",
			Status:                   []string{"ACTIVE", "VALIDATED"},
			ModifiedTime:             1700000000,
			LastTenantValidationTime: 1700001000,
			LastValidationMsg: &c2c_incident_receiver.LastValidationMsg{
				ErrorMsg:  "none",
				ErrorCode: "0",
			},
			LastModifiedBy: &ziacommon.IDNameExtensions{ID: 1, Name: "admin@example.com"},
			OnboardableEntity: &c2c_incident_receiver.OnboardableEntity{
				ID:                 50,
				Name:               "AWS Entity",
				Type:               "AWS",
				EnterpriseTenantID: "tenant-123",
				Application:        "S3",
				LastValidationMsg: c2c_incident_receiver.LastValidationMsg{
					ErrorCode: "OK",
				},
				TenantAuthorizationInfo: c2c_incident_receiver.TenantAuthorizationInfo{
					Type:                 "AWS",
					RoleArn:              "arn:aws:iam::123:role/zscaler",
					QuarantineBucketName: "quarantine-bucket",
					FeaturesSupported:    []string{"DLP", "MALWARE"},
				},
				ZscalerAppTenantID: &ziacommon.IDNameExtensions{ID: 99, Name: "zs-tenant"},
			},
		}

		data, err := json.Marshal(receiver)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"AWS S3 Receiver"`)
		assert.Contains(t, string(data), `"lastValidationMsg"`)
		assert.Contains(t, string(data), `"onboardableEntity"`)
		assert.Contains(t, string(data), `"roleArn"`)
	})

	t.Run("C2CIncidentReceiver JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 200,
			"name": "Box Receiver",
			"status": ["PENDING"],
			"modifiedTime": 1700100000,
			"lastValidationMsg": {
				"errorMsg": "validation pending",
				"errorCode": "PENDING"
			},
			"onboardableEntity": {
				"id": 60,
				"name": "Box Tenant",
				"type": "BOX",
				"tenantAuthorizationInfo": {
					"clientId": "client-id",
					"clientSecret": "secret",
					"subdomain": "corp",
					"smirBucketConfig": [
						{
							"configName": "cfg1",
							"metadataBucketName": "meta-bucket",
							"dataBucketName": "data-bucket",
							"id": 1
						}
					]
				}
			}
		}`

		var receiver c2c_incident_receiver.C2CIncidentReceiver
		err := json.Unmarshal([]byte(jsonData), &receiver)
		require.NoError(t, err)

		assert.Equal(t, 200, receiver.ID)
		assert.Equal(t, "Box Receiver", receiver.Name)
		require.NotNil(t, receiver.LastValidationMsg)
		assert.Equal(t, "PENDING", receiver.LastValidationMsg.ErrorCode)
		require.NotNil(t, receiver.OnboardableEntity)
		assert.Equal(t, "client-id", receiver.OnboardableEntity.TenantAuthorizationInfo.ClientID)
		assert.Len(t, receiver.OnboardableEntity.TenantAuthorizationInfo.SmirBucketConfig, 1)
		assert.Equal(t, "meta-bucket", receiver.OnboardableEntity.TenantAuthorizationInfo.SmirBucketConfig[0].MetadataBucketName)
	})

	t.Run("LastValidationMsg JSON marshaling", func(t *testing.T) {
		msg := c2c_incident_receiver.LastValidationMsg{
			ErrorMsg:  "invalid credentials",
			ErrorCode: "AUTH_FAILED",
		}

		data, err := json.Marshal(msg)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"errorMsg":"invalid credentials"`)
		assert.Contains(t, string(data), `"errorCode":"AUTH_FAILED"`)
	})

	t.Run("SmirBucketConfig JSON marshaling", func(t *testing.T) {
		cfg := c2c_incident_receiver.SmirBucketConfig{
			ConfigName:         "primary",
			MetadataBucketName: "meta",
			DataBucketName:     "data",
			ID:                 42,
		}

		data, err := json.Marshal(cfg)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"configName":"primary"`)
		assert.Contains(t, string(data), `"id":42`)
	})

	t.Run("TenantAuthorizationInfo JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"accessToken": "token",
			"botToken": "bot-token",
			"redirectUrl": "https://redirect.example.com",
			"type": "SLACK",
			"env": "production",
			"qtnInfoCleared": true,
			"featuresSupported": ["DLP"]
		}`

		var info c2c_incident_receiver.TenantAuthorizationInfo
		err := json.Unmarshal([]byte(jsonData), &info)
		require.NoError(t, err)

		assert.Equal(t, "token", info.AccessToken)
		assert.Equal(t, "SLACK", info.Type)
		assert.True(t, info.QtnInfoCleared)
		assert.Len(t, info.FeaturesSupported, 1)
	})
}

func TestC2CIncidentReceiver_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse incident receivers list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Receiver 1", "status": ["ACTIVE"]},
			{"id": 2, "name": "Receiver 2", "status": ["INACTIVE"]}
		]`

		var receivers []c2c_incident_receiver.C2CIncidentReceiver
		err := json.Unmarshal([]byte(jsonResponse), &receivers)
		require.NoError(t, err)

		assert.Len(t, receivers, 2)
		assert.Equal(t, "ACTIVE", receivers[0].Status[0])
	})

	t.Run("Parse count integer response", func(t *testing.T) {
		jsonResponse := `12`

		var count int
		err := json.Unmarshal([]byte(jsonResponse), &count)
		require.NoError(t, err)

		assert.Equal(t, 12, count)
	})
}
