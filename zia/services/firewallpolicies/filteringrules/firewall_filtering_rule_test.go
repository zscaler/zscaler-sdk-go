package filteringrules

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/ipdestinationgroups"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/ipsourcegroups"
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

func TestFirewallFilteringRule(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateDescription := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	cleanupTasks := []func(){}

	defer func() {
		for _, task := range cleanupTasks {
			task()
		}
	}()

	// create ip source group for testing
	sourceIPGroupService := ipsourcegroups.New(client)
	sourceIPGroup, err := sourceIPGroupService.Create(&ipsourcegroups.IPSourceGroups{
		Name:        name,
		Description: name,
		IPAddresses: []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"},
	})
	if err != nil {
		t.Fatalf("Error creating source ip group for testing server group: %v", err)
	}

	cleanupTasks = append(cleanupTasks, func() {
		_, err := sourceIPGroupService.Delete(sourceIPGroup.ID)
		if err != nil && !strings.Contains(err.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
			t.Errorf("Error deleting source ip group: %v", err)
		}
	})

	// create ip destination group for testing
	dstIPGroupService := ipdestinationgroups.New(client)
	dstIPGroup, err := dstIPGroupService.Create(&ipdestinationgroups.IPDestinationGroups{
		Name:        name,
		Description: name,
		Type:        "DSTN_FQDN",
		Addresses:   []string{"test1.acme.com", "test2.acme.com", "test3.acme.com"},
	})
	if err != nil {
		t.Fatalf("Error creating ip destination group for testing server group: %v", err)
	}

	cleanupTasks = append(cleanupTasks, func() {
		_, err := dstIPGroupService.Delete(dstIPGroup.ID)
		if err != nil && !strings.Contains(err.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
			t.Errorf("Error deleting destination ip group: %v", err)
		}
	})

	service := New(client)
	rule := FirewallFilteringRules{
		Name:           name,
		Description:    name,
		Order:          5,
		Rank:           7,
		Action:         "ALLOW",
		DestCountries:  []string{"COUNTRY_CA", "COUNTRY_US", "COUNTRY_MX", "COUNTRY_AU", "COUNTRY_GB"},
		NwApplications: []string{"APNS", "GARP", "PERFORCE", "WINDOWS_MARKETPLACE", "DIAMETER"},
		SrcIpGroups: []common.IDNameExtensions{
			{
				ID: sourceIPGroup.ID,
			},
		},
		DestIpGroups: []common.IDNameExtensions{
			{
				ID: dstIPGroup.ID,
			},
		},
	}

	var createdResource *FirewallFilteringRules

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, err = service.Create(&rule)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	defer func() {
		// Delete the main rule first
		err = retryOnConflict(func() error {
			_, delErr := service.Delete(createdResource.ID)
			return delErr
		})
		if err != nil && !strings.Contains(err.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
			t.Errorf("Error deleting rule: %v", err)
		}

		// Then delete secondary resources
		for _, task := range cleanupTasks {
			task()
		}
	}()

	// Add the cleanup task for the main resource
	cleanupTasks = append(cleanupTasks, func() {
		_, err := service.Delete(createdResource.ID)
		if err != nil && !strings.Contains(err.Error(), `"code":"RESOURCE_NOT_FOUND"`) {
			t.Errorf("Error deleting rule: %v", err)
		}
	})

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
		_, err = service.Update(createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}
	updatedResource, err := service.Get(createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.Description != updateDescription {
		t.Errorf("Expected updated resource description '%s', but got '%s'", updateDescription, updatedResource.Description)
	}
	if updatedResource.Name != updateName { // Added this line
		t.Errorf("Expected updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, err = service.GetByName(updateName)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource description '%s', but got '%s'", updateName, createdResource.Name)
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
	// alreadyDeleted = true // Set the flag to true after successfully deleting the rule

	_, err = service.Get(createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *Service, id int) (*FirewallFilteringRules, error) {
	var result *FirewallFilteringRules
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		result, lastErr = s.Get(id)
		if lastErr == nil {
			return result, nil
		}

		time.Sleep(retryInterval)
	}

	return nil, lastErr
}
