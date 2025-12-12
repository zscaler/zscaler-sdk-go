// Package services provides unit tests for ZTW services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/workload_groups"
)

func TestWorkloadGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WorkloadGroup JSON marshaling", func(t *testing.T) {
		group := workload_groups.WorkloadGroup{
			ID:               12345,
			Name:             "Production Workloads",
			Description:      "All production environment workloads",
			Expression:       "env=production AND region=us-east",
			LastModifiedTime: 1699000000,
			WorkloadTagExpression: workload_groups.WorkloadTagExpression{
				ExpressionContainers: []workload_groups.ExpressionContainer{
					{
						TagType:  "AWS",
						Operator: "AND",
						TagContainer: workload_groups.TagContainer{
							Tags: []workload_groups.Tags{
								{Key: "env", Value: "production"},
								{Key: "region", Value: "us-east"},
							},
							Operator: "AND",
						},
					},
				},
			},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Production Workloads"`)
		assert.Contains(t, string(data), `"expression":"env=production AND region=us-east"`)
	})

	t.Run("WorkloadGroup JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Development Workloads",
			"description": "All development environment workloads",
			"expression": "env=dev OR env=staging",
			"lastModifiedTime": 1699500000,
			"lastModifiedBy": {
				"id": 100,
				"name": "admin@company.com"
			},
			"expressionJson": {
				"expressionContainers": [
					{
						"tagType": "AZURE",
						"operator": "OR",
						"tagContainer": {
							"tags": [
								{"key": "env", "value": "dev"},
								{"key": "env", "value": "staging"}
							],
							"operator": "OR"
						}
					}
				]
			}
		}`

		var group workload_groups.WorkloadGroup
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Equal(t, "Development Workloads", group.Name)
		assert.NotNil(t, group.LastModifiedBy)
		assert.Equal(t, "admin@company.com", group.LastModifiedBy.Name)
		assert.Len(t, group.WorkloadTagExpression.ExpressionContainers, 1)
		assert.Equal(t, "AZURE", group.WorkloadTagExpression.ExpressionContainers[0].TagType)
	})

	t.Run("WorkloadTagExpression JSON marshaling", func(t *testing.T) {
		expr := workload_groups.WorkloadTagExpression{
			ExpressionContainers: []workload_groups.ExpressionContainer{
				{
					TagType:  "GCP",
					Operator: "AND",
					TagContainer: workload_groups.TagContainer{
						Tags: []workload_groups.Tags{
							{Key: "project", Value: "my-project"},
						},
						Operator: "AND",
					},
				},
			},
		}

		data, err := json.Marshal(expr)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"tagType":"GCP"`)
		assert.Contains(t, string(data), `"project"`)
	})

	t.Run("Tags JSON marshaling", func(t *testing.T) {
		tags := []workload_groups.Tags{
			{Key: "environment", Value: "production"},
			{Key: "team", Value: "platform"},
			{Key: "cost-center", Value: "CC-123"},
		}

		data, err := json.Marshal(tags)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"key":"environment"`)
		assert.Contains(t, string(data), `"value":"production"`)
		assert.Contains(t, string(data), `"cost-center"`)
	})
}

func TestWorkloadGroups_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse workload groups list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"name": "Production",
				"description": "Production workloads",
				"expression": "env=prod"
			},
			{
				"id": 2,
				"name": "Development",
				"description": "Development workloads",
				"expression": "env=dev"
			},
			{
				"id": 3,
				"name": "Staging",
				"description": "Staging workloads",
				"expression": "env=staging"
			}
		]`

		var groups []workload_groups.WorkloadGroup
		err := json.Unmarshal([]byte(jsonResponse), &groups)
		require.NoError(t, err)

		assert.Len(t, groups, 3)
		assert.Equal(t, "Production", groups[0].Name)
		assert.Equal(t, "env=prod", groups[0].Expression)
	})

	t.Run("Parse complex expression", func(t *testing.T) {
		jsonResponse := `{
			"id": 100,
			"name": "Complex Workload Group",
			"expressionJson": {
				"expressionContainers": [
					{
						"tagType": "AWS",
						"operator": "AND",
						"tagContainer": {
							"tags": [
								{"key": "env", "value": "production"}
							],
							"operator": "AND"
						}
					},
					{
						"tagType": "CUSTOM",
						"operator": "OR",
						"tagContainer": {
							"tags": [
								{"key": "team", "value": "platform"},
								{"key": "team", "value": "infra"}
							],
							"operator": "OR"
						}
					}
				]
			}
		}`

		var group workload_groups.WorkloadGroup
		err := json.Unmarshal([]byte(jsonResponse), &group)
		require.NoError(t, err)

		assert.Len(t, group.WorkloadTagExpression.ExpressionContainers, 2)
		assert.Equal(t, "AWS", group.WorkloadTagExpression.ExpressionContainers[0].TagType)
		assert.Equal(t, "CUSTOM", group.WorkloadTagExpression.ExpressionContainers[1].TagType)
		assert.Len(t, group.WorkloadTagExpression.ExpressionContainers[1].TagContainer.Tags, 2)
	})
}

