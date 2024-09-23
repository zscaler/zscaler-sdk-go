package ipdestinationgroups

/*
import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

const (
	maxRetries    = 3
	retryInterval = 2 * time.Second
)

// Constants for conflict retries
const (
	maxConflictRetries    = 5
	conflictRetryInterval = 1 * time.Second
)

func retryOnConflict(operation func() error) error {
	var lastErr error
	for i := 0; i < maxConflictRetries; i++ {
		lastErr = operation()
		if lastErr == nil {
			return nil
		}

		if strings.Contains(lastErr.Error(), `"code":"EDIT_LOCK_NOT_AVAILABLE"`) {
			log.Printf("Conflict error detected, retrying in %v... (Attempt %d/%d)", conflictRetryInterval, i+1, maxConflictRetries)
			time.Sleep(conflictRetryInterval)
			continue
		}

		return lastErr
	}
	return lastErr
}

func TestFWFilteringIPDestGroups(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}

	group := IPDestinationGroups{
		Name:        name,
		Description: name,
		Type:        "DSTN_IP",
		Addresses:   []string{"3.217.228.0-3.217.231.255"},
	}

	// Test resource creation
	var createdResource *IPDestinationGroups
	err = retryOnConflict(func() error {
		createdResource, err = Create(service, &group)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}
	log.Printf("Created resource with ID: %d", createdResource.ID) // Added log statement for the created resource

	if createdResource.ID == 0 {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}

	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved dlp engine '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Name = updateName
	// Ensure that the Addresses field (and others if needed) remains unchanged during the update
	retrievedResource.Addresses = group.Addresses

	err = retryOnConflict(func() error {
		_, _, err = Update(service, createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := Get(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
		return
	}
	if updatedResource == nil {
		t.Error("Updated resource is nil")
		return
	}

	// Continue checking
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, err = GetByName(service, updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
		return
	}
	if retrievedResource == nil {
		t.Error("Retrieved resource by name is nil")
		return
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.Name)
	}

	// Test resources retrieval
	allResources, err := GetAll(service)
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
	}
	if len(allResources) == 0 {
		t.Fatal("Expected retrieved resources to be non-empty, but got empty slice")
	}

	// check if the created resource is in the list
	found := false
	for _, resource := range allResources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved resources to contain created resource '%d', but it didn't", createdResource.ID)
	}

	// Introduce a delay before deleting
	time.Sleep(5 * time.Second) // sleep for 5 seconds

	// Test resource removal
	err = retryOnConflict(func() error {
		// First, attempt to delete the resource
		_, delErr := Delete(service, createdResource.ID)
		if delErr != nil {
			if strings.Contains(delErr.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
				// If we get a RESOURCE_NOT_FOUND error during deletion, it's already deleted.
				log.Printf("Resource %d already deleted.", createdResource.ID)
				return nil
			}
			// If deletion returns another error, return that error.
			return delErr
		}

		// Confirm deletion by trying to get the deleted resource.
		_, getErr := Get(service, createdResource.ID)
		if getErr != nil {
			if strings.Contains(getErr.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
				// If we get a RESOURCE_NOT_FOUND error, it confirms successful deletion.
				return nil
			}
			// If the get operation returns another error, return that error.
			return getErr
		}

		// If get operation does not return an error, it means deletion was not successful.
		return fmt.Errorf("Resource %d was not deleted", createdResource.ID)
	})
	if err != nil {
		t.Fatalf("Unexpected error during deletion: %v", err)
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *zscaler.Service, id int) (*IPDestinationGroups, error) {
	var resource *IPDestinationGroups
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = Get(s, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}

func TestRetrieveNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}

	_, err = Get(service, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}

	_, err = Delete(service, 0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}

	_, _, err = Update(service, 0, &IPDestinationGroups{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}

	_, err = GetByName(service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
*/
