package forwarding_rules

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/forwarding_control_policy/zpa_gateways"
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

func TestForwardingRules(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
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

	ziaClient, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// create segment group for testing
	zpaGatewayService := zpa_gateways.New(ziaClient)
	zpaGateway, err := zpaGatewayService.Create(&zpa_gateways.ZPAGateways{
		Name:        name,
		Description: name,
		Type:        "ZPA",
		ZPAServerGroup: zpa_gateways.ZPAServerGroup{
			ExternalID: serverGroup.ID,
			Name:       serverGroup.Name,
		},
		ZPAAppSegments: []zpa_gateways.ZPAAppSegments{
			{
				ExternalID: appSegment.ID,
				Name:       appSegment.Name,
			},
		},
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating zpa gateway for testing forwarding rule: %v", err)
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, getErr := zpaGatewayService.Get(zpaGateway.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := zpaGatewayService.Delete(zpaGateway.ID)
			if err != nil {
				t.Errorf("Error deleting zpa gateway: %v", err)
			}
		}
	}()

	service := New(ziaClient)
	// create ForwardingRule
	rule := ForwardingRules{
		Name:          name,
		Description:   name,
		Order:         1,
		Rank:          7,
		State:         "ENABLED",
		Type:          "FORWARDING",
		ForwardMethod: "ZPA",
		ZPAGateway: ZPAGateway{
			ID:   zpaGateway.ID,
			Name: zpaGateway.Name,
		},
		ZPAAppSegments: []ZPAAppSegments{
			{
				ExternalID: appSegment.ID,
				Name:       appSegment.Name,
			},
		},
	}

	// Inside Forwarding Control Rule function
	createdResource, err := service.Create(&rule)
	if err != nil {
		t.Fatalf("Error creating Forwarding Control Rule resource: %v", err)
	}

	defer func() {
		// Attempt to delete the resource
		_, delErr := service.Delete(createdResource.ID)
		if delErr != nil {
			// If the error indicates the resource is already deleted, log it as information
			if strings.Contains(delErr.Error(), "409") || strings.Contains(delErr.Error(), "RESOURCE_NOT_FOUND") {
				t.Logf("Resource with ID %d not found (already deleted).", createdResource.ID)
			} else {
				// If the deletion error is not due to the resource being missing, log it as an actual error
				t.Errorf("Error deleting Forwarding Control Rule resource: %v", delErr)
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
		t.Errorf("Expected retrieved dlp engine '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Name = updateName
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
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedByNameResource, err := service.GetByName(updateName)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedByNameResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedByNameResource.ID)
	}
	if retrievedByNameResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, retrievedByNameResource.Name)
	}

	// Test resources retrieval
	allResources, err := service.GetAll()
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
	time.Sleep(5 * time.Second) // sleep for 5 seconds

	// Test resource removal
	err = retryOnConflict(func() error {
		_, getErr := service.Get(createdResource.ID)
		if getErr != nil {
			return fmt.Errorf("Resource %d may have already been deleted: %v", createdResource.ID, getErr)
		}
		_, delErr := service.Delete(createdResource.ID)
		return delErr
	})
	_, err = service.Get(createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(s *Service, id int) (*ForwardingRules, error) {
	var resource *ForwardingRules
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
