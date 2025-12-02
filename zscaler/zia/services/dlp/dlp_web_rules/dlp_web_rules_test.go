package dlp_web_rules

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
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

func TestDLPWebRule(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "dlp_web_rules", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	name := tests.GetTestName("tests-dlpweb")
	updateName := tests.GetTestName("tests-dlpweb")

	rule := WebDLPRules{
		Name:                     name,
		Description:              name,
		Order:                    1,
		Rank:                     7,
		State:                    "ENABLED",
		Action:                   "BLOCK",
		WithoutContentInspection: false,
		Severity:                 "RULE_SEVERITY_HIGH",
		Protocols:                []string{"FTP_RULE", "HTTPS_RULE", "HTTP_RULE"},
		CloudApplications:        []string{"SALESFORCE", "GOOGLEANALYTICS", "OTHER_OFFICE365"},
		UserRiskScoreLevels:      []string{"LOW", "MEDIUM", "HIGH", "CRITICAL"},
		URLCategories: []common.IDNameExtensions{
			{
				ID: 2,
			},
		},
	}

	var createdResource *WebDLPRules

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = Create(context.Background(), service, &rule)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

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
		t.Errorf("Expected retrieved dlp engine '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update - create a clean update object to avoid sending read-only fields
	updateRule := WebDLPRules{
		ID:                       createdResource.ID,
		Name:                     updateName,
		Description:              updateName,
		Order:                    retrievedResource.Order,
		Rank:                     retrievedResource.Rank,
		State:                    retrievedResource.State,
		Action:                   retrievedResource.Action,
		WithoutContentInspection: retrievedResource.WithoutContentInspection,
		Severity:                 retrievedResource.Severity,
		Protocols:                retrievedResource.Protocols,
		CloudApplications:        retrievedResource.CloudApplications,
		UserRiskScoreLevels:      retrievedResource.UserRiskScoreLevels,
		URLCategories:            retrievedResource.URLCategories,
	}

	err = retryOnConflict(func() error {
		_, err = Update(context.Background(), service, createdResource.ID, &updateRule)
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

	// Test resource retrieval by name (use the updated name)
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
	time.Sleep(5 * time.Second)

	// Test resource removal
	err = retryOnConflict(func() error {
		_, getErr := Get(context.Background(), service, createdResource.ID)
		if getErr != nil {
			return fmt.Errorf("Resource %d may have already been deleted: %v", createdResource.ID, getErr)
		}
		_, delErr := Delete(context.Background(), service, createdResource.ID)
		return delErr
	})
	if err != nil {
		t.Fatalf("Error deleting resource: %v", err)
	}

	_, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}

	// Test error cases
	t.Run("RetrieveNonExistentResource", func(t *testing.T) {
		_, err := Get(context.Background(), service, 999999999)
		if err == nil {
			t.Error("Expected error retrieving non-existent resource, but got nil")
		}
	})

	t.Run("DeleteNonExistentResource", func(t *testing.T) {
		_, err := Delete(context.Background(), service, 999999999)
		if err == nil {
			t.Error("Expected error deleting non-existent resource, but got nil")
		}
	})

	t.Run("UpdateNonExistentResource", func(t *testing.T) {
		_, err := Update(context.Background(), service, 999999999, &WebDLPRules{})
		if err == nil {
			t.Error("Expected error updating non-existent resource, but got nil")
		}
	})

	t.Run("GetByNameNonExistentResource", func(t *testing.T) {
		_, err := GetByName(context.Background(), service, "non_existent_name")
		if err == nil {
			t.Error("Expected error retrieving resource by non-existent name, but got nil")
		}
	})
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *zscaler.Service, id int) (*WebDLPRules, error) {
	var resource *WebDLPRules
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
