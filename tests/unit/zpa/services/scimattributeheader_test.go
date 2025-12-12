// Package unit provides unit tests for ZPA SCIM Attribute Header service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimattributeheader"
)

func TestScimAttributeHeader_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ScimAttributeHeader JSON marshaling", func(t *testing.T) {
		header := scimattributeheader.ScimAttributeHeader{
			ID:     "sah-123",
			Name:   "email",
			IdpID:  "idp-001",
			DataType: "String",
		}

		data, err := json.Marshal(header)
		require.NoError(t, err)

		var unmarshaled scimattributeheader.ScimAttributeHeader
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, header.ID, unmarshaled.ID)
		assert.Equal(t, header.Name, unmarshaled.Name)
	})
}

func TestScimAttributeHeader_MockServerOperations(t *testing.T) {
	t.Run("GET attribute header by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "sah-123", "name": "Mock Attribute"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/scimAttributeHeader")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
