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

func TestZWACommon_AdditionalStructures(t *testing.T) {
	t.Parallel()

	t.Run("ApplicationInfo JSON marshaling", func(t *testing.T) {
		app := common.ApplicationInfo{
			URL:                   "https://app.example.com/data",
			Category:              "Cloud Storage",
			Name:                  "Example Cloud App",
			HostnameOrApplication: "app.example.com",
			AdditionalInfo:        "Enterprise tier",
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"url":"https://app.example.com/data"`)
		assert.Contains(t, string(data), `"category":"Cloud Storage"`)
		assert.Contains(t, string(data), `"hostnameOrApplication":"app.example.com"`)
	})

	t.Run("ContentInfo JSON marshaling", func(t *testing.T) {
		content := common.ContentInfo{
			FileName:       "confidential_report.pdf",
			FileType:       "application/pdf",
			AdditionalInfo: "Contains SSN data",
		}

		data, err := json.Marshal(content)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"fileName":"confidential_report.pdf"`)
		assert.Contains(t, string(data), `"fileType":"application/pdf"`)
	})

	t.Run("NetworkInfo JSON marshaling", func(t *testing.T) {
		network := common.NetworkInfo{
			Source:      "10.0.0.100",
			Destination: "external.server.com",
		}

		data, err := json.Marshal(network)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"source":"10.0.0.100"`)
		assert.Contains(t, string(data), `"destination":"external.server.com"`)
	})

	t.Run("AssignedAdmin JSON marshaling", func(t *testing.T) {
		admin := common.AssignedAdmin{
			Email: "security-admin@company.com",
		}

		data, err := json.Marshal(admin)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"email":"security-admin@company.com"`)
	})

	t.Run("LastNotifiedUser JSON marshaling", func(t *testing.T) {
		user := common.LastNotifiedUser{
			Role:  "Manager",
			Email: "manager@company.com",
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"role":"Manager"`)
		assert.Contains(t, string(data), `"email":"manager@company.com"`)
	})

	t.Run("IncidentGroup JSON marshaling", func(t *testing.T) {
		group := common.IncidentGroup{
			ID:          42,
			Name:        "PCI-DSS Compliance",
			Description: "Payment Card Industry Data Security Standard incidents",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":42`)
		assert.Contains(t, string(data), `"name":"PCI-DSS Compliance"`)
	})

	t.Run("DLPIncidentTicket JSON marshaling", func(t *testing.T) {
		ticket := common.DLPIncidentTicket{
			TicketType:          "JIRA",
			TicketingSystemName: "Corporate JIRA",
			ProjectID:           "SEC",
			ProjectName:         "Security Team",
			TicketInfo: common.TicketInfo{
				TicketID:  "SEC-12345",
				TicketURL: "https://jira.company.com/browse/SEC-12345",
				State:     "IN_PROGRESS",
			},
		}

		data, err := json.Marshal(ticket)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"ticketType":"JIRA"`)
		assert.Contains(t, string(data), `"ticketId":"SEC-12345"`)
		assert.Contains(t, string(data), `"state":"IN_PROGRESS"`)
	})

	t.Run("Note JSON marshaling", func(t *testing.T) {
		note := common.Note{
			Body:          "Investigation started by security team",
			CreatedAt:     "2024-01-15T10:00:00Z",
			LastUpdatedAt: "2024-01-15T11:00:00Z",
			CreatedBy:     1001,
			LastUpdatedBy: 1002,
		}

		data, err := json.Marshal(note)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"body":"Investigation started by security team"`)
		assert.Contains(t, string(data), `"createdBy":1001`)
	})

	t.Run("Engine and Rule JSON marshaling", func(t *testing.T) {
		engine := common.Engine{
			Name: "DLP Detection Engine",
			Rule: "SSN Detection Rule",
		}

		rule := common.Rule{
			Name: "PCI-DSS Rule",
		}

		engineData, err := json.Marshal(engine)
		require.NoError(t, err)
		assert.Contains(t, string(engineData), `"name":"DLP Detection Engine"`)

		ruleData, err := json.Marshal(rule)
		require.NoError(t, err)
		assert.Contains(t, string(ruleData), `"name":"PCI-DSS Rule"`)
	})

	t.Run("Dictionary JSON marshaling", func(t *testing.T) {
		dict := common.Dictionary{
			Name:           "Credit Card Numbers",
			MatchCount:     25,
			NameMatchCount: "25",
		}

		data, err := json.Marshal(dict)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"Credit Card Numbers"`)
		assert.Contains(t, string(data), `"matchCount":25`)
	})

	t.Run("Fields JSON marshaling", func(t *testing.T) {
		fields := common.Fields{
			Name:  "severity",
			Value: []string{"HIGH", "CRITICAL", "MEDIUM"},
		}

		data, err := json.Marshal(fields)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"severity"`)
		assert.Contains(t, string(data), `["HIGH","CRITICAL","MEDIUM"]`)
	})

	t.Run("TimeRange JSON marshaling", func(t *testing.T) {
		timeRange := common.TimeRange{
			StartTime: "2024-01-01T00:00:00Z",
			EndTime:   "2024-12-31T23:59:59Z",
		}

		data, err := json.Marshal(timeRange)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"startTime":"2024-01-01T00:00:00Z"`)
		assert.Contains(t, string(data), `"endTime":"2024-12-31T23:59:59Z"`)
	})

	t.Run("ManagerInfo JSON marshaling", func(t *testing.T) {
		manager := common.ManagerInfo{
			ID:    500,
			Name:  "John Manager",
			Email: "john.manager@company.com",
		}

		data, err := json.Marshal(manager)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":500`)
		assert.Contains(t, string(data), `"name":"John Manager"`)
		assert.Contains(t, string(data), `"email":"john.manager@company.com"`)
	})
}

func TestZWACommon_ComplexParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse full IncidentDetails response", func(t *testing.T) {
		jsonResponse := `{
			"internalId": "inc-complex-001",
			"integrationType": "ZIA",
			"transactionId": "txn-987654",
			"sourceType": "EMAIL",
			"sourceSubType": "OUTBOUND",
			"sourceActions": ["BLOCK", "QUARANTINE", "NOTIFY"],
			"severity": "CRITICAL",
			"priority": "P1",
			"matchingPolicies": {
				"engines": [
					{"name": "DLP Engine V2", "rule": "Exact Data Match"}
				],
				"rules": [
					{"name": "SSN Detection"},
					{"name": "Credit Card Detection"}
				],
				"dictionaries": [
					{"name": "SSN Dictionary", "matchCount": 15, "nameMatchCount": "15"},
					{"name": "CC Dictionary", "matchCount": 8, "nameMatchCount": "8"}
				]
			},
			"matchCount": 23,
			"createdAt": "2024-03-15T10:00:00Z",
			"lastUpdatedAt": "2024-03-15T12:30:00Z",
			"userInfo": {
				"name": "Jane Doe",
				"email": "jane.doe@company.com",
				"clientIP": "192.168.1.100",
				"userId": 12345,
				"department": "Finance",
				"homeCountry": "US",
				"managerInfo": {
					"id": 100,
					"name": "Finance Manager",
					"email": "finance.manager@company.com"
				}
			},
			"applicationInfo": {
				"url": "https://email.service.com/send",
				"category": "Email",
				"name": "Corporate Email"
			},
			"contentInfo": {
				"fileName": "q1_financials.xlsx",
				"fileType": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				"additionalInfo": "Contains SSN data in column D"
			},
			"networkInfo": {
				"source": "internal.company.com",
				"destination": "external.recipient.com"
			},
			"status": "OPEN",
			"assignedAdmin": {
				"email": "security@company.com"
			},
			"notes": [
				{
					"body": "Initial investigation started",
					"createdAt": "2024-03-15T11:00:00Z",
					"createdBy": 1001
				}
			],
			"incidentGroupIds": [1, 2, 5],
			"labels": [
				{"key": "priority", "value": "urgent"},
				{"key": "team", "value": "security"},
				{"key": "compliance", "value": "pci-dss"}
			]
		}`

		var incident common.IncidentDetails
		err := json.Unmarshal([]byte(jsonResponse), &incident)
		require.NoError(t, err)

		assert.Equal(t, "inc-complex-001", incident.InternalID)
		assert.Equal(t, "CRITICAL", incident.Severity)
		assert.Equal(t, "P1", incident.Priority)
		assert.Len(t, incident.SourceActions, 3)
		assert.Len(t, incident.MatchingPolicies.Engines, 1)
		assert.Len(t, incident.MatchingPolicies.Rules, 2)
		assert.Len(t, incident.MatchingPolicies.Dictionaries, 2)
		assert.Equal(t, 23, incident.MatchCount)
		assert.Equal(t, "Jane Doe", incident.UserInfo.Name)
		assert.Equal(t, "Finance Manager", incident.UserInfo.ManagerInfo.Name)
		assert.Equal(t, "q1_financials.xlsx", incident.ContentInfo.FileName)
		assert.Len(t, incident.Notes, 1)
		assert.Len(t, incident.Labels, 3)
		assert.Len(t, incident.IncidentGroupIDs, 3)
	})
}

