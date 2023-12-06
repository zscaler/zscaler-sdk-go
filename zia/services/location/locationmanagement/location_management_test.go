package locationmanagement

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
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

	client, err := tests.NewZiaClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	service := New(client)
	resources, err := service.GetAll()
	if err != nil {
		log.Printf("Error retrieving resources during cleanup: %v", err)
		return
	}

	for _, r := range resources {
		if strings.HasPrefix(r.Name, "tests-") {
			_, err := service.Delete(r.ID)
			if err != nil {
				log.Printf("Error deleting resource %d: %v", r.ID, err)
			}
		}
	}
}

func TestLocationManagement(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("104.239.236.0/24")
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateDescription := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	cleanupTasks := []func(){}

	defer func() {
		for _, task := range cleanupTasks {
			task()
		}
	}()

	// static ip for location management testing
	staticipsService := staticips.New(client)
	staticIP, _, err := staticipsService.Create(&staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   "testing static ip for location management",
	})
	if err != nil {
		t.Fatalf("creating static ip failed: %v", err)
		return
	}

	cleanupTasks = append(cleanupTasks, func() {
		_, err := staticipsService.Delete(staticIP.ID)
		if err != nil && !strings.Contains(err.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
			t.Errorf("Error deleting static ip: %v", err)
		}
	})

	service := New(client)
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
		IPAddresses:       []string{ipAddress},
	}

	var createdResource *Locations

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = service.Create(&location)
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
			_, delErr := service.Delete(createdResource.ID)
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
		_, _, err = service.Update(createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}
	updatedResource, err := service.GetLocation(createdResource.ID)
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
	retrievedResource, err = service.GetLocationByName(updateName)
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
	resources, err := service.GetAll()
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
		_, delErr := service.Delete(createdResource.ID)
		return delErr
	})
	if err != nil && !strings.Contains(err.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
		t.Fatalf("Error deleting resource: %v", err)
	}
	alreadyDeleted = true // Set the flag to true after successfully deleting the rule

	_, err = service.GetLocation(createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *Service, id int) (*Locations, error) {
	var result *Locations
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		result, lastErr = s.GetLocation(id)
		if lastErr == nil {
			return result, nil
		}

		time.Sleep(retryInterval)
	}

	return nil, lastErr
}
