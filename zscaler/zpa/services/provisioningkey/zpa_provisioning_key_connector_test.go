package provisioningkey

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/enrollmentcert"
)

const (
	connGrpAssociationType = "CONNECTOR_GRP"
)

func TestProvisioningKeyConnectorGroup(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	appConnGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	// service, err := tests.NewOneAPIClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	service, err := tests.NewZPAClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	appGroup := appconnectorgroup.AppConnectorGroup{
		Name:                     appConnGroupName,
		Description:              appConnGroupName,
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
	}

	createdAppConnGroup, _, err := appconnectorgroup.Create(context.Background(), service, appGroup)
	if err != nil || createdAppConnGroup == nil || createdAppConnGroup.ID == "" {
		t.Fatalf("Error creating application connector group or ID is empty")
		return
	}

	defer func() {
		if createdAppConnGroup != nil && createdAppConnGroup.ID != "" {
			existingGroup, _, errCheck := appconnectorgroup.Get(context.Background(), service, createdAppConnGroup.ID)
			if errCheck == nil && existingGroup != nil {
				_, errDelete := appconnectorgroup.Delete(context.Background(), service, createdAppConnGroup.ID)
				if errDelete != nil {
					t.Errorf("Error deleting application connector group: %v", errDelete)
				}
			}
		}
	}()

	enrollmentCert, _, err := enrollmentcert.GetByName(context.Background(), service, "Connector")
	if err != nil {
		t.Errorf("Error getting enrollment cert: %v", err)
		return
	}

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
	createdResource, _, err := Create(context.Background(), service, connGrpAssociationType, &resource)
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
	retrievedResource, _, err := Get(context.Background(), service, connGrpAssociationType, createdResource.ID)
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
	_, err = Update(context.Background(), service, connGrpAssociationType, createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := Get(context.Background(), service, connGrpAssociationType, createdResource.ID)
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
	retrievedResource, _, err = GetByName(context.Background(), service, connGrpAssociationType, updateName)
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
	resources, err := GetAll(context.Background(), service)
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

	// Additional Tests for missing functions

	// Test GetByNameAllAssociations
	retrievedResource, assocType, _, err := GetByNameAllAssociations(context.Background(), service, updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name across all associations: %v", err)
	}
	if retrievedResource == nil {
		t.Fatalf("Expected retrieved resource, but got nil")
	}
	if assocType != connGrpAssociationType {
		t.Errorf("Expected association type '%s', but got '%s'", connGrpAssociationType, assocType)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}

	// Test GetByIDAllAssociations
	retrievedResource, assocType, _, err = GetByIDAllAssociations(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource by ID across all associations: %v", err)
	}
	if retrievedResource == nil {
		t.Fatalf("Expected retrieved resource, but got nil")
	}
	if assocType != connGrpAssociationType {
		t.Errorf("Expected association type '%s', but got '%s'", connGrpAssociationType, assocType)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}

	// Test GetAllByAssociationType
	associationTypeResources, err := GetAllByAssociationType(context.Background(), service, connGrpAssociationType)
	if err != nil {
		t.Errorf("Error retrieving resources by association type: %v", err)
	}
	if len(associationTypeResources) == 0 {
		t.Error("Expected non-empty list of resources by association type, but got empty list")
	}
	// check if the created resource is in the list
	found = false
	for _, resource := range associationTypeResources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved resources by association type to contain created resource '%s', but it didn't", createdResource.ID)
	}

	// Test resource removal
	_, err = Delete(context.Background(), service, connGrpAssociationType, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = Get(context.Background(), service, connGrpAssociationType, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}
