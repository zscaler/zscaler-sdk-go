package applicationsegment

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/segmentgroup"
)

// clean all resources
func init() {
	log.Printf("init cleaning test")
	shouldCleanAllResources, _ := strconv.ParseBool(os.Getenv("ZSCALER_SDK_TEST_SWEEP"))
	if !shouldCleanAllResources {
		return
	}
	client, err := tests.NewZpaClient()
	if err != nil {
		panic(fmt.Sprintf("Error creating client: %v", err))
	}
	service := New(client)
	resources, _, _ := service.GetAll()
	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		_, _ = service.Delete(r.ID)
	}
}

func TestApplicationSegment(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	segmentGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rPort := strconv.Itoa(acctest.RandIntRange(1000, 9999))
	updatedPort := strconv.Itoa(acctest.RandIntRange(1000, 9999))
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	// create application segment group for testing
	appGroupService := segmentgroup.New(client)
	appGroup := segmentgroup.SegmentGroup{
		Name:        segmentGroupName,
		Description: segmentGroupName,
	}
	createdAppGroup, _, err := appGroupService.Create(&appGroup)
	if err != nil {
		t.Errorf("Error creating application segment group: %v", err)
		return
	}
	defer func() {
		_, err := appGroupService.Delete(createdAppGroup.ID)
		if err != nil {
			t.Errorf("Error deleting application segment group: %v", err)
		}
	}()

	service := New(client)

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
		DomainNames:           []string{"test.example.com"},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: rPort,
				To:   rPort,
			},
		},
	}
	// Test resource creation
	createdResource, _, err := service.Create(appSegment)
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
	retrievedResource.TCPAppPortRange = []common.NetworkPorts{
		{
			From: updatedPort,
			To:   updatedPort,
		},
	}
	_, err = service.Update(createdResource.ID, *retrievedResource)
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
	if len(updatedResource.TCPPortRanges) != 2 || updatedResource.TCPPortRanges[0] != updatedPort || updatedResource.TCPPortRanges[1] != updatedPort {
		t.Errorf("Expected created resource port '%s-%s', but got '%s'", updatedPort, updatedPort, updatedResource.TCPPortRanges)
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
	_, _, err = service.Get(createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}
