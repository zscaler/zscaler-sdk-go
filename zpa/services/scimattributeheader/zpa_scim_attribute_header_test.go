package scimattributeheader

import (
	"net/http"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/idpcontroller"
)

func getTestIdpId(t *testing.T) string {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return ""
	}

	idpService := idpcontroller.New(client)
	idpList, _, err := idpService.GetAll()
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
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	testIdpId := getTestIdpId(t)
	scimAttributeService := New(client)

	// Test GetAllByIdpId function
	scimAttribute, resp, err := scimAttributeService.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Fatalf("Error getting all SCIM Attribute Header by IdP ID: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// If attribute list is empty, skip the subsequent logic
	if len(scimAttribute) == 0 {
		t.Log("No SCIM Attribute Header found, skipping further tests.")
		return
	}

	// Use the first SCIM attribute headers's name from the list for testing
	scimName := scimAttribute[0].Name
	_, _, err = scimAttributeService.GetByName(scimName, testIdpId)
	if err != nil {
		t.Fatalf("Error getting SCIM Attribute Headers by name: %v", err)
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	testIdpId := getTestIdpId(t)
	service := New(client)

	groups, _, err := service.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Errorf("Error getting SCIM Attribute Header: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No SCIM Attribute Header found")
		return
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
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	testIdpId := getTestIdpId(t)
	service := New(client)
	_, _, err = service.GetByName("NonExistentName", testIdpId)
	if err == nil {
		t.Errorf("Expected error when getting non-existent SCIM Attribute Header by name, got none")
	}
}

func TestEmptyResponse(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	testIdpId := getTestIdpId(t)
	groups, _, err := service.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Errorf("Error getting SCIM Attribute Header: %v", err)
		return
	}
	if groups == nil {
		t.Errorf("Received nil response for SCIM Attribute Header")
		return
	}
}

func TestGetSCIMAttributeHeaderByID(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	testIdpId := getTestIdpId(t)
	attributes, _, err := service.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Errorf("Error getting all SCIM Attribute Headers: %v", err)
		return
	}

	if len(attributes) == 0 {
		t.Errorf("No SCIM Attribute Header found")
		return
	}

	specificID := attributes[0].ID
	group, _, err := service.Get(testIdpId, specificID)
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
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)
	testIdpId := getTestIdpId(t)

	// Retrieve all attributes for the IdP
	attributes, _, err := service.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Fatalf("Error getting all SCIM Attribute Header: %v", err)
	}
	if len(attributes) == 0 {
		t.Fatalf("No SCIM Attribute Header found")
	}

	// Use the ID of the first attribute for GetValues
	attributeID := attributes[0].ID
	values, err := service.GetValues(testIdpId, attributeID)
	if err != nil {
		t.Fatalf("Error getting values: %v", err)
	}
	if len(values) == 0 {
		t.Errorf("No values found")
	}
	// ... add more assertions as needed
}

func TestAllFieldsOfSCIMAttributeHeaders(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	testIdpId := getTestIdpId(t)
	attributes, _, err := service.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Errorf("Error getting all SCIM Attribute Header: %v", err)
		return
	}

	if len(attributes) == 0 {
		t.Errorf("No SCIM Attribute Header found")
		return
	}

	specificID := attributes[0].ID
	attribute, _, err := service.Get(testIdpId, specificID)
	if err != nil {
		t.Errorf("Error getting SCIM Attribute Header by ID: %v", err)
		return
	}

	// Now check each field
	if attribute.CreationTime == "" {
		t.Errorf("CreationTime is empty")
	}
	if attribute.ModifiedTime == "" {
		t.Errorf("ModifiedTime is empty")
	}
	if attribute.ModifiedBy == "" {
		t.Errorf("ModifiedBy is empty")
	}
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

func TestResponseHeadersAndFormat(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	testIdpId := getTestIdpId(t)
	_, resp, err := service.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Errorf("Error getting SCIM Attribute Header: %v", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		t.Errorf("Expected content type to start with 'application/json', got %s", contentType)
	}
}
