// Package services provides unit tests for ZWA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services/dlp_incidents"
)

func TestDLPIncidents_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CreateNoteRequest JSON marshaling", func(t *testing.T) {
		request := dlp_incidents.CreateNoteRequest{
			Notes: "This incident has been reviewed and confirmed as legitimate.",
		}

		data, err := json.Marshal(request)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"notes":"This incident has been reviewed and confirmed as legitimate."`)
	})

	t.Run("ResolutionDetailsRequest JSON marshaling", func(t *testing.T) {
		request := dlp_incidents.ResolutionDetailsRequest{
			ResolutionLabel: common.Label{
				Key:   "resolution",
				Value: "false_positive",
			},
			ResolutionCode: "FP_001",
			Notes:          "Confirmed as false positive after investigation",
		}

		data, err := json.Marshal(request)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"resolutionLabel"`)
		assert.Contains(t, string(data), `"resolutionCode":"FP_001"`)
		assert.Contains(t, string(data), `"notes"`)
	})

	t.Run("LabelsRequest JSON marshaling", func(t *testing.T) {
		request := dlp_incidents.LabelsRequest{
			Labels: []common.Label{
				{Key: "priority", Value: "high"},
				{Key: "team", Value: "security"},
				{Key: "reviewed", Value: "true"},
			},
		}

		data, err := json.Marshal(request)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"labels"`)
		assert.Contains(t, string(data), `"priority"`)
		assert.Contains(t, string(data), `"security"`)
	})

	t.Run("IncidentGroupRequest JSON marshaling", func(t *testing.T) {
		request := dlp_incidents.IncidentGroupRequest{
			IncidentGroupIDs: []int{1, 2, 3, 4, 5},
		}

		data, err := json.Marshal(request)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"incidentGroupIds":[1,2,3,4,5]`)
	})

	t.Run("IncidentGroup JSON marshaling", func(t *testing.T) {
		group := dlp_incidents.IncidentGroup{
			ID:                              123,
			Name:                            "PCI-DSS Incidents",
			Description:                     "Incidents related to PCI-DSS compliance",
			Status:                          "ACTIVE",
			IncidentGroupType:               "COMPLIANCE",
			IsDLPIncidentGroupAlreadyMapped: true,
			IsDLPAdminConfigAlreadyMapped:   false,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":123`)
		assert.Contains(t, string(data), `"name":"PCI-DSS Incidents"`)
		assert.Contains(t, string(data), `"incidentGroupType":"COMPLIANCE"`)
	})

	t.Run("IncidentHistoryResponse JSON marshaling", func(t *testing.T) {
		response := dlp_incidents.IncidentHistoryResponse{
			IncidentID: "inc-12345",
			StartDate:  "2024-01-01T00:00:00Z",
			EndDate:    "2024-01-31T23:59:59Z",
			ChangeHistory: []dlp_incidents.ChangeHistory{
				{
					ChangeType:    "STATUS_CHANGE",
					ChangedAt:     "2024-01-15T10:00:00Z",
					ChangedByName: "admin@company.com",
					ChangeData: dlp_incidents.ChangeData{
						Before: "OPEN",
						After:  "IN_PROGRESS",
					},
					Comment: "Started investigation",
				},
			},
		}

		data, err := json.Marshal(response)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"incidentId":"inc-12345"`)
		assert.Contains(t, string(data), `"changeHistory"`)
		assert.Contains(t, string(data), `"STATUS_CHANGE"`)
	})

	t.Run("Ticket JSON marshaling", func(t *testing.T) {
		ticket := dlp_incidents.Ticket{
			TicketType:          "JIRA",
			TicketingSystemName: "Company JIRA",
			ProjectID:           "SEC",
			ProjectName:         "Security Project",
			TicketInfo: dlp_incidents.TicketInfo{
				TicketID:  "SEC-12345",
				TicketURL: "https://jira.company.com/browse/SEC-12345",
				State:     "OPEN",
			},
		}

		data, err := json.Marshal(ticket)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"ticketType":"JIRA"`)
		assert.Contains(t, string(data), `"ticketId":"SEC-12345"`)
		assert.Contains(t, string(data), `"ticketUrl"`)
	})

	t.Run("DLPIncidentEvidence JSON marshaling", func(t *testing.T) {
		evidence := dlp_incidents.DLPIncidentEvidence{
			FileName:       "confidential_data.xlsx",
			FileType:       "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			AdditionalInfo: "Contains SSN data in column C",
			EvidenceURL:    "https://evidence.zscaler.com/download/abc123",
		}

		data, err := json.Marshal(evidence)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"fileName":"confidential_data.xlsx"`)
		assert.Contains(t, string(data), `"evidenceURL"`)
	})

	t.Run("DLPIncidentTriggerData JSON marshaling", func(t *testing.T) {
		triggers := dlp_incidents.DLPIncidentTriggerData{
			"dictionary_match":   "Credit Card Numbers",
			"pattern_detected":   "4111-****-****-1111",
			"confidence_score":   "95",
			"detection_location": "Email Body",
		}

		data, err := json.Marshal(triggers)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"dictionary_match":"Credit Card Numbers"`)
		assert.Contains(t, string(data), `"confidence_score":"95"`)
	})
}

func TestDLPIncidents_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse incident groups response", func(t *testing.T) {
		jsonResponse := `{
			"incidentGroups": [
				{
					"id": 1,
					"name": "PCI-DSS",
					"description": "Payment Card Industry compliance incidents",
					"status": "ACTIVE",
					"incidentGroupType": "COMPLIANCE"
				},
				{
					"id": 2,
					"name": "HIPAA",
					"description": "Healthcare compliance incidents",
					"status": "ACTIVE",
					"incidentGroupType": "COMPLIANCE"
				},
				{
					"id": 3,
					"name": "General DLP",
					"description": "General data loss prevention incidents",
					"status": "ACTIVE",
					"incidentGroupType": "GENERAL"
				}
			]
		}`

		var response dlp_incidents.IncidentGroupsResponse
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Len(t, response.IncidentGroups, 3)
		assert.Equal(t, "PCI-DSS", response.IncidentGroups[0].Name)
		assert.Equal(t, "COMPLIANCE", response.IncidentGroups[0].IncidentGroupType)
	})

	t.Run("Parse incident history response", func(t *testing.T) {
		jsonResponse := `{
			"incidentId": "inc-99999",
			"startDate": "2024-01-01T00:00:00Z",
			"endDate": "2024-01-31T23:59:59Z",
			"changeHistory": [
				{
					"changeType": "CREATED",
					"changedAt": "2024-01-10T08:00:00Z",
					"changedByName": "system",
					"changeData": {"before": "", "after": "OPEN"},
					"comment": "Incident created automatically"
				},
				{
					"changeType": "ASSIGNED",
					"changedAt": "2024-01-10T09:00:00Z",
					"changedByName": "admin@company.com",
					"changeData": {"before": "", "after": "analyst@company.com"},
					"comment": "Assigned to security analyst"
				},
				{
					"changeType": "STATUS_CHANGE",
					"changedAt": "2024-01-11T10:00:00Z",
					"changedByName": "analyst@company.com",
					"changeData": {"before": "OPEN", "after": "CLOSED"},
					"comment": "Confirmed as false positive"
				}
			]
		}`

		var response dlp_incidents.IncidentHistoryResponse
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, "inc-99999", response.IncidentID)
		assert.Len(t, response.ChangeHistory, 3)
		assert.Equal(t, "CREATED", response.ChangeHistory[0].ChangeType)
		assert.Equal(t, "ASSIGNED", response.ChangeHistory[1].ChangeType)
		assert.Equal(t, "STATUS_CHANGE", response.ChangeHistory[2].ChangeType)
	})

	t.Run("Parse tickets response", func(t *testing.T) {
		jsonResponse := `{
			"tickets": [
				{
					"ticketType": "JIRA",
					"ticketingSystemName": "Corporate JIRA",
					"projectId": "SEC",
					"projectName": "Security",
					"ticketInfo": {
						"ticketId": "SEC-001",
						"ticketUrl": "https://jira.corp.com/SEC-001",
						"state": "CLOSED"
					}
				},
				{
					"ticketType": "SERVICENOW",
					"ticketingSystemName": "IT ServiceNow",
					"projectId": "INC",
					"projectName": "Incidents",
					"ticketInfo": {
						"ticketId": "INC0012345",
						"ticketUrl": "https://corp.service-now.com/INC0012345",
						"state": "RESOLVED"
					}
				}
			]
		}`

		var response dlp_incidents.DLPIncidentTicketsResponse
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Len(t, response.Tickets, 2)
		assert.Equal(t, "JIRA", response.Tickets[0].TicketType)
		assert.Equal(t, "SERVICENOW", response.Tickets[1].TicketType)
		assert.Equal(t, "CLOSED", response.Tickets[0].TicketInfo.State)
	})
}

