package cloudappcontrol

import (
	"bytes"
	"encoding/json"
	"fmt"
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

		if strings.Contains(lastErr.Error(), `"code":"EDIT_LOCK_NOT_AVAILABLE"`) {
			log.Printf("Conflict error detected, retrying in %v... (Attempt %d/%d)", conflictRetryInterval, i+1, maxConflictRetries)
			time.Sleep(conflictRetryInterval)
			continue
		}

		return lastErr
	}
	return lastErr
}

func TestCloudAppControlRule(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	ziaClient, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(ziaClient)
	rule := WebApplicationRules{
		Name:         name,
		Description:  name,
		Order:        1,
		Rank:         7,
		State:        "ENABLED",
		Type:         "STREAMING_MEDIA",
		Applications: []string{"YOUTUBE", "GOOGLE_STREAMING"},
		Actions:      []string{"ALLOW_STREAMING_VIEW_LISTEN", "ALLOW_STREAMING_UPLOAD"},
	}

	createdResource, err := Create(service, rule.Type, &rule)
	time.Sleep(1 * time.Second) // Adding delay
	if err != nil {
		t.Fatalf("Error creating Web Application Rule resource: %v", err)
	}

	defer func() {
		time.Sleep(1 * time.Second) // Adding delay
		// Check if the resource exists before attempting to delete
		_, getErr := GetByRuleID(service, createdResource.Type, createdResource.ID)
		if getErr == nil {
			_, delErr := Delete(service, createdResource.Type, createdResource.ID)
			if delErr != nil {
				// If the error indicates the resource is already deleted, log it as information
				if strings.Contains(delErr.Error(), "409") || strings.Contains(delErr.Error(), "RESOURCE_NOT_FOUND") {
					t.Logf("Resource with ID %d not found (already deleted).", createdResource.ID)
				} else {
					// If the deletion error is not due to the resource being missing, log it as an actual error
					t.Errorf("Error deleting Web Application Rule resource: %v", delErr)
				}
			}
		} else {
			t.Logf("Resource with ID %d not found (already deleted).", createdResource.ID)
		}
	}()

	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(service, rule.Type, createdResource.ID)
	time.Sleep(1 * time.Second) // Adding delay
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved rule name '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Name = updateName
	err = retryOnConflict(func() error {
		_, err = Update(service, rule.Type, createdResource.ID, retrievedResource)
		time.Sleep(1 * time.Second) // Adding delay
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := GetByRuleID(service, rule.Type, createdResource.ID)
	time.Sleep(1 * time.Second) // Adding delay
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resources retrieval
	allResources, err := GetByRuleType(service, rule.Type)
	time.Sleep(1 * time.Second) // Adding delay
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
	time.Sleep(5 * time.Second) // sleep for 5 seconds

	// Test resource removal
	err = retryOnConflict(func() error {
		_, getErr := GetByRuleID(service, rule.Type, createdResource.ID)
		if getErr != nil {
			return fmt.Errorf("Resource %d may have already been deleted: %v", createdResource.ID, getErr)
		}
		_, delErr := Delete(service, rule.Type, createdResource.ID)
		time.Sleep(1 * time.Second) // Adding delay
		return delErr
	})
	_, err = GetByRuleID(service, rule.Type, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *services.Service, ruleType string, id int) (*WebApplicationRules, error) {
	var resource *WebApplicationRules
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = GetByRuleID(s, ruleType, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}

func TestAllAvailableActions(t *testing.T) {
	ziaClient, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(ziaClient)
	ruleType := "STREAMING_MEDIA" // Adjust as necessary for your specific use case

	// Set up the payload
	payload := AvailableActionsRequest{
		CloudApps: []string{"DROPBOX"},
		Type:      "ANY",
	}

	// Call the AllAvailableActions function
	actions, err := AllAvailableActions(service, ruleType, payload)
	time.Sleep(1 * time.Second) // Adding delay
	if err != nil {
		t.Fatalf("Error getting available actions: %v", err)
	}

	// Verify the response
	expectedActions := []string{
		"ALLOW_FILE_SHARE_CREATE",
		"ALLOW_FILE_SHARE_DELETE",
		"ALLOW_FILE_SHARE_DOWNLOAD",
		"ALLOW_FILE_SHARE_EDIT",
		"ALLOW_FILE_SHARE_INVITE",
		"ALLOW_FILE_SHARE_RENAME",
		"ALLOW_FILE_SHARE_SHARE",
		"DENY_FILE_SHARE_CREATE",
		"DENY_FILE_SHARE_DELETE",
		"DENY_FILE_SHARE_DOWNLOAD",
		"DENY_FILE_SHARE_EDIT",
		"DENY_FILE_SHARE_INVITE",
		"DENY_FILE_SHARE_RENAME",
		"DENY_FILE_SHARE_SHARE",
	}

	if len(actions) != len(expectedActions) {
		t.Errorf("Expected %d actions, but got %d", len(expectedActions), len(actions))
	}

	for i, action := range actions {
		if action != expectedActions[i] {
			t.Errorf("Expected action %s, but got %s", expectedActions[i], action)
		}
	}
}

func TestRuleTypeMapping(t *testing.T) {
	ziaClient, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(ziaClient)

	// Invoke the GetRuleTypeMapping function
	mappings, err := GetRuleTypeMapping(service)
	if err != nil {
		t.Fatalf("Error getting rule type mappings: %v", err)
	}

	// Check if the response is not empty
	if len(mappings) == 0 {
		t.Error("Expected rule type mappings, but got an empty response")
	}

	// Optionally, check if specific keys are present in the response
	expectedKeys := []string{
		"Webmail", "Social Networking", "Finance", "Legal", "AI & ML Applications",
	}

	for _, key := range expectedKeys {
		if _, exists := mappings[key]; !exists {
			t.Errorf("Expected key %s in the response, but it was not found", key)
		}
	}

	// Convert the map back to JSON for logging with HTML escaping disabled
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(mappings)
	if err != nil {
		t.Fatalf("Error marshalling rule type mappings to JSON: %v", err)
	}

	t.Logf("Rule Type Mappings: %s", buffer.String())
}

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)
	ruleType := "STREAMING_MEDIA" // Adjust as necessary for your specific use case

	_, err = GetByRuleID(service, ruleType, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
	time.Sleep(1 * time.Second) // Adding delay
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)
	ruleType := "STREAMING_MEDIA" // Adjust as necessary for your specific use case

	_, err = Delete(service, ruleType, 0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
	time.Sleep(1 * time.Second) // Adding delay
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)
	ruleType := "STREAMING_MEDIA" // Adjust as necessary for your specific use case

	_, err = Update(service, ruleType, 0, &WebApplicationRules{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
	time.Sleep(1 * time.Second) // Adding delay
}
