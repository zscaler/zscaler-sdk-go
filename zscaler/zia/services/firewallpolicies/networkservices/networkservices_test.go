package networkservices

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestFWFilteringNWServices(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	networkServices := NetworkServices{
		Name:        name,
		Description: name,
		Type:        "CUSTOM",
		SrcTCPPorts: []NetworkPorts{
			{
				Start: 5000,
			},
			{
				Start: 5001,
			},
			{
				Start: 5002,
				End:   5005,
			},
		},
		DestTCPPorts: []NetworkPorts{
			{
				Start: 5000,
			},
			{
				Start: 5001,
			},
			{
				Start: 5003,
				End:   5005,
			},
		},
	}
	// Test resource creation
	createdResource, err := Create(context.Background(), service, &networkServices)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}

	// Test resource retrieval
	retrievedResource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource update
	retrievedResource.Name = updateName
	retrievedResource.SrcTCPPorts = []NetworkPorts{
		{
			Start: 5000,
		},
		{
			Start: 5001,
		},
		{
			Start: 5002,
			End:   5005,
		},
	}
	retrievedResource.DestTCPPorts = []NetworkPorts{
		{
			Start: 5001,
		},
		{
			Start: 5003,
			End:   5005,
		},
	}
	_, _, err = Update(context.Background(), service, createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}
	// Test resource retrieval by name
	retrievedResource, err = GetByName(context.Background(), service, updateName, nil, nil)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.Name)
	}
	// Test resources retrieval
	resources, err := GetAllNetworkServices(context.Background(), service, nil, nil)
	if err != nil {
		t.Errorf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		t.Error("Expected retrieved resources to be non-empty, but got empty slice")
	}
	// check if the created resource is in the list
	found := false
	for _, resource := range resources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved resources to contain created resource '%d', but it didn't", createdResource.ID)
	}
	// Test resource removal
	_, err = Delete(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}

func TestRetrieveNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	_, err = Get(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	_, err = Delete(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	_, _, err = Update(context.Background(), service, 0, &NetworkServices{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	_, err = GetByName(context.Background(), service, "non_existent_name", nil, nil)
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
