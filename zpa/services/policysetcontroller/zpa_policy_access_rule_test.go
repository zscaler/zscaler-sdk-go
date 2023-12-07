package policysetcontroller

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/idpcontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/postureprofile"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/samlattribute"
)

func TestPolicyAccessRule(t *testing.T) {
	policyType := "ACCESS_POLICY"
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
	postureService := postureprofile.New(client)
	postureList, _, err := postureService.GetAll()
	if err != nil {
		t.Errorf("Error getting posture profiles: %v", err)
		return
	}
	if len(postureList) == 0 {
		t.Error("Expected retrieved posture profiles to be non-empty, but got empty slice")
	}
	service := New(client)
	accessPolicySet, _, err := service.GetByPolicyType(policyType)
	if err != nil {
		t.Errorf("Error getting access policy set: %v", err)
		return
	}

	for i := 0; i < 3; i++ {
		// Generate a unique name for each iteration
		name := fmt.Sprintf("tests-%s-%d", acctest.RandStringFromCharSet(10, acctest.CharSetAlpha), i)

		accessPolicyRule := PolicyRule{
			Name:        name,
			Description: name,
			PolicySetID: accessPolicySet.ID,
			Action:      "ALLOW",
			Conditions: []Conditions{
				{
					Operator: "OR",
					Operands: []Operands{
						{
							ObjectType: "POSTURE",
							LHS:        postureList[0].PostureudID,
							RHS:        "false",
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
		if err != nil {
			t.Errorf("Error making POST request: %v", err)
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
}
