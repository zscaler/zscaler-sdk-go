// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/rule_labels"
)

func TestRuleLabels_Structure(t *testing.T) {
	t.Parallel()

	t.Run("RuleLabels JSON marshaling", func(t *testing.T) {
		label := rule_labels.RuleLabels{
			ID:                  12345,
			Name:                "Security Policy Label",
			Description:         "Label for security policies",
			LastModifiedTime:    1699000000,
			ReferencedRuleCount: 15,
		}

		data, err := json.Marshal(label)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Security Policy Label"`)
		assert.Contains(t, string(data), `"referencedRuleCount":15`)
	})

	t.Run("RuleLabels JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Compliance Label",
			"description": "Label for compliance rules",
			"lastModifiedTime": 1699500000,
			"lastModifiedBy": {
				"id": 100,
				"name": "admin@company.com"
			},
			"createdBy": {
				"id": 101,
				"name": "creator@company.com"
			},
			"referencedRuleCount": 25
		}`

		var label rule_labels.RuleLabels
		err := json.Unmarshal([]byte(jsonData), &label)
		require.NoError(t, err)

		assert.Equal(t, 54321, label.ID)
		assert.Equal(t, "Compliance Label", label.Name)
		assert.Equal(t, 25, label.ReferencedRuleCount)
		assert.NotNil(t, label.LastModifiedBy)
		assert.NotNil(t, label.CreatedBy)
	})
}

func TestRuleLabels_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse rule labels list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Label 1", "referencedRuleCount": 5},
			{"id": 2, "name": "Label 2", "referencedRuleCount": 10},
			{"id": 3, "name": "Label 3", "referencedRuleCount": 0}
		]`

		var labels []rule_labels.RuleLabels
		err := json.Unmarshal([]byte(jsonResponse), &labels)
		require.NoError(t, err)

		assert.Len(t, labels, 3)
		assert.Equal(t, 10, labels[1].ReferencedRuleCount)
	})
}

