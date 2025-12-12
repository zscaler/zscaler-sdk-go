// Package unit provides unit tests for ZPA CBI Certificate Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
)

func TestCBICertificateController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CBICertificate JSON marshaling", func(t *testing.T) {
		cert := cbicertificatecontroller.CBICertificate{
			ID:          "cert-123",
			Name:        "Test Certificate",
			PEM:         "-----BEGIN CERTIFICATE-----",
			IsDefault:   false,
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		var unmarshaled cbicertificatecontroller.CBICertificate
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, cert.ID, unmarshaled.ID)
		assert.Equal(t, cert.Name, unmarshaled.Name)
	})
}

func TestCBICertificateController_MockServerOperations(t *testing.T) {
	t.Run("GET certificate by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "cert-123", "name": "Mock Cert"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbiCert")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
