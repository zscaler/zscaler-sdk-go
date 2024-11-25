package scimattributeheader

import (
	"context"
	"net/http"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
)

func getTestIdpId(t *testing.T) string {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	idpList, _, err := idpcontroller.GetAll(context.Background(), service)
	if err != nil {
		t.Fatalf("Error getting idps: %v", err)
		return ""
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
		t.Fatalf("No IdP with ssoType USER found")
		return ""
	}

	return testIdpId
}

func TestSCIMAttributeHeader(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	testIdpId := getTestIdpId(t)

	// Test GetAllByIdpId function
	scimAttribute, resp, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Fatalf("Error getting all SCIM Attribute Header by IdP ID: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Instead of failing the test, log a message and return successfully if no groups are found.
	if len(scimAttribute) == 0 {
		t.Logf("No SCIM Attribute Header found, skipping further tests.")
		return // Return successfully since the absence of SCIM Groups is not considered a failure condition.
	}

	// Use the first SCIM attribute headers's name from the list for testing
	scimName := scimAttribute[0].Name
	_, _, err = GetByName(context.Background(), service, scimName, testIdpId)
	if err != nil {
		t.Fatalf("Error getting SCIM Attribute Headers by name: %v", err)
	}
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	testIdpId := getTestIdpId(t)

	groups, _, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Errorf("Error getting SCIM Attribute Header: %v", err)
		return
	}

	// Instead of failing the test, log a message and return successfully if no groups are found.
	if len(groups) == 0 {
		t.Logf("No SCIM Attribute Header found")
		return // Return successfully since the absence of SCIM Groups is not considered a failure condition.
	}

	// Validate each group
	for _, group := range groups {
		// Checking if essential fields are not empty
		if group.ID == "" {
			t.Errorf("SCIM Attribute Header ID is empty")
		}
		if group.Name == "" {
			t.Errorf("SCIM Attribute Header Name is empty")
		}
	}
}

func TestNonExistentSCIMAttributeHeaderName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	testIdpId := getTestIdpId(t)
	_, _, err = GetByName(context.Background(), service, "NonExistentName", testIdpId)
	if err == nil {
		t.Errorf("Expected error when getting non-existent SCIM Attribute Header by name, got none")
	}
}

func TestEmptyResponse(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	testIdpId := getTestIdpId(t)
	groups, _, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Errorf("Error getting SCIM Attribute Header: %v", err)
		return
	}

	// Simplified check for an empty response
	if len(groups) == 0 {
		t.Logf("Received an empty response for SCIM Attribute Header for IdP ID: %s. This may be expected if no SCIM groups are configured.", testIdpId)
	} else {
		t.Logf("Received response for SCIM Attribute Header for IdP ID: %s with %d groups.", testIdpId, len(groups))
	}

}

func TestGetSCIMAttributeHeaderByID(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	testIdpId := getTestIdpId(t)

	attributes, _, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Errorf("Error getting all SCIM Attribute Headers: %v", err)
		return
	}

	// Instead of failing the test, log a message and return successfully if no groups are found.
	if len(attributes) == 0 {
		t.Logf("No SCIM Attribute Header found")
		return // Return successfully since the absence of SCIM Groups is not considered a failure condition.
	}

	specificID := attributes[0].ID
	group, _, err := Get(context.Background(), service, testIdpId, specificID)
	if err != nil {
		t.Errorf("Error getting SCIM Attribute Header by ID: %v", err)
		return
	}
	if group.ID != specificID {
		t.Errorf("Mismatch in attribute header ID: expected '%s', got %s", specificID, group.ID)
		return
	}
}

func TestSCIMAttributeHeaderGetValues(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	testIdpId := getTestIdpId(t)

	// Retrieve all attributes for the IdP
	attributes, _, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Fatalf("Error getting all SCIM Attribute Header: %v", err)
	}

	// Instead of failing the test, log a message and return successfully if no groups are found.
	if len(attributes) == 0 {
		t.Logf("No SCIM Attribute Header found")
		return // Return successfully since the absence of SCIM Groups is not considered a failure condition.
	}

	// Use the ID of the first attribute for GetValues
	attributeID := attributes[0].ID
	values, err := GetValues(context.Background(), service, testIdpId, attributeID)
	if err != nil {
		t.Fatalf("Error getting values for attribute ID %s: %v", attributeID, err)
	}
	if len(values) == 0 {
		t.Logf("No values found for attribute ID %s, but proceeding with the test.", attributeID)
		return // Proceed with the test despite no values found
	}

	// Add any additional assertions here if you have values
}

func TestAllFieldsOfSCIMAttributeHeaders(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	testIdpId := getTestIdpId(t)
	attributes, _, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Errorf("Error getting all SCIM Attribute Header: %v", err)
		return
	}

	// Instead of failing the test, log a message and return successfully if no groups are found.
	if len(attributes) == 0 {
		t.Logf("No SCIM Attribute Header found")
		return // Return successfully since the absence of SCIM Groups is not considered a failure condition.
	}

	specificID := attributes[0].ID
	attribute, _, err := Get(context.Background(), service, testIdpId, specificID)
	if err != nil {
		t.Errorf("Error getting SCIM Attribute Header by ID: %v", err)
		return
	}

	// Now check each field
	if attribute.ID == "" {
		t.Errorf("ID is empty")
	}
	if attribute.Name == "" {
		t.Errorf("Name is empty")
	}
	if attribute.IdpID == "" {
		t.Errorf("IdpID is empty")
	}
	if attribute.DataType == "" {
		t.Errorf("DataType is empty")
	}
	if attribute.SchemaURI == "" {
		t.Errorf("SchemaURI is empty")
	}
}

/*
func TestResponseHeadersAndFormat(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	testIdpId := getTestIdpId(t)
	_, resp, err :=  GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Errorf("Error getting SCIM Attribute Header: %v", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	t.Logf("Content-Type header received: %s", contentType)
	if !strings.HasPrefix(contentType, "application/json") {
		t.Errorf("Expected content type to start with 'application/json', got %s", contentType)
	}
}
*/
