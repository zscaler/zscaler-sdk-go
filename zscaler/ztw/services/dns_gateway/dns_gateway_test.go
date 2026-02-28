package dnsgateway

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
)

const (
	maxRetries    = 3
	retryInterval = 2 * time.Second
)

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

func TestDNSGatewayCRUD(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	gateway := DNSGateway{
		Name:                         name,
		DNSGatewayType:               "EC_DNS_GW",
		ECDnsGatewayOptionsPrimary:   "LAN_PRI_DNS_AS_PRI",
		ECDnsGatewayOptionsSecondary: "LAN_SEC_DNS_AS_SEC",
		FailureBehavior:              "FAIL_RET_ERR",
	}

	var createdResource *DNSGateway

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = Create(context.Background(), service, &gateway)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-zero, but got 0")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created DNS gateway name '%s', but got '%s'", name, createdResource.Name)
	}

	defer func() {
		err = retryOnConflict(func() error {
			_, delErr := Delete(context.Background(), service, createdResource.ID)
			return delErr
		})
		if err != nil {
			t.Errorf("Error deleting DNS gateway: %v", err)
		}
	}()

	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved DNS gateway name '%s', but got '%s'", name, retrievedResource.Name)
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
	retrievedByName, err := GetByName(context.Background(), service, updateName)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedByName.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedByName.ID)
	}

	// Test GetAll
	resources, err := GetAll(context.Background(), service)
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		t.Fatal("Expected retrieved resources to be non-empty, but got empty slice")
	}
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

	// Test GetAllLite
	resourcesLite, err := GetAllLite(context.Background(), service)
	if err != nil {
		t.Fatalf("Error retrieving lite resources: %v", err)
	}
	found = false
	for _, resource := range resourcesLite {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected GetAllLite to contain created resource '%d', but it didn't", createdResource.ID)
	}
}

func tryRetrieveResource(ctx context.Context, service *zscaler.Service, id int) (*DNSGateway, error) {
	var resource *DNSGateway
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
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = Get(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = Delete(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, _, err = Update(context.Background(), service, 0, &DNSGateway{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = GetByName(context.Background(), service, "non_existent_dns_gateway")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}

func TestDNSGatewayGetAll(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	gateways, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting DNS gateways: %v", err)
		return
	}
	if len(gateways) == 0 {
		t.Log("No DNS gateways found. Moving on with other tests.")
	} else {
		name := gateways[0].Name
		t.Log("Getting DNS gateway by name: " + name)
		gw, err := GetByName(context.Background(), service, name)
		if err != nil {
			t.Errorf("Error getting DNS gateway by name: %v", err)
			return
		}
		if gw.Name != name {
			t.Errorf("DNS gateway name does not match: expected %s, got %s", name, gw.Name)
		}
	}

	nonExistentName := "ThisDNSGatewayDoesNotExist"
	_, err = GetByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
	} else {
		t.Log("Correctly received error when attempting to get non-existent DNS gateway")
	}
}

func TestDNSGatewayCaseSensitivity(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	gateways, err := GetAll(context.Background(), service)
	if err != nil {
		t.Fatalf("Error getting DNS gateways: %v", err)
	}
	if len(gateways) == 0 {
		t.Log("No DNS gateways found. Skipping case sensitivity test.")
		return
	}

	knownName := gateways[0].Name
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
	}

	for _, variation := range variations {
		t.Run(fmt.Sprintf("variation_%s", variation), func(t *testing.T) {
			gw, err := GetByName(context.Background(), service, variation)
			if err != nil {
				t.Errorf("Error getting DNS gateway with name variation '%s': %v", variation, err)
				return
			}
			if gw.Name != knownName {
				t.Errorf("Expected DNS gateway name '%s' for variation '%s', but got '%s'", knownName, variation, gw.Name)
			}
		})
	}
}
