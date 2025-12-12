// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/workloadgroups"
)

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

