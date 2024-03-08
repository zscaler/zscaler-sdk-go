package zpath

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestGetAllAltClouds(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Failed to create ZPA client: %v", err)
	}
	service := New(client)

	altClouds, resp, err := service.GetAltCloud()
	if err != nil {
		t.Fatalf("Failed to get alternate clouds: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	if len(altClouds) == 0 {
		t.Fatalf("Expected at least one alternate cloud endpoint, got none")
	}
}

func TestStatusCodesErrorResponse(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Failed to create ZPA client: %v", err)
	}
	service := New(client)

	altClouds, httpResponse, err := service.GetAltCloud()
	if err != nil {
		t.Fatalf("Error retrieving alternate clouds: %v", err)
	}

	var errResponse ErrorResponse
	if httpResponse.StatusCode == http.StatusOK {
		defer httpResponse.Body.Close()
		if err := json.NewDecoder(httpResponse.Body).Decode(&errResponse); err == nil && errResponse.Error != "" {
			t.Fatalf("Received error response despite 200 OK status: %s", errResponse.Error)
		}
	}

	if len(altClouds) == 0 {
		t.Fatalf("Expected at least one alternate cloud endpoint, got none")
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
