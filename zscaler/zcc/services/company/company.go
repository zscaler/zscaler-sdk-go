package company

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	getCompanyInfoEndpoint = "/zcc/papi/public/v1/getCompanyInfo"
)

// GetAdminUserSyncInfo retrieves synchronization information for admin users.
func GetCompanyInfo(ctx context.Context, service *zscaler.Service) error {
	// Make the GET request
	resp, err := service.Client.NewZccRequestDo(ctx, "GET", getCompanyInfoEndpoint, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve company info: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP response codes
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to retrieve company info: received status code %d", resp.StatusCode)
	}

	// Since the API returns an empty JSON, simply return nil to indicate success
	return nil
}
