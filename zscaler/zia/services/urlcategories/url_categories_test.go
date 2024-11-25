package urlcategories

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
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
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

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
		createdResource, err = CreateURLCategories(context.Background(), service, &urlCategories)
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
		_, _, err = UpdateURLCategories(context.Background(), service, createdResource.ID, retrievedResource)
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
		t.Errorf("Expected retrieved updated resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Description != updateDescription {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateDescription, updatedResource.ConfiguredName)
	}

	// Test resources retrieval
	resources, err := GetAll(context.Background(), service)
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

	// Test the GetIncludeOnlyUrlKeyWordCounts function with both parameters
	keywordCountResource, err := GetCustomURLCategories(context.Background(), service, name, true, true)
	if err != nil {
		t.Errorf("Error retrieving URL category with includeOnlyUrlKeywordCounts and customOnly: %v", err)
		return
	}
	if keywordCountResource == nil {
		t.Errorf("Expected non-nil keywordCountResource, but got nil")
	} else if keywordCountResource.ConfiguredName != name {
		t.Errorf("Expected keywordCountResource name '%s', but got '%s'", name, keywordCountResource.ConfiguredName)
	}

	// Test resource removal
	err = retryOnConflict(func() error {
		_, delErr := DeleteURLCategories(context.Background(), service, createdResource.ID)
		return delErr
	})
	_, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with a retry mechanism.
func tryRetrieveResource(s *zscaler.Service, id string) (*URLCategory, error) {
	var resource *URLCategory
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = Get(context.Background(), s, id)
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

func TestGetURLQuota(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Call the GetURLQuota function
	quota, err := GetURLQuota(context.Background(), service)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	// Check if the returned quota is valid
	if quota == nil {
		t.Fatal("Expected valid quota, but got nil")
	}

	// Example assertions on the response values
	if quota.UniqueUrlsProvisioned < 0 {
		t.Errorf("Unexpected value for UniqueUrlsProvisioned: %d", quota.UniqueUrlsProvisioned)
	}
	if quota.RemainingUrlsQuota < 0 {
		t.Errorf("Unexpected value for RemainingUrlsQuota: %d", quota.RemainingUrlsQuota)
	}
}

func TestGetURLLookup(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	urls := []string{"google.com"}
	lookupResults, err := GetURLLookup(context.Background(), service, urls)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if lookupResults == nil {
		t.Fatal("Expected valid lookup results, but got nil")
	}
	found := false
	for _, result := range lookupResults {
		if result.URL == "google.com" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected lookup results to contain 'google.com', but it was not found")
	}
	for _, result := range lookupResults {
		if len(result.URLClassifications) == 0 {
			t.Errorf("Expected URLClassifications to be non-empty for URL: %s", result.URL)
		}
	}
}

func TestGetAllLite(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Call the GetAllLite function
	urlCategories, err := GetAllLite(context.Background(), service)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	// Check if the returned URL categories are valid
	if urlCategories == nil {
		t.Fatal("Expected valid URL categories, but got nil")
	}

	// Check if the list is not empty and contains valid items
	if len(urlCategories) == 0 {
		t.Fatalf("Expected at least one URL category, but got none")
	}

	// Only retrieve the first item off the list
	firstCategory := urlCategories[0]
	if firstCategory.ID == "" {
		t.Errorf("Expected first URL category to have a non-empty ID, but got ''")
	}
	if firstCategory.Type == "" {
		t.Errorf("Expected first URL category to have a non-empty Type, but got ''")
	}

	// Log the first category
	t.Logf("First URL Category: %+v", firstCategory)
}

/*
// Test for Domain Review Methods
func TestURLCategoriesDomainReview(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	name1 := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	name2 := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	// Define the URL categories
	urlCategories1 := URLCategory{
		SuperCategory:  "USER_DEFINED",
		ConfiguredName: name1,
		Description:    name1,
		CustomCategory: true,
		Type:           "URL_CATEGORY",
		// UrlsRetainingParentCategoryToAdd: []string{".microsoft.com"},
		Urls: []string{".microsoft.com", "teams.microsoft.com", ".teams1.microsoft.com"},
	}

	urlCategories2 := URLCategory{
		SuperCategory:                    "USER_DEFINED",
		ConfiguredName:                   name2,
		Description:                      name2,
		CustomCategory:                   true,
		Type:                             "URL_CATEGORY",
		UrlsRetainingParentCategoryToAdd: []string{".microsoft.com"},
		Urls:                             []string{".microsoft.com", "teams.microsoft.com", ".teams1.microsoft.com"},
	}

	// Create the first URL category
	createdResource1, err := CreateURLCategories(service, &urlCategories1)
	if err != nil {
		t.Fatalf("Error creating first URL category: %v", err)
	}
	t.Logf("Created first URL category with ID: %s", createdResource1.ID)

	// Create the second URL category
	createdResource2, err := CreateURLCategories(service, &urlCategories2)
	if err != nil {
		t.Fatalf("Error creating second URL category: %v", err)
	}
	t.Logf("Created second URL category with ID: %s", createdResource2.ID)

	// Define the URLs to be reviewed
	urlsToReview := []string{"teams.microsoft.com"}

	// Create URL Review
	createResults, err := CreateURLReview(service, urlsToReview)
	if err != nil {
		t.Fatalf("Error creating URL review: %v", err)
	}
	t.Logf("Create URL Review Results: %+v", createResults)

	// Pause for 2 seconds
	time.Sleep(2 * time.Second)

	// Update URL Review
	err = UpdateURLReview(service, createResults)
	if err != nil {
		t.Fatalf("Error updating URL review: %v", err)
	}
	t.Logf("Update URL Review successfully completed")

	// Cleanup: Delete the created URL categories
	_, err = DeleteURLCategories(service, createdResource1.ID)
	if err != nil {
		t.Errorf("Error deleting first URL category: %v", err)
	} else {
		t.Logf("Deleted first URL category with ID: %s", createdResource1.ID)
	}

	_, err = DeleteURLCategories(service, createdResource2.ID)
	if err != nil {
		t.Errorf("Error deleting second URL category: %v", err)
	} else {
		t.Logf("Deleted second URL category with ID: %s", createdResource2.ID)
	}
}
*/
