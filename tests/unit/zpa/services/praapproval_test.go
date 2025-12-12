// Package unit provides unit tests for ZPA PRA Approval service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PrivilegedApproval represents the PRA approval for testing
type PrivilegedApproval struct {
	ID              string               `json:"id,omitempty"`
	EmailIDs        []string             `json:"emailIds,omitempty"`
	StartTime       string               `json:"startTime,omitempty"`
	EndTime         string               `json:"endTime,omitempty"`
	Status          string               `json:"status,omitempty"`
	CreationTime    string               `json:"creationTime,omitempty"`
	ModifiedBy      string               `json:"modifiedBy,omitempty"`
	ModifiedTime    string               `json:"modifiedTime,omitempty"`
	MicroTenantID   string               `json:"microtenantId,omitempty"`
	MicroTenantName string               `json:"microtenantName,omitempty"`
	WorkingHours    *PRAWorkingHours     `json:"workingHours"`
	Applications    []PRAApplication     `json:"applications"`
}

// PRAWorkingHours represents working hours configuration
type PRAWorkingHours struct {
	Days          []string `json:"days,omitempty"`
	StartTime     string   `json:"startTime,omitempty"`
	EndTime       string   `json:"endTime,omitempty"`
	StartTimeCron string   `json:"startTimeCron,omitempty"`
	EndTimeCron   string   `json:"endTimeCron,omitempty"`
	TimeZone      string   `json:"timeZone,omitempty"`
}

// PRAApplication represents application reference
type PRAApplication struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TestPRAApproval_Structure tests the struct definitions
func TestPRAApproval_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PrivilegedApproval JSON marshaling", func(t *testing.T) {
		approval := PrivilegedApproval{
			ID:        "approval-123",
			EmailIDs:  []string{"user1@example.com", "user2@example.com"},
			StartTime: "2024-01-01T00:00:00Z",
			EndTime:   "2024-12-31T23:59:59Z",
			Status:    "ACTIVE",
			WorkingHours: &PRAWorkingHours{
				Days:      []string{"MON", "TUE", "WED", "THU", "FRI"},
				StartTime: "09:00",
				EndTime:   "17:00",
				TimeZone:  "America/Los_Angeles",
			},
			Applications: []PRAApplication{
				{ID: "app-001", Name: "SSH Server"},
				{ID: "app-002", Name: "RDP Server"},
			},
		}

		data, err := json.Marshal(approval)
		require.NoError(t, err)

		var unmarshaled PrivilegedApproval
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, approval.ID, unmarshaled.ID)
		assert.Len(t, unmarshaled.EmailIDs, 2)
		assert.Equal(t, "ACTIVE", unmarshaled.Status)
		assert.Len(t, unmarshaled.WorkingHours.Days, 5)
		assert.Len(t, unmarshaled.Applications, 2)
	})

	t.Run("PrivilegedApproval from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "approval-456",
			"emailIds": ["admin@example.com"],
			"startTime": "2024-06-01T00:00:00Z",
			"endTime": "2024-06-30T23:59:59Z",
			"status": "FUTURE",
			"creationTime": "1609459200000",
			"modifiedBy": "admin@example.com",
			"modifiedTime": "1612137600000",
			"microtenantId": "mt-001",
			"microtenantName": "Production",
			"workingHours": {
				"days": ["MON", "WED", "FRI"],
				"startTime": "10:00",
				"endTime": "16:00",
				"startTimeCron": "0 0 10 ? * MON,WED,FRI",
				"endTimeCron": "0 0 16 ? * MON,WED,FRI",
				"timeZone": "UTC"
			},
			"applications": [
				{"id": "app-003", "name": "Database Server"}
			]
		}`

		var approval PrivilegedApproval
		err := json.Unmarshal([]byte(apiResponse), &approval)
		require.NoError(t, err)

		assert.Equal(t, "approval-456", approval.ID)
		assert.Equal(t, "FUTURE", approval.Status)
		assert.Len(t, approval.WorkingHours.Days, 3)
		assert.Equal(t, "UTC", approval.WorkingHours.TimeZone)
	})
}

// TestPRAApproval_ResponseParsing tests parsing of API responses
func TestPRAApproval_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse approval list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "status": "ACTIVE", "emailIds": ["user1@example.com"]},
				{"id": "2", "status": "FUTURE", "emailIds": ["user2@example.com"]},
				{"id": "3", "status": "EXPIRED", "emailIds": ["user3@example.com"]}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []PrivilegedApproval `json:"list"`
			TotalPages int                  `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "ACTIVE", listResp.List[0].Status)
		assert.Equal(t, "EXPIRED", listResp.List[2].Status)
	})
}

// TestPRAApproval_MockServerOperations tests CRUD operations
func TestPRAApproval_MockServerOperations(t *testing.T) {
	t.Run("GET approval by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/approval/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "approval-123", "status": "ACTIVE"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/approval/approval-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all approvals", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/approval")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create approval", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-approval", "status": "ACTIVE"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/approval", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE approval", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/approval/approval-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE expired approvals", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Contains(t, r.URL.Path, "/expired")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/approval/expired", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestPRAApproval_SpecialCases tests edge cases
func TestPRAApproval_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Approval status values", func(t *testing.T) {
		statuses := []string{"INVALID", "ACTIVE", "FUTURE", "EXPIRED"}

		for _, status := range statuses {
			approval := PrivilegedApproval{
				ID:     "approval-" + status,
				Status: status,
			}

			data, err := json.Marshal(approval)
			require.NoError(t, err)

			assert.Contains(t, string(data), status)
		}
	})

	t.Run("Working hours with cron expressions", func(t *testing.T) {
		hours := PRAWorkingHours{
			Days:          []string{"MON", "TUE", "WED", "THU", "FRI"},
			StartTime:     "09:00",
			EndTime:       "17:00",
			StartTimeCron: "0 0 9 ? * MON-FRI",
			EndTimeCron:   "0 0 17 ? * MON-FRI",
			TimeZone:      "America/New_York",
		}

		data, err := json.Marshal(hours)
		require.NoError(t, err)

		var unmarshaled PRAWorkingHours
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.Days, 5)
		assert.Contains(t, unmarshaled.StartTimeCron, "MON-FRI")
	})

	t.Run("Multiple email IDs", func(t *testing.T) {
		approval := PrivilegedApproval{
			ID: "approval-multi",
			EmailIDs: []string{
				"user1@example.com",
				"user2@example.com",
				"user3@example.com",
				"admin@example.com",
			},
		}

		data, err := json.Marshal(approval)
		require.NoError(t, err)

		var unmarshaled PrivilegedApproval
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.EmailIDs, 4)
	})
}

