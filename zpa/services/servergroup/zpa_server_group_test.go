package servergroup

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appservercontroller"
)

// clean all resources
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources() // clean up at the beginning
}

func teardown() {
	cleanResources() // clean up at the end
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	return !present || (present && (val == "" || val == "true")) // simplified for clarity
}

func cleanResources() {
	if !shouldClean() {
		return
	}

	client, err := tests.NewZpaClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	service := New(client)
	resources, _, _ := service.GetAll()
	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, _ = service.Delete(r.ID)
	}
}

func TestServerGroup(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}
	// create app connector group for testing
	appConnGroupService := appconnectorgroup.New(client)
	appConnGroup, _, err := appConnGroupService.Create(appconnectorgroup.AppConnectorGroup{
		Name:                     name,
		Description:              name,
		Enabled:                  true,
		CityCountry:              "San Jose, US",
		Latitude:                 "37.3382082",
		Longitude:                "-121.8863286",
		Location:                 "San Jose, CA, USA",
		UpgradeDay:               "SUNDAY",
		UpgradeTimeInSecs:        "66600",
		OverrideVersionProfile:   true,
		VersionProfileName:       "Default",
		VersionProfileID:         "0",
		DNSQueryType:             "IPV4_IPV6",
		PRAEnabled:               false,
		WAFDisabled:              true,
		TCPQuickAckApp:           true,
		TCPQuickAckAssistant:     true,
		TCPQuickAckReadAssistant: true,
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating app connector group for testing server group: %v", err)
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := appConnGroupService.Get(appConnGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := appConnGroupService.Delete(appConnGroup.ID)
			if err != nil {
				t.Errorf("Error deleting app connector group: %v", err)
			}
		}
	}()

	// create app server for testing
	appServerService := appservercontroller.New(client)
	appServer, _, err := appServerService.Create(appservercontroller.ApplicationServer{
		Name:        name,
		Description: name,
		Address:     "192.168.1.1",
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating app server for testing server group: %v", err)
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := appServerService.Get(appServer.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := appServerService.Delete(appServer.ID)
			if err != nil {
				t.Errorf("Error deleting app server: %v", err)
			}
		}
	}()

	// creat

	service := New(client)

	appGroup := ServerGroup{
		Name:        name,
		Description: name,
		AppConnectorGroups: []AppConnectorGroups{
			{
				ID: appConnGroup.ID,
			},
		},
		Servers: []ApplicationServer{
			{
				ID: appServer.ID,
			},
		},
	}

	// Test resource creation
	createdResource, _, err := service.Create(&appGroup)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource retrieval
	retrievedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource update
	retrievedResource.Name = updateName
	_, err = service.Update(createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, _, err = service.GetByName(updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.Name)
	}
	// Test resources retrieval
	resources, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		t.Error("Expected retrieved resources to be non-empty, but got empty slice")
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
		t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.ID)
	}
	// Test resource removal
	_, err = service.Delete(createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	retrievedAfterDelete, _, err := service.Get(createdResource.ID)
	if err != nil {
		// Check if the error implies the resource doesn't exist.
		// Note: This is a basic check. You may want to improve the accuracy based on the exact error message or code.
		if strings.Contains(err.Error(), "resource.not.found") {
			t.Logf("Resource with ID %s has been deleted as expected.", createdResource.ID)
		} else {
			t.Errorf("Unexpected error retrieving resource after delete: %v", err)
		}
		return
	}

	// If no error and the resource still exists, this is unexpected.
	if retrievedAfterDelete != nil && retrievedAfterDelete.ID == createdResource.ID {
		t.Errorf("Expected resource with ID %s to be deleted, but it still exists.", createdResource.ID)
	}
}
