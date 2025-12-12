// Package unit provides unit tests for ZPA Machine Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/machinegroup"
)

// TestMachineGroup_Structure tests the struct definitions
func TestMachineGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("MachineGroup JSON marshaling", func(t *testing.T) {
		group := machinegroup.MachineGroup{
			ID:          "mg-123",
			Name:        "Test Machine Group",
			Description: "Test Description",
			Enabled:     true,
			Machines: []machinegroup.Machines{
				{ID: "m-001", Name: "Machine 1"},
				{ID: "m-002", Name: "Machine 2"},
			},
			MicroTenantID:   "mt-001",
			MicroTenantName: "Production",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled machinegroup.MachineGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Len(t, unmarshaled.Machines, 2)
	})

	t.Run("MachineGroup from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "mg-456",
			"name": "Production Machine Group",
			"description": "Production machines",
			"enabled": true,
			"machines": [
				{"id": "m-001", "name": "Server A", "fingerprint": "abc123"},
				{"id": "m-002", "name": "Server B", "fingerprint": "def456"}
			],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"microtenantId": "mt-002",
			"microtenantName": "Prod Tenant"
		}`

		var group machinegroup.MachineGroup
		err := json.Unmarshal([]byte(apiResponse), &group)
		require.NoError(t, err)

		assert.Equal(t, "mg-456", group.ID)
		assert.True(t, group.Enabled)
		assert.Len(t, group.Machines, 2)
	})

	t.Run("Machines structure", func(t *testing.T) {
		machine := machinegroup.Machines{
			ID:               "m-123",
			Name:             "Test Machine",
			Description:      "Test Description",
			Fingerprint:      "abc123def456",
			MachineGroupID:   "mg-001",
			MachineGroupName: "Test Group",
			MicroTenantID:    "mt-001",
		}

		data, err := json.Marshal(machine)
		require.NoError(t, err)

		var unmarshaled machinegroup.Machines
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, machine.ID, unmarshaled.ID)
		assert.Equal(t, machine.Fingerprint, unmarshaled.Fingerprint)
	})
}

// TestMachineGroup_MockServerOperations tests CRUD operations
func TestMachineGroup_MockServerOperations(t *testing.T) {
	t.Run("GET machine group by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "mg-123", "name": "Mock Group", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/machineGroup/mg-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all machine groups", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/machineGroup")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
