package policysetcontrollerv2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/postureprofile"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/samlattribute"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
)

func TestAccessTimeoutPolicyV2(t *testing.T) {
	segmentGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	policyType := "TIMEOUT_POLICY"
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	appGroupService := segmentgroup.New(client)
	idpService := idpcontroller.New(client)
	samlService := samlattribute.New(client)
	postureService := postureprofile.New(client)
	policyService := policysetcontroller.New(client)
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

	postureList, _, err := postureService.GetAll()
	if err != nil {
		t.Errorf("Error getting posture profiles: %v", err)
		return
	}
	if len(postureList) == 0 {
		t.Error("Expected retrieved posture profiles to be non-empty, but got empty slice")
	}

	appGroup := segmentgroup.SegmentGroup{
		Name:        segmentGroupName,
		Description: segmentGroupName,
	}
	createdAppGroup, _, err := appGroupService.Create(&appGroup)
	if err != nil {
		t.Errorf("Error creating segment group: %v", err)
		return
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := appGroupService.Get(createdAppGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := appGroupService.Delete(createdAppGroup.ID)
			if err != nil {
				t.Errorf("Error deleting segment group: %v", err)
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
			Name:              name,
			Description:       name,
			PolicySetID:       accessPolicySet.ID,
			Action:            "RE_AUTH",
			ReauthIdleTimeout: "600",
			ReauthTimeout:     "172800",
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
					Operator: "OR",
					Operands: []PolicyRuleResourceOperands{
						{
							ObjectType:        "POSTURE",
							EntryValuesLHSRHS: []OperandsResourceLHSRHSValue{{LHS: postureList[0].PostureudID, RHS: "true"}, {LHS: postureList[1].PostureudID, RHS: "true"}},
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
							Values:     []string{"zpn_client_type_zapp"},
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
