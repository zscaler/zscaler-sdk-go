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
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/certificate/" + certID

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

	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/certificate"

	server.On("GET", path, common.SuccessResponse([]cbicertificatecontroller.CBICertificate{{ID: "cert-001"}, {ID: "cert-002"}}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cbicertificatecontroller.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
