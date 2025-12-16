// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/workloadgroups"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestWorkloadGroups_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	workloadID := 12345
	path := "/zia/api/v1/workloadGroups/12345"

	server.On("GET", path, common.SuccessResponse(workloadgroups.WorkloadGroup{
		ID:          workloadID,
		Name:        "Production Workloads",
		Description: "Workload group for production servers",
		Expression:  "(TAG.environment = 'production')",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := workloadgroups.Get(context.Background(), service, workloadID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, workloadID, result.ID)
	assert.Equal(t, "Production Workloads", result.Name)
}

func TestWorkloadGroups_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	workloadName := "Production Workloads"
	path := "/zia/api/v1/workloadGroups"

	server.On("GET", path, common.SuccessResponse([]workloadgroups.WorkloadGroup{
		{ID: 1, Name: "Development Workloads", Description: "Dev workloads"},
		{ID: 2, Name: workloadName, Description: "Prod workloads"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := workloadgroups.GetByName(context.Background(), service, workloadName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, workloadName, result.Name)
}

func TestWorkloadGroups_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/workloadGroups"

	server.On("GET", path, common.SuccessResponse([]workloadgroups.WorkloadGroup{
		{ID: 1, Name: "Workload 1", Description: "Description 1"},
		{ID: 2, Name: "Workload 2", Description: "Description 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := workloadgroups.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestWorkloadGroups_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/workloadGroups"

	server.On("POST", path, common.SuccessResponse(workloadgroups.WorkloadGroup{
		ID:          99999,
		Name:        "New Workload Group",
		Description: "New description",
		Expression:  "(TAG.env = 'staging')",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newWorkload := &workloadgroups.WorkloadGroup{
		Name:        "New Workload Group",
		Description: "New description",
		Expression:  "(TAG.env = 'staging')",
	}

	result, _, err := workloadgroups.Create(context.Background(), service, newWorkload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestWorkloadGroups_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	workloadID := 12345
	path := "/zia/api/v1/workloadGroups/12345"

	server.On("PUT", path, common.SuccessResponse(workloadgroups.WorkloadGroup{
		ID:          workloadID,
		Name:        "Updated Workload Group",
		Description: "Updated description",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateWorkload := &workloadgroups.WorkloadGroup{
		ID:          workloadID,
		Name:        "Updated Workload Group",
		Description: "Updated description",
	}

	result, _, err := workloadgroups.Update(context.Background(), service, workloadID, updateWorkload)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Workload Group", result.Name)
}

func TestWorkloadGroups_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	workloadID := 12345
	path := "/zia/api/v1/workloadGroups/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = workloadgroups.Delete(context.Background(), service, workloadID)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests
// =====================================================

func TestWorkloadGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WorkloadGroup JSON marshaling", func(t *testing.T) {
		group := workloadgroups.WorkloadGroup{
			ID:               12345,
			Name:             "Production Workloads",
			Description:      "Workload group for production servers",
			Expression:       "(TAG.environment = 'production') AND (TAG.tier = 'web')",
			LastModifiedTime: 1699000000,
			WorkloadTagExpression: workloadgroups.WorkloadTagExpression{
				ExpressionContainers: []workloadgroups.ExpressionContainer{
					{
						TagType:  "AWS_TAG",
						Operator: "AND",
						TagContainer: workloadgroups.TagContainer{
							Operator: "AND",
							Tags: []workloadgroups.Tags{
								{Key: "environment", Value: "production"},
								{Key: "tier", Value: "web"},
							},
						},
					},
				},
			},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Production Workloads"`)
		assert.Contains(t, string(data), `"expressionJson"`)
	})

	t.Run("WorkloadGroup JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Development Workloads",
			"description": "Workload group for dev servers",
			"expression": "(TAG.environment = 'development')",
			"lastModifiedTime": 1699500000,
			"lastModifiedBy": {
				"id": 100,
				"name": "admin@company.com"
			},
			"expressionJson": {
				"expressionContainers": [
					{
						"tagType": "AWS_TAG",
						"operator": "OR",
						"tagContainer": {
							"operator": "OR",
							"tags": [
								{"key": "environment", "value": "development"},
								{"key": "environment", "value": "staging"}
							]
						}
					}
				]
			}
		}`

		var group workloadgroups.WorkloadGroup
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.NotNil(t, group.LastModifiedBy)
		assert.Len(t, group.WorkloadTagExpression.ExpressionContainers, 1)
		assert.Len(t, group.WorkloadTagExpression.ExpressionContainers[0].TagContainer.Tags, 2)
	})

	t.Run("WorkloadTagExpression JSON marshaling", func(t *testing.T) {
		expr := workloadgroups.WorkloadTagExpression{
			ExpressionContainers: []workloadgroups.ExpressionContainer{
				{
					TagType:  "AZURE_TAG",
					Operator: "AND",
					TagContainer: workloadgroups.TagContainer{
						Operator: "OR",
						Tags: []workloadgroups.Tags{
							{Key: "app", Value: "api"},
						},
					},
				},
			},
		}

		data, err := json.Marshal(expr)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"tagType":"AZURE_TAG"`)
	})

	t.Run("Tags JSON marshaling", func(t *testing.T) {
		tag := workloadgroups.Tags{
			Key:   "environment",
			Value: "production",
		}

		data, err := json.Marshal(tag)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"key":"environment"`)
		assert.Contains(t, string(data), `"value":"production"`)
	})
}

func TestWorkloadGroups_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse workload groups list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Production", "expression": "(TAG.env = 'prod')"},
			{"id": 2, "name": "Development", "expression": "(TAG.env = 'dev')"},
			{"id": 3, "name": "Staging", "expression": "(TAG.env = 'staging')"}
		]`

		var groups []workloadgroups.WorkloadGroup
		err := json.Unmarshal([]byte(jsonResponse), &groups)
		require.NoError(t, err)

		assert.Len(t, groups, 3)
	})

	t.Run("Parse complex expression", func(t *testing.T) {
		jsonResponse := `{
			"id": 100,
			"name": "Complex Workload",
			"expressionJson": {
				"expressionContainers": [
					{
						"tagType": "AWS_TAG",
						"operator": "AND",
						"tagContainer": {
							"operator": "AND",
							"tags": [
								{"key": "env", "value": "prod"}
							]
						}
					},
					{
						"tagType": "GCP_TAG",
						"operator": "OR",
						"tagContainer": {
							"operator": "OR",
							"tags": [
								{"key": "region", "value": "us-west-1"},
								{"key": "region", "value": "us-east-1"}
							]
						}
					}
				]
			}
		}`

		var group workloadgroups.WorkloadGroup
		err := json.Unmarshal([]byte(jsonResponse), &group)
		require.NoError(t, err)

		assert.Len(t, group.WorkloadTagExpression.ExpressionContainers, 2)
	})
}

