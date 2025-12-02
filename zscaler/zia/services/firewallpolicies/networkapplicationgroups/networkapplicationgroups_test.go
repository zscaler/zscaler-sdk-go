package networkapplicationgroups

import (
	"context"
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

func TestNetworkApplicationGroups(t *testing.T) {
	tests.ResetTestNameCounter()
	name := tests.GetTestName("tests-netsvc")
	updateName := tests.GetTestName("tests-netsvc")
	client, err := tests.NewVCRTestClient(t, "networkapplicationgroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	nwAppgroup := NetworkApplicationGroups{
		Name:                name,
		Description:         name,
		NetworkApplications: []string{"APNS", "APPSTORE", "DICT"},
	}

	var createdResource *NetworkApplicationGroups

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = Create(context.Background(), service, &nwAppgroup)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

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
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, createdResource.Name)
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

	updatedResource, err := GetNetworkApplicationGroups(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, err = GetNetworkApplicationGroupsByName(context.Background(), service, updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.Name)
	}
	// Test resources retrieval
	resources, err := GetAllNetworkApplicationGroups(context.Background(), service)
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
	// Test resource removal
	err = retryOnConflict(func() error {
		_, delErr := Delete(context.Background(), service, createdResource.ID)
		return delErr
	})
	_, err = GetNetworkApplicationGroups(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *zscaler.Service, id int) (*NetworkApplicationGroups, error) {
	var resource *NetworkApplicationGroups
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = GetNetworkApplicationGroups(context.Background(), s, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}

func TestRetrieveNonExistentResource(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "networkapplicationgroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, err = GetNetworkApplicationGroups(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "networkapplicationgroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, err = Delete(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "networkapplicationgroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, _, err = Update(context.Background(), service, 0, &NetworkApplicationGroups{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "networkapplicationgroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, err = GetNetworkApplicationGroupsByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
