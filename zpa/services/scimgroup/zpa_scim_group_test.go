package scimgroup

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/idpcontroller"
)

func TestSCIMGroup(t *testing.T) {
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

	// Find an IdP with ssoType USER
	var testIdpId string
	for _, idp := range idpList {
		for _, ssoType := range idp.SsoType {
			if ssoType == "USER" {
				testIdpId = idp.ID
				break
			}
		}
		if testIdpId != "" {
			break
		}
	}

	if testIdpId == "" {
		t.Error("No IdP with ssoType USER found")
		return
	}

	scimGroupService := New(client)

	// Test GetAllByIdpId function
	scimGroups, resp, err := scimGroupService.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Errorf("Error getting all SCIM groups by IdP ID: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		return
	}

	// If SCIM groups are present, test the GetByName function
	if len(scimGroups) > 0 {
		// Use the first SCIM group's name from the list for testing
		scimName := scimGroups[0].Name
		_, _, err = scimGroupService.GetByName(scimName, testIdpId)
		if err != nil {
			t.Errorf("Error getting SCIM group by name: %v", err)
		}
	} else {
		t.Logf("No SCIM groups retrieved for IdP ID: %s", testIdpId)
	}
}
