package lssconfigcontroller

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

func TestGetAllClientTypes(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Failed to create ZPA client: %v", err)
	}
	service := services.New(client)

	clientTypes, _, err := GetClientTypes(service)
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
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Failed to create ZPA client: %v", err)
	}
	service := services.New(client)

	// Fetch the client types and the HTTP response
	_, httpResponse, err := GetClientTypes(service)
	if err != nil {
		t.Fatalf("Error retrieving client types: %v", err)
	}

	// Check for 200 OK status code
	if httpResponse.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200 OK, but got: %d", httpResponse.StatusCode)
	}
}

func TestClientTypesErrorResponse(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Failed to create ZPA client: %v", err)
	}
	service := services.New(client)

	// Fetch the client types and the HTTP response
	clientTypes, httpResponse, err := GetClientTypes(service)
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
