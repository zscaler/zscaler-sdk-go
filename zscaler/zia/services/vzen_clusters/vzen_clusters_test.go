package vzen_clusters

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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/vzen_nodes"
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

func TestVZENClusters(t *testing.T) {
	tests.ResetTestNameCounter()
	// VZEN cluster names must be alphanumeric only (no hyphens, underscores, or spaces)
	name := fmt.Sprintf("VZENCluster%02d", tests.GetTestNameCounter())

	client, err := tests.NewVCRTestClient(t, "vzen_clusters", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// First, fetch available VZen nodes
	availableNodes, err := vzen_nodes.GetAll(context.Background(), service)
	if err != nil {
		t.Fatalf("Error fetching VZen nodes: %v", err)
	}

	// Need at least 2 nodes to create a cluster
	if len(availableNodes) < 2 {
		t.Skip("Skipping test: Need at least 2 VZen nodes to create a cluster")
	}

	// Use the first two available nodes
	var virtualZenNodes []common.IDNameExternalID
	for i := 0; i < 2 && i < len(availableNodes); i++ {
		virtualZenNodes = append(virtualZenNodes, common.IDNameExternalID{
			ID: availableNodes[i].ID,
		})
	}

	vzenCluster := VZENClusters{
		Name:            name,
		Status:          "ENABLED",
		Type:            "VIP",
		IpAddress:       "192.168.90.7",
		SubnetMask:      "255.255.255.0",
		DefaultGateway:  "192.168.90.254",
		IpSecEnabled:    false,
		VirtualZenNodes: virtualZenNodes,
	}

	var createdResource *VZENClusters

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, _, err = Create(context.Background(), service, &vzenCluster)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-zero, but got 0")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created vzen cluster '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved vzen cluster '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Status = "DISABLED"
	retrievedResource.Name = name // Ensure name is preserved in the update payload
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

	// Test resource retrieval by name
	retrievedResource, err = GetClusterByName(context.Background(), service, name)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
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
		t.Errorf("Expected retrieved resources to contain created resource '%d', but it didn't", createdResource.ID)
	}
	// Test resource removal
	err = retryOnConflict(func() error {
		_, delErr := Delete(context.Background(), service, createdResource.ID)
		return delErr
	})
	_, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(ctx context.Context, service *zscaler.Service, id int) (*VZENClusters, error) {
	var resource *VZENClusters
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = Get(ctx, service, id)
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
	client, err := tests.NewVCRTestClient(t, "vzen_clusters", "zia")
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
	client, err := tests.NewVCRTestClient(t, "vzen_clusters", "zia")
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
	client, err := tests.NewVCRTestClient(t, "vzen_clusters", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, _, err = Update(context.Background(), service, 0, &VZENClusters{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "vzen_clusters", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	_, err = GetClusterByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
