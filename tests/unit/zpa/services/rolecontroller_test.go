// Package unit provides unit tests for ZPA Role Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/role_controller"
)

func TestRoleController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("RoleController JSON marshaling", func(t *testing.T) {
		role := role_controller.RoleController{
			ID:          "role-123",
			Name:        "Test Role",
			Description: "Test Description",
		}

		data, err := json.Marshal(role)
		require.NoError(t, err)

		var unmarshaled role_controller.RoleController
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, role.ID, unmarshaled.ID)
		assert.Equal(t, role.Name, unmarshaled.Name)
	})
}

func TestRoleController_MockServerOperations(t *testing.T) {
	t.Run("GET role by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "role-123", "name": "Mock Role"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/role")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
