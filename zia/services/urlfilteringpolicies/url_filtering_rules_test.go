package urlfilteringpolicies

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/cloudbrowserisolation"
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

func retryGetByName(service *Service, name string) (*URLFilteringRule, error) {
	var rule *URLFilteringRule
	var err error

	for i := 0; i < maxRetries; i++ {
		rule, err = service.GetByName(name)
		if err == nil {
			return rule, nil
		}
		time.Sleep(retryInterval)
	}
	return nil, err
}

func TestURLFilteringRuleIsolation(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)

	cbiService := cloudbrowserisolation.New(client)
	cbiProfileList, err := cbiService.GetAll()
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
		Action:         "ISOLATE",
		URLCategories:  []string{"ANY"},
		Protocols:      []string{"HTTPS_RULE", "HTTP_RULE"},
		RequestMethods: []string{"CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"},
		UserAgentTypes: []string{"OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE"},
		CBIProfile: CBIProfile{
			ProfileSeq: 0,
			ID:         cbiProfileList[0].ID,
			Name:       cbiProfileList[0].Name,
			URL:        cbiProfileList[0].URL,
		},
	}

	var createdResource *URLFilteringRule

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = service.Create(&rule)
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
		t.Errorf("Expected retrieved rule '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Name = updateName
	err = retryOnConflict(func() error {
		_, _, err = service.Update(createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := service.Get(createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name with retry mechanism
	retrievedByNameResource, err := retryGetByName(service, updateName)
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
	allResources, err := service.GetAll()
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
	}
	if len(allResources) == 0 {
		t.Fatal("Expected retrieved resources to be non-empty, but got empty slice")
	}

	// Check if the created resource is in the list
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
	time.Sleep(5 * time.Second) // Sleep for 5 seconds

	// Test resource removal
	err = retryOnConflict(func() error {
		_, getErr := service.Get(createdResource.ID)
		if getErr != nil {
			return fmt.Errorf("Resource %d may have already been deleted: %v", createdResource.ID, getErr)
		}
		_, delErr := service.Delete(createdResource.ID)
		return delErr
	})
	_, err = service.Get(createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

func tryRetrieveResource(service *Service, id int) (*URLFilteringRule, error) {
	var resource *URLFilteringRule
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = service.Get(id)
		if err == nil {
			return resource, nil
		}
		time.Sleep(retryInterval)
	}
	return nil, err
}

// / Testing URL Filtering Rule with BLOCK action
func TestURLFilteringRuleBlock(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-updated-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	rule := URLFilteringRule{
		Name:           name,
		Description:    name,
		Order:          1,
		Rank:           7,
		State:          "ENABLED",
		Action:         "BLOCK",
		URLCategories:  []string{"ANY"},
		Protocols:      []string{"ANY_RULE"},
		RequestMethods: []string{"CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"},
	}

	// Create a URL Filtering Rule
	createdResource, err := createURLFilteringRule(t, service, &rule)
	if err != nil {
		t.Fatalf("Error creating URL Filtering Rule: %v", err)
	}

	defer cleanupURLFilteringRule(t, service, createdResource.ID)

	// Retrieve and check the URL Filtering Rule
	retrievedResource, err := tryRetrieveResource(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving URL Filtering Rule: %v", err)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected name '%s', got '%s'", name, retrievedResource.Name)
	}

	// Update the URL Filtering Rule
	retrievedResource.Name = updateName
	err = updateURLFilteringRule(t, service, retrievedResource)
	if err != nil {
		t.Fatalf("Error updating URL Filtering Rule: %v", err)
	}

	// Retrieve and check the updated URL Filtering Rule
	updatedResource, err := service.Get(createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving updated URL Filtering Rule: %v", err)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected updated name '%s', got '%s'", updateName, updatedResource.Name)
	}
}

func createURLFilteringRule(t *testing.T, service *Service, rule *URLFilteringRule) (*URLFilteringRule, error) {
	var createdResource *URLFilteringRule
	err := retryOnConflict(func() error {
		var err error
		createdResource, err = service.Create(rule)
		if err != nil {
			t.Logf("Error creating URL Filtering Rule: %v", err) // Use t for logging
		}
		return err
	})
	return createdResource, err
}

func updateURLFilteringRule(t *testing.T, service *Service, resource *URLFilteringRule) error {
	return retryOnConflict(func() error {
		_, _, err := service.Update(resource.ID, resource)
		if err != nil {
			t.Logf("Error updating URL Filtering Rule: %v", err) // Use t for logging
		}
		return err
	})
}

func cleanupURLFilteringRule(t *testing.T, service *Service, id int) {
	err := retryOnConflict(func() error {
		_, err := service.Delete(id)
		return err
	})
	if err != nil {
		t.Errorf("Failed to cleanup URL Filtering Rule: %v", err)
	}
}

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Get(0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Delete(0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.Update(0, &URLFilteringRule{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.GetByName("non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
