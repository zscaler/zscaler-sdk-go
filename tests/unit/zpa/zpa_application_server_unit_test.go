package unit

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appservercontroller"
)

func TestApplicationServer_Get(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/server/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Server 1"}`))
	})

	// Make the GET request
	appServer, _, err := appservercontroller.Get(service, "123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	// Check if the server ID and name match the expected values
	if appServer.ID != "123" {
		t.Errorf("Expected server ID '123', but got '%s'", appServer.ID)
	}
	if appServer.Name != "Server 1" {
		t.Errorf("Expected group name 'Server 1', but got '%s'", appServer.Name)
	}
}

func TestApplicationServer_Create(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/server", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Server 1"}`))
	})

	// Create a sample group
	appServer := appservercontroller.ApplicationServer{
		ID:          "123",
		Name:        "Server 1",
		Description: "Server 1",
		Address:     "192.168.1.1",
	}

	// Make the POST request
	createdAppServer, _, err := appservercontroller.Create(service, appServer)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	// Check if the created server ID and name match the expected values
	if createdAppServer.ID != "123" {
		t.Errorf("Expected created server ID '123', but got '%s'", createdAppServer.ID)
	}
	if createdAppServer.Name != "Server 1" {
		t.Errorf("Expected created server name 'Server 1', but got '%s'", createdAppServer.Name)
	}
}

// You can write similar tests for other functions like GetByName, Update, Delete, and GetAll.

func TestApplicationServer_GetByName(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/server", func(w http.ResponseWriter, r *http.Request) {
		// Get the query parameter "name" from the request
		query := r.URL.Query()
		name := query.Get("search")

		// Check if the name matches the expected value
		if name == "Server1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"list":[
					{"id": "123", "name": "Server1"}
				],
				"totalPages":1
				}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Server not found"}`))
		}
	})

	// Make the GetByName request
	appServer, _, err := appservercontroller.GetByName(service, "Server1")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetByName request: %v", err)
	}

	// Check if the server ID and name match the expected values
	if appServer.ID != "123" {
		t.Errorf("Expected server ID '123', but got '%s'", appServer.ID)
	}
	if appServer.Name != "Server1" {
		t.Errorf("Expected group name 'Server1', but got '%s'", appServer.Name)
	}
}

func TestApplicationServer_Update(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/server/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	})

	appServer := appservercontroller.ApplicationServer{
		ID:          "123",
		Name:        "Server 1",
		Description: "Server 1",
		Address:     "192.168.1.1",
	}

	// Make the Update request
	_, err := appservercontroller.Update(service, "123", appServer)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Update request: %v", err)
	}
}

func TestApplicationServer_Delete(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/server/123", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	})

	// Make the Delete request
	_, err := appservercontroller.Delete(service, "123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Delete request: %v", err)
	}
}

func TestApplicationServer_GetAll(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/server", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
				{"id": "123", "name": "Server 1"},
				{"id": "456", "name": "Server 2"}
			],
			"totalPages":1
			}`))
	})

	// Make the GetAll request
	appServers, _, err := appservercontroller.GetAll(service)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check the returned groups
	expectedAppServers := []*appservercontroller.ApplicationServer{
		{ID: "123", Name: "Server 1"},
		{ID: "456", Name: "Server 2"},
	}
	if len(appServers) != len(expectedAppServers) {
		t.Errorf("Expected %d groups, but got %d", len(expectedAppServers), len(appServers))
	}
	for i, expectedAppServer := range expectedAppServers {
		appServer := appServers[i]
		if appServer.ID != expectedAppServer.ID {
			t.Errorf("Expected group ID '%s', but got '%s'", expectedAppServer.ID, appServer.ID)
		}
		if appServer.Name != expectedAppServer.Name {
			t.Errorf("Expected group name '%s', but got '%s'", expectedAppServer.Name, appServer.Name)
		}
	}
}
