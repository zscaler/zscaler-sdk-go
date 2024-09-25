package locationmanagement

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/staticips"
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

func TestLocationManagement(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("104.239.245.0/24")
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateDescription := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	comment := "tests-" + acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	cleanupTasks := []func(){}

	defer func() {
		for _, task := range cleanupTasks {
			task()
		}
	}()

	staticIP, _, err := staticips.Create(context.Background(), service, &staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   comment,
	})
	if err != nil {
		t.Fatalf("creating static ip failed: %v", err)
		return
	}

	cleanupTasks = append(cleanupTasks, func() {
		_, err := staticips.Delete(context.Background(), service, staticIP.ID)
		if err != nil && !strings.Contains(err.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
			t.Errorf("Error deleting static ip: %v", err)
		}
	})

	location := Locations{
		Name:              name,
		Description:       name,
		Country:           "UNITED_STATES",
		TZ:                "UNITED_STATES_AMERICA_LOS_ANGELES",
		AuthRequired:      true,
		IdleTimeInMinutes: 720,
		DisplayTimeUnit:   "HOUR",
		SurrogateIP:       true,
		XFFForwardEnabled: true,
		OFWEnabled:        true,
		IPSControl:        true,
		Profile:           "CORPORATE",
		IPAddresses:       []string{ipAddress},
	}

	var createdResource *Locations

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = Create(context.Background(), service, &location)
		return err
	})
	// Check if the request was successful
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	alreadyDeleted := false

	// This deferred deletion should be declared right after confirming resource creation
	defer func() {
		if alreadyDeleted {
			return
		}
		err = retryOnConflict(func() error {
			_, delErr := Delete(context.Background(), service, createdResource.ID)
			return delErr
		})
		if err != nil {
			t.Errorf("Error deleting rule: %v", err)
		}
	}()

	if createdResource.ID == 0 {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}

	// Pause the test for 30 seconds before creating the sub-location
	time.Sleep(15 * time.Second)

	// Create sub-location
	subLocation := Locations{
		Name:              name + "-sub",
		Description:       name + "-sub",
		Country:           "UNITED_STATES",
		TZ:                "UNITED_STATES_AMERICA_LOS_ANGELES",
		AuthRequired:      true,
		IdleTimeInMinutes: 720,
		DisplayTimeUnit:   "HOUR",
		SurrogateIP:       true,
		XFFForwardEnabled: true,
		OFWEnabled:        true,
		IPSControl:        true,
		IPAddresses:       []string{"10.6.0.0-10.6.255.255"},
		ParentID:          createdResource.ID,
		Profile:           "CORPORATE",
	}

	var createdSubLocation *Locations

	// Test sub-location creation
	err = retryOnConflict(func() error {
		createdSubLocation, err = Create(context.Background(), service, &subLocation)
		return err
	})
	// Check if the request was successful
	if err != nil {
		t.Fatalf("Error making POST request for sub-location: %v", err)
	}

	if createdSubLocation.ID == 0 {
		t.Error("Expected created sub-location ID to be non-empty, but got ''")
	}
	if createdSubLocation.Name != subLocation.Name {
		t.Errorf("Expected created sub-location name '%s', but got '%s'", subLocation.Name, createdSubLocation.Name)
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
	retrievedResource.Description = updateDescription
	retrievedResource.Name = updateName // Added this line
	err = retryOnConflict(func() error {
		_, _, err = Update(context.Background(), service, createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}
	updatedResource, err := GetLocation(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.Description != updateDescription {
		t.Errorf("Expected updated resource description '%s', but got '%s'", updateDescription, updatedResource.Description)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, err = GetLocationByName(context.Background(), service, updateName)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.Name)
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

	// Test sub-location retrieval by ID
	retrievedSubLocation, err := tryRetrieveResource(service, createdSubLocation.ID)
	if err != nil {
		t.Fatalf("Error retrieving sub-location: %v", err)
	}
	if retrievedSubLocation.ID != createdSubLocation.ID {
		t.Errorf("Expected retrieved sub-location ID '%d', but got '%d'", createdSubLocation.ID, retrievedSubLocation.ID)
	}
	if retrievedSubLocation.Name != subLocation.Name {
		t.Errorf("Expected retrieved sub-location name '%s', but got '%s'", subLocation.Name, retrievedSubLocation.Name)
	}

	// Test GetSubLocation
	subLocationByParent, err := GetSubLocation(context.Background(), service, createdResource.ID, createdSubLocation.ID)
	if err != nil {
		t.Fatalf("Error getting sub-location by parent ID: %v", err)
	}
	if subLocationByParent.ID != createdSubLocation.ID {
		t.Errorf("Expected sub-location ID '%d', but got '%d'", createdSubLocation.ID, subLocationByParent.ID)
	}

	// Test GetSubLocationBySubID
	subLocationBySubID, err := GetSubLocationBySubID(context.Background(), service, createdSubLocation.ID)
	if err != nil {
		t.Fatalf("Error getting sub-location by sub ID: %v", err)
	}
	if subLocationBySubID.ID != createdSubLocation.ID {
		t.Errorf("Expected sub-location ID '%d', but got '%d'", createdSubLocation.ID, subLocationBySubID.ID)
	}

	// Test GetAllSublocations
	allSubLocations, err := GetAllSublocations(context.Background(), service)
	if err != nil {
		t.Fatalf("Error getting all sub-locations: %v", err)
	}
	if len(allSubLocations) == 0 {
		t.Fatalf("Expected sub-locations, but got none")
	}

	// Test GetSubLocationByName
	subLocationByName, err := GetSubLocationByName(context.Background(), service, subLocation.Name)
	if err != nil {
		t.Fatalf("Error getting sub-location by name: %v", err)
	}
	if subLocationByName.ID != createdSubLocation.ID {
		t.Errorf("Expected sub-location ID '%d', but got '%d'", createdSubLocation.ID, subLocationByName.ID)
	}

	// Test GetSubLocationByNames
	subLocationByNames, err := GetSubLocationByNames(context.Background(), service, updateName, subLocation.Name)
	if err != nil {
		t.Fatalf("Error getting sub-location by names: %v", err)
	}
	if subLocationByNames.ID != createdSubLocation.ID {
		t.Errorf("Expected sub-location ID '%d', but got '%d'", createdSubLocation.ID, subLocationByNames.ID)
	}

	// Test GetLocationOrSublocationByName
	locationOrSubLocation, err := GetLocationOrSublocationByName(context.Background(), service, subLocation.Name)
	if err != nil {
		t.Fatalf("Error getting location or sub-location by name: %v", err)
	}
	if locationOrSubLocation.ID != createdSubLocation.ID {
		t.Errorf("Expected sub-location ID '%d', but got '%d'", createdSubLocation.ID, locationOrSubLocation.ID)
	}

	// Test resource removal
	err = retryOnConflict(func() error {
		_, delErr := Delete(context.Background(), service, createdResource.ID)
		return delErr
	})
	if err != nil && !strings.Contains(err.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
		t.Fatalf("Error deleting resource: %v", err)
	}
	alreadyDeleted = true // Set the flag to true after successfully deleting the rule

	_, err = GetLocation(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *zscaler.Service, id int) (*Locations, error) {
	var result *Locations
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		result, lastErr = GetLocationOrSublocationByID(context.Background(), s, id)
		if lastErr == nil {
			return result, nil
		}

		time.Sleep(retryInterval)
	}

	return nil, lastErr
}

func TestRetrieveNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = GetLocation(context.Background(), service, 0)
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

	_, _, err = Update(context.Background(), service, 0, &Locations{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = GetLocationByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
