package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/organization_details"
)

const (
	subscriptionsPath      = "/zia/api/v1/subscriptions"
	orgInformationPath     = "/zia/api/v1/orgInformation"
	orgInformationLitePath = "/zia/api/v1/orgInformation/lite"
)

func TestOrganizationDetails_GetSubscriptions_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", subscriptionsPath, common.SuccessResponse([]organization_details.Subscription{
		{
			ID:         "ZIA",
			Status:     "SUBSCRIBED",
			State:      "ACTIVE",
			Licenses:   1000,
			SKU:        "ZIA-ENT",
			Subscribed: true,
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := organization_details.GetSubscriptions(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "ZIA", result[0].ID)
	assert.True(t, result[0].Subscribed)
}

func TestOrganizationDetails_GetOrgInformation_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", orgInformationPath, common.SuccessResponse(organization_details.Organization{
		OrgID:            123456,
		Name:             "Acme Corp",
		HQLocation:       "San Jose, US",
		Domains:          []string{"acme.com", "acme.net"},
		GeoLocation:      "AMER",
		IndustryVertical: "TECHNOLOGY",
		Country:          "US",
		Language:         "ENGLISH",
		Timezone:         "GMT",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := organization_details.GetOrgInformation(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Acme Corp", result.Name)
	assert.Equal(t, 123456, result.OrgID)
	assert.Contains(t, result.Domains, "acme.com")
}

func TestOrganizationDetails_GetOrgInformationLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", orgInformationLitePath, common.SuccessResponse(organization_details.OrganizationInfoLite{
		OrgID:     123456,
		Name:      "Acme Corp",
		CloudName: "zscaler.net",
		Domains:   []string{"acme.com"},
		Language:  "ENGLISH",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := organization_details.GetOrgInformationLite(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Acme Corp", result.Name)
	assert.Equal(t, "zscaler.net", result.CloudName)
}

func TestOrganizationDetails_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Organization JSON marshaling", func(t *testing.T) {
		org := organization_details.Organization{
			OrgID:   123456,
			Name:    "Acme Corp",
			Country: "US",
			City:    "San Jose",
		}

		data, err := json.Marshal(org)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"Acme Corp"`)
		assert.Contains(t, string(data), `"orgId":123456`)
	})
}
