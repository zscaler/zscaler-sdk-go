package lssconfigcontroller

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
)

func TestGetAllStatusCodes(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Failed to create ZPA client: %v", err)
	}
	service := New(client)

	statusCodes, resp, err := service.GetStatusCodes()
	if err != nil {
		t.Fatalf("Failed to get status codes: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Check if the returned mappings for each log are non-empty
	if len(statusCodes.ZPNAuthLog) == 0 {
		t.Error("ZPNAuthLog is empty")
	}
	if len(statusCodes.ZPNAstAuthLog) == 0 {
		t.Error("ZPNAstAuthLog is empty")
	}
	if len(statusCodes.ZPNTransLog) == 0 {
		t.Error("ZPNTransLog is empty")
	}
	if len(statusCodes.ZPNSysAuthLog) == 0 {
		t.Error("ZPNSysAuthLog is empty")
	}
}
