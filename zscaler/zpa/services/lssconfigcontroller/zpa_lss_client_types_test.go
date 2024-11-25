package lssconfigcontroller

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestGetAllClientTypes(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	clientTypes, _, err := GetClientTypes(context.Background(), service)
	if err != nil {
		t.Fatalf("Failed to get client types: %v", err)
	}

	// Check that each client type field is not empty.
	if clientTypes.ZPNClientTypeExporter == "" {
		t.Error("ZPNClientTypeExporter is empty")
	}
	if clientTypes.ZPNClientTypeMachineTunnel == "" {
		t.Error("ZPNClientTypeMachineTunnel is empty")
	}
	if clientTypes.ZPNClientTypeIPAnchoring == "" {
		t.Error("ZPNClientTypeIPAnchoring is empty")
	}
	if clientTypes.ZPNClientTypeEdgeConnector == "" {
		t.Error("ZPNClientTypeEdgeConnector is empty")
	}
	if clientTypes.ZPNClientTypeZAPP == "" {
		t.Error("ZPNClientTypeZAPP is empty")
	}
	if clientTypes.ZPNClientTypeSlogger == "" {
		t.Error("ZPNClientTypeSlogger is empty")
	}
}

func TestClientTypesStatusCodeCheck(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Fetch the client types and the HTTP response
	_, httpResponse, err := GetClientTypes(context.Background(), service)
	if err != nil {
		t.Fatalf("Error retrieving client types: %v", err)
	}

	// Check for 200 OK status code
	if httpResponse.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200 OK, but got: %d", httpResponse.StatusCode)
	}
}

func TestClientTypesErrorResponse(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Fetch the client types and the HTTP response
	clientTypes, httpResponse, err := GetClientTypes(context.Background(), service)
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
	if clientTypes.ZPNClientTypeExporter == "" ||
		clientTypes.ZPNClientTypeMachineTunnel == "" ||
		clientTypes.ZPNClientTypeIPAnchoring == "" ||
		clientTypes.ZPNClientTypeEdgeConnector == "" ||
		clientTypes.ZPNClientTypeZAPP == "" {
		t.Error("One or more client type fields are empty")
	}
	// Since ZPNClientTypeSlogger is omitempty, we only check if it's present.
	if clientTypes.ZPNClientTypeSlogger == "" {
		t.Error("ZPNClientTypeSlogger is empty")
	}
}
