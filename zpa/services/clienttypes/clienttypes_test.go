package clienttypes

import (
	"reflect"
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

	// Test case: Normal scenario
	t.Run("TestGetAllClientTypesNormal", func(t *testing.T) {
		clientTypes, resp, err := GetAllClientTypes(service)
		if err != nil {
			t.Fatalf("Failed to fetch client types: %v", err)
		}

		if resp.StatusCode >= 400 {
			t.Errorf("Received an HTTP error %d when fetching client types", resp.StatusCode)
		}

		if clientTypes == nil {
			t.Fatal("Client types are nil, expected a valid response")
		}

		tests := map[string]string{
			"zpn_client_type_exporter":          "Web Browser",
			"zpn_client_type_exporter_noauth":   "Web Browser Unauthenticated",
			"zpn_client_type_browser_isolation": "Cloud Browser",
			"zpn_client_type_machine_tunnel":    "Machine Tunnel",
			"zpn_client_type_ip_anchoring":      "ZIA Service Edge",
			"zpn_client_type_edge_connector":    "Cloud Connector",
			"zpn_client_type_zapp":              "Client Connector",
			"zpn_client_type_slogger":           "ZPA LSS",
			"zpn_client_type_branch_connector":  "Branch Connector",
			// "zpn_client_type_zapp_partner":      "Client Connector Partner",
		}

		clientTypeValues := getValuesByTags(clientTypes)
		for jsonTag, expectedValue := range tests {
			actualValue, found := clientTypeValues[jsonTag]
			if !found || actualValue != expectedValue {
				t.Errorf("Expected %s but got %s for json tag %s", expectedValue, actualValue, jsonTag)
			}
		}
	})

	// Test case: Error scenario
	t.Run("TestGetAllClientTypesError", func(t *testing.T) {
		// Temporarily change the client configuration to trigger an error
		service.Client.Config.CustomerID = "invalid_customer_id"
		_, _, err := GetAllClientTypes(service)
		if err == nil {
			t.Errorf("Expected error while fetching client types with invalid customer ID, got nil")
		}
		// Reset the customer ID to avoid affecting other tests
		service.Client.Config.CustomerID = client.Config.CustomerID
	})
}

func getValuesByTags(types *ClientTypes) map[string]string {
	values := make(map[string]string)
	r := reflect.ValueOf(types).Elem()
	for i := 0; i < r.NumField(); i++ {
		fieldTag := r.Type().Field(i).Tag.Get("json")
		values[fieldTag] = r.Field(i).String()
	}
	return values
}
