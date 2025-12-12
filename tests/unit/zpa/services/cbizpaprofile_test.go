// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbizpaprofile"
)

func TestCBIZPAProfile_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := "zpa-profile-12345"
	// Note: Get uses GetAll internally and finds by ID, so we mock GetAll
	// Correct path: /zpa/cbiconfig/cbi/api/customers/{customerId}/zpaprofiles (plural, no ID in path for GetAll)
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/zpaprofiles"

	server.On("GET", path, common.SuccessResponse([]cbizpaprofile.ZPAProfiles{
		{ID: profileID, Name: "Test ZPA Profile"},
		{ID: "other-profile", Name: "Other Profile"},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cbizpaprofile.Get(context.Background(), service, profileID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, profileID, result.ID)
}

func TestCBIZPAProfile_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Correct path: /zpa/cbiconfig/cbi/api/customers/{customerId}/zpaprofiles (plural)
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/zpaprofiles"

	server.On("GET", path, common.SuccessResponse([]cbizpaprofile.ZPAProfiles{{ID: "zpa-001"}, {ID: "zpa-002"}}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cbizpaprofile.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
