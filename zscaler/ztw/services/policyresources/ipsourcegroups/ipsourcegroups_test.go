package ipsourcegroups

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

func TestFWFileringIPSourceGroups(t *testing.T) {
	tests.ResetTestNameCounter()
	ipAddress1 := tests.GetTestIP("192.168.1.0/24")
	ipAddress2 := tests.GetTestIP("192.168.2.0/24")
	name := tests.GetTestName("tests-ipsrc")
	updateName := tests.GetTestName("tests-ipsrc")

	client, err := tests.NewVCRTestClient(t, "ipsourcegroups", "ztw")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	group := IPSourceGroups{
		Name:        name,
		Description: name,
		IPAddresses: []string{ipAddress1, ipAddress2},
	}

	// Test resource creation
	var createdResource *IPSourceGroups
	err = retryOnConflict(func() error {
		createdResource, _, err = Create(context.Background(), service, &group)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}
	log.Printf("Created resource with ID: %d", createdResource.ID)

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
	retrievedResource.IPAddresses = group.IPAddresses

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
	if updatedResource == nil {
		t.Fatal("Updated resource is nil")
	}

	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, err = GetByName(context.Background(), service, updateName)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource == nil {
		t.Fatal("Retrieved resource by name is nil")
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.Name)
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
			if strings.Contains(getErr.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
				log.Printf("Resource %d already deleted.", createdResource.ID)
				return nil
			}
			return fmt.Errorf("Error retrieving resource %d: %v", createdResource.ID, getErr)
		}
		_, delErr := Delete(context.Background(), service, createdResource.ID)
		if delErr != nil {
			if strings.Contains(delErr.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
				log.Printf("Resource %d already deleted.", createdResource.ID)
				return nil
			}
			return delErr
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Unexpected error during deletion: %v", err)
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *zscaler.Service, id int) (*IPSourceGroups, error) {
	var resource *IPSourceGroups
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
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "ipsourcegroups", "ztw")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, err = Get(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "ipsourcegroups", "ztw")
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
	client, err := tests.NewVCRTestClient(t, "ipsourcegroups", "ztw")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, _, err = Update(context.Background(), service, 0, &IPSourceGroups{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "ipsourcegroups", "ztw")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
