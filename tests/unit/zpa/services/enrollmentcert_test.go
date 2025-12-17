// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/enrollmentcert"
)

func TestEnrollmentCert_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := "cert-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/enrollmentCert/" + certID

	server.On("GET", path, common.SuccessResponse(enrollmentcert.EnrollmentCert{
		ID:          certID,
		Name:        "Test Certificate",
		Description: "Test description",
		AllowSigning: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := enrollmentcert.Get(context.Background(), service, certID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, certID, result.ID)
	assert.Equal(t, "Test Certificate", result.Name)
}

func TestEnrollmentCert_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certName := "Production Certificate"
	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/enrollmentCert"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []enrollmentcert.EnrollmentCert{
			{ID: "cert-001", Name: "Other Cert"},
			{ID: "cert-002", Name: certName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := enrollmentcert.GetByName(context.Background(), service, certName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "cert-002", result.ID)
	assert.Equal(t, certName, result.Name)
}

func TestEnrollmentCert_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/enrollmentCert"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []enrollmentcert.EnrollmentCert{
			{ID: "cert-001", Name: "Cert 1"},
			{ID: "cert-002", Name: "Cert 2"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := enrollmentcert.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestEnrollmentCert_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/enrollmentCert/" + certID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := enrollmentcert.Get(context.Background(), service, certID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestEnrollmentCert_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/enrollmentCert"

	server.On("POST", path, common.SuccessResponse(enrollmentcert.EnrollmentCert{
		ID:           "new-cert-123",
		Name:         "New Certificate",
		AllowSigning: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newCert := &enrollmentcert.EnrollmentCert{
		Name:         "New Certificate",
		AllowSigning: true,
	}

	result, _, err := enrollmentcert.Create(context.Background(), service, newCert)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-cert-123", result.ID)
}

func TestEnrollmentCert_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := "cert-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/enrollmentCert/" + certID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateCert := &enrollmentcert.EnrollmentCert{
		ID:           certID,
		Name:         "Updated Certificate",
		AllowSigning: false,
	}

	resp, err := enrollmentcert.Update(context.Background(), service, certID, updateCert)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestEnrollmentCert_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := "cert-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/enrollmentCert/" + certID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := enrollmentcert.Delete(context.Background(), service, certID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestEnrollmentCert_GenerateCSR_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/enrollmentCert/csr"

	server.On("POST", path, common.SuccessResponse(enrollmentcert.GenerateEnrollmentCSR{
		Name: "CSR Request",
		CSR:  "CSR-CONTENT-HERE",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	csrRequest := &enrollmentcert.GenerateEnrollmentCSR{
		Name: "CSR Request",
	}

	result, _, err := enrollmentcert.GenerateCSR(context.Background(), service, csrRequest)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.CSR)
}

func TestEnrollmentCert_GenerateSelfSigned_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/enrollmentCert/selfsigned"

	server.On("POST", path, common.SuccessResponse(enrollmentcert.GenerateSelfSignedCert{
		Name:        "Self-Signed Certificate",
		Certificate: "SELF-SIGNED-CERT-HERE",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	certRequest := &enrollmentcert.GenerateSelfSignedCert{
		Name: "Self-Signed Certificate",
	}

	result, _, err := enrollmentcert.GenerateSelfSigned(context.Background(), service, certRequest)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Certificate)
}
