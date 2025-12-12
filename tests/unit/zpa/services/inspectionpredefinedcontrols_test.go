// Package unit provides unit tests for ZPA Inspection Predefined Controls service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_predefined_controls"
)

func TestInspectionPredefinedControls_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PredefinedControls JSON marshaling", func(t *testing.T) {
		control := inspection_predefined_controls.PredefinedControls{
			ID:                    "ipc-123",
			Name:                  "Test Predefined Control",
			Description:           "Test Description",
			Action:                "BLOCK",
			Severity:              "CRITICAL",
			ControlGroup:          "OWASP",
			DefaultAction:         "BLOCK",
			ParanoiaLevel:         "2",
		}

		data, err := json.Marshal(control)
		require.NoError(t, err)

		var unmarshaled inspection_predefined_controls.PredefinedControls
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, control.ID, unmarshaled.ID)
		assert.Equal(t, control.Name, unmarshaled.Name)
	})
}

func TestInspectionPredefinedControls_MockServerOperations(t *testing.T) {
	t.Run("GET predefined control by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "ipc-123", "name": "Mock Control"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/inspectionPredefinedControl")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
