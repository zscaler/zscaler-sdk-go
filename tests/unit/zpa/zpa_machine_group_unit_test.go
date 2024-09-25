package unit

/*
import (
	"context"
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/machinegroup"
)

func TestMachineGroup_Get(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/machineGroup/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "machineGroup 1"}`))
	})

	// Make the GET request
	group, _, err := machinegroup.Get(context.Background(), service, "123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	// Check if the idp ID and name match the expected values
	if group.ID != "123" {
		t.Errorf("Expected machineGroup ID '123', but got '%s'", group.ID)
	}
	if group.Name != "machineGroup 1" {
		t.Errorf("Expected machineGroup name 'machineGroup 1', but got '%s'", group.Name)
	}
}

// You can write similar tests for other functions like GetByName, Update, Delete, and GetAll.
/*
func TestMachineGroup_GetByName(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/machineGroup", func(w http.ResponseWriter, r *http.Request) {
		// Get the query parameter "name" from the request
		query := r.URL.Query()
		name := query.Get("search")

		// Check if the name matches the expected value
		if name == "machineGroup1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"list":[
					{"id": "123", "name": "machineGroup1"}
				],
				"totalPages":1
				}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Machine Group not found"}`))
		}
	})
	service := &machinegroup.Service{
		Client: client,
	}

	// Make the GetByName request
	group, _, err := service.GetByName("machineGroup1")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetByName request: %v", err)
	}

	// Check if the machine group ID and name match the expected values
	if group.ID != "123" {
		t.Errorf("Expected machine group ID '123', but got '%s'", group.ID)
	}
	if group.Name != "machineGroup1" {
		t.Errorf("Expected machine name 'Idp1', but got '%s'", group.Name)
	}
}

func TestMachineGroup_GetAll(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/machineGroup", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
				{"id": "123", "name": "machineGroup 1"},
				{"id": "456", "name": "machineGroup 2"}
			],
			"totalPages":1
			}`))
	})

	// Make the GetAll request
	groups, _, err := machinegroup.GetAll(context.Background(), service)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check the returned idps
	expectedMachineGroups := []*machinegroup.MachineGroup{
		{ID: "123", Name: "machineGroup 1"},
		{ID: "456", Name: "machineGroup 2"},
	}
	if len(groups) != len(expectedMachineGroups) {
		t.Errorf("Expected %d machine groups, but got %d", len(expectedMachineGroups), len(groups))
	}
	for i, expectedMachineGroup := range expectedMachineGroups {
		group := groups[i]
		if group.ID != expectedMachineGroup.ID {
			t.Errorf("Expected machine group ID '%s', but got '%s'", expectedMachineGroup.ID, group.ID)
		}
		if group.Name != expectedMachineGroup.Name {
			t.Errorf("Expected machine group name '%s', but got '%s'", expectedMachineGroup.Name, group.Name)
		}
	}
}
*/
