package locationtemplate

import (
	"context"
	"fmt"
	"log"
	"os"
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

func TestMain(m *testing.M) {
	tests.ResetTestNameCounter()
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources()
}

func teardown() {
	cleanResources()
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	return !present || (present && (val == "" || val == "true")) // simplified for clarity
}

func cleanResources() {
	if !shouldClean() {
		return
	}

	service, err := tests.NewOneAPIClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	resources, err := GetAll(context.Background(), service)
	if err != nil {
		log.Printf("Error retrieving resources during cleanup: %v", err)
		return
	}

	for _, r := range resources {
		if strings.HasPrefix(r.Name, "tests-") {
			_, err := Delete(context.Background(), service, r.ID)
			if err != nil {
				log.Printf("Error deleting resource %d: %v", r.ID, err)
			}
		}
	}
}

func TestZConLocationTemplate(t *testing.T) {
	tests.ResetTestNameCounter()
	name := tests.GetTestName("tests-locati")
	description := tests.GetTestName("tests-locati")
	updateName := tests.GetTestName("tests-locati")
	client, err := tests.NewVCRTestClient(t, "locationtemplate", "ztw")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	template := LocationTemplate{
		Name:        name,
		Description: description,
		LocationTemplateDetails: &LocationTemplateDetails{
			TemplatePrefix:          "zcon-prefix",
			AupEnabled:              true,
			AupTimeoutInDays:        30,
			AuthRequired:            true,
			IPSControl:              true,
			OFWEnabled:              true,
			XFFForwardEnabled:       true,
			EnforceBandwidthControl: true,
			UpBandwidth:             10,
			DnBandwidth:             10,
			// DisplayTimeUnit:         "DAY",
		},
	}

	var createdResource *LocationTemplate
	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = Create(context.Background(), service, &template)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	// Other assertions based on the creation result
	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-empty, but got ''")
	}
	expectedName := name
	if createdResource.Name != expectedName {
		t.Errorf("Expected created resource name '%s', but got '%s'", expectedName, createdResource.Name)
	}

	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != expectedName {
		t.Errorf("Expected retrieved location template '%s', but got '%s'", expectedName, retrievedResource.Name)
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
func tryRetrieveResource(ctx context.Context, s *zscaler.Service, id int) (*LocationTemplate, error) {
	var resource *LocationTemplate
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = Get(ctx, s, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}
