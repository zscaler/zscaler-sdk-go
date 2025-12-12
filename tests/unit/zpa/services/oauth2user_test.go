// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/oauth2_user"
)

func TestOAuth2User_VerifyUserCodes_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	associationType := "CONNECTOR_GRP"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/" + associationType + "/usercodes"

	server.On("POST", path, common.SuccessResponse(oauth2_user.OauthUser{
		TenantID:     "tenant-123",
		ZcomponentID: "zcomp-456",
		UserCodes:    []string{"code1", "code2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	oauthUser := &oauth2_user.OauthUser{
		UserCodes: []string{"code1", "code2"},
	}

	result, _, err := oauth2_user.VerifyUserCodes(context.Background(), service, associationType, oauthUser)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tenant-123", result.TenantID)
	assert.Len(t, result.UserCodes, 2)
}

func TestOAuth2User_VerifyUserCodeStatus_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	associationType := "SERVICE_EDGE_GRP"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/" + associationType + "/usercodes/status"

	server.On("POST", path, common.SuccessResponse(oauth2_user.OauthUser{
		TenantID:             "tenant-789",
		NonceAssociationType: associationType,
		UserCodes:            []string{"validcode1"},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	userCodes := []string{"validcode1"}
	result, _, err := oauth2_user.VerifyUserCodeStatus(context.Background(), service, associationType, userCodes)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tenant-789", result.TenantID)
}
