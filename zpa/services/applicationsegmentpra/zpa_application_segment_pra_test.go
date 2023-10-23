package applicationsegmentpra

/*
import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
)

// Declare the global variable
var createdResourceID string
var createdResourceName string

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

func setupTest(t *testing.T) *segmentgroup.SegmentGroup {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	appGroupService := segmentgroup.New(client)
	segmentGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	appGroup := segmentgroup.SegmentGroup{
		Name:        segmentGroupName,
		Description: segmentGroupName,
	}
	createdAppGroup, _, err := appGroupService.Create(&appGroup)
	if err != nil {
		t.Fatalf("Error creating application segment group: %v", err)
	}

	return createdAppGroup
}

func cleanupTest(t *testing.T, createdAppGroup *segmentgroup.SegmentGroup) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	appGroupService := segmentgroup.New(client)
	_, err = appGroupService.Delete(createdAppGroup.ID)
	if err != nil {
		t.Errorf("Error deleting application segment group: %v", err)
	}
}

func TestApplicationSegmentPRACreate(t *testing.T) {
	createdAppGroup := setupTest(t)
	// defer cleanupTest(t, createdAppGroup)

	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	initialPort := "3389"

	appSegment := AppSegmentPRA{
		Name:             name,
		Description:      name,
		Enabled:          true,
		SegmentGroupID:   createdAppGroup.ID,
		SegmentGroupName: createdAppGroup.Name,
		IsCnameEnabled:   true,
		BypassType:       "NEVER",
		IcmpAccessType:   "PING_TRACEROUTING",
		HealthReporting:  "ON_ACCESS",
		HealthCheckType:  "DEFAULT",
		TCPKeepAlive:     "1",
		DomainNames:      []string{"rdp_pra.example.com"},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: initialPort,
				To:   initialPort,
			},
		},
		CommonAppsDto: CommonAppsDto{
			AppsConfig: []AppsConfig{
				{
					Name:                name,
					Description:         name,
					Enabled:             true,
					AppTypes:            []string{"SECURE_REMOTE_ACCESS"},
					ApplicationPort:     initialPort,
					ApplicationProtocol: "RDP",
					ConnectionSecurity:  "ANY",
					Domain:              "rdp_pra.example.com",
				},
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
	if len(createdResource.TCPPortRanges) != 2 || createdResource.TCPPortRanges[0] != initialPort || createdResource.TCPPortRanges[1] != initialPort {
		t.Errorf("Expected created resource port '%s-%s', but got '%s'", initialPort, initialPort, createdResource.TCPPortRanges)
	}
	// Save the created resource's ID to the global variable
	createdResourceID = createdResource.ID
	createdResourceName = name
}

func TestApplicationSegmentPRAGet(t *testing.T) {
	createdAppGroup := setupTest(t)
	defer cleanupTest(t, createdAppGroup)

	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	resourceID := createdResourceID
	if resourceID == "" {
		t.Fatal("Resource ID is empty, might be due to TestApplicationSegmentPRACreate not executed or failed.")
	}

	retrievedResource, _, err := service.Get(resourceID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource == nil {
		t.Fatalf("Expected a resource but got nil")
	}
	if retrievedResource.ID != resourceID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", resourceID, retrievedResource.ID)
	}
}

func TestApplicationSegmentPRAGetByName(t *testing.T) {
	createdAppGroup := setupTest(t)
	defer cleanupTest(t, createdAppGroup)

	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	resourceName := createdResourceName
	if resourceName == "" {
		t.Fatal("Resource name is empty, might be due to TestApplicationSegmentPRACreate not executed or failed.")
	}

	retrievedResource, _, err := service.GetByName(resourceName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
		return
	}

	if retrievedResource == nil {
		t.Fatal("Retrieved resource is nil")
		return
	}

	if retrievedResource.Name != resourceName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", resourceName, retrievedResource.Name)
	}
}

/*
	func TestApplicationSegmentPRAUpdate(t *testing.T) {
		createdAppGroup := setupTest(t)
		defer cleanupTest(t, createdAppGroup)

		client, err := tests.NewZpaClient()
		if err != nil {
			t.Errorf("Error creating client: %v", err)
			return
		}

		service := New(client)
		resourceID := createdResourceID
		if resourceID == "" {
			t.Fatal("Resource ID is empty, might be due to TestApplicationSegmentPRACreate not executed or failed.")
		}
		updatedPort := "3389"

		retrievedResource, _, err := service.Get(resourceID)
		if err != nil {
			t.Fatalf("Error retrieving resource: %v", err)
		}

		initialAppPort := ""
		if len(retrievedResource.SRAAppsDto) > 0 {
			initialAppPort = retrievedResource.SRAAppsDto[0].ApplicationPort
		}

		// Check if there's a change in ApplicationPort
		if initialAppPort != updatedPort {
			// Delete the old resource
			_, err = service.Delete(resourceID)
			if err != nil {
				t.Fatalf("Error deleting resource: %v", err)
			}

			// Create new resource with the updated configurations
			// NOTE: This assumes that the service.Create function returns the created resource
			// Create new resource with the updated configurations
			newResource, _, err := service.Create(*retrievedResource)
			if err != nil {
				t.Fatalf("Error creating new resource: %v", err)
			}
			retrievedResource = newResource   // Updating the reference to point to the newly created resource
			resourceID = retrievedResource.ID // Assuming the newResource (or retrievedResource after the assignment) has an ID field

		} else {
			// Your existing update logic for properties that can be updated in place.
			retrievedResource.SegmentGroupID = createdAppGroup.ID
			retrievedResource.SegmentGroupName = createdAppGroup.Name
			retrievedResource.Name = createdResourceName
			retrievedResource.TCPAppPortRange = []common.NetworkPorts{
				{
					From: updatedPort,
					To:   updatedPort,
				},
			}

			if len(retrievedResource.SRAAppsDto) > 0 {
				retrievedResource.SRAAppsDto[0].ApplicationPort = updatedPort
			} else {
				retrievedResource.CommonAppsDto.AppsConfig = []AppsConfig{
					{
						Name:                createdResourceName,
						Description:         createdResourceName,
						Enabled:             true,
						AppTypes:            []string{"SECURE_REMOTE_ACCESS"},
						ApplicationPort:     updatedPort,
						ApplicationProtocol: "RDP",
						ConnectionSecurity:  "ANY",
						Domain:              "rdp_pra.example.com",
					},
				}
			}

			// Update the resource
			_, err = service.Update(resourceID, retrievedResource)
			if err != nil {
				t.Errorf("Error updating resource: %v", err)
			}
		}

		// Delay to give some time for the update to propagate (if needed)
		time.Sleep(time.Second * 5)

		// Fetch the updated resource again
		updatedResource, _, err := service.Get(resourceID)
		if err != nil {
			t.Fatalf("Error retrieving updated resource: %v", err)
		}

		// Assertions based on your requirements
		if updatedResource.TCPAppPortRange[0].From != updatedPort || updatedResource.TCPAppPortRange[0].To != updatedPort {
			t.Errorf("Expected updated resource port '%s-%s', but got '%s-%s'", updatedPort, updatedPort, updatedResource.TCPAppPortRange[0].From, updatedResource.TCPAppPortRange[0].To)
		}

		if len(updatedResource.CommonAppsDto.AppsConfig) > 0 && updatedResource.CommonAppsDto.AppsConfig[0].ApplicationPort != updatedPort {
			t.Errorf("Expected updated ApplicationPort '%s', but got '%s'", updatedPort, updatedResource.CommonAppsDto.AppsConfig[0].ApplicationPort)
		}
	}

func TestApplicationSegmentPRADelete(t *testing.T) {
	createdAppGroup := setupTest(t)
	defer cleanupTest(t, createdAppGroup)

	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	// Read the ID from the global variable
	resourceID := createdResourceID
	if resourceID == "" {
		t.Fatal("Resource ID is empty, might be due to TestApplicationSegmentPRACreate not executed or failed.")
	}

	// Delete the resource
	_, err = service.Delete(resourceID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Try fetching the deleted resource to ensure it's deleted
	_, _, err = service.Get(resourceID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}
*/
