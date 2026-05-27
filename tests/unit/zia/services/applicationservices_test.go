package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/applicationservices"
)

const appServicesLitePath = "/zia/api/v1/appServices/lite"

func TestApplicationServices_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serviceName := "HTTP"
	server.On("GET", appServicesLitePath, common.SuccessResponse([]applicationservices.ApplicationServicesLite{
		{ID: 1, Name: "DNS"},
		{ID: 2, Name: serviceName, NameL10nTag: true},
	}))

	svc, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := applicationservices.GetByName(context.Background(), svc, serviceName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, serviceName, result.Name)
}

func TestApplicationServices_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", appServicesLitePath, common.SuccessResponse([]applicationservices.ApplicationServicesLite{
		{ID: 1, Name: "HTTP"},
		{ID: 2, Name: "HTTPS"},
	}))

	svc, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := applicationservices.GetAll(context.Background(), svc)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestApplicationServices_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		appSvc := applicationservices.ApplicationServicesLite{
			ID:          1,
			Name:        "HTTP",
			NameL10nTag: true,
		}

		data, err := json.Marshal(appSvc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"HTTP"`)
		assert.Contains(t, string(data), `"nameL10nTag":true`)
	})
}
