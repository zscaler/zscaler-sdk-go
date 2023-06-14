package unit

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appconnectorgroup"
)

func TestAppConnectorGroup_Get(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/appConnectorGroup/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Group 1"}`))
	})
	service := &appconnectorgroup.Service{
		Client: client,
	}

	// Make the GET request
	group, _, err := service.Get("123")

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	// Check if the group ID and name match the expected values
	if group.ID != "123" {
		t.Errorf("Expected group ID '123', but got '%s'", group.ID)
	}
	if group.Name != "Group 1" {
		t.Errorf("Expected group name 'Group 1', but got '%s'", group.Name)
	}
}

func TestAppConnectorGroup_Create(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/appConnectorGroup", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Group 1"}`))
	})

	service := &appconnectorgroup.Service{
		Client: client,
	}
	// Create a sample group
	group := appconnectorgroup.AppConnectorGroup{
		ID:   "123",
		Name: "Group 1",
	}

	// Make the POST request
	createdGroup, _, err := service.Create(group)

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	// Check if the created group ID and name match the expected values
	if createdGroup.ID != "123" {
		t.Errorf("Expected created group ID '123', but got '%s'", createdGroup.ID)
	}
	if createdGroup.Name != "Group 1" {
		t.Errorf("Expected created group name 'Group 1', but got '%s'", createdGroup.Name)
	}
}

// You can write similar tests for other functions like GetByName, Update, Delete, and GetAll.

func TestAppConnectorGroup_GetByName(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/appConnectorGroup", func(w http.ResponseWriter, r *http.Request) {
		// Get the query parameter "name" from the request
		query := r.URL.Query()
		name := query.Get("search")

		// Check if the name matches the expected value
		if name == "Group1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"list":[
					{"id": "123", "name": "Group1"}
				],
				"totalPages":1
				}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Group not found"}`))
		}
	})
	service := &appconnectorgroup.Service{
		Client: client,
	}

	// Make the GetByName request
	group, _, err := service.GetByName("Group1")

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetByName request: %v", err)
	}

	// Check if the group ID and name match the expected values
	if group.ID != "123" {
		t.Errorf("Expected group ID '123', but got '%s'", group.ID)
	}
	if group.Name != "Group1" {
		t.Errorf("Expected group name 'Group1', but got '%s'", group.Name)
	}
}

func TestAppConnectorGroup_Update(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/appConnectorGroup/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	})
	service := &appconnectorgroup.Service{
		Client: client,
	}
	group := appconnectorgroup.AppConnectorGroup{
		ID:   "123",
		Name: "Group 1",
	}

	// Make the Update request
	_, err := service.Update("123", &group)

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Update request: %v", err)
	}
}

func TestAppConnectorGroup_Delete(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/appConnectorGroup/123", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	})
	service := &appconnectorgroup.Service{
		Client: client,
	}

	// Make the Delete request
	_, err := service.Delete("123")

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Delete request: %v", err)
	}
}

func TestAppConnectorGroup_GetAll(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/appConnectorGroup", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
				{"id": "123", "name": "Group 1"},
				{"id": "456", "name": "Group 2"}
			],
			"totalPages":1
			}`))
	})
	service := &appconnectorgroup.Service{
		Client: client,
	}

	// Make the GetAll request
	groups, _, err := service.GetAll()

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check the returned groups
	expectedGroups := []*appconnectorgroup.AppConnectorGroup{
		{ID: "123", Name: "Group 1"},
		{ID: "456", Name: "Group 2"},
	}
	if len(groups) != len(expectedGroups) {
		t.Errorf("Expected %d groups, but got %d", len(expectedGroups), len(groups))
	}
	for i, expectedGroup := range expectedGroups {
		group := groups[i]
		if group.ID != expectedGroup.ID {
			t.Errorf("Expected group ID '%s', but got '%s'", expectedGroup.ID, group.ID)
		}
		if group.Name != expectedGroup.Name {
			t.Errorf("Expected group name '%s', but got '%s'", expectedGroup.Name, group.Name)
		}
	}
}
