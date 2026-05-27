// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	ziacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloud_app_instances"
)

const cloudAppInstancesPath = "/zia/api/v1/cloudApplicationInstances"

// =====================================================
// SDK Function Tests
// =====================================================

func TestCloudAppInstances_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	instanceID := 5001
	path := "/zia/api/v1/cloudApplicationInstances/5001"

	server.On("GET", path, common.SuccessResponse(cloud_app_instances.CloudApplicationInstances{
		InstanceID:   instanceID,
		InstanceType: "SALESFORCE",
		InstanceName: "Corp Salesforce",
		ModifiedAt:   1700000000,
		ModifiedBy:   &ziacommon.IDNameExtensions{ID: 1, Name: "admin@example.com"},
		InstanceIdentifiers: []cloud_app_instances.InstanceIdentifiers{
			{
				InstanceID:             instanceID,
				InstanceIdentifier:     "sf-123",
				InstanceIdentifierName: "Primary Org",
				IdentifierType:         "ORG_ID",
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloud_app_instances.Get(context.Background(), service, instanceID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, instanceID, result.InstanceID)
	assert.Equal(t, "Corp Salesforce", result.InstanceName)
	assert.Len(t, result.InstanceIdentifiers, 1)
}

func TestCloudAppInstances_Get_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/cloudApplicationInstances/9999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloud_app_instances.Get(context.Background(), service, 9999)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCloudAppInstances_GetInstanceByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	instanceName := "Office 365 Tenant"
	server.On("GET", cloudAppInstancesPath, common.SuccessResponse([]cloud_app_instances.CloudApplicationInstances{
		{InstanceID: 1, InstanceName: "Other Instance", InstanceType: "BOX"},
		{InstanceID: 2, InstanceName: instanceName, InstanceType: "OFFICE365"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloud_app_instances.GetInstanceByName(context.Background(), service, instanceName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, instanceName, result.InstanceName)
	assert.Equal(t, 2, result.InstanceID)
}

func TestCloudAppInstances_GetInstanceByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cloudAppInstancesPath, common.SuccessResponse([]cloud_app_instances.CloudApplicationInstances{
		{InstanceID: 10, InstanceName: "GitHub Enterprise", InstanceType: "GITHUB"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloud_app_instances.GetInstanceByName(context.Background(), service, "github enterprise")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "GitHub Enterprise", result.InstanceName)
}

func TestCloudAppInstances_GetInstanceByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cloudAppInstancesPath, common.SuccessResponse([]cloud_app_instances.CloudApplicationInstances{
		{InstanceID: 1, InstanceName: "Existing"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloud_app_instances.GetInstanceByName(context.Background(), service, "Missing Instance")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no cloud instance found with name: Missing Instance")
}

func TestCloudAppInstances_GetInstanceByName_APIError_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cloudAppInstancesPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloud_app_instances.GetInstanceByName(context.Background(), service, "Any")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCloudAppInstances_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", cloudAppInstancesPath, common.SuccessResponse(cloud_app_instances.CloudApplicationInstances{
		InstanceID:   88888,
		InstanceType: "SALESFORCE",
		InstanceName: "New Salesforce Instance",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newInstance := &cloud_app_instances.CloudApplicationInstances{
		InstanceType: "SALESFORCE",
		InstanceName: "New Salesforce Instance",
		InstanceIdentifiers: []cloud_app_instances.InstanceIdentifiers{
			{InstanceIdentifier: "org-abc", IdentifierType: "ORG_ID"},
		},
	}

	result, _, err := cloud_app_instances.Create(context.Background(), service, newInstance)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 88888, result.InstanceID)
	assert.Equal(t, "New Salesforce Instance", result.InstanceName)
}

func TestCloudAppInstances_Create_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", cloudAppInstancesPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newInstance := &cloud_app_instances.CloudApplicationInstances{
		InstanceName: "Fail Instance",
	}

	result, _, err := cloud_app_instances.Create(context.Background(), service, newInstance)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCloudAppInstances_Create_UnexpectedResponseType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Non-JSON / empty POST response returns *http.Response from Create,
	// which fails the *CloudApplicationInstances type assertion.
	server.On("POST", cloudAppInstancesPath, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newInstance := &cloud_app_instances.CloudApplicationInstances{
		InstanceName: "No Body Instance",
	}

	result, _, err := cloud_app_instances.Create(context.Background(), service, newInstance)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "object returned from api was not a cloud instance pointer")
}

func TestCloudAppInstances_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	instanceID := 5001
	path := "/zia/api/v1/cloudApplicationInstances/5001"

	server.On("PUT", path, common.SuccessResponse(cloud_app_instances.CloudApplicationInstances{
		InstanceID:   instanceID,
		InstanceType: "SALESFORCE",
		InstanceName: "Updated Instance Name",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateInstance := &cloud_app_instances.CloudApplicationInstances{
		InstanceID:   instanceID,
		InstanceName: "Updated Instance Name",
	}

	result, _, err := cloud_app_instances.Update(context.Background(), service, instanceID, updateInstance)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Instance Name", result.InstanceName)
}

func TestCloudAppInstances_Update_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/cloudApplicationInstances/5001"
	server.On("PUT", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateInstance := &cloud_app_instances.CloudApplicationInstances{
		InstanceName: "Updated",
	}

	result, _, err := cloud_app_instances.Update(context.Background(), service, 5001, updateInstance)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCloudAppInstances_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	instanceID := 5001
	path := "/zia/api/v1/cloudApplicationInstances/5001"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = cloud_app_instances.Delete(context.Background(), service, instanceID)

	require.NoError(t, err)
}

func TestCloudAppInstances_Delete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/cloudApplicationInstances/5001"
	server.On("DELETE", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = cloud_app_instances.Delete(context.Background(), service, 5001)

	require.Error(t, err)
}

func TestCloudAppInstances_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cloudAppInstancesPath, common.SuccessResponse([]cloud_app_instances.CloudApplicationInstances{
		{InstanceID: 1, InstanceName: "Instance A", InstanceType: "BOX"},
		{InstanceID: 2, InstanceName: "Instance B", InstanceType: "GITHUB"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloud_app_instances.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCloudAppInstances_GetAll_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cloudAppInstancesPath, common.SuccessResponse([]cloud_app_instances.CloudApplicationInstances{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloud_app_instances.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestCloudAppInstances_GetAll_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cloudAppInstancesPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloud_app_instances.GetAll(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// Structure Tests
// =====================================================

func TestCloudAppInstances_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CloudApplicationInstances JSON marshaling", func(t *testing.T) {
		instance := cloud_app_instances.CloudApplicationInstances{
			InstanceID:   100,
			InstanceType: "OFFICE365",
			InstanceName: "Corp O365",
			ModifiedAt:   1700000000,
			ModifiedBy:   &ziacommon.IDNameExtensions{ID: 5, Name: "admin@example.com"},
			InstanceIdentifiers: []cloud_app_instances.InstanceIdentifiers{
				{
					InstanceID:             100,
					InstanceIdentifier:     "tenant-guid",
					InstanceIdentifierName: "Primary Tenant",
					IdentifierType:         "TENANT_ID",
					ModifiedAt:             1700000100,
					ModifiedBy:             &ziacommon.IDNameExtensions{ID: 5, Name: "admin@example.com"},
				},
			},
		}

		data, err := json.Marshal(instance)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"instanceId":100`)
		assert.Contains(t, string(data), `"instanceType":"OFFICE365"`)
		assert.Contains(t, string(data), `"instanceName":"Corp O365"`)
		assert.Contains(t, string(data), `"instanceIdentifiers"`)
	})

	t.Run("CloudApplicationInstances JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"instanceId": 200,
			"instanceType": "SALESFORCE",
			"instanceName": "SF Prod",
			"modifiedAt": 1700200000,
			"modifiedBy": {"id": 10, "name": "ops@example.com"},
			"instanceIdentifiers": [
				{
					"instanceId": 200,
					"instanceIdentifier": "sf-org-1",
					"instanceIdentifierName": "Production Org",
					"identifierType": "ORG_ID"
				}
			]
		}`

		var instance cloud_app_instances.CloudApplicationInstances
		err := json.Unmarshal([]byte(jsonData), &instance)
		require.NoError(t, err)

		assert.Equal(t, 200, instance.InstanceID)
		assert.Equal(t, "SALESFORCE", instance.InstanceType)
		require.NotNil(t, instance.ModifiedBy)
		assert.Equal(t, "ops@example.com", instance.ModifiedBy.Name)
		assert.Len(t, instance.InstanceIdentifiers, 1)
		assert.Equal(t, "sf-org-1", instance.InstanceIdentifiers[0].InstanceIdentifier)
	})

	t.Run("InstanceIdentifiers JSON marshaling", func(t *testing.T) {
		identifier := cloud_app_instances.InstanceIdentifiers{
			InstanceID:             50,
			InstanceIdentifier:     "id-xyz",
			InstanceIdentifierName: "Secondary",
			IdentifierType:         "CUSTOM",
			ModifiedAt:             1700300000,
		}

		data, err := json.Marshal(identifier)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"instanceIdentifier":"id-xyz"`)
		assert.Contains(t, string(data), `"identifierType":"CUSTOM"`)
	})

	t.Run("InstanceIdentifiers JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"instanceId": 75,
			"instanceIdentifier": "box-tenant",
			"instanceIdentifierName": "Box Tenant",
			"identifierType": "TENANT",
			"modifiedAt": 1700400000
		}`

		var identifier cloud_app_instances.InstanceIdentifiers
		err := json.Unmarshal([]byte(jsonData), &identifier)
		require.NoError(t, err)

		assert.Equal(t, 75, identifier.InstanceID)
		assert.Equal(t, "box-tenant", identifier.InstanceIdentifier)
		assert.Equal(t, "TENANT", identifier.IdentifierType)
	})
}

func TestCloudAppInstances_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse cloud application instances list", func(t *testing.T) {
		jsonResponse := `[
			{"instanceId": 1, "instanceName": "Inst 1", "instanceType": "BOX"},
			{"instanceId": 2, "instanceName": "Inst 2", "instanceType": "GITHUB"}
		]`

		var instances []cloud_app_instances.CloudApplicationInstances
		err := json.Unmarshal([]byte(jsonResponse), &instances)
		require.NoError(t, err)

		assert.Len(t, instances, 2)
		assert.Equal(t, "BOX", instances[0].InstanceType)
	})
}
