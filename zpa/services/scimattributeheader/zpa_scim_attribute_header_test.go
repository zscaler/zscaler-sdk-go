package scimattributeheader

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/idpcontroller"
)

func TestSCIMAttributeHeader(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	idpService := idpcontroller.New(client)
	idpList, _, err := idpService.GetAll()
	if err != nil {
		t.Fatalf("Error getting idps: %v", err)
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
		t.Fatal("No IdP with ssoType USER found")
	}

	scimAttributeService := New(client)

	// Test GetAllByIdpId function
	scimAttribute, resp, err := scimAttributeService.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Fatalf("Error getting all SCIM attribute headers by IdP ID: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// If attribute list is empty, skip the subsequent logic
	if len(scimAttribute) == 0 {
		t.Log("No SCIM attribute headers found, skipping further tests.")
		return
	}

	// Use the first SCIM attribute headers's name from the list for testing
	scimName := scimAttribute[0].Name
	_, _, err = scimAttributeService.GetByName(scimName, testIdpId)
	if err != nil {
		t.Fatalf("Error getting SCIM attribute headers by name: %v", err)
	}
}
