package policysetcontrollerv2

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/postureprofile"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/samlattribute"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/trustednetwork"
)

func TestAccessForwardingPolicyV2(t *testing.T) {
	segmentGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	policyType := "CLIENT_FORWARDING_POLICY"
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	idpList, _, err := idpcontroller.GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting idps: %v", err)
		return
	}
	if len(idpList) == 0 {
		t.Skip("No IdPs retrieved, skipping test as it requires at least one IdP")
		return
	}

	samlsList, _, err := samlattribute.GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting saml attributes: %v", err)
		return
	}
	if len(samlsList) == 0 {
		t.Error("Expected retrieved saml attributes to be non-empty, but got empty slice")
	}

	postureList, _, err := postureprofile.GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting posture profiles: %v", err)
		return
	}
	if len(postureList) == 0 {
		t.Error("Expected retrieved posture profiles to be non-empty, but got empty slice")
	}

	trustedNetworkList, _, err := trustednetwork.GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting trusted networks: %v", err)
		return
	}
	if len(postureList) == 0 {
		t.Error("Expected retrieved trusted networks to be non-empty, but got empty slice")
	}

	appGroup := segmentgroup.SegmentGroup{
		Name:        segmentGroupName,
		Description: segmentGroupName,
	}
	createdAppGroup, _, err := segmentgroup.Create(context.Background(), service, &appGroup)
	if err != nil {
		t.Errorf("Error creating segment group: %v", err)
		return
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := segmentgroup.Get(context.Background(), service, createdAppGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := segmentgroup.Delete(context.Background(), service, createdAppGroup.ID)
			if err != nil {
				t.Errorf("Error deleting segment group: %v", err)
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
			Action:      "INTERCEPT",
			Conditions: []PolicyRuleResourceConditions{
				{
					Operator: "OR",
					Operands: []PolicyRuleResourceOperands{
						{
							ObjectType: "APP_GROUP",
							Values:     []string{createdAppGroup.ID},
						},
					},
				},
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
						{
							ObjectType:        "TRUSTED_NETWORK",
							EntryValuesLHSRHS: []OperandsResourceLHSRHSValue{{LHS: trustedNetworkList[0].NetworkID, RHS: "true"}, {LHS: trustedNetworkList[1].NetworkID, RHS: "true"}},
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
