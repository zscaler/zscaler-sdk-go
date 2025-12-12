// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/intermediatecacertificates"
)

func TestIntermediateCACertificates_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IntermediateCACertificate JSON marshaling", func(t *testing.T) {
		cert := intermediatecacertificates.IntermediateCACertificate{
			ID:                         12345,
			Name:                       "Zscaler Intermediate CA",
			Description:                "Primary intermediate CA certificate",
			Type:                       "ZSCALER",
			Region:                     "US_WEST",
			Status:                     "ENABLED",
			DefaultCertificate:         true,
			CertStartDate:              1699000000,
			CertExpDate:                1730536000,
			CurrentState:               "READY",
			KeyGenerationTime:          1698900000,
			HSMAttestationVerifiedTime: 1698950000,
			CSRFileName:                "intermediate_ca.csr",
			CSRGenerationTime:          1698960000,
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"ZSCALER"`)
		assert.Contains(t, string(data), `"defaultCertificate":true`)
	})

	t.Run("IntermediateCACertificate JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Custom Intermediate CA",
			"description": "Custom cloud HSM protected certificate",
			"type": "CLOUD_HSM",
			"region": "EU_WEST",
			"status": "ENABLED",
			"defaultCertificate": false,
			"certStartDate": 1699000000,
			"certExpDate": 1730536000,
			"currentState": "CSR_GENERATED",
			"publicKey": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...",
			"keyGenerationTime": 1698900000,
			"csrFileName": "custom_ca.csr",
			"csrGenerationTime": 1698960000
		}`

		var cert intermediatecacertificates.IntermediateCACertificate
		err := json.Unmarshal([]byte(jsonData), &cert)
		require.NoError(t, err)

		assert.Equal(t, 54321, cert.ID)
		assert.Equal(t, "CLOUD_HSM", cert.Type)
		assert.False(t, cert.DefaultCertificate)
	})

	t.Run("CertSigningRequest JSON marshaling", func(t *testing.T) {
		csr := intermediatecacertificates.CertSigningRequest{
			CertID:               12345,
			CSRFileName:          "csr_request.csr",
			CommName:             "company.com",
			ORGName:              "Company Inc.",
			DeptName:             "IT Security",
			City:                 "San Jose",
			State:                "California",
			Country:              "US",
			KeySize:              2048,
			SignatureAlgorithm:   "SHA256",
			PathLengthConstraint: 0,
		}

		data, err := json.Marshal(csr)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"certId":12345`)
		assert.Contains(t, string(data), `"commName":"company.com"`)
		assert.Contains(t, string(data), `"keySize":2048`)
	})

	t.Run("CertSigningRequest JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"certId": 67890,
			"csrFileName": "corporate_csr.csr",
			"commName": "example.com",
			"orgName": "Example Corp",
			"deptName": "InfoSec",
			"city": "Seattle",
			"state": "Washington",
			"country": "US",
			"keySize": 4096,
			"signatureAlgorithm": "SHA384",
			"pathLengthConstraint": 1
		}`

		var csr intermediatecacertificates.CertSigningRequest
		err := json.Unmarshal([]byte(jsonData), &csr)
		require.NoError(t, err)

		assert.Equal(t, 67890, csr.CertID)
		assert.Equal(t, 4096, csr.KeySize)
		assert.Equal(t, "SHA384", csr.SignatureAlgorithm)
	})
}

func TestIntermediateCACertificates_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse certificates list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Zscaler CA", "type": "ZSCALER", "status": "ENABLED", "defaultCertificate": true},
			{"id": 2, "name": "Custom CA 1", "type": "SOFTWARE", "status": "ENABLED", "defaultCertificate": false},
			{"id": 3, "name": "Custom CA 2", "type": "CLOUD_HSM", "status": "DISABLED", "defaultCertificate": false}
		]`

		var certs []intermediatecacertificates.IntermediateCACertificate
		err := json.Unmarshal([]byte(jsonResponse), &certs)
		require.NoError(t, err)

		assert.Len(t, certs, 3)
		assert.True(t, certs[0].DefaultCertificate)
		assert.Equal(t, "CLOUD_HSM", certs[2].Type)
	})
}

