package lssconfigcontroller

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
)

func TestGetAllStatusCodes(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	statusCodes, resp, err := GetStatusCodes(context.Background(), service)
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

func TestStatusCodesErrorResponse(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Fetch the client types and the HTTP response
	statusCodes, httpResponse, err := GetStatusCodes(context.Background(), service)
	if err != nil {
		t.Fatalf("Error retrieving client types: %v", err)
	}

	// Check if the status code indicates success but there's an error message in the body
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	var errResponse ErrorResponse
	if httpResponse.StatusCode == http.StatusOK {
		if err := json.NewDecoder(httpResponse.Body).Decode(&errResponse); err == nil && errResponse.Error != "" {
			t.Fatalf("Received error response despite 200 OK status: %s", errResponse.Error)
		}
	}

	// Continue with other validations...
	if len(statusCodes.ZPNAuthLog) == 0 ||
		len(statusCodes.ZPNAstAuthLog) == 0 ||
		len(statusCodes.ZPNTransLog) == 0 ||
		len(statusCodes.ZPNSysAuthLog) == 0 {
		t.Error("One or more status code fields are empty")
	}
}
