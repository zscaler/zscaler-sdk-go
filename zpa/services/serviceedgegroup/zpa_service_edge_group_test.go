package serviceedgegroup

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

func TestServiceEdgeGroup_Create(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	serviceEdgeGroup := ServiceEdgeGroup{
		Name:                   name,
		Description:            name,
		Enabled:                true,
		Latitude:               "37.33874",
		Longitude:              "-121.8852525",
		Location:               "San Jose, CA, USA",
		IsPublic:               "TRUE",
		UpgradeDay:             "SUNDAY",
		UpgradeTimeInSecs:      "66600",
		OverrideVersionProfile: true,
		VersionProfileName:     "Default",
		VersionProfileID:       "0",
		GraceDistanceEnabled:   true,
		GraceDistanceValue:     "10",
		GraceDistanceValueUnit: "MILES",
	}

	// Test resource creation
	createdResource, _, err := Create(service, &serviceEdgeGroup)

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource retrieval
	retrievedResource, _, err := Get(service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource update
	retrievedResource.Name = updateName
	_, err = Update(service, createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := Get(service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, _, err = GetByName(service, updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.Name)
	}
	// Test resources retrieval
	resources, _, err := GetAll(service)
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
		t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.ID)
	}
	// Test resource removal
	_, err = Delete(service, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = Get(service, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}

// func TestRetrieveNonExistentResource(t *testing.T) {
// 	client, err := tests.NewZpaClient()
// 	if err != nil {
// 		t.Fatalf("Error creating client: %v", err)
// 	}
// 	service := services.New(client)

// 	_, _, err = Get(service, "non_existent_id")
// 	if err == nil {
// 		t.Error("Expected error retrieving non-existent resource, but got nil")
// 	}
// }

// func TestDeleteNonExistentResource(t *testing.T) {
// 	client, err := tests.NewZpaClient()
// 	if err != nil {
// 		t.Fatalf("Error creating client: %v", err)
// 	}
// 	service := services.New(client)

// 	_, err = Delete(service, "non_existent_id")
// 	if err == nil {
// 		t.Error("Expected error deleting non-existent resource, but got nil")
// 	}
// }

// func TestUpdateNonExistentResource(t *testing.T) {
// 	client, err := tests.NewZpaClient()
// 	if err != nil {
// 		t.Fatalf("Error creating client: %v", err)
// 	}
// 	service := services.New(client)

// 	_, err = Update(service, "non_existent_id", &ServiceEdgeGroup{})
// 	if err == nil {
// 		t.Error("Expected error updating non-existent resource, but got nil")
// 	}
// }

// func TestGetByNameNonExistentResource(t *testing.T) {
// 	client, err := tests.NewZpaClient()
// 	if err != nil {
// 		t.Fatalf("Error creating client: %v", err)
// 	}
// 	service := services.New(client)

// 	_, _, err = GetByName(service, "non_existent_name")
// 	if err == nil {
// 		t.Error("Expected error retrieving resource by non-existent name, but got nil")
// 	}
// }
