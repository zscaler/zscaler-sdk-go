package integration

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/serviceedgegroup"
)

func TestServiceEdgeGroup_Create(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := serviceedgegroup.New(client)

	// Create a sample group
	group := serviceedgegroup.ServiceEdgeGroup{
		Name:                   "Group1",
		Description:            "Group1",
		Enabled:                true,
		Latitude:               "37.3861",
		Longitude:              "-122.0839",
		Location:               "Mountain View, CA",
		IsPublic:               "TRUE",
		UpgradeDay:             "SUNDAY",
		UpgradeTimeInSecs:      "66600",
		OverrideVersionProfile: true,
		VersionProfileName:     "New Release",
		VersionProfileID:       "2",
	}

	// Test resource creation
	createdResource, _, err := service.Create(group)

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == "" {
		t.Error("Expected created group ID to be non-empty, but got ''")
	}
	if createdResource.Name != "Group1" {
		t.Errorf("Expected created group name 'Group1', but got '%s'", createdResource.Name)
	}

	// Test resource retrieval
	retrievedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving group: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved group ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != "Group1" {
		t.Errorf("Expected retrieved group name 'Group1', but got '%s'", createdResource.Name)
	}
	// Test resource update
	retrievedResource.Name = "Group1-Updated"
	_, err = service.Update(createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating group: %v", err)
	}
	updatedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving group: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated group ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != "Group1-Updated" {
		t.Errorf("Expected retrieved updated group name 'Group1', but got '%s'", updatedResource.Name)
	}
	// Test resource retrieval by name
	retrievedResource, _, err = service.GetByName("Group1-Updated")
	if err != nil {
		t.Errorf("Error retrieving group by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved group ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != "Group1-Updated" {
		t.Errorf("Expected retrieved group name 'Group1', but got '%s'", createdResource.Name)
	}
	// Test resources retrieval
	resources, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error retrieving groups: %v", err)
	}
	if len(resources) == 0 {
		t.Error("Expected retrieved groups to be non-empty, but got empty slice")
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
		t.Errorf("Expected retrieved groups to contain created group '%s', but it didn't", createdResource.ID)
	}
	// Test resource removal
	_, err = service.Delete(createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting group: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = service.Get(createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted group, but got nil")
	}

}
