package policysetcontrollerv2

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
)

func TestAccessRedirectionPolicyV2(t *testing.T) {
	policyType := "REDIRECTION_POLICY"
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	// create service edge group for testing
	svcEdgeGroup, _, err := serviceedgegroup.Create(context.Background(), service, serviceedgegroup.ServiceEdgeGroup{
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
		_, _, getErr := serviceedgegroup.Get(context.Background(), service, svcEdgeGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := serviceedgegroup.Delete(context.Background(), service, svcEdgeGroup.ID)
			if err != nil {
				t.Errorf("Error deleting service edge group: %v", err)
			}
		}
	}()

	accessPolicySet, _, err := GetByPolicyType(context.Background(), service, policyType)
	if err != nil {
		t.Errorf("Error getting access policy set: %v", err)
		return
	}

	var ruleIDs []string // Store the IDs of the created rules

	for i := 0; i < 1; i++ {
		// Generate a unique name for each iteration
		name := fmt.Sprintf("tests-%s-%d", acctest.RandStringFromCharSet(10, acctest.CharSetAlpha), i)

		accessPolicyRule := PolicyRule{
			Name:        name,
			Description: name,
			PolicySetID: accessPolicySet.ID,
			Action:      "REDIRECT_PREFERRED",
			ServiceEdgeGroups: []serviceedgegroup.ServiceEdgeGroup{
				{
					ID:   svcEdgeGroup.ID,
					Name: svcEdgeGroup.Name,
				},
			},
			Conditions: []PolicyRuleResourceConditions{
				{
					Operator: "OR",
					Operands: []PolicyRuleResourceOperands{
						{
							ObjectType:        "COUNTRY_CODE",
							EntryValuesLHSRHS: []OperandsResourceLHSRHSValue{{LHS: "CA", RHS: "true"}, {LHS: "US", RHS: "true"}},
						},
					},
				},
				{
					Operator: "OR",
					Operands: []PolicyRuleResourceOperands{
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_machine_tunnel", "zpn_client_type_branch_connector", "zpn_client_type_edge_connector", "zpn_client_type_zapp"},
						},
					},
				},
			},
		}

		// Test resource creation
		createdResource, _, err := CreateRule(context.Background(), service, &accessPolicyRule)

		if err != nil {
			t.Errorf("Error making POST request: %v", err)
			continue
		}
		if createdResource.ID == "" {
			t.Error("Expected created resource ID to be non-empty, but got ''")
			continue
		}
		// if err == nil {
		// 	jsonBytes, _ := json.Marshal(createdResource)
		// 	fmt.Println(string(jsonBytes)) // This prints the JSON response
		// }
		ruleIDs = append(ruleIDs, createdResource.ID) // Collect rule ID for reordering

		// Update the rule name
		updatedName := name + "-updated"
		accessPolicyRule.Name = updatedName
		_, updateErr := UpdateRule(context.Background(), service, accessPolicySet.ID, createdResource.ID, &accessPolicyRule)

		if updateErr != nil {
			t.Errorf("Error updating rule: %v", updateErr)
			continue
		}

		// Retrieve and print the updated resource as JSON
		updatedResource, _, getErr := GetPolicyRule(context.Background(), service, accessPolicySet.ID, createdResource.ID)
		if getErr != nil {
			t.Errorf("Error retrieving updated resource: %v", getErr)
			continue
		}
		if updatedResource.Name != updatedName {
			t.Errorf("Expected updated resource name '%s', but got '%s'", updatedName, updatedResource.Name)
		}

		// Test resource retrieval by name
		updatedResource, _, err = GetByNameAndType(context.Background(), service, policyType, updatedName)
		if err != nil {
			t.Errorf("Error retrieving resource by name: %v", err)
		}
		if updatedResource.ID != createdResource.ID {
			t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
		}
		if updatedResource.Name != updatedName {
			t.Errorf("Expected retrieved resource name '%s', but got '%s'", updatedName, updatedResource.Name)
		}
		time.Sleep(5 * time.Second)
	}

	// Reorder the rules after all have been created and updated
	ruleIdToOrder := make(map[string]int)
	for i, id := range ruleIDs {
		ruleIdToOrder[id] = len(ruleIDs) - i // Reverse the order
	}

	_, err = BulkReorder(context.Background(), service, policyType, ruleIdToOrder)
	if err != nil {
		t.Errorf("Error reordering rules: %v", err)
	}

	// Optionally verify the new order of rules here

	// Clean up: Delete the rules
	for _, ruleID := range ruleIDs {
		_, err = Delete(context.Background(), service, accessPolicySet.ID, ruleID)
		if err != nil {
			t.Errorf("Error deleting resource: %v", err)
		}
	}
}
