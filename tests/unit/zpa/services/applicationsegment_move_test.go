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
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID + "/move"

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	req := applicationsegment_move.AppSegmentMicrotenantMoveRequest{
		TargetSegmentGroupID: "sg-456",
	}

	_, err = applicationsegment_move.AppSegmentMicrotenantMove(context.Background(), service, appID, req)
	require.NoError(t, err)
}
