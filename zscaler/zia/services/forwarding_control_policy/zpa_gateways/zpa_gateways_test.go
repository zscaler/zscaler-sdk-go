package zpa_gateways

/*
import (
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
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

func TestZPAGateways(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	description := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateDescription := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rPort := strconv.Itoa(acctest.RandIntRange(1000, 9999))

	zpaClient, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	// create app connector group for testing
	service := services.New(zpaClient)
	appConnGroup, _, err := appconnectorgroup.Create(context.Background(), service, appconnectorgroup.AppConnectorGroup{
		Name:                     name,
		Description:              name,
		Enabled:                  true,
		CityCountry:              "San Jose, US",
		Latitude:                 "37.33874",
		Longitude:                "-121.8852525",
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
		_, _, getErr := appconnectorgroup.Get(service, appConnGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := appconnectorgroup.Delete(context.Background(), service, appConnGroup.ID)
			if err != nil {
				t.Errorf("Error deleting app connector group: %v", err)
			}
		}
	}()

	// create app server for testing
	serverGroup, _, err := servergroup.Create(context.Background(), service, &servergroup.ServerGroup{
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
		_, _, getErr := servergroup.Get(service, serverGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := servergroup.Delete(context.Background(), service, serverGroup.ID)
			if err != nil {
				t.Errorf("Error deleting server group: %v", err)
			}
		}
	}()

	// create segment group for testing
	segmentGroup, _, err := segmentgroup.Create(context.Background(), service, &segmentgroup.SegmentGroup{
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
		_, _, getErr := segmentgroup.Get(service, segmentGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := segmentgroup.Delete(context.Background(), service, segmentGroup.ID)
			if err != nil {
				t.Errorf("Error deleting segment group: %v", err)
			}
		}
	}()

	// create segment group for testing
	appSegment, _, err := applicationsegment.Create(context.Background(), service, applicationsegment.ApplicationSegmentResource{
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
		_, _, getErr := applicationsegment.Get(service, appSegment.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := applicationsegment.Delete(context.Background(), service, appSegment.ID)
			if err != nil {
				t.Errorf("Error deleting application segment: %v", err)
			}
		}
	}()

	ziaService, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

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
	createdResource, err := ziaService.Create(&zpaGateways)
	if err != nil {
		t.Fatalf("Error creating ZPAGateways resource: %v", err)
	}

	defer func() {
		// Attempt to delete the resource
		_, delErr := ziaService.Delete(createdResource.ID)
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
	retrievedResource, err := tryRetrieveResource(ziaService, createdResource.ID)
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

	// Convert the retrievedResource to JSON and log it before the update
	// var jsonRepresentation []byte
	// jsonRepresentation, err = json.MarshalIndent(retrievedResource, "", "  ")
	// if err != nil {
	// 	t.Fatalf("Error converting retrievedResource to JSON before update: %v", err)
	// }
	// t.Logf("JSON Payload being sent for update:\n%s", string(jsonRepresentation))

	err = retryOnConflict(func() error {
		_, err = ziaService.Update(createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := ziaService.Get(createdResource.ID)
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
	retrievedResource, err = ziaService.GetByName(name)
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
	resources, err := ziaService.GetAll()
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
		_, delErr := ziaService.Delete(createdResource.ID)
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

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Get(0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Delete(0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Update(0, &ZPAGateways{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.GetByName("non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}

*/
