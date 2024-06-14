package urlcategories

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
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

		if strings.Contains(lastErr.Error(), `"code":"EDIT_LOCK_NOT_AVAILABLE"`) ||
			strings.Contains(lastErr.Error(), `"code":"INVALID_OPERATION"`) {
			log.Printf("Conflict error detected, retrying in %v... (Attempt %d/%d)", conflictRetryInterval, i+1, maxConflictRetries)
			time.Sleep(conflictRetryInterval)
			continue
		}

		return lastErr
	}
	return lastErr
}

func TestURLCategories(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateDescription := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	urlCategories := URLCategory{
		SuperCategory:     "USER_DEFINED",
		ConfiguredName:    name,
		Description:       name,
		Keywords:          []string{"microsoft"},
		CustomCategory:    true,
		DBCategorizedUrls: []string{".creditkarma.com", ".youku.com"},
		Type:              "URL_CATEGORY",
		Urls: []string{
			".coupons.com",
		},
		IPRanges: []string{
			"3.217.228.0/25",
			"3.235.112.0/24",
		},
		IPRangesRetainingParentCategory: []string{
			"13.107.6.152/31",
		},
	}

	var createdResource *URLCategory

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = CreateURLCategories(service, &urlCategories)
		return err
	})
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
		return
	}

	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.ConfiguredName != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.ConfiguredName)
	}
	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
		return
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.ConfiguredName != name {
		t.Errorf("Expected retrieved dlp engine '%s', but got '%s'", name, retrievedResource.ConfiguredName)
	}

	// Test resource update
	retrievedResource.Description = updateDescription
	err = retryOnConflict(func() error {
		_, _, err = UpdateURLCategories(service, createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := Get(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Description != updateDescription {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateDescription, updatedResource.ConfiguredName)
	}

	// Test resource retrieval by name
	retrievedResource, err = GetCustomURLCategories(service, name)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Description != updateDescription {
		t.Errorf("Expected retrieved resource description '%s', but got '%s'", updateDescription, createdResource.ConfiguredName)
	}

	// Test resources retrieval
	resources, err := GetAll(service)
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
		t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.ID)
	}
	// Test resource removal
	err = retryOnConflict(func() error {
		_, delErr := DeleteURLCategories(service, createdResource.ID)
		return delErr
	})
	_, err = Get(service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with a retry mechanism.
func tryRetrieveResource(s *services.Service, id string) (*URLCategory, error) {
	var resource *URLCategory
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = Get(s, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		if strings.Contains(err.Error(), "404") { // handle 404 error
			return nil, err
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = Get(service, "non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = DeleteURLCategories(service, "non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, _, err = UpdateURLCategories(service, "non_existent_id", &URLCategory{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = GetCustomURLCategories(service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
