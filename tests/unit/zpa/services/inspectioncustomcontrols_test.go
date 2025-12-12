// Package unit provides unit tests for ZPA Inspection Custom Controls service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_custom_controls"
)

func TestInspectionCustomControls_Structure(t *testing.T) {
	t.Parallel()

	t.Run("InspectionCustomControl JSON marshaling", func(t *testing.T) {
		control := inspection_custom_controls.InspectionCustomControl{
			ID:          "icc-123",
			Name:        "Test Custom Control",
			Description: "Test Description",
			Action:      "BLOCK",
			Severity:    "CRITICAL",
			Type:        "REQUEST",
		}

		data, err := json.Marshal(control)
		require.NoError(t, err)

		var unmarshaled inspection_custom_controls.InspectionCustomControl
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, control.ID, unmarshaled.ID)
		assert.Equal(t, control.Name, unmarshaled.Name)
	})
}

func TestInspectionCustomControls_MockServerOperations(t *testing.T) {
	t.Run("GET custom control by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "icc-123", "name": "Mock Control"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/inspectionCustomControl")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
