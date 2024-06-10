package unit

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/trustednetwork"
)

func TestTrustedNetworks_Get(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/network/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "trustedNetwork 1"}`))
	})

	// Make the GET request
	network, _, err := trustednetwork.Get(service, "123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	// Check if the trusted network ID and name match the expected values
	if network.ID != "123" {
		t.Errorf("Expected network ID '123', but got '%s'", network.ID)
	}
	if network.Name != "trustedNetwork 1" {
		t.Errorf("Expected network name 'trustedNetwork 1', but got '%s'", network.Name)
	}
}

// You can write similar tests for other functions like GetByName, Update, Delete, and GetAll.
func TestTrustedNetworks_GetByName(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/network", func(w http.ResponseWriter, r *http.Request) {
		// Get the query parameter "name" from the request
		query := r.URL.Query()
		name := query.Get("search")

		// Check if the name matches the expected value
		if name == "trustedNetwork1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"list":[
					{"id": "123", "name": "trustedNetwork1"}
				],
				"totalPages":1
				}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Trusted Network not found"}`))
		}
	})

	// Make the GetByName request
	network, _, err := trustednetwork.GetByName(service, "trustedNetwork1")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetByName request: %v", err)
	}

	// Check if the network ID and name match the expected values
	if network.ID != "123" {
		t.Errorf("Expected network ID '123', but got '%s'", network.ID)
	}
	if network.Name != "trustedNetwork1" {
		t.Errorf("Expected network name 'trustedNetwork1', but got '%s'", network.Name)
	}
}

func TestTrustedNetworks_GetAll(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/network", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
				{"id": "123", "name": "trustedNetwork 1"},
				{"id": "456", "name": "trustedNetwork 2"}
			],
			"totalPages":1
			}`))
	})

	// Make the GetAll request
	networks, _, err := trustednetwork.GetAll(service)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check the returned trusted networks
	expectedNetworks := []*trustednetwork.TrustedNetwork{
		{ID: "123", Name: "trustedNetwork 1"},
		{ID: "456", Name: "trustedNetwork 2"},
	}
	if len(networks) != len(expectedNetworks) {
		t.Errorf("Expected %d trusted networks, but got %d", len(expectedNetworks), len(networks))
	}
	for i, expectedNetwork := range expectedNetworks {
		network := networks[i]
		if network.ID != expectedNetwork.ID {
			t.Errorf("Expected network ID '%s', but got '%s'", expectedNetwork.ID, network.ID)
		}
		if network.Name != expectedNetwork.Name {
			t.Errorf("Expected trusted network name '%s', but got '%s'", expectedNetwork.Name, network.Name)
		}
	}
}
