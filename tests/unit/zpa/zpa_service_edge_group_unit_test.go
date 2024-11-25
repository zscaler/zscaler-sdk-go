package unit

/*
import (
	"context"
	"net/http"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
)

func TestServiceEdgeGroup_Get(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serviceEdgeGroup/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Group 1"}`))
	})

	// Make the GET request
	group, _, err := serviceedgegroup.Get(context.Background(), service, "123")
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

func TestServiceEdgeGroup_Create(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serviceEdgeGroup", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Group 1"}`))
	})

	// Create a sample group
	group := serviceedgegroup.ServiceEdgeGroup{
		ID:   "123",
		Name: "Group 1",
	}

	// Make the POST request
	createdGroup, _, err := serviceedgegroup.Create(context.Background(), service, group)
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
func TestServiceEdgeGroup_GetByName(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serviceEdgeGroup", func(w http.ResponseWriter, r *http.Request) {
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

	// Make the GetByName request
	group, _, err := serviceedgegroup.GetByName(context.Background(), service, "Group1")
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

func TestServiceEdgeGroup_Update(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serviceEdgeGroup/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	})

	group := serviceedgegroup.ServiceEdgeGroup{
		ID:   "123",
		Name: "Group 1",
	}

	// Make the Update request
	_, err := serviceedgegroup.Update(context.Background(), service, "123", &group)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Update request: %v", err)
	}
}

func TestServiceEdgeGroup_Delete(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serviceEdgeGroup/123", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	})

	// Make the Delete request
	_, err := serviceedgegroup.Delete(context.Background(), service, "123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Delete request: %v", err)
	}
}

func TestServiceEdgeGroup_GetAll(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serviceEdgeGroup", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
				{"id": "123", "name": "Group 1"},
				{"id": "456", "name": "Group 2"}
			],
			"totalPages":1
			}`))
	})

	// Make the GetAll request
	groups, _, err := serviceedgegroup.GetAll(context.Background(), service)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check the returned groups
	expectedGroups := []*serviceedgegroup.ServiceEdgeGroup{
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
*/
