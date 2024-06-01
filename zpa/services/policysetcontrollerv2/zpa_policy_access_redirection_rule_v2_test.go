package policysetcontrollerv2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/serviceedgegroup"
)

func TestAccessRedirectionPolicyV2(t *testing.T) {
	policyType := "REDIRECTION_POLICY"
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	policyService := policysetcontroller.New(client)
	policyServiceV2 := New(client)

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

	accessPolicySet, _, err := policyService.GetByPolicyType(policyType)
	if err != nil {
		t.Errorf("Error getting access policy set: %v", err)
		return
	}

	var ruleIDs []string // Store the IDs of the created rules

	for i := 0; i < 5; i++ {
		// Generate a unique name for each iteration
		name := fmt.Sprintf("tests-%s-%d", acctest.RandStringFromCharSet(10, acctest.CharSetAlpha), i)

		accessPolicyRule := PolicyRule{
			Name:        name,
			Description: name,
			PolicySetID: accessPolicySet.ID,
			Action:      "REDIRECT_PREFERRED",
			ServiceEdgeGroups: []ServiceEdgeGroups{
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
		createdResource, _, err := policyServiceV2.CreateRule(&accessPolicyRule)

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
		_, updateErr := policyServiceV2.UpdateRule(accessPolicySet.ID, createdResource.ID, &accessPolicyRule)

		if updateErr != nil {
			t.Errorf("Error updating rule: %v", updateErr)
			continue
		}

		// Retrieve and print the updated resource as JSON
		updatedResource, _, getErr := policyService.GetPolicyRule(accessPolicySet.ID, createdResource.ID)
		if getErr != nil {
			t.Errorf("Error retrieving updated resource: %v", getErr)
			continue
		}
		if updatedResource.Name != updatedName {
			t.Errorf("Expected updated resource name '%s', but got '%s'", updatedName, updatedResource.Name)
		}
		// Print the updated resource as JSON
		// updatedJson, _ := json.Marshal(updatedResource)
		// fmt.Println(string(updatedJson))

		// Introduce a delay to prevent rate limit issues
		time.Sleep(10 * time.Second)
	}

	// Reorder the rules after all have been created and updated
	ruleIdToOrder := make(map[string]int)
	for i, id := range ruleIDs {
		ruleIdToOrder[id] = len(ruleIDs) - i // Reverse the order
	}

	_, err = policyService.BulkReorder(policyType, ruleIdToOrder)
	if err != nil {
		t.Errorf("Error reordering rules: %v", err)
	}

	// Optionally verify the new order of rules here

	// Clean up: Delete the rules
	for _, ruleID := range ruleIDs {
		_, err = policyService.Delete(accessPolicySet.ID, ruleID)
		if err != nil {
			t.Errorf("Error deleting resource: %v", err)
		}
	}
}
