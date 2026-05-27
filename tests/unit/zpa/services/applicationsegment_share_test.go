// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment_share"
)

func TestApplicationSegmentShare_Share_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "app-123"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID + "/share"

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	req := applicationsegment_share.AppSegmentSharedToMicrotenant{
		ApplicationID:       appID,
		ShareToMicrotenants: []string{"mt-001", "mt-002"},
	}

	_, err = applicationsegment_share.AppSegmentMicrotenantShare(context.Background(), service, appID, req)
	require.NoError(t, err)
}

func TestApplicationSegmentShare_Share_Error_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	appID := "app-123"
	path := common.ZPAPath(api.CustomerID, "application", appID, "share")
	api.On("PUT", path, common.NotFoundResponse())

	_, err := applicationsegment_share.AppSegmentMicrotenantShare(context.Background(), api.Service, appID, applicationsegment_share.AppSegmentSharedToMicrotenant{
		ApplicationID: appID,
	})
	require.Error(t, err)
}
