package applicationsegment_move

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig         = "/mgmtconfig/v1/admin/customers/"
	appSegmentEndpoint = "/application"
)

type AppSegmentMicrotenantMove struct {
	ApplicationID        string `json:"applicationId,omitempty"`
	MicroTenantID        string `json:"microtenantId,omitempty"`
	TargetSegmentGroupID string `json:"targetSegmentGroupId,omitempty"`
	TargetMicrotenantID  string `json:"targetMicrotenantId,omitempty"`
	TargetServerGroupID  string `json:"targetServerGroupId,omitempty"`
}

func (service *Service) AppSegmentMicrotenantMove(applicationID string, move AppSegmentMicrotenantMove) (*http.Response, error) {
	// Check if a microtenant ID was provided in the move struct, else use the one from the service
	microTenantID := move.MicroTenantID
	if microTenantID == "" && service.microTenantID != nil {
		microTenantID = *service.microTenantID
	}
	// Corrected URL format to include the applicationID before /move
	relativeURL := fmt.Sprintf("%s%s%s/%s/move", mgmtConfig, service.Client.Config.CustomerID, appSegmentEndpoint, applicationID)
	// Add microTenantID to the filter if it's provided
	filter := common.Filter{}
	if microTenantID != "" {
		filter.MicroTenantID = &microTenantID
	}
	resp, err := service.Client.NewRequestDo("POST", relativeURL, filter, move, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
