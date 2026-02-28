package email_profiles

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
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

func TestEmailProfiles(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	description := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateDescription := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	label := EmailProfiles{
		Name:        name,
		Description: description,
		Emails:      []string{"test@example.com", "test2@example.com"},
	}

	var createdResource *EmailProfiles

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, _, err = Create(context.Background(), service, &label)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-zero, but got 0")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created rule label '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved rule label '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Description = updateDescription
	err = retryOnConflict(func() error {
		_, _, err = Update(context.Background(), service, createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Description != updateDescription {
		t.Errorf("Expected retrieved updated resource description '%s', but got '%s'", updateDescription, updatedResource.Description)
	}

	// Test resource retrieval by name
	retrievedResource, err = GetEmailProfileByName(context.Background(), service, name)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Description != updateDescription {
		t.Errorf("Expected retrieved resource comment '%s', but got '%s'", updateDescription, createdResource.Description)
	}
	// Test resources retrieval
	resources, err := GetAll(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		t.Fatal("Expected retrieved resources to be non-empty, but got empty slice")
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

	// Test GetAll with name filter
	nameFilter := name
	resourcesFiltered, err := GetAll(context.Background(), service, &GetAllFilterOptions{Name: &nameFilter})
	if err != nil {
		t.Fatalf("Error retrieving resources with name filter: %v", err)
	}
	found = false
	for _, resource := range resourcesFiltered {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected GetAll with name filter to contain created resource '%d', but it didn't", createdResource.ID)
	}

	// Test GetAllLite without filters
	resourcesLite, err := GetAllLite(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error retrieving lite resources: %v", err)
	}
	found = false
	for _, resource := range resourcesLite {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected GetAllLite to contain created resource '%d', but it didn't", createdResource.ID)
	}

	// Test GetAllLite with name filter
	resourcesLiteFiltered, err := GetAllLite(context.Background(), service, &GetAllLiteFilterOptions{Name: &nameFilter})
	if err != nil {
		t.Fatalf("Error retrieving lite resources with name filter: %v", err)
	}
	found = false
	for _, resource := range resourcesLiteFiltered {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected GetAllLite with name filter to contain created resource '%d', but it didn't", createdResource.ID)
	}

	// Test GetCount without filters
	totalCount, err := GetCount(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error retrieving count: %v", err)
	}
	if totalCount < 1 {
		t.Errorf("Expected count to be at least 1, got %d", totalCount)
	}

	// Test GetCount with name filter
	countWithFilter, err := GetCount(context.Background(), service, &GetAllFilterOptions{Name: &nameFilter})
	if err != nil {
		t.Fatalf("Error retrieving count with name filter: %v", err)
	}
	if countWithFilter < 1 {
		t.Errorf("Expected count with name filter to be at least 1, got %d", countWithFilter)
	}

	// Test resource removal
	err = retryOnConflict(func() error {
		_, delErr := Delete(context.Background(), service, createdResource.ID)
		return delErr
	})
	_, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(ctx context.Context, service *zscaler.Service, id int) (*EmailProfiles, error) {
	var resource *EmailProfiles
	var err error

	for i := 0; i < maxRetries; i++ {
		// Use the passed context (ctx) instead of context.Background()
		resource, err = Get(ctx, service, id)
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
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = Get(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	_, err = Delete(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, _, err = Update(context.Background(), service, 0, &EmailProfiles{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = GetEmailProfileByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
