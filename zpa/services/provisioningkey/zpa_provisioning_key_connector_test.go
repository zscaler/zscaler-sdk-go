package provisioningkey

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/enrollmentcert"
)

const (
	connGrpAssociationType          = "CONNECTOR_GRP"
	serviceEdgeGroupAssociationType = "SERVICE_EDGE_GRP"
)

func TestProvisiongKeyConnectorGroup(t *testing.T) {
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
	if err != nil || createdAppConnGroup == nil || createdAppConnGroup.ID == "" {
		t.Fatalf("Error creating application connector group or ID is empty")
		return
	}

	defer func() {
		if createdAppConnGroup != nil && createdAppConnGroup.ID != "" {
			existingGroup, _, errCheck := appConnectorGroupService.Get(createdAppConnGroup.ID)
			if errCheck == nil && existingGroup != nil {
				_, errDelete := appConnectorGroupService.Delete(createdAppConnGroup.ID)
				if errDelete != nil {
					t.Errorf("Error deleting application connector group: %v", errDelete)
				}
			}
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
	if err != nil || createdResource == nil || createdResource.ID == "" {
		t.Fatalf("Error making POST request or created resource is nil/empty: %v", err)
		return
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
		t.Fatalf("Error retrieving resource by name: %v", err)
		return
	}
	if retrievedResource == nil {
		t.Fatalf("Error: retrievedResource is nil")
		return
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, retrievedResource.Name)
	}

	// Test resources retrieval
	resources, err := service.GetAll()
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
		return
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

/*
// ######################################################################################################################################
// ############################################## TEST SERVICE EDGE GROUP PROVISIONING KEY ##############################################
// ######################################################################################################################################
func TestProvisiongKeyServiceEdgeGroup(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	edgeGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	// create service edge group for testing
	serviceEdgeGroupService := serviceedgegroup.New(client)
	edgeGroup := serviceedgegroup.ServiceEdgeGroup{
		Name:                   edgeGroupName,
		Description:            edgeGroupName,
		Enabled:                true,
		Latitude:               "37.3861",
		Longitude:              "-122.0839",
		Location:               "Mountain View, CA",
		IsPublic:               "TRUE",
		UpgradeDay:             "SUNDAY",
		UpgradeTimeInSecs:      "66600",
		OverrideVersionProfile: true,
		VersionProfileName:     "Default",
		VersionProfileID:       "0",
	}

	createdEdgeGroup, _, err := serviceEdgeGroupService.Create(edgeGroup)
	if err != nil || createdEdgeGroup == nil || createdEdgeGroup.ID == "" {
		t.Fatalf("Error creating service edge group or ID is empty")
		return
	}

	defer func() {
		if createdEdgeGroup != nil && createdEdgeGroup.ID != "" {
			existingGroup, _, errCheck := serviceEdgeGroupService.Get(createdEdgeGroup.ID)
			if errCheck == nil && existingGroup != nil {
				_, errDelete := serviceEdgeGroupService.Delete(createdEdgeGroup.ID)
				if errDelete != nil {
					t.Errorf("Error deleting service edge group: %v", errDelete)
				}
			}
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
		AssociationType:       serviceEdgeGroupAssociationType,
		Name:                  name,
		AppConnectorGroupID:   createdEdgeGroup.ID,
		AppConnectorGroupName: createdEdgeGroup.Name,
		EnrollmentCertID:      enrollmentCert.ID,
		ZcomponentID:          createdEdgeGroup.ID,
		MaxUsage:              "10",
	}
	// Test resource creation
	createdResource, _, err := service.Create(serviceEdgeGroupAssociationType, &resource)
	if err != nil || createdResource == nil || createdResource.ID == "" {
		t.Fatalf("Error making POST request or created resource is nil/empty: %v", err)
		return
	}

	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource retrieval
	retrievedResource, _, err := service.Get(serviceEdgeGroupAssociationType, createdResource.ID)
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
	_, err = service.Update(serviceEdgeGroupAssociationType, createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := service.Get(serviceEdgeGroupAssociationType, createdResource.ID)
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
	retrievedResource, _, err = service.GetByName(serviceEdgeGroupAssociationType, updateName)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
		return
	}
	if retrievedResource == nil {
		t.Fatalf("Error: retrievedResource is nil")
		return
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, retrievedResource.Name)
	}

	// Test resources retrieval
	resources, err := service.GetAll()
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
		return
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
	_, err = service.Delete(serviceEdgeGroupAssociationType, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = service.Get(serviceEdgeGroupAssociationType, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.Get(connGrpAssociationType, "non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}

	_, _, err = service.Get(serviceEdgeGroupAssociationType, "non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Delete(connGrpAssociationType, "non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}

	_, err = service.Delete(serviceEdgeGroupAssociationType, "non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Update(connGrpAssociationType, "non_existent_id", &ProvisioningKey{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
	_, err = service.Update(serviceEdgeGroupAssociationType, "non_existent_id", &ProvisioningKey{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.GetByName(connGrpAssociationType, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
	_, _, err = service.GetByName(serviceEdgeGroupAssociationType, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
*/
