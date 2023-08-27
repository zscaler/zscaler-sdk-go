package provisioningkey

import (
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/enrollmentcert"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/serviceedgegroup"
)

const (
	connGrpAssociationType        = "CONNECTOR_GRP"
	serviceEdgeGrpAssociationType = "SERVICE_EDGE_GRP"
)

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

// func setup() {
// 	cleanResources() // clean up at the beginning
// }

func teardown() {
	cleanResources() // clean up at the end
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	if !present {
		return true // default value
	}
	shouldClean, err := strconv.ParseBool(val)
	if err != nil {
		return true // default to cleaning if the value is not parseable
	}
	log.Printf("ZSCALER_SDK_TEST_SWEEP value: %v", shouldClean)
	return shouldClean
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
	resources, _ := service.GetAll()
	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)

		// Determine which type to use based on the AssociationType of the resource
		var associationType string
		if r.AssociationType == connGrpAssociationType {
			associationType = connGrpAssociationType
		} else if r.AssociationType == serviceEdgeGrpAssociationType {
			associationType = serviceEdgeGrpAssociationType
		} else {
			// If the AssociationType is neither of the two we expect, skip to the next resource
			continue
		}

		_, _ = service.Delete(associationType, r.ID)
	}
}

func TestAppConnectorGroupProvisiongKey(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	appConnGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	// create application connector group for testing
	appConnectorGroupService := appconnectorgroup.New(client)
	appGroup := appconnectorgroup.AppConnectorGroup{
		Name:                     appConnGroupName,
		Description:              appConnGroupName,
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
	}
	createdAppConnGroup, _, err := appConnectorGroupService.Create(appGroup)
	if err != nil {
		t.Errorf("Error creating application connector group: %v", err)
		return
	}
	defer func() {
		_, err := appConnectorGroupService.Delete(createdAppConnGroup.ID)
		if err != nil {
			t.Errorf("Error deleting application connector group: %v", err)
		}
	}()
	// get enrollment cert for testing
	enrollmentCertService := enrollmentcert.New(client)
	enrollmentCert, _, err := enrollmentCertService.GetByName("Connector")
	if err != nil {
		t.Errorf("Error getting enrollment cert: %v", err)
		return
	}

	service := New(client)

	resource := ProvisioningKey{
		AssociationType:       connGrpAssociationType,
		Name:                  name,
		AppConnectorGroupID:   createdAppConnGroup.ID,
		AppConnectorGroupName: createdAppConnGroup.Name,
		EnrollmentCertID:      enrollmentCert.ID,
		ZcomponentID:          createdAppConnGroup.ID,
		MaxUsage:              "10",
	}
	// Test resource creation
	createdResource, _, err := service.Create(connGrpAssociationType, &resource)
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
	retrievedResource, _, err := service.Get(connGrpAssociationType, createdResource.ID)
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
	_, err = service.Update(connGrpAssociationType, createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := service.Get(connGrpAssociationType, createdResource.ID)
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
	retrievedResource, _, err = service.GetByName(connGrpAssociationType, updateName)
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
		t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.ID)
	}
	// Test resource removal
	_, err = service.Delete(connGrpAssociationType, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = service.Get(connGrpAssociationType, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}

func TestServiceEdgeGroupProvisiongKey(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	serviceEdgeGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	// create application connector group for testing
	serviceEdgeGroupService := serviceedgegroup.New(client)
	edgeGroup := serviceedgegroup.ServiceEdgeGroup{
		Name:                   serviceEdgeGroupName,
		Description:            serviceEdgeGroupName,
		Enabled:                true,
		CityCountry:            "San Jose, US",
		Latitude:               "37.3382082",
		Longitude:              "-121.8863286",
		Location:               "San Jose, CA, USA",
		UpgradeDay:             "SUNDAY",
		UpgradeTimeInSecs:      "66600",
		OverrideVersionProfile: true,
		VersionProfileName:     "Default",
		VersionProfileID:       "0",
	}
	createdServiceEdgeGroup, _, err := serviceEdgeGroupService.Create(edgeGroup)
	if err != nil {
		t.Errorf("Error creating service edge group: %v", err)
		return
	}
	defer func() {
		_, err := serviceEdgeGroupService.Delete(createdServiceEdgeGroup.ID)
		if err != nil {
			t.Errorf("Error deleting service edge group: %v", err)
		}
	}()
	// get enrollment cert for testing
	enrollmentCertService := enrollmentcert.New(client)
	enrollmentCert, _, err := enrollmentCertService.GetByName("Service Edge")
	if err != nil {
		t.Errorf("Error getting enrollment cert: %v", err)
		return
	}

	service := New(client)

	resource := ProvisioningKey{
		AssociationType:       serviceEdgeGrpAssociationType,
		Name:                  name,
		AppConnectorGroupID:   createdServiceEdgeGroup.ID,
		AppConnectorGroupName: createdServiceEdgeGroup.Name,
		EnrollmentCertID:      enrollmentCert.ID,
		ZcomponentID:          createdServiceEdgeGroup.ID,
		MaxUsage:              "10",
	}
	// Test resource creation
	createdResource, _, err := service.Create(serviceEdgeGrpAssociationType, &resource)
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
	retrievedResource, _, err := service.Get(serviceEdgeGrpAssociationType, createdResource.ID)
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
	_, err = service.Update(serviceEdgeGrpAssociationType, createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := service.Get(serviceEdgeGrpAssociationType, createdResource.ID)
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
	retrievedResource, _, err = service.GetByName(serviceEdgeGrpAssociationType, updateName)
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
		t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.ID)
	}
	// Test resource removal
	_, err = service.Delete(serviceEdgeGrpAssociationType, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = service.Get(serviceEdgeGrpAssociationType, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}
