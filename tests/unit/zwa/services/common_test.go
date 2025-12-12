// Package services provides unit tests for ZWA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services/common"
)

func TestZWACommon_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IncidentDetails JSON marshaling", func(t *testing.T) {
		incident := common.IncidentDetails{
			InternalID:            "inc-12345",
			IntegrationType:       "ZIA",
			TransactionID:         "txn-67890",
			SourceType:            "EMAIL",
			SourceSubType:         "OUTBOUND",
			SourceActions:         []string{"BLOCK", "QUARANTINE"},
			Severity:              "HIGH",
			Priority:              "P1",
			MatchCount:            5,
			CreatedAt:             "2024-01-15T10:30:00Z",
			LastUpdatedAt:         "2024-01-15T11:00:00Z",
			Status:                "OPEN",
			Resolution:            "",
			IncidentGroupIDs:      []int{1, 2, 3},
		}

		data, err := json.Marshal(incident)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"internalId":"inc-12345"`)
		assert.Contains(t, string(data), `"severity":"HIGH"`)
		assert.Contains(t, string(data), `"status":"OPEN"`)
	})

	t.Run("IncidentDetails JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"internalId": "inc-99999",
			"integrationType": "ZPA",
			"transactionId": "txn-11111",
			"sourceType": "WEB",
			"sourceSubType": "DOWNLOAD",
			"sourceActions": ["ALLOW", "LOG"],
			"severity": "MEDIUM",
			"priority": "P2",
			"matchCount": 3,
			"createdAt": "2024-02-01T08:00:00Z",
			"lastUpdatedAt": "2024-02-01T09:00:00Z",
			"status": "CLOSED",
			"resolution": "FALSE_POSITIVE",
			"userInfo": {
				"name": "John Doe",
				"email": "john.doe@company.com",
				"clientIP": "192.168.1.100",
				"userId": 12345,
				"department": "Engineering"
			},
			"applicationInfo": {
				"url": "https://example.com/download",
				"category": "Cloud Storage",
				"name": "Example Cloud"
			},
			"labels": [
				{"key": "reviewed", "value": "true"},
				{"key": "team", "value": "security"}
			]
		}`

		var incident common.IncidentDetails
		err := json.Unmarshal([]byte(jsonData), &incident)
		require.NoError(t, err)

		assert.Equal(t, "inc-99999", incident.InternalID)
		assert.Equal(t, "MEDIUM", incident.Severity)
		assert.Equal(t, "CLOSED", incident.Status)
		assert.Equal(t, "John Doe", incident.UserInfo.Name)
		assert.Len(t, incident.Labels, 2)
	})

	t.Run("MatchingPolicies JSON marshaling", func(t *testing.T) {
		policies := common.MatchingPolicies{
			Engines: []common.Engine{
				{Name: "DLP Engine", Rule: "SSN Detection"},
			},
			Rules: []common.Rule{
				{Name: "PCI-DSS Rule"},
				{Name: "HIPAA Rule"},
			},
			Dictionaries: []common.Dictionary{
				{Name: "Credit Card Numbers", MatchCount: 10, NameMatchCount: "10"},
			},
		}

		data, err := json.Marshal(policies)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"engines"`)
		assert.Contains(t, string(data), `"DLP Engine"`)
		assert.Contains(t, string(data), `"PCI-DSS Rule"`)
	})

	t.Run("UserInfo JSON marshaling", func(t *testing.T) {
		user := common.UserInfo{
			Name:             "Jane Smith",
			Email:            "jane.smith@company.com",
			ClientIP:         "10.0.0.50",
			UniqueIdentifier: "uid-12345",
			UserID:           67890,
			Department:       "Sales",
			HomeCountry:      "US",
			ManagerInfo: common.ManagerInfo{
				ID:    100,
				Name:  "Manager Name",
				Email: "manager@company.com",
			},
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"Jane Smith"`)
		assert.Contains(t, string(data), `"department":"Sales"`)
		assert.Contains(t, string(data), `"managerInfo"`)
	})

	t.Run("Label JSON marshaling", func(t *testing.T) {
		label := common.Label{
			Key:   "priority",
			Value: "high",
		}

		data, err := json.Marshal(label)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"key":"priority"`)
		assert.Contains(t, string(data), `"value":"high"`)
	})

	t.Run("Cursor JSON marshaling", func(t *testing.T) {
		cursor := common.Cursor{
			TotalPages:        10,
			CurrentPageNumber: 3,
			CurrentPageSize:   100,
			PageID:            "page-abc-123",
			TotalElements:     950,
		}

		data, err := json.Marshal(cursor)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"totalPages":10`)
		assert.Contains(t, string(data), `"currentPageNumber":3`)
		assert.Contains(t, string(data), `"totalElements":950`)
	})

	t.Run("CommonDLPIncidentFiltering JSON marshaling", func(t *testing.T) {
		filter := common.CommonDLPIncidentFiltering{
			Fields: []common.Fields{
				{Name: "severity", Value: []string{"HIGH", "CRITICAL"}},
				{Name: "status", Value: []string{"OPEN"}},
			},
			TimeRange: common.TimeRange{
				StartTime: "2024-01-01T00:00:00Z",
				EndTime:   "2024-01-31T23:59:59Z",
			},
		}

		data, err := json.Marshal(filter)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"fields"`)
		assert.Contains(t, string(data), `"severity"`)
		assert.Contains(t, string(data), `"timeRange"`)
	})

	t.Run("PaginationParams pointer helpers", func(t *testing.T) {
		page := common.IntPtr(5)
		assert.Equal(t, 5, *page)

		pageSize := common.GetPageSize()
		assert.Equal(t, 1000, pageSize)
	})
}

func TestZWACommon_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse incidents list with cursor", func(t *testing.T) {
		jsonResponse := `{
			"logs": [
				{
					"internalId": "inc-1",
					"severity": "HIGH",
					"status": "OPEN"
				},
				{
					"internalId": "inc-2",
					"severity": "MEDIUM",
					"status": "CLOSED"
				}
			],
			"cursor": {
				"totalPages": 5,
				"currentPageNumber": 1,
				"currentPageSize": 2,
				"totalElements": 10
			}
		}`

		var response struct {
			Logs   []common.IncidentDetails `json:"logs"`
			Cursor common.Cursor            `json:"cursor"`
		}

		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Len(t, response.Logs, 2)
		assert.Equal(t, "inc-1", response.Logs[0].InternalID)
		assert.Equal(t, 5, response.Cursor.TotalPages)
	})
}

