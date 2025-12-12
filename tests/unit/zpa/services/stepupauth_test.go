// Package unit provides unit tests for ZPA Step Up Auth service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// StepAuthLevel represents the step auth level for testing
type StepAuthLevel struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description,omitempty"`
	Delta                string `json:"delta,omitempty"`
	IamAuthLevelID       string `json:"iamAuthLevelId,omitempty"`
	ParentIamAuthLevelID string `json:"parentIamAuthLevelId,omitempty"`
	UserMessage          string `json:"userMessage,omitempty"`
	MicrotenantID        string `json:"microtenantId,omitempty"`
	MicrotenantName      string `json:"microtenantName,omitempty"`
	CreationTime         string `json:"creationTime,omitempty"`
	ModifiedTime         string `json:"modifiedTime,omitempty"`
}

// TestStepUpAuth_Structure tests the struct definitions
func TestStepUpAuth_Structure(t *testing.T) {
	t.Parallel()

	t.Run("StepAuthLevel JSON marshaling", func(t *testing.T) {
		auth := StepAuthLevel{
			ID:             "sal-123",
			Name:           "High Security",
			Description:    "High security step-up authentication",
			Delta:          "30",
			IamAuthLevelID: "iam-001",
			UserMessage:    "Please complete additional verification",
		}

		data, err := json.Marshal(auth)
		require.NoError(t, err)

		var unmarshaled StepAuthLevel
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, auth.ID, unmarshaled.ID)
		assert.Equal(t, auth.Name, unmarshaled.Name)
	})

	t.Run("StepAuthLevel from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "sal-456",
			"name": "MFA Required",
			"description": "Multi-factor authentication required",
			"delta": "60",
			"iamAuthLevelId": "iam-002",
			"parentIamAuthLevelId": "iam-001",
			"userMessage": "MFA verification required",
			"microtenantId": "mt-001",
			"microtenantName": "Production",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var auth StepAuthLevel
		err := json.Unmarshal([]byte(apiResponse), &auth)
		require.NoError(t, err)

		assert.Equal(t, "sal-456", auth.ID)
		assert.Equal(t, "MFA Required", auth.Name)
		assert.Equal(t, "60", auth.Delta)
	})
}

// TestStepUpAuth_MockServerOperations tests operations
func TestStepUpAuth_MockServerOperations(t *testing.T) {
	t.Run("GET step up auth levels", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/stepupauthlevel")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `["level1", "level2", "level3"]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/stepupauthlevel")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestStepUpAuth_SpecialCases tests edge cases
func TestStepUpAuth_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Auth level names", func(t *testing.T) {
		levels := []string{
			"DEFAULT",
			"MFA_REQUIRED",
			"HIGH_SECURITY",
			"CERTIFICATE_REQUIRED",
		}

		for _, level := range levels {
			auth := StepAuthLevel{
				ID:   "sal-" + level,
				Name: level,
			}

			data, err := json.Marshal(auth)
			require.NoError(t, err)

			assert.Contains(t, string(data), level)
		}
	})

	t.Run("Delta values", func(t *testing.T) {
		deltas := []string{"15", "30", "60", "120", "240"}

		for _, delta := range deltas {
			auth := StepAuthLevel{
				ID:    "sal-delta",
				Name:  "Test Level",
				Delta: delta,
			}

			data, err := json.Marshal(auth)
			require.NoError(t, err)

			assert.Contains(t, string(data), delta)
		}
	})
}

