package applicationsegment_move

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig         = "/zpa/mgmtconfig/v1/admin/customers/"
	appSegmentEndpoint = "/application"
)

type AppSegmentMicrotenantMoveRequest struct {
	ApplicationID        string `json:"applicationId,omitempty"`
	MicroTenantID        string `json:"microtenantId,omitempty"`
	TargetSegmentGroupID string `json:"targetSegmentGroupId,omitempty"`
	TargetMicrotenantID  string `json:"targetMicrotenantId,omitempty"`
	TargetServerGroupID  string `json:"targetServerGroupId,omitempty"`
}

func AppSegmentMicrotenantMove(ctx context.Context, service *zscaler.Service, applicationID string, move AppSegmentMicrotenantMoveRequest) (*http.Response, error) {
	// Check if a microtenant ID was provided in the move struct, else use the one from the service
	microTenantID := move.MicroTenantID
	if microTenantID == "" && service.MicroTenantID() != nil {
		microTenantID = *service.MicroTenantID()
	}
	// Corrected URL format to include the applicationID before /move
	relativeURL := fmt.Sprintf("%s%s%s/%s/move", mgmtConfig, service.Client.GetCustomerID(), appSegmentEndpoint, applicationID)
	// Add microTenantID to the filter if it's provided
	filter := common.Filter{}
	if microTenantID != "" {
		filter.MicroTenantID = &microTenantID
	}
	resp, err := service.Client.NewRequestDo(ctx, "POST", relativeURL, filter, move, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
