// Package unit provides unit tests for ZPA Inspection Profile service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_profile"
)

func TestInspectionProfile_Structure(t *testing.T) {
	t.Parallel()

	t.Run("InspectionProfile JSON marshaling", func(t *testing.T) {
		profile := inspection_profile.InspectionProfile{
			ID:          "ip-123",
			Name:        "Test Profile",
			Description: "Test Description",
			ParanoiaLevel: "2",
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled inspection_profile.InspectionProfile
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, profile.ID, unmarshaled.ID)
		assert.Equal(t, profile.Name, unmarshaled.Name)
	})
}

func TestInspectionProfile_MockServerOperations(t *testing.T) {
	t.Run("GET profile by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "ip-123", "name": "Mock Profile"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/inspectionProfile")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
