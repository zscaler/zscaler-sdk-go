package unit

/*
import (
	"context"
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

func TestSegmentGroup_Get(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/segmentGroup/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"applications": [],
			"configSpace": "",
			"creationTime": "",
			"description": "",
			"enabled": true,
			"id": "123",
			"modifiedBy": "",
			"modifiedTime": "",
			"name": "Group1",
			"policyMigrated": false,
			"tcpKeepAliveEnabled": ""
		}`))
	})

	// Make the Get request
	group, _, err := segmentgroup.Get(context.Background(), service, "123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Get request: %v", err)
	}

	// Check if the group ID and name match the expected values
	if group.ID != "123" {
		t.Errorf("Expected group ID '123', but got '%s'", group.ID)
	}
	if group.Name != "Group1" {
		t.Errorf("Expected group name 'Group 1', but got '%s'", group.Name)
	}
}

/*
	func TestSegmentGroup_GetByName(t *testing.T) {
		client, mux, server := tests.NewOneAPIClientMock()
		defer server.Close()
		mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/segmentGroup", func(w http.ResponseWriter, r *http.Request) {
			// Get the query parameter "name" from the request
			query := r.URL.Query()
			name := query.Get("search")

			// Check if the name matches the expected value
			if name == "Group1" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"list": [
						{
							"applications": [],
							"configSpace": "",
							"creationTime": "",
							"description": "",
							"enabled": true,
							"id": "123",
							"modifiedBy": "",
							"modifiedTime": "",
							"name": "Group1",
							"policyMigrated": false,
							"tcpKeepAliveEnabled": ""
						}
					],
					"totalPages": 1
				}`))
			} else {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error": "Group not found"}`))
			}
		})
		service := &segmentgroup.Service{
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
			t.Errorf("Expected group name 'Group 1', but got '%s'", group.Name)
		}
	}

func TestSegmentGroup_Create(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/segmentGroup", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"applications": [],
			"configSpace": "",
			"creationTime": "",
			"description": "",
			"enabled": true,
			"id": "123",
			"modifiedBy": "",
			"modifiedTime": "",
			"name": "Group1",
			"policyMigrated": false,
			"tcpKeepAliveEnabled": ""
		}`))
	})

	// Create a sample group
	group := &segmentgroup.SegmentGroup{
		Name:    "Group1",
		Enabled: true,
	}

	// Make the Create request
	createdGroup, _, err := segmentgroup.Create(context.Background(), service, group)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Create request: %v", err)
	}

	// Check if the created group ID and name match the expected values
	if createdGroup.ID != "123" {
		t.Errorf("Expected created group ID '123', but got '%s'", createdGroup.ID)
	}
	if createdGroup.Name != "Group1" {
		t.Errorf("Expected created group name 'Group 1', but got '%s'", createdGroup.Name)
	}
}

func TestSegmentGroup_Update(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/segmentGroup/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"applications": [],
			"configSpace": "",
			"creationTime": "",
			"description": "",
			"enabled": true,
			"id": "123",
			"modifiedBy": "",
			"modifiedTime": "",
			"name": "Group1",
			"policyMigrated": false,
			"tcpKeepAliveEnabled": ""
		}`))
	})

	// Update a sample group
	group := &segmentgroup.SegmentGroup{
		ID:      "123",
		Name:    "Group1",
		Enabled: true,
	}

	// Make the Update request
	_, err := segmentgroup.Update(context.Background(), service, "123", group)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Update request: %v", err)
	}
}

func TestSegmentGroup_Delete(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/segmentGroup/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
	})

	// Make the Delete request
	_, err := segmentgroup.Delete(context.Background(), service, "123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Delete request: %v", err)
	}
}

func TestSegmentGroup_GetAll(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/segmentGroup", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list": [
				{
					"applications": [],
					"configSpace": "",
					"creationTime": "",
					"description": "",
					"enabled": true,
					"id": "123",
					"modifiedBy": "",
					"modifiedTime": "",
					"name": "Group1",
					"policyMigrated": false,
					"tcpKeepAliveEnabled": ""
				},
				{
					"applications": [],
					"configSpace": "",
					"creationTime": "",
					"description": "",
					"enabled": true,
					"id": "456",
					"modifiedBy": "",
					"modifiedTime": "",
					"name": "Group 2",
					"policyMigrated": false,
					"tcpKeepAliveEnabled": ""
				}
			],
			"totalPages": 1
		}`))
	})

	// Make the GetAll request
	groups, _, err := segmentgroup.GetAll(context.Background(), service)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check if the number of groups returned is correct
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, but got %d", len(groups))
	}

	// Check if the group IDs and names match the expected values
	if groups[0].ID != "123" {
		t.Errorf("Expected group 1 ID '123', but got '%s'", groups[0].ID)
	}
	if groups[0].Name != "Group1" {
		t.Errorf("Expected group 1 name 'Group 1', but got '%s'", groups[0].Name)
	}
	if groups[1].ID != "456" {
		t.Errorf("Expected group 2 ID '456', but got '%s'", groups[1].ID)
	}
	if groups[1].Name != "Group 2" {
		t.Errorf("Expected group 2 name 'Group 2', but got '%s'", groups[1].Name)
	}
}
*/
