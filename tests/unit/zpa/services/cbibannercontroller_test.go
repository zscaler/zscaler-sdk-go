// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbibannercontroller"
)

func TestCBIBannerController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	bannerID := "banner-12345"
	// Correct path: /zpa/cbiconfig/cbi/api/customers/{customerId}/banners/{id} (plural "banners")
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/banners/" + bannerID

	server.On("GET", path, common.SuccessResponse(cbibannercontroller.CBIBannerController{
		ID:   bannerID,
		Name: "Test Banner",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cbibannercontroller.Get(context.Background(), service, bannerID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, bannerID, result.ID)
}

func TestCBIBannerController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Correct path: /zpa/cbiconfig/cbi/api/customers/{customerId}/banners (plural)
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/banners"

	// Note: GetAll returns a raw array, not paginated
	server.On("GET", path, common.SuccessResponse([]cbibannercontroller.CBIBannerController{{ID: "banner-001"}, {ID: "banner-002"}}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cbibannercontroller.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCBIBannerController_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Create uses singular endpoint "/banner"
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/banner"

	server.On("POST", path, common.SuccessResponse(cbibannercontroller.CBIBannerController{
		ID:   "new-banner-123",
		Name: "New Banner",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newBanner := &cbibannercontroller.CBIBannerController{
		Name: "New Banner",
	}

	result, _, err := cbibannercontroller.Create(context.Background(), service, newBanner)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-banner-123", result.ID)
}

func TestCBIBannerController_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	bannerID := "banner-12345"
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/banners/" + bannerID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateBanner := &cbibannercontroller.CBIBannerController{
		ID:   bannerID,
		Name: "Updated Banner",
	}

	resp, err := cbibannercontroller.Update(context.Background(), service, bannerID, updateBanner)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestCBIBannerController_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	bannerID := "banner-12345"
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/banners/" + bannerID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := cbibannercontroller.Delete(context.Background(), service, bannerID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
