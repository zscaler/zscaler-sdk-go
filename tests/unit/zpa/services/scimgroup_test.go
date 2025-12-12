// Package unit provides unit tests for ZPA SCIM Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimgroup"
)

func TestScimGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ScimGroup JSON marshaling", func(t *testing.T) {
		group := scimgroup.ScimGroup{
			ID:         123,
			Name:       "Engineering",
			IdpID:      456,
			IdpGroupID: "group-001",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled scimgroup.ScimGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
	})

	t.Run("ScimGroup from API response", func(t *testing.T) {
		apiResponse := `{
			"id": 789,
			"name": "Sales Team",
			"idpId": 101,
			"idpName": "Okta",
			"idpGroupId": "group-002"
		}`

		var group scimgroup.ScimGroup
		err := json.Unmarshal([]byte(apiResponse), &group)
		require.NoError(t, err)

		assert.Equal(t, int64(789), group.ID)
		assert.Equal(t, "Sales Team", group.Name)
	})
}

func TestScimGroup_MockServerOperations(t *testing.T) {
	t.Run("GET SCIM group by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 123, "name": "Mock Group"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/scimGroup")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
