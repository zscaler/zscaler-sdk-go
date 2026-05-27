package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/appservicegroups"
)

const appServiceGroupsLitePath = "/zia/api/v1/appServiceGroups/lite"

func TestAppServiceGroups_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "Web Services"
	server.On("GET", appServiceGroupsLitePath, common.SuccessResponse([]appservicegroups.ApplicationServicesGroupLite{
		{ID: 1, Name: "DNS Services"},
		{ID: 2, Name: groupName, NameL10nTag: true},
	}))

	svc, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := appservicegroups.GetByName(context.Background(), svc, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupName, result.Name)
}

func TestAppServiceGroups_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", appServiceGroupsLitePath, common.SuccessResponse([]appservicegroups.ApplicationServicesGroupLite{
		{ID: 1, Name: "Web Services"},
		{ID: 2, Name: "Mail Services"},
	}))

	svc, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := appservicegroups.GetAll(context.Background(), svc)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestAppServiceGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		group := appservicegroups.ApplicationServicesGroupLite{
			ID:          1,
			Name:        "Web Services",
			NameL10nTag: true,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"Web Services"`)
	})
}
