package applicationsegment

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

func TestApplicationSegment(t *testing.T) {
	ctx := context.Background()
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	segmentGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rPort := strconv.Itoa(acctest.RandIntRange(1000, 9999))
	updatedPort := strconv.Itoa(acctest.RandIntRange(1000, 9999))
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	appGroup := segmentgroup.SegmentGroup{
		Name:        segmentGroupName,
		Description: segmentGroupName,
		Enabled:     true,
	}
	createdAppGroup, _, err := segmentgroup.Create(ctx, service, &appGroup)
	if err != nil {
		t.Errorf("Error creating application segment group: %v", err)
		return
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := segmentgroup.Get(ctx, service, createdAppGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := segmentgroup.Delete(ctx, service, createdAppGroup.ID)
			if err != nil {
				t.Errorf("Error deleting application segment group: %v", err)
			}
		}
	}()

	appSegment := ApplicationSegmentResource{
		Name:                  name,
		Description:           "New application segment",
		Enabled:               true,
		SegmentGroupID:        createdAppGroup.ID,
		SegmentGroupName:      createdAppGroup.Name,
		IsCnameEnabled:        true,
		BypassType:            "NEVER",
		IcmpAccessType:        "PING_TRACEROUTING",
		HealthReporting:       "ON_ACCESS",
		HealthCheckType:       "DEFAULT",
		TCPKeepAlive:          "1",
		InspectTrafficWithZia: false,
		MatchStyle:            "EXCLUSIVE",
		DomainNames:           []string{"test.example.com"},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: rPort,
				To:   rPort,
			},
		},
	}
	// Test resource creation
	createdResource, _, err := Create(ctx, service, appSegment)
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
	if len(createdResource.TCPPortRanges) != 2 || createdResource.TCPPortRanges[0] != rPort || createdResource.TCPPortRanges[1] != rPort {
		t.Errorf("Expected created resource port '%s-%s', but got '%s'", rPort, rPort, createdResource.TCPPortRanges)
	}
	// Test resource retrieval
	retrievedResource, _, err := Get(ctx, service, createdResource.ID)
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
	retrievedResource.TCPAppPortRange = []common.NetworkPorts{
		{
			From: updatedPort,
			To:   updatedPort,
		},
	}
	_, err = Update(ctx, service, createdResource.ID, *retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := Get(ctx, service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}
	if len(updatedResource.TCPPortRanges) != 2 || updatedResource.TCPPortRanges[0] != updatedPort || updatedResource.TCPPortRanges[1] != updatedPort {
		t.Errorf("Expected created resource port '%s-%s', but got '%s'", updatedPort, updatedPort, updatedResource.TCPPortRanges)
	}
	// Test resource retrieval by name
	retrievedResource, _, err = GetByName(ctx, service, updateName)
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
	resources, _, err := GetAll(ctx, service)
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
	_, err = Delete(ctx, service, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = Get(ctx, service, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}

func TestRetrieveNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, _, err = Get(context.Background(), service, "non-existent-id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = Delete(context.Background(), service, "non-existent-id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = Update(context.Background(), service, "non-existent-id", ApplicationSegmentResource{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, _, err = GetByName(context.Background(), service, "non-existent-name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
