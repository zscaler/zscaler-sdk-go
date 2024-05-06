package policysetcontroller

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/serviceedgegroup"
)

func TestAccessRedirectionPolicy(t *testing.T) {
	policyType := "REDIRECTION_POLICY"
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "updated_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// create service edge group for testing
	svcEdgeGroupService := serviceedgegroup.New(client)
	svcEdgeGroup, _, err := svcEdgeGroupService.Create(serviceedgegroup.ServiceEdgeGroup{
		Name:                   name,
		Description:            name,
		Enabled:                true,
		Latitude:               "37.33874",
		Longitude:              "-121.8852525",
		Location:               "San Jose, CA, USA",
		IsPublic:               "TRUE",
		UpgradeDay:             "SUNDAY",
		UpgradeTimeInSecs:      "66600",
		OverrideVersionProfile: true,
		VersionProfileName:     "Default",
		VersionProfileID:       "0",
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating service edge group for testing server group: %v", err)
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := svcEdgeGroupService.Get(svcEdgeGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := svcEdgeGroupService.Delete(svcEdgeGroup.ID)
			if err != nil {
				t.Errorf("Error deleting service edge group: %v", err)
			}
		}
	}()

	service := New(client)
	accessPolicySet, _, err := service.GetByPolicyType(policyType)
	if err != nil {
		t.Errorf("Error getting redirection access policy set: %v", err)
		return
	}
	redirectionPolicyRule := PolicyRule{
		Name:        name,
		Description: name,
		PolicySetID: accessPolicySet.ID,
		Action:      "REDIRECT_PREFERRED",
		ServiceEdgeGroups: []ServiceEdgeGroups{
			{
				ID: svcEdgeGroup.ID,
			},
		},
		Conditions: []Conditions{
			{
				Operator: "OR",
				Operands: []Operands{
					{
						ObjectType: "CLIENT_TYPE",
						LHS:        "id",
						RHS:        "zpn_client_type_machine_tunnel",
					},
					{
						ObjectType: "CLIENT_TYPE",
						LHS:        "id",
						RHS:        "zpn_client_type_branch_connector",
					},
					{
						ObjectType: "CLIENT_TYPE",
						LHS:        "id",
						RHS:        "zpn_client_type_edge_connector",
					},
					{
						ObjectType: "CLIENT_TYPE",
						LHS:        "id",
						RHS:        "zpn_client_type_zapp",
					},
				},
			},
		},
	}
	// Test resource creation
	createdResource, _, err := service.CreateRule(&redirectionPolicyRule)
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
	_, err = service.UpdateRule(accessPolicySet.ID, createdResource.ID, retrievedResource)
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
