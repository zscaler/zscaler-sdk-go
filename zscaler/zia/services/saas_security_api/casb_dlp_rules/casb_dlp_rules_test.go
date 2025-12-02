package casb_dlp_rules

/*
import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

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

func TestCasbDLPRules(t *testing.T) {
	tests.ResetTestNameCounter()
	name := tests.GetTestName("tests-casbdlp")
	updateName := tests.GetTestName("tests-casbdlp")
	client, err := tests.NewVCRTestClient(t, "casb_dlp_rules", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Create a minimal rule without tenant-specific IDs that may not exist
	rule := CasbDLPRules{
		Name:                     name,
		Description:              name,
		Order:                    1,
		Rank:                     7,
		State:                    "ENABLED",
		Type:                     "OFLCASB_DLP_ITSM",
		Severity:                 "RULE_SEVERITY_HIGH",
		Action:                   "OFLCASB_DLP_BLOCK",
		FileTypes:                []string{"FTCATEGORY_JAVASCRIPT", "FTCATEGORY_TAR"},
		WithoutContentInspection: true,
		DeviceTrustLevels:        []string{"UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"},
		Components:               []string{"COMPONENT_ITSM_ATTACHMENTS", "COMPONENT_ITSM_OBJECTS"},
		CollaborationScope: []string{"COLLABORATION_SCOPE_EXTERNAL_COLLAB_EDIT",
			"COLLABORATION_SCOPE_EXTERNAL_COLLAB_VIEW",
			"COLLABORATION_SCOPE_INTERNAL_COLLAB_EDIT",
			"COLLABORATION_SCOPE_INTERNAL_COLLAB_VIEW",
			"COLLABORATION_SCOPE_PRIVATE_EDIT",
			"COLLABORATION_SCOPE_PRIVATE"},
	}

	createdResource, err := Create(context.Background(), service, &rule)
	time.Sleep(1 * time.Second)
	if err != nil {
		t.Fatalf("Error creating CASB DLP Rule resource: %v", err)
	}

	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(service, rule.Type, createdResource.ID)
	time.Sleep(1 * time.Second)
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
		_, err = Update(context.Background(), service, createdResource.ID, retrievedResource)
		time.Sleep(1 * time.Second)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := GetByRuleID(context.Background(), service, rule.Type, createdResource.ID)
	time.Sleep(1 * time.Second)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	allRules, err := GetAll(context.Background(), service)
	time.Sleep(1 * time.Second)
	if err != nil {
		t.Fatalf("Error retrieving all rules via GetAll: %v", err)
	}
	if len(allRules) == 0 {
		t.Fatal("Expected GetAll to return non-empty list of rules")
	}

	// Introduce a delay before deleting
	time.Sleep(5 * time.Second)

	// Test resources retrieval
	allResources, err := GetByRuleType(context.Background(), service, rule.Type)
	time.Sleep(1 * time.Second)
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
	time.Sleep(5 * time.Second)

	// Test resource removal
	err = retryOnConflict(func() error {
		_, getErr := GetByRuleID(context.Background(), service, rule.Type, createdResource.ID)
		if getErr != nil {
			return fmt.Errorf("Resource %d may have already been deleted: %v", createdResource.ID, getErr)
		}
		_, delErr := Delete(context.Background(), service, rule.Type, createdResource.ID)
		time.Sleep(1 * time.Second)
		return delErr
	})
	if err != nil {
		t.Fatalf("Error deleting resource: %v", err)
	}

	_, err = GetByRuleID(context.Background(), service, rule.Type, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *zscaler.Service, ruleType string, id int) (*CasbDLPRules, error) {
	var resource *CasbDLPRules
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = GetByRuleID(context.Background(), s, ruleType, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}
*/
