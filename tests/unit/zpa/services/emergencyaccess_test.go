// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/emergencyaccess"
)

func TestEmergencyAccess_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	userID := "user-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/emergencyAccess/user/" + userID

	server.On("GET", path, common.SuccessResponse(emergencyaccess.EmergencyAccess{
		UserID:    userID,
		EmailID:   "user@example.com",
		FirstName: "Test",
		LastName:  "User",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := emergencyaccess.Get(context.Background(), service, userID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
}

// Note: GetAll uses complex pagination with pageId and iterates through pages.
// It also uses a different endpoint (/emergencyAccess/users with 's').
// Testing GetAll would require mocking multiple paginated responses.
// For comprehensive GetAll testing, use integration tests.
