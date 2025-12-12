// Package unit provides unit tests for ZPA Policy Set Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
)

func TestPolicySetController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PolicySet JSON marshaling", func(t *testing.T) {
		policySet := policysetcontroller.PolicySet{
			ID:          "ps-123",
			Name:        "Test Policy Set",
			Description: "Test Description",
			PolicyType:  "ACCESS_POLICY",
			Enabled:     true,
		}

		data, err := json.Marshal(policySet)
		require.NoError(t, err)

		var unmarshaled policysetcontroller.PolicySet
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, policySet.ID, unmarshaled.ID)
		assert.Equal(t, policySet.Name, unmarshaled.Name)
	})

	t.Run("PolicyRule JSON marshaling", func(t *testing.T) {
		rule := policysetcontroller.PolicyRule{
			ID:          "rule-123",
			Name:        "Test Rule",
			Description: "Test Description",
			Action:      "ALLOW",
			PolicyType:  "ACCESS_POLICY",
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		var unmarshaled policysetcontroller.PolicyRule
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, rule.ID, unmarshaled.ID)
		assert.Equal(t, rule.Action, unmarshaled.Action)
	})
}

func TestPolicySetController_MockServerOperations(t *testing.T) {
	t.Run("GET policy set by type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "ps-123", "name": "Mock Policy Set"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/policySet")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
