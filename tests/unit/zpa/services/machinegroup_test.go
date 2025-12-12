// Package unit provides unit tests for ZPA Machine Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MachineGroup represents the machine group structure for testing
type MachineGroup struct {
	ID              string    `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Description     string    `json:"description,omitempty"`
	Enabled         bool      `json:"enabled"`
	CreationTime    string    `json:"creationTime,omitempty"`
	ModifiedBy      string    `json:"modifiedBy,omitempty"`
	ModifiedTime    string    `json:"modifiedTime,omitempty"`
	MicroTenantID   string    `json:"microtenantId,omitempty"`
	MicroTenantName string    `json:"microtenantName,omitempty"`
	Machines        []Machine `json:"machines,omitempty"`
}

// Machine represents a machine in a machine group
type Machine struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Description         string `json:"description,omitempty"`
	MachineGroupID      string `json:"machineGroupId,omitempty"`
	MachineGroupName    string `json:"machineGroupName,omitempty"`
	MachineTokenID      string `json:"machineTokenId,omitempty"`
	Fingerprint         string `json:"fingerprint,omitempty"`
	IssuedCertID        string `json:"issuedCertId,omitempty"`
	SigningCert         string `json:"signingCert,omitempty"`
	CreationTime        string `json:"creationTime,omitempty"`
	ModifiedBy          string `json:"modifiedBy,omitempty"`
	ModifiedTime        string `json:"modifiedTime,omitempty"`
}

// TestMachineGroup_Structure tests the struct definitions
func TestMachineGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("MachineGroup JSON marshaling", func(t *testing.T) {
		group := MachineGroup{
			ID:          "mg-123",
			Name:        "Engineering Machines",
			Description: "Machine group for engineering team",
			Enabled:     true,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled MachineGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
	})

	t.Run("MachineGroup JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "mg-456",
			"name": "Production Machines",
			"description": "Production environment machines",
			"enabled": true,
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-001",
			"microtenantName": "Production",
			"machines": [
				{
					"id": "m-001",
					"name": "Workstation 1",
					"fingerprint": "ABC123"
				},
				{
					"id": "m-002",
					"name": "Workstation 2",
					"fingerprint": "DEF456"
				}
			]
		}`

		var group MachineGroup
		err := json.Unmarshal([]byte(apiResponse), &group)
		require.NoError(t, err)

		assert.Equal(t, "mg-456", group.ID)
		assert.Equal(t, "Production Machines", group.Name)
		assert.True(t, group.Enabled)
		assert.Len(t, group.Machines, 2)
	})

	t.Run("Machine structure", func(t *testing.T) {
		machine := Machine{
			ID:               "m-001",
			Name:             "Developer Workstation",
			Description:      "Primary developer machine",
			MachineGroupID:   "mg-001",
			MachineGroupName: "Engineering",
			Fingerprint:      "ABC123DEF456",
			IssuedCertID:     "cert-001",
		}

		data, err := json.Marshal(machine)
		require.NoError(t, err)

		var unmarshaled Machine
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, machine.ID, unmarshaled.ID)
		assert.Equal(t, machine.Fingerprint, unmarshaled.Fingerprint)
	})
}

// TestMachineGroup_MockServerOperations tests CRUD operations with mock server
func TestMachineGroup_MockServerOperations(t *testing.T) {
	t.Run("GET machine group by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/machineGroup/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "mg-123",
				"name": "Mock Machine Group",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/machineGroup/mg-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all machine groups", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Group A", "enabled": true},
					{"id": "2", "name": "Group B", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/machineGroup")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestMachineGroup_ErrorHandling tests error scenarios
func TestMachineGroup_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Machine Group Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "Machine group not found"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/machineGroup/nonexistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// TestMachineGroup_SpecialCases tests edge cases
func TestMachineGroup_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Machine group with multiple machines", func(t *testing.T) {
		group := MachineGroup{
			ID:      "mg-123",
			Name:    "Large Group",
			Enabled: true,
			Machines: []Machine{
				{ID: "m-1", Name: "Machine 1", Fingerprint: "fp-1"},
				{ID: "m-2", Name: "Machine 2", Fingerprint: "fp-2"},
				{ID: "m-3", Name: "Machine 3", Fingerprint: "fp-3"},
			},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled MachineGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.Machines, 3)
	})

	t.Run("Disabled machine group", func(t *testing.T) {
		group := MachineGroup{
			ID:      "mg-123",
			Name:    "Disabled Group",
			Enabled: false,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"enabled":false`)
	})
}

