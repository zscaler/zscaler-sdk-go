// Package services provides unit tests for ZWA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services/customeraudit"
)

func TestCustomerAudit_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AuditLogsResponse JSON marshaling", func(t *testing.T) {
		response := customeraudit.AuditLogsResponse{
			Cursor: common.Cursor{
				TotalPages:        5,
				CurrentPageNumber: 1,
				CurrentPageSize:   100,
				TotalElements:     450,
			},
			Logs: []customeraudit.AuditLog{
				{
					Action:     customeraudit.Action{Action: "CREATE"},
					Module:     "DLP",
					Resource:   "DLP Policy",
					ChangedAt:  "2024-01-15T10:30:00Z",
					ChangedBy:  "admin@company.com",
					OldRowJSON: "{}",
					NewRowJSON: `{"name": "New Policy"}`,
					ChangeNote: "Created new DLP policy",
				},
			},
		}

		data, err := json.Marshal(response)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"cursor"`)
		assert.Contains(t, string(data), `"logs"`)
		assert.Contains(t, string(data), `"module":"DLP"`)
	})

	t.Run("AuditLog JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"action": {"action": "UPDATE"},
			"module": "FIREWALL",
			"resource": "Firewall Rule",
			"changedAt": "2024-01-20T14:00:00Z",
			"changedBy": "security@company.com",
			"oldRowJson": "{\"enabled\": false}",
			"newRowJson": "{\"enabled\": true}",
			"changeNote": "Enabled firewall rule"
		}`

		var log customeraudit.AuditLog
		err := json.Unmarshal([]byte(jsonData), &log)
		require.NoError(t, err)

		assert.Equal(t, "UPDATE", log.Action.Action)
		assert.Equal(t, "FIREWALL", log.Module)
		assert.Equal(t, "Firewall Rule", log.Resource)
		assert.Equal(t, "security@company.com", log.ChangedBy)
	})

	t.Run("Action JSON marshaling", func(t *testing.T) {
		action := customeraudit.Action{
			Action: "DELETE",
		}

		data, err := json.Marshal(action)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"action":"DELETE"`)
	})
}

func TestCustomerAudit_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse audit logs response", func(t *testing.T) {
		jsonResponse := `{
			"cursor": {
				"totalPages": 3,
				"currentPageNumber": 1,
				"currentPageSize": 50,
				"totalElements": 125
			},
			"logs": [
				{
					"action": {"action": "CREATE"},
					"module": "DLP",
					"resource": "DLP Dictionary",
					"changedAt": "2024-01-10T09:00:00Z",
					"changedBy": "admin1@company.com"
				},
				{
					"action": {"action": "UPDATE"},
					"module": "URL_FILTERING",
					"resource": "URL Category",
					"changedAt": "2024-01-11T10:00:00Z",
					"changedBy": "admin2@company.com"
				},
				{
					"action": {"action": "DELETE"},
					"module": "FIREWALL",
					"resource": "Firewall Rule",
					"changedAt": "2024-01-12T11:00:00Z",
					"changedBy": "admin3@company.com"
				}
			]
		}`

		var response customeraudit.AuditLogsResponse
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, 3, response.Cursor.TotalPages)
		assert.Equal(t, 125, response.Cursor.TotalElements)
		assert.Len(t, response.Logs, 3)
		assert.Equal(t, "CREATE", response.Logs[0].Action.Action)
		assert.Equal(t, "UPDATE", response.Logs[1].Action.Action)
		assert.Equal(t, "DELETE", response.Logs[2].Action.Action)
	})

	t.Run("Parse audit log with detailed changes", func(t *testing.T) {
		jsonResponse := `{
			"action": {"action": "UPDATE"},
			"module": "SSL_INSPECTION",
			"resource": "SSL Policy",
			"changedAt": "2024-02-01T08:00:00Z",
			"changedBy": "sslteam@company.com",
			"oldRowJson": "{\"inspectAll\": false, \"bypassCategories\": [\"finance\", \"health\"]}",
			"newRowJson": "{\"inspectAll\": true, \"bypassCategories\": [\"health\"]}",
			"changeNote": "Expanded SSL inspection scope, removed finance bypass"
		}`

		var log customeraudit.AuditLog
		err := json.Unmarshal([]byte(jsonResponse), &log)
		require.NoError(t, err)

		assert.Equal(t, "SSL_INSPECTION", log.Module)
		assert.Contains(t, log.OldRowJSON, "inspectAll")
		assert.Contains(t, log.NewRowJSON, "inspectAll")
		assert.Equal(t, "Expanded SSL inspection scope, removed finance bypass", log.ChangeNote)
	})
}

