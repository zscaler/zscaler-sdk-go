package customeraudit

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services/common"
)

const (
	customerAuditEndpoint = "/dlp/v1/customer/audit"
)

type AuditLogsResponse struct {
	Cursor common.Cursor `json:"cursor"` // Use the shared Cursor struct
	Logs   []AuditLog    `json:"logs"`
}

type AuditLog struct {
	Action     Action `json:"action"`
	Module     string `json:"module"`
	Resource   string `json:"resource"`
	ChangedAt  string `json:"changedAt"`
	ChangedBy  string `json:"changedBy"`
	OldRowJSON string `json:"oldRowJson"`
	NewRowJSON string `json:"newRowJson"`
	ChangeNote string `json:"changeNote"`
}

type Action struct {
	Action string `json:"action"`
}

func GetCustomerAudit(ctx context.Context, service *services.Service, filters common.CommonDLPIncidentFiltering, paginationParams *common.PaginationParams) ([]AuditLog, *common.Cursor, error) {
	// Construct the endpoint URL
	endpoint := customerAuditEndpoint

	// Read all pages of audit logs using POST
	allResults, cursor, err := common.ReadAllPages[AuditLog](ctx, service.Client, http.MethodPost, endpoint, paginationParams, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch customer audit logs: %w", err)
	}

	return allResults, cursor, nil
}
