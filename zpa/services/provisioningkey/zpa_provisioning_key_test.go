package provisioningkey

import (
	"fmt"
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
	// Uncomment setup if needed in the future
	// setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func teardown() {
	err := cleanResources()
	if err != nil {
		log.Fatalf("Error during cleanup: %v", err)
	}
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

func cleanResources() error {
	if !shouldClean() {
		return nil
	}

	client, err := tests.NewZpaClient()
	if err != nil {
		return fmt.Errorf("Error creating client: %v", err)
	}
	service := New(client)
	resources, err := service.GetAll()
	if err != nil {
		return fmt.Errorf("Error getting all resources: %v", err)
	}
	for _, r := range resources {
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)

		var associationType string
		if r.AssociationType == connGrpAssociationType {
			associationType = connGrpAssociationType
		} else if r.AssociationType == serviceEdgeGrpAssociationType {
			associationType = serviceEdgeGrpAssociationType
		} else {
			continue
		}

		_, err := service.Delete(associationType, r.ID)
		if err != nil {
			return fmt.Errorf("Error deleting resource with ID: %s, Name: %s, Error: %v", r.ID, r.Name, err)
		}
	}
	return nil
}

func TestAppConnectorGroupProvisiongKey(t *testing.T) {
	randStr := func() string { return "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) }
	assertNoError := func(err error, errMsg string) {
		if err != nil {
			t.Fatalf(errMsg, err)
			return
		}
	}

	name := randStr()
	updateName := randStr()
	appConnGroupName := randStr()

	client, err := tests.NewZpaClient()
	assertNoError(err, "Error creating client: %v")

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
	assertNoError(err, "Error creating application connector group: %v")
	defer func() {
		_, err := appConnectorGroupService.Delete(createdAppConnGroup.ID)
		if err != nil {
			t.Errorf("Error deleting application connector group: %v", err)
		}
	}()
	// get enrollment cert for testing
	enrollmentCertService := enrollmentcert.New(client)
	enrollmentCert, _, err := enrollmentCertService.GetByName("Connector")
	assertNoError(err, "Error getting enrollment cert: %v")

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
	createdResource, _, err := service.Create(connGrpAssociationType, &resource)
	assertNoError(err, "Error making POST request: %v")

	assertNoError(checkResourceFields(*createdResource, name), "")

	retrievedResource, _, err := service.Get(connGrpAssociationType, createdResource.ID)
	assertNoError(err, "Error retrieving resource: %v")
	assertNoError(checkResourceFields(*retrievedResource, name), "")

	// Test resource update
	retrievedResource.Name = updateName
	_, err = service.Update(connGrpAssociationType, createdResource.ID, retrievedResource)
	assertNoError(err, "Error updating resource: %v")

	updatedResource, _, err := service.Get(connGrpAssociationType, createdResource.ID)
	assertNoError(err, "Error retrieving resource: %v")
	assertNoError(checkResourceFields(*updatedResource, updateName), "")

	retrievedResource, _, err = service.GetByName(connGrpAssociationType, updateName)
	assertNoError(err, "Error retrieving resource by name: %v")
	assertNoError(checkResourceFields(*retrievedResource, updateName), "")

	resources, err := service.GetAll()
	assertNoError(err, "Error retrieving resources: %v")
	if len(resources) == 0 {
		t.Error("Expected retrieved resources to be non-empty, but got empty slice")
	}

	if !resourceInSlice(*createdResource, resources) {
		t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.ID)
	}

	_, err = service.Delete(connGrpAssociationType, createdResource.ID)
	assertNoError(err, "Error deleting resource: %v")

	_, _, err = service.Get(connGrpAssociationType, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}

func TestServiceEdgeGroupProvisiongKey(t *testing.T) {
	randStr := func() string { return "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha) }
	assertNoError := func(err error, errMsg string) {
		if err != nil {
			t.Fatalf(errMsg, err)
			return
		}
	}

	name := randStr()
	updateName := randStr()
	serviceEdgeGroupName := randStr()

	client, err := tests.NewZpaClient()
	assertNoError(err, "Error creating client: %v")
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
	assertNoError(err, "Error creating service edge group: %v")
	defer func() {
		_, err := serviceEdgeGroupService.Delete(createdServiceEdgeGroup.ID)
		if err != nil {
			t.Errorf("Error deleting service edge group: %v", err)
		}
	}()
	// get enrollment cert for testing
	enrollmentCertService := enrollmentcert.New(client)
	enrollmentCert, _, err := enrollmentCertService.GetByName("Service Edge")
	assertNoError(err, "Error getting enrollment cert: %v")

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
	createdResource, _, err := service.Create(serviceEdgeGrpAssociationType, &resource)
	assertNoError(err, "Error making POST request: %v")

	assertNoError(checkResourceFields(*createdResource, name), "")

	retrievedResource, _, err := service.Get(serviceEdgeGrpAssociationType, createdResource.ID)
	assertNoError(err, "Error retrieving resource: %v")
	assertNoError(checkResourceFields(*retrievedResource, name), "")

	// Test resource update
	retrievedResource.Name = updateName
	_, err = service.Update(serviceEdgeGrpAssociationType, createdResource.ID, retrievedResource)
	assertNoError(err, "Error updating resource: %v")

	updatedResource, _, err := service.Get(serviceEdgeGrpAssociationType, createdResource.ID)
	assertNoError(err, "Error retrieving resource: %v")
	assertNoError(checkResourceFields(*updatedResource, updateName), "")

	retrievedResource, _, err = service.GetByName(serviceEdgeGrpAssociationType, updateName)
	assertNoError(err, "Error retrieving resource by name: %v")
	assertNoError(checkResourceFields(*retrievedResource, updateName), "")

	resources, err := service.GetAll()
	assertNoError(err, "Error retrieving resources: %v")
	if len(resources) == 0 {
		t.Error("Expected retrieved resources to be non-empty, but got empty slice")
	}

	if !resourceInSlice(*createdResource, resources) {
		t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.ID)
	}

	_, err = service.Delete(serviceEdgeGrpAssociationType, createdResource.ID)
	assertNoError(err, "Error deleting resource: %v")

	_, _, err = service.Get(serviceEdgeGrpAssociationType, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}

func checkResourceFields(resource ProvisioningKey, expectedName string) error {
	if resource.ID == "" {
		return fmt.Errorf("Expected created resource ID to be non-empty, but got ''")
	}
	if resource.Name != expectedName {
		return fmt.Errorf("Expected created resource name '%s', but got '%s'", expectedName, resource.Name)
	}
	return nil
}

func resourceInSlice(target ProvisioningKey, resourceList []ProvisioningKey) bool {
	for _, resource := range resourceList {
		if resource.ID == target.ID {
			return true
		}
	}
	return false
}
