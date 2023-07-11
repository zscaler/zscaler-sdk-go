package locationgroups

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
)

func TestLocationGroups_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	// static ip for vpn credentials testing

	service := New(client)

	// Test resources retrieval
	resources, err := service.GetAll()
	if err != nil {
		t.Errorf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		return
	}
	name := resources[0].Name
	resourceByName, err := service.GetLocationGroupByName(name)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	// Test resource retrieval after deletion
	_, err = service.GetLocationGroup(resourceByName.ID)
	if err != nil {
		t.Errorf("expected resource to exist: %v", err)
	}
}

/*
func TestLocationGroups(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("1.1.1.1/24")
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	locationName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	// static ip for vpn credentials testing
	staticipsService := staticips.New(client)
	// Test resource creation
	staticIP, _, err := staticipsService.Create(&staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   "testing static ip for location management",
	})
	if err != nil {
		t.Errorf("creating static ip failed: %v", err)
		return
	}
	defer func() {
		_, err := staticipsService.Delete(staticIP.ID)
		if err != nil {
			t.Errorf("deleting static ip failed: %v", err)
		}
	}()

	locationmanagementService := locationmanagement.New(client)

	// Test resource creation
	location, err := locationmanagementService.Create(&locationmanagement.Locations{
		Name:              locationName,
		Description:       locationName,
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
	})

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}
	defer func() {
		_, err := locationmanagementService.Delete(location.ID)
		if err != nil {
			t.Errorf("deleting static ip failed: %v", err)
		}
	}()

	service := locationgroups.New(client)

	locationGroup := locationgroups.LocationGroup{
		Name:      name,
		GroupType: "Static",
		Locations: []common.IDNameExtensions{
			{
				ID: location.ID,
			},
		},
	}

	// Test resource creation
	createdResource, err := service.CreateLocationGroup(&locationGroup)

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource retrieval
	retrievedResource, err := service.GetLocationGroup(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource update
	retrievedResource.Name = updateName
	_, _, err = service.UpdateLocationGroup(createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, err := service.GetLocationGroup(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resources retrieval
	resources, err := service.GetAll()
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
		t.Errorf("Expected retrieved resources to contain created resource '%d', but it didn't", createdResource.ID)
	}
	// Test resource removal
	_, err = service.DeleteLocationGroup(createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, err = service.GetLocationGroup(createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}
*/
