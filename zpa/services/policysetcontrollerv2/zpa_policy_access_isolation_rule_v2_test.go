package policysetcontrollerv2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/isolationprofile"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/samlattribute"
)

func TestAccessIsolationPolicyIsolateV2(t *testing.T) {
	policyType := "ISOLATION_POLICY"
	isolationProfileID := "BD_SA_Profile1"
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	idpService := idpcontroller.New(client)
	samlService := samlattribute.New(client)
	policyServiceV2 := New(client)

	idpList, _, err := idpService.GetAll()
	if err != nil {
		t.Errorf("Error getting idps: %v", err)
		return
	}
	if len(idpList) == 0 {
		t.Skip("No IdPs retrieved, skipping test as it requires at least one IdP")
		return
	}

	samlsList, _, err := samlService.GetAll()
	if err != nil {
		t.Errorf("Error getting saml attributes: %v", err)
		return
	}
	if len(samlsList) == 0 {
		t.Error("Expected retrieved saml attributes to be non-empty, but got empty slice")
	}

	profile := isolationprofile.New(client)
	profileID, _, err := profile.GetByName(isolationProfileID)
	if err != nil {
		t.Errorf("Error getting isolation profile id set: %v", err)
		return
	}

	accessPolicySet, _, err := policyServiceV2.GetByPolicyType(policyType)
	if err != nil {
		t.Errorf("Error getting access policy set: %v", err)
		return
	}

	var ruleIDs []string // Store the IDs of the created rules

	for i := 0; i < 5; i++ {
		// Generate a unique name for each iteration
		name := fmt.Sprintf("tests-%s-%d", acctest.RandStringFromCharSet(10, acctest.CharSetAlpha), i)

		accessPolicyRule := PolicyRule{
			Name:                  name,
			Description:           name,
			PolicySetID:           accessPolicySet.ID,
			ZpnIsolationProfileID: profileID.ID,
			Action:                "ISOLATE",
			Conditions: []PolicyRuleResourceConditions{
				{
					Operator: "OR",
					Operands: []PolicyRuleResourceOperands{
						{
							ObjectType:        "SAML",
							EntryValuesLHSRHS: []OperandsResourceLHSRHSValue{{LHS: samlsList[0].ID, RHS: "user1@acme.com"}},
						},
					},
				},
				{
					Operands: []PolicyRuleResourceOperands{
						{
							ObjectType:        "PLATFORM",
							EntryValuesLHSRHS: []OperandsResourceLHSRHSValue{{LHS: "linux", RHS: "true"}, {LHS: "windows", RHS: "true"}},
						},
					},
				},
				{
					Operator: "OR",
					Operands: []PolicyRuleResourceOperands{
						{
							ObjectType: "CLIENT_TYPE",
							Values:     []string{"zpn_client_type_exporter"},
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
		updatedResource, _, getErr := policyServiceV2.GetPolicyRule(accessPolicySet.ID, createdResource.ID)
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

	_, err = policyServiceV2.BulkReorder(policyType, ruleIdToOrder)
	if err != nil {
		t.Errorf("Error reordering rules: %v", err)
	}

	// Optionally verify the new order of rules here

	// Clean up: Delete the rules
	for _, ruleID := range ruleIDs {
		_, err = policyServiceV2.Delete(accessPolicySet.ID, ruleID)
		if err != nil {
			t.Errorf("Error deleting resource: %v", err)
		}
	}
}
