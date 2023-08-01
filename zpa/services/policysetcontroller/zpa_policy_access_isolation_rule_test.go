package policysetcontroller

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/cloudbrowserisolation/isolationprofile"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/samlattribute"
)

func TestAccessIsolationPolicy(t *testing.T) {
	policyType := "ISOLATION_POLICY"
	isolationProfileID := "BD_SA_Profile1"
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	idpService := idpcontroller.New(client)
	idpList, _, err := idpService.GetAll()
	if err != nil {
		t.Errorf("Error getting idps: %v", err)
		return
	}
	if len(idpList) == 0 {
		t.Error("Expected retrieved idps to be non-empty, but got empty slice")
	}
	samlService := samlattribute.New(client)
	samlsList, _, err := samlService.GetAll()
	if err != nil {
		t.Errorf("Error getting saml attributes: %v", err)
		return
	}
	if len(samlsList) == 0 {
		t.Error("Expected retrieved saml attributes to be non-empty, but got empty slice")
	}
	service := New(client)
	accessPolicySet, _, err := service.GetByPolicyType(policyType)
	if err != nil {
		t.Errorf("Error getting access isolation policy set: %v", err)
		return
	}
	profile := isolationprofile.New(client)
	profileID, _, err := profile.GetByName(isolationProfileID)
	if err != nil {
		t.Errorf("Error getting isolation profile id set: %v", err)
		return
	}
	accessPolicyRule := PolicyRule{
		Name:                  name,
		Description:           "New application segment",
		RuleOrder:             "1",
		PolicySetID:           accessPolicySet.ID,
		ZpnIsolationProfileID: profileID.ID,
		Action:                "ISOLATE",
		Conditions: []Conditions{
			{
				Operator: "OR",
				Operands: []Operands{
					{
						ObjectType: "CLIENT_TYPE",
						LHS:        "id",
						RHS:        "zpn_client_type_exporter",
					},
					{
						ObjectType: "SAML",
						LHS:        samlsList[0].ID,
						RHS:        "user1@acme.com",
						IdpID:      idpList[0].ID,
					},
				},
			},
		},
	}
	// Test resource creation
	createdResource, _, err := service.Create(&accessPolicyRule)
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
	retrievedResource, _, err := service.GetPolicyRule(accessPolicySet.ID, createdResource.ID)
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
	_, err = service.Update(accessPolicySet.ID, createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := service.GetPolicyRule(accessPolicySet.ID, createdResource.ID)
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
	retrievedResource, _, err = service.GetByNameAndType(policyType, updateName)
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
	resources, _, err := service.GetAllByType(policyType)
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
	_, err = service.Delete(accessPolicySet.ID, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = service.GetPolicyRule(accessPolicySet.ID, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}
*/
