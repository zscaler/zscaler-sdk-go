// Package services provides unit tests for ZWA services
package services

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testcommon "github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
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

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

// mockZWAAuth adds the ZWA authentication endpoint mock to the server
func mockZWAAuth(server *testcommon.TestServer) {
	server.On("POST", "/v1/auth/api-key/token", testcommon.SuccessResponse(map[string]interface{}{
		"token_type": "Bearer",
		"token":      "mock-zwa-access-token",
		"expires_in": 3600,
	}))
}

func TestDLPIncidents_GetDLPIncident_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	// Mock authentication endpoint
	mockZWAAuth(server)

	incidentID := "inc-12345"
	path := "/dlp/v1/incidents/" + incidentID

	server.On("GET", path, testcommon.SuccessResponse(common.IncidentDetails{
		InternalID: incidentID,
		Status:     "OPEN",
		Severity:   "HIGH",
		Priority:   "P1",
		UserInfo: common.UserInfo{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
	}))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	result, resp, err := dlp_incidents.GetDLPIncident(context.Background(), service, incidentID, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, incidentID, result.InternalID)
	assert.Equal(t, "OPEN", result.Status)
	assert.Equal(t, "HIGH", result.Severity)
}

func TestDLPIncidents_CreateNotes_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()
	mockZWAAuth(server)

	incidentID := 12345
	path := "/dlp/v1/incidents/notes/12345"

	server.On("POST", path, testcommon.SuccessResponse(common.IncidentDetails{
		InternalID: "inc-12345",
		Status:     "IN_PROGRESS",
		Notes: []common.Note{
			{Body: "Investigation started"},
		},
	}))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	result, resp, err := dlp_incidents.CreateNotes(context.Background(), service, incidentID, "Investigation started")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "IN_PROGRESS", result.Status)
}

func TestDLPIncidents_CreateNotes_EmptyNote(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	result, resp, err := dlp_incidents.CreateNotes(context.Background(), service, 12345, "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "note is required")
}

func TestDLPIncidents_UpdateIncidentStatus_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()
	mockZWAAuth(server)

	incidentID := "inc-12345"
	path := "/dlp/v1/incidents/" + incidentID + "/close"

	server.On("POST", path, testcommon.SuccessResponse(common.IncidentDetails{
		InternalID: incidentID,
		Status:     "CLOSED",
		ClosedCode: "RESOLVED",
	}))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	result, resp, err := dlp_incidents.UpdateIncidentStatus(context.Background(), service, incidentID, "Resolved after investigation")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "CLOSED", result.Status)
}

func TestDLPIncidents_UpdateIncidentStatus_EmptyID(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	result, resp, err := dlp_incidents.UpdateIncidentStatus(context.Background(), service, "", "close")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "DLP incident ID is required")
}

func TestDLPIncidents_AssignLabels_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()
	mockZWAAuth(server)

	incidentID := "inc-12345"
	path := "/dlp/v1/incidents/" + incidentID + "/labels"

	server.On("POST", path, testcommon.SuccessResponse(common.IncidentDetails{
		InternalID: incidentID,
		Labels: []common.Label{
			{Key: "priority", Value: "high"},
			{Key: "team", Value: "security"},
		},
	}))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	labels := []common.Label{
		{Key: "priority", Value: "high"},
		{Key: "team", Value: "security"},
	}

	result, resp, err := dlp_incidents.AssignLabels(context.Background(), service, incidentID, labels)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Labels, 2)
}

func TestDLPIncidents_AssignLabels_EmptyLabels(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	result, resp, err := dlp_incidents.AssignLabels(context.Background(), service, "inc-12345", []common.Label{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "labels are required")
}

func TestDLPIncidents_DeleteDLPIncident_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()
	mockZWAAuth(server)

	incidentID := "inc-12345"
	path := "/dlp/v1/incidents/" + incidentID

	server.On("DELETE", path, testcommon.SuccessResponseWithStatus(http.StatusNoContent, nil))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	resp, err := dlp_incidents.DeleteDLPIncident(context.Background(), service, incidentID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestDLPIncidents_DeleteDLPIncident_EmptyID(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	resp, err := dlp_incidents.DeleteDLPIncident(context.Background(), service, "")
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "DLP incident ID is required")
}

func TestDLPIncidents_HistoryDLPIncident_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()
	mockZWAAuth(server)

	incidentID := "inc-12345"
	path := "/dlp/v1/incidents/" + incidentID + "/change-history"

	server.On("GET", path, testcommon.SuccessResponse(dlp_incidents.IncidentHistoryResponse{
		IncidentID: incidentID,
		StartDate:  "2024-01-01T00:00:00Z",
		EndDate:    "2024-01-31T23:59:59Z",
		ChangeHistory: []dlp_incidents.ChangeHistory{
			{
				ChangeType:    "STATUS_CHANGE",
				ChangedAt:     "2024-01-15T10:00:00Z",
				ChangedByName: "admin@company.com",
				ChangeData: dlp_incidents.ChangeData{
					Before: "OPEN",
					After:  "CLOSED",
				},
			},
		},
	}))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	result, resp, err := dlp_incidents.HistoryDLPIncident(context.Background(), service, incidentID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, incidentID, result.IncidentID)
	assert.Len(t, result.ChangeHistory, 1)
}

func TestDLPIncidents_GetDLPIncidentTriggers_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()
	mockZWAAuth(server)

	incidentID := "inc-12345"
	path := "/dlp/v1/incidents/" + incidentID + "/triggers"

	server.On("GET", path, testcommon.SuccessResponse(dlp_incidents.DLPIncidentTriggerData{
		"dictionary_match": "Credit Card Numbers",
		"pattern_detected": "4111-****-****-1111",
		"confidence_score": "95",
	}))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	result, resp, err := dlp_incidents.GetDLPIncidentTriggers(context.Background(), service, incidentID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Credit Card Numbers", result["dictionary_match"])
}

func TestDLPIncidents_GetDLPIncidentEvidence_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()
	mockZWAAuth(server)

	incidentID := "inc-12345"
	path := "/dlp/v1/incidents/" + incidentID + "/evidence"

	server.On("GET", path, testcommon.SuccessResponse(dlp_incidents.DLPIncidentEvidence{
		FileName:       "confidential.xlsx",
		FileType:       "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		AdditionalInfo: "Contains SSN data",
		EvidenceURL:    "https://evidence.example.com/abc123",
	}))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	result, resp, err := dlp_incidents.GetDLPIncidentEvidence(context.Background(), service, incidentID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "confidential.xlsx", result.FileName)
}

func TestDLPIncidents_AssignIncidentGroups_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()
	mockZWAAuth(server)

	incidentID := 12345
	path := "/dlp/v1/incidents/12345/incident-groups/search"

	server.On("POST", path, testcommon.SuccessResponse(dlp_incidents.IncidentGroupsResponse{
		IncidentGroups: []dlp_incidents.IncidentGroup{
			{ID: 1, Name: "PCI-DSS", Status: "ACTIVE"},
			{ID: 2, Name: "HIPAA", Status: "ACTIVE"},
		},
	}))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	result, resp, err := dlp_incidents.AssignIncidentGroups(context.Background(), service, incidentID, []int{1, 2})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.IncidentGroups, 2)
}

func TestDLPIncidents_AssignIncidentGroups_EmptyGroups(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	result, resp, err := dlp_incidents.AssignIncidentGroups(context.Background(), service, 12345, []int{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "incident group IDs are required")
}

func TestDLPIncidents_FilterIncidentSearch_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()
	mockZWAAuth(server)

	path := "/dlp/v1/incidents/search"

	server.On("POST", path, testcommon.SuccessResponse(map[string]interface{}{
		"logs": []common.IncidentDetails{
			{InternalID: "inc-1", Status: "OPEN", Severity: "HIGH"},
			{InternalID: "inc-2", Status: "CLOSED", Severity: "MEDIUM"},
		},
		"cursor": common.Cursor{
			TotalPages:        1,
			CurrentPageNumber: 1,
			CurrentPageSize:   100,
			TotalElements:     2,
		},
	}))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	filters := common.CommonDLPIncidentFiltering{
		TimeRange: common.TimeRange{
			StartTime: "2024-01-01T00:00:00Z",
			EndTime:   "2024-01-31T23:59:59Z",
		},
	}

	results, cursor, err := dlp_incidents.FilterIncidentSearch(context.Background(), service, filters, nil)
	require.NoError(t, err)
	require.NotNil(t, results)
	require.NotNil(t, cursor)
	assert.Len(t, results, 2)
}

func TestDLPIncidents_GetIncidentTransactions_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()
	mockZWAAuth(server)

	transactionID := "txn-12345"
	path := "/dlp/v1/incidents/transactions/" + transactionID

	server.On("GET", path, testcommon.SuccessResponse(map[string]interface{}{
		"logs": []common.IncidentDetails{
			{InternalID: "inc-1", TransactionID: transactionID, Status: "OPEN"},
		},
		"cursor": common.Cursor{
			TotalPages:        1,
			CurrentPageNumber: 1,
			CurrentPageSize:   100,
			TotalElements:     1,
		},
	}))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	results, cursor, err := dlp_incidents.GetIncidentTransactions(context.Background(), service, transactionID, nil)
	require.NoError(t, err)
	require.NotNil(t, results)
	require.NotNil(t, cursor)
	assert.Len(t, results, 1)
}

func TestDLPIncidents_GetIncidentTransactions_EmptyID(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	results, cursor, err := dlp_incidents.GetIncidentTransactions(context.Background(), service, "", nil)
	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Nil(t, cursor)
	assert.Contains(t, err.Error(), "transaction ID is required")
}

func TestDLPIncidents_GetDLPIncidentTickets_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()
	mockZWAAuth(server)

	incidentID := "inc-12345"
	path := "/dlp/v1/incidents/tickets/" + incidentID

	server.On("GET", path, testcommon.SuccessResponse(map[string]interface{}{
		"logs": []dlp_incidents.Ticket{
			{
				TicketType:          "JIRA",
				TicketingSystemName: "Company JIRA",
				ProjectID:           "SEC",
				TicketInfo: dlp_incidents.TicketInfo{
					TicketID:  "SEC-123",
					TicketURL: "https://jira.company.com/SEC-123",
					State:     "OPEN",
				},
			},
		},
		"cursor": common.Cursor{
			TotalPages:        1,
			CurrentPageNumber: 1,
			CurrentPageSize:   100,
			TotalElements:     1,
		},
	}))

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	results, cursor, err := dlp_incidents.GetDLPIncidentTickets(context.Background(), service, incidentID, nil)
	require.NoError(t, err)
	require.NotNil(t, results)
	require.NotNil(t, cursor)
	assert.Len(t, results, 1)
}

func TestDLPIncidents_GetDLPIncidentTickets_EmptyID(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	service, err := testcommon.CreateZWATestService(context.Background(), server)
	require.NoError(t, err)

	results, cursor, err := dlp_incidents.GetDLPIncidentTickets(context.Background(), service, "", nil)
	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Nil(t, cursor)
	assert.Contains(t, err.Error(), "DLP incident ID is required")
}

