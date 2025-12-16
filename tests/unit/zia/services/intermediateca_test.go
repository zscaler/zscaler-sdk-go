// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/intermediatecacertificates"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestIntermediateCA_GetCertificate_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := 12345
	path := "/zia/api/v1/intermediateCaCertificate/12345"

	server.On("GET", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:   certID,
		Name: "Zscaler Intermediate CA",
		Type: "ZSCALER",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := intermediatecacertificates.GetCertificate(context.Background(), service, certID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, certID, result.ID)
}

func TestIntermediateCA_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certName := "Zscaler Intermediate CA"
	path := "/zia/api/v1/intermediateCaCertificate"

	server.On("GET", path, common.SuccessResponse([]intermediatecacertificates.IntermediateCACertificate{
		{ID: 1, Name: "Other CA", Type: "SOFTWARE"},
		{ID: 2, Name: certName, Type: "ZSCALER"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := intermediatecacertificates.GetByName(context.Background(), service, certName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
}

func TestIntermediateCA_GetIntCAReadyToUse_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/intermediateCaCertificate/readyToUse"

	server.On("GET", path, common.SuccessResponse([]intermediatecacertificates.IntermediateCACertificate{
		{ID: 1, Name: "CA 1", CurrentState: "READY"},
		{ID: 2, Name: "CA 2", CurrentState: "READY"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := intermediatecacertificates.GetIntCAReadyToUse(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestIntermediateCA_GetShowCert_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := 12345
	path := "/zia/api/v1/intermediateCaCertificate/showCert/12345"

	server.On("GET", path, common.SuccessResponse(intermediatecacertificates.CertSigningRequest{
		CertID:   certID,
		CommName: "company.com",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := intermediatecacertificates.GetShowCert(context.Background(), service, certID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, certID, result.CertID)
}

func TestIntermediateCA_GetShowCSR_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := 12345
	path := "/zia/api/v1/intermediateCaCertificate/showCsr/12345"

	server.On("GET", path, common.SuccessResponse(intermediatecacertificates.CertSigningRequest{
		CertID:      certID,
		CSRFileName: "csr_request.csr",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := intermediatecacertificates.GetShowCSR(context.Background(), service, certID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, certID, result.CertID)
}

func TestIntermediateCA_GetDownloadAttestation_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := 12345
	path := "/zia/api/v1/intermediateCaCertificate/downloadAttestation/12345"

	server.On("GET", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:   certID,
		Name: "Attestation Data",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := intermediatecacertificates.GetDownloadAttestation(context.Background(), service, certID)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntermediateCA_GetDownloadCSR_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := 12345
	path := "/zia/api/v1/intermediateCaCertificate/downloadCsr/12345"

	server.On("GET", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID: certID,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := intermediatecacertificates.GetDownloadCSR(context.Background(), service, certID)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntermediateCA_GetDownloadPublicKey_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := 12345
	path := "/zia/api/v1/intermediateCaCertificate/downloadPublicKey/12345"

	server.On("GET", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:        certID,
		PublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := intermediatecacertificates.GetDownloadPublicKey(context.Background(), service, certID)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntermediateCA_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/intermediateCaCertificate"

	server.On("GET", path, common.SuccessResponse([]intermediatecacertificates.IntermediateCACertificate{
		{ID: 1, Name: "Zscaler CA", Type: "ZSCALER"},
		{ID: 2, Name: "Custom CA", Type: "SOFTWARE"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := intermediatecacertificates.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestIntermediateCA_CreateIntCACertificate_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/intermediateCaCertificate"

	server.On("POST", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:   99999,
		Name: "New Intermediate CA",
		Type: "SOFTWARE",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newCert := &intermediatecacertificates.IntermediateCACertificate{
		Name: "New Intermediate CA",
		Type: "SOFTWARE",
	}

	result, err := intermediatecacertificates.CreateIntCACertificate(context.Background(), service, newCert)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestIntermediateCA_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := 12345
	path := "/zia/api/v1/intermediateCaCertificate/12345"

	server.On("PUT", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:   certID,
		Name: "Updated CA",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateCert := &intermediatecacertificates.IntermediateCACertificate{
		ID:   certID,
		Name: "Updated CA",
	}

	result, err := intermediatecacertificates.Update(context.Background(), service, certID, updateCert)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated CA", result.Name)
}

func TestIntermediateCA_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := 12345
	path := "/zia/api/v1/intermediateCaCertificate/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = intermediatecacertificates.Delete(context.Background(), service, certID)

	require.NoError(t, err)
}

func TestIntermediateCA_UpdateMakeDefault_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := 12345
	path := "/zia/api/v1/intermediateCaCertificate/makeDefault/12345"

	server.On("PUT", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:                 certID,
		Name:               "Default CA",
		DefaultCertificate: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateCert := &intermediatecacertificates.IntermediateCACertificate{
		Name: "Default CA",
	}

	result, err := intermediatecacertificates.UpdateMakeDefault(context.Background(), service, certID, updateCert)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.DefaultCertificate)
}

func TestIntermediateCA_CreateIntCAGenerateCSR_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/intermediateCaCertificate/generateCsr"

	server.On("POST", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:   99999,
		Name: "CSR Generated",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	csrReq := &intermediatecacertificates.IntermediateCACertificate{
		Name: "CSR Cert",
	}

	result, err := intermediatecacertificates.CreateIntCAGenerateCSR(context.Background(), service, csrReq)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntermediateCA_CreateIntCAKeyPair_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/intermediateCaCertificate/keyPair"

	server.On("POST", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:   99999,
		Name: "Key Pair Generated",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	keyPairReq := &intermediatecacertificates.IntermediateCACertificate{
		Name: "Key Pair Cert",
		Type: "CLOUD_HSM",
	}

	result, err := intermediatecacertificates.CreateIntCAKeyPair(context.Background(), service, keyPairReq)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntermediateCA_CreateIntCAFinalizeCert_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/intermediateCaCertificate/finalizeCert"

	server.On("POST", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:           99999,
		Name:         "Finalized Cert",
		CurrentState: "READY",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	finalizeReq := &intermediatecacertificates.IntermediateCACertificate{
		Name: "Finalize Cert",
	}

	result, err := intermediatecacertificates.CreateIntCAFinalizeCert(context.Background(), service, finalizeReq)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntermediateCA_CreateUploadCert_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/intermediateCaCertificate/uploadCert"

	server.On("POST", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:   99999,
		Name: "Uploaded Cert",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	uploadReq := &intermediatecacertificates.IntermediateCACertificate{
		Name: "Upload Cert",
	}

	result, err := intermediatecacertificates.CreateUploadCert(context.Background(), service, uploadReq)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntermediateCA_CreateUploadCertChain_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/intermediateCaCertificate/uploadCertChain"

	server.On("POST", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:   99999,
		Name: "Uploaded Cert Chain",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	uploadReq := &intermediatecacertificates.IntermediateCACertificate{
		Name: "Upload Cert Chain",
	}

	result, err := intermediatecacertificates.CreateUploadCertChain(context.Background(), service, uploadReq)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestIntermediateCA_CreateVerifyKeyAttestation_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/intermediateCaCertificate/verifyKeyAttestation"

	server.On("POST", path, common.SuccessResponse(intermediatecacertificates.IntermediateCACertificate{
		ID:                         99999,
		Name:                       "Verified Attestation",
		HSMAttestationVerifiedTime: 1699000000,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	verifyReq := &intermediatecacertificates.IntermediateCACertificate{
		Name: "Verify Attestation",
	}

	result, err := intermediatecacertificates.CreateVerifyKeyAttestation(context.Background(), service, verifyReq)

	require.NoError(t, err)
	require.NotNil(t, result)
}

// =====================================================
// Structure Tests
// =====================================================

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

