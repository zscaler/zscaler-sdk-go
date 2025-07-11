package urlfilteringpolicies

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/browser_isolation"
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

func TestURLFilteringRuleIsolation(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	cbiProfileList, err := browser_isolation.GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting cbi profile: %v", err)
		return
	}
	if len(cbiProfileList) == 0 {
		t.Error("Expected retrieved cbi profile to be non-empty, but got empty slice")
	}

	rule := URLFilteringRule{
		Name:           name,
		Description:    name,
		Order:          1,
		Rank:           7,
		State:          "ENABLED",
		Action:         "BLOCK",
		URLCategories:  []string{"ANY"},
		Protocols:      []string{"HTTPS_RULE", "HTTP_RULE"},
		RequestMethods: []string{"CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"},
		UserAgentTypes: []string{"OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE"},
		CBIProfile: browser_isolation.CBIProfile{
			ID:   cbiProfileList[0].ID,
			Name: cbiProfileList[0].Name,
			URL:  cbiProfileList[0].URL,
		},
	}

	var createdResource *URLFilteringRule

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = Create(context.Background(), service, &rule)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	// Other assertions based on the creation result
	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-empty, but got ''")
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
		t.Errorf("Expected retrieved url filtering rule '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Name = updateName
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
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}
	// Test resource retrieval by name
	retrievedByNameResource, err := GetByName(context.Background(), service, updateName)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedByNameResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedByNameResource.ID)
	}
	if retrievedByNameResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, retrievedByNameResource.Name)
	}

	// Test resources retrieval
	allResources, err := GetAll(context.Background(), service)
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
		_, getErr := Get(context.Background(), service, createdResource.ID)
		if getErr != nil {
			return fmt.Errorf("Resource %d may have already been deleted: %v", createdResource.ID, getErr)
		}
		_, delErr := Delete(context.Background(), service, createdResource.ID)
		return delErr
	})
	_, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *zscaler.Service, id int) (*URLFilteringRule, error) {
	var resource *URLFilteringRule
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = Get(context.Background(), s, id)
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

	_, _, err = Update(context.Background(), service, 0, &URLFilteringRule{})
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
	_, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}

func TestUrlAndAppSettings(t *testing.T) {
	// Initialize the API client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Context for API calls
	ctx := context.Background()

	// Step 1: Retrieve existing URL and App Settings
	initialSettings, err := GetUrlAndAppSettings(ctx, service)
	if err != nil {
		t.Fatalf("Error retrieving URL and App settings: %v", err)
	}

	t.Logf("Initial settings retrieved: %+v", initialSettings)

	// Step 2: Update the URL and App Settings with test changes
	updatedSettings := *initialSettings // Start with initial settings to ensure a complete payload
	updatedSettings.EnableDynamicContentCat = true
	updatedSettings.ConsiderEmbeddedSites = false
	updatedSettings.EnforceSafeSearch = false
	// updatedSettings.EnableOffice365 = false
	updatedSettings.EnableNewlyRegisteredDomains = false
	updatedSettings.EnableBlockOverrideForNonAuthUser = false

	t.Logf("Payload to be sent: %+v", updatedSettings)

	// Send the update request
	_, _, err = UpdateUrlAndAppSettings(ctx, service, updatedSettings)
	if err != nil {
		t.Fatalf("Error updating URL and App settings: %v", err)
	}
	t.Logf("Updated settings sent successfully")

	// Step 3: Retrieve the updated settings
	retrievedSettings, err := GetUrlAndAppSettings(ctx, service)
	if err != nil {
		t.Fatalf("Error retrieving updated settings: %v", err)
	}
	t.Logf("Retrieved updated settings: %+v", retrievedSettings)

	// Step 4: Validate changes
	if retrievedSettings.EnableDynamicContentCat != updatedSettings.EnableDynamicContentCat {
		t.Errorf("EnableDynamicContentCat mismatch: Expected: %v, Got: %v", updatedSettings.EnableDynamicContentCat, retrievedSettings.EnableDynamicContentCat)
	}
	if retrievedSettings.ConsiderEmbeddedSites != updatedSettings.ConsiderEmbeddedSites {
		t.Errorf("ConsiderEmbeddedSites mismatch: Expected: %v, Got: %v", updatedSettings.ConsiderEmbeddedSites, retrievedSettings.ConsiderEmbeddedSites)
	}
	if retrievedSettings.EnforceSafeSearch != updatedSettings.EnforceSafeSearch {
		t.Errorf("EnforceSafeSearch mismatch: Expected: %v, Got: %v", updatedSettings.EnforceSafeSearch, retrievedSettings.EnforceSafeSearch)
	}
	if retrievedSettings.EnableOffice365 != updatedSettings.EnableOffice365 {
		t.Errorf("EnableOffice365 mismatch: Expected: %v, Got: %v", updatedSettings.EnableOffice365, retrievedSettings.EnableOffice365)
	}
	if retrievedSettings.EnableNewlyRegisteredDomains != updatedSettings.EnableNewlyRegisteredDomains {
		t.Errorf("EnableNewlyRegisteredDomains mismatch: Expected: %v, Got: %v", updatedSettings.EnableNewlyRegisteredDomains, retrievedSettings.EnableNewlyRegisteredDomains)
	}
	if retrievedSettings.EnableBlockOverrideForNonAuthUser != updatedSettings.EnableBlockOverrideForNonAuthUser {
		t.Errorf("EnableBlockOverrideForNonAuthUser mismatch: Expected: %v, Got: %v", updatedSettings.EnableBlockOverrideForNonAuthUser, retrievedSettings.EnableBlockOverrideForNonAuthUser)
	}

	// Step 5: Revert to the original settings
	_, _, err = UpdateUrlAndAppSettings(ctx, service, *initialSettings)
	if err != nil {
		t.Fatalf("Error reverting URL and App settings to original values: %v", err)
	}
	t.Logf("Reverted settings to original values: %+v", initialSettings)
}
