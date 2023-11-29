package zpa_gateways

import (
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/applicationsegment"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/servergroup"
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

func TestZPAGateways(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	description := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateDescription := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rPort := strconv.Itoa(acctest.RandIntRange(1000, 9999))

	zpaClient, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	// create app connector group for testing
	appConnGroupService := appconnectorgroup.New(zpaClient)
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
	serverGroupService := servergroup.New(zpaClient)
	serverGroup, _, err := serverGroupService.Create(&servergroup.ServerGroup{
		Name:             name,
		Description:      name,
		Enabled:          true,
		DynamicDiscovery: true,
		AppConnectorGroups: []servergroup.AppConnectorGroups{
			{
				ID: appConnGroup.ID,
			},
		},
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating app server for testing server group: %v", err)
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := serverGroupService.Get(serverGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := serverGroupService.Delete(serverGroup.ID)
			if err != nil {
				t.Errorf("Error deleting server group: %v", err)
			}
		}
	}()

	// create segment group for testing
	segmentGroupService := segmentgroup.New(zpaClient)
	segmentGroup, _, err := segmentGroupService.Create(&segmentgroup.SegmentGroup{
		Name:        name,
		Description: name,
		Enabled:     true,
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating segment for testing: %v", err)
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := segmentGroupService.Get(segmentGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := segmentGroupService.Delete(segmentGroup.ID)
			if err != nil {
				t.Errorf("Error deleting segment group: %v", err)
			}
		}
	}()

	// create segment group for testing
	appSegmentService := applicationsegment.New(zpaClient)
	appSegment, _, err := appSegmentService.Create(applicationsegment.ApplicationSegmentResource{
		Name:                  name,
		Description:           name,
		Enabled:               true,
		IpAnchored:            true,
		SegmentGroupID:        segmentGroup.ID,
		IsCnameEnabled:        true,
		BypassType:            "NEVER",
		IcmpAccessType:        "PING_TRACEROUTING",
		HealthReporting:       "ON_ACCESS",
		HealthCheckType:       "DEFAULT",
		TCPKeepAlive:          "1",
		InspectTrafficWithZia: false,
		DomainNames:           []string{"test.example.com"},
		ServerGroups: []applicationsegment.AppServerGroups{
			{
				ID: serverGroup.ID,
			},
		},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: rPort,
				To:   rPort,
			},
		},
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating application segment for testing: %v", err)
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := appSegmentService.Get(appSegment.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := appSegmentService.Delete(appSegment.ID)
			if err != nil {
				t.Errorf("Error deleting application segment: %v", err)
			}
		}
	}()

	// Initialize ZIA client for creating ZPA Gateways
	ziaClient, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating ZIA client: %v", err)
		return
	}
	service := New(ziaClient)

	// Initialize ZPAGateways resource
	zpaGateways := ZPAGateways{
		Name:        name,
		Description: description,
		Type:        "ZPA",
		ZPAServerGroup: ZPAServerGroup{
			ExternalID: serverGroup.ID, // Assigning int value
			Name:       serverGroup.Name,
		},
		ZPAAppSegments: []ZPAAppSegments{
			{
				ExternalID: appSegment.ID, // Assigning int value
				Name:       appSegment.Name,
			},
		},
	}

	// Inside TestZPAGateways function
	createdResource, err := service.Create(&zpaGateways)
	if err != nil {
		t.Fatalf("Error creating ZPAGateways resource: %v", err)
	}

	defer func() {
		// Attempt to delete the resource
		_, delErr := service.Delete(createdResource.ID)
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
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved zpa gateway '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Description = updateDescription
	// Ensure type is retained
	retrievedResource.Type = "ZPA" // or whatever the correct value should be

	// Remove metadata fields before sending update request
	retrievedResource.LastModifiedBy = nil
	retrievedResource.LastModifiedTime = 0

	//Convert the retrievedResource to JSON and log it before the update
	// var jsonRepresentation []byte
	// jsonRepresentation, err = json.MarshalIndent(retrievedResource, "", "  ")
	// if err != nil {
	// 	t.Fatalf("Error converting retrievedResource to JSON before update: %v", err)
	// }
	// t.Logf("JSON Payload being sent for update:\n%s", string(jsonRepresentation))

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
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Description != updateDescription {
		t.Errorf("Expected retrieved updated resource description '%s', but got '%s'", updateDescription, updatedResource.Description)
	}

	// Test resource retrieval by name
	retrievedResource, err = service.GetByName(name)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Description != updateDescription {
		t.Errorf("Expected retrieved resource comment '%s', but got '%s'", updateDescription, createdResource.Description)
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
func tryRetrieveResource(s *Service, id int) (*ZPAGateways, error) {
	var resource *ZPAGateways
	var err error

	for i := 0; i < maxRetries; i++ {
		resource, err = s.Get(id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}
