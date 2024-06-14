package gretunnels

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
	virtualipaddress "github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/virtualipaddress"
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

func TestGRETunnel(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("104.239.238.0/24")
	comment := "tests-" + acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	updateComment := "tests-" + acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	// create static IP for testing
	service := services.New(client)

	staticIP, _, err := staticips.Create(service, &staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   comment,
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating static IP for testing: %v", err)
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, getErr := staticips.Get(service, staticIP.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := staticips.Delete(service, staticIP.ID)
			if err != nil {
				t.Errorf("Error deleting static ip: %v", err)
			}
		}
	}()

	vipRecommendedList, err := virtualipaddress.GetAll(service, ipAddress)
	if err != nil {
		t.Errorf("Error getting recommended vip: %v", err)
		return
	}
	if len(vipRecommendedList) == 0 {
		t.Error("Expected retrieved recommended vip to be non-empty, but got empty slice")
	}

	withinCountry := true // Create a boolean variable
	greTunnel := GreTunnels{
		SourceIP:      staticIP.IpAddress,
		Comment:       comment,
		WithinCountry: &withinCountry,
		IPUnnumbered:  true,
		PrimaryDestVip: &PrimaryDestVip{
			ID:        vipRecommendedList[0].ID,
			VirtualIP: vipRecommendedList[0].VirtualIp,
		},
		SecondaryDestVip: &SecondaryDestVip{
			ID:        vipRecommendedList[1].ID,
			VirtualIP: vipRecommendedList[1].VirtualIp,
		},
	}

	// Inside TestZPAGateways function
	createdResource, _, err := CreateGreTunnels(service, &greTunnel)
	if err != nil {
		t.Fatalf("Error creating GRE Tunnel resource: %v", err)
	}

	defer func() {
		// Attempt to delete the resource
		_, delErr := DeleteGreTunnels(service, createdResource.ID)
		if delErr != nil {
			// If the error indicates the resource is already deleted, log it as information
			if strings.Contains(delErr.Error(), "404") || strings.Contains(delErr.Error(), "RESOURCE_NOT_FOUND") {
				t.Logf("Resource with ID %d not found (already deleted).", createdResource.ID)
			} else {
				// If the deletion error is not due to the resource being missing, log it as an actual error
				t.Errorf("Error deleting ZPAGateways resource: %v", delErr)
			}
		}
	}()

	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Comment != comment {
		t.Errorf("Expected retrieved zpa gateway '%s', but got '%s'", comment, retrievedResource.Comment)
	}

	// Test resource update
	retrievedResource.Comment = updateComment

	err = retryOnConflict(func() error {
		_, _, err = UpdateGreTunnels(service, createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := GetGreTunnels(service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Comment != updateComment {
		t.Errorf("Expected retrieved updated resource description '%s', but got '%s'", updateComment, updatedResource.Comment)
	}

	// Test resource retrieval by name
	retrievedResource, err = GetByIPAddress(service, ipAddress)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Comment != updateComment {
		t.Errorf("Expected retrieved resource comment '%s', but got '%s'", updateComment, createdResource.Comment)
	}
	// Test resources retrieval
	resources, err := GetAll(service)
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
		_, delErr := DeleteGreTunnels(service, createdResource.ID)
		if delErr != nil {
			// Check if the error is due to the resource not being found (i.e., already deleted)
			if strings.Contains(delErr.Error(), "404") || strings.Contains(delErr.Error(), "RESOURCE_NOT_FOUND") {
				log.Printf("Resource with ID %d not found (already deleted).", createdResource.ID)
				return nil // No error, as the resource is already deleted
			}
			return delErr // Return the actual error for other cases
		}
		return nil // No error, deletion successful
	})
	if err != nil {
		t.Errorf("Error occurred during resource deletion: %v", err)
	} else {
		t.Logf("Resource deleted successfully.")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *services.Service, id int) (*GreTunnels, error) {
	var resource *GreTunnels
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = GetGreTunnels(s, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = GetGreTunnels(service, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = DeleteGreTunnels(service, 0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, _, err = UpdateGreTunnels(service, 0, &GreTunnels{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByIPAddressNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, err = GetByIPAddress(service, "non-existent-ip-address")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent ip address, but got nil")
	}
}
