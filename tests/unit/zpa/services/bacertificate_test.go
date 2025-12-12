// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
)

func TestBACertificate_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	certID := "cert-12345"
	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/certificate/{id}
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/certificate/" + certID

	server.On("GET", path, common.SuccessResponse(bacertificate.BaCertificate{
		ID:   certID,
		Name: "Test Certificate",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := bacertificate.Get(context.Background(), service, certID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, certID, result.ID)
}

func TestBACertificate_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Correct path: /zpa/mgmtconfig/v2/admin/customers/{customerId}/clientlessCertificate/issued
	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/clientlessCertificate/issued"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []bacertificate.BaCertificate{{ID: "cert-001"}, {ID: "cert-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := bacertificate.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
