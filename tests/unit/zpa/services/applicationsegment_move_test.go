// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment_move"
)

func TestApplicationSegmentMove_Move_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "app-123"
	// Correct path: POST to /zpa/mgmtconfig/v1/admin/customers/{customerId}/application/{appId}/move
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID + "/move"

	// SDK uses POST, not PUT
	server.On("POST", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	req := applicationsegment_move.AppSegmentMicrotenantMoveRequest{
		TargetSegmentGroupID: "sg-456",
	}

	_, err = applicationsegment_move.AppSegmentMicrotenantMove(context.Background(), service, appID, req)
	require.NoError(t, err)
}

func TestApplicationSegmentMove_Move_Error_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	appID := "app-123"
	path := common.ZPAPath(api.CustomerID, "application", appID, "move")
	api.On("POST", path, common.NotFoundResponse())

	_, err := applicationsegment_move.AppSegmentMicrotenantMove(context.Background(), api.Service, appID, applicationsegment_move.AppSegmentMicrotenantMoveRequest{})
	require.Error(t, err)
}
