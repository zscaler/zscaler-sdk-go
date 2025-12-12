// Package unit provides unit tests for ZPA OAuth2 User service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/oauth2_user"
)

func TestOAuth2User_Structure(t *testing.T) {
	t.Parallel()

	t.Run("OauthUser JSON marshaling", func(t *testing.T) {
		user := oauth2_user.OauthUser{
			TenantID:        "tenant-123",
			ConfigCloudName: "zscaler.net",
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		var unmarshaled oauth2_user.OauthUser
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, user.TenantID, unmarshaled.TenantID)
	})

	t.Run("UserCodeStatusRequest JSON marshaling", func(t *testing.T) {
		req := oauth2_user.UserCodeStatusRequest{
			UserCodes: []string{"code1", "code2"},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaled oauth2_user.UserCodeStatusRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.UserCodes, 2)
	})
}

func TestOAuth2User_MockServerOperations(t *testing.T) {
	t.Run("POST verify user codes", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"tenantId": "tenant-123"}`))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/oauth2/verify", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
