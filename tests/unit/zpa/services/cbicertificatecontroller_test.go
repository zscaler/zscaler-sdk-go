// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
)

func TestCBICertificateController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := "cert-12345"
	// Correct path: /zpa/cbiconfig/cbi/api/customers/{customerId}/certificates/{id} (plural)
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/certificates/" + certID

	server.On("GET", path, common.SuccessResponse(cbicertificatecontroller.CBICertificate{
		ID:   certID,
		Name: "Test Certificate",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cbicertificatecontroller.Get(context.Background(), service, certID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, certID, result.ID)
}

func TestCBICertificateController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Correct path: /zpa/cbiconfig/cbi/api/customers/{customerId}/certificates (plural)
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/certificates"

	server.On("GET", path, common.SuccessResponse([]cbicertificatecontroller.CBICertificate{{ID: "cert-001"}, {ID: "cert-002"}}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cbicertificatecontroller.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCBICertificateController_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Create uses singular endpoint
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/certificate"

	server.On("POST", path, common.SuccessResponse(cbicertificatecontroller.CBICertificate{
		ID:   "new-cert-123",
		Name: "New Certificate",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newCert := &cbicertificatecontroller.CBICertificate{
		Name: "New Certificate",
	}

	result, _, err := cbicertificatecontroller.Create(context.Background(), service, newCert)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-cert-123", result.ID)
}

func TestCBICertificateController_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := "cert-12345"
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/certificates/" + certID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateCert := &cbicertificatecontroller.CBICertificate{
		ID:   certID,
		Name: "Updated Certificate",
	}

	resp, err := cbicertificatecontroller.Update(context.Background(), service, certID, updateCert)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestCBICertificateController_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := "cert-12345"
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/certificates/" + certID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := cbicertificatecontroller.Delete(context.Background(), service, certID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
