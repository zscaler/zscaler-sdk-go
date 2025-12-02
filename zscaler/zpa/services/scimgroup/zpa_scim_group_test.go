package scimgroup

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
)

func getTestIdpId(t *testing.T, service *zscaler.Service) string {
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

func TestSCIMGroup(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "scimgroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	idpList, _, err := idpcontroller.GetAll(context.Background(), service)
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

	// Test GetAllByIdpId function
	scimGroups, resp, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Logf("Error getting all SCIM groups by IdP ID: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Logf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		return
	}

	// If SCIM groups are present, test the GetByName function
	if len(scimGroups) > 0 {
		// Use the first SCIM group's name from the list for testing
		scimName := scimGroups[0].Name
		_, _, err = GetByName(context.Background(), service, scimName, testIdpId)
		if err != nil {
			t.Logf("Error getting SCIM group by name: %v", err)
		}
	} else {
		t.Logf("No SCIM groups retrieved for IdP ID: %s", testIdpId)
	}
}

func TestSCIMGroupGetByNameWithSort(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "scimgroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	testIdpId := getTestIdpId(t, service)

	// Retrieve a list of SCIM groups
	scimGroups, _, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Fatalf("Error retrieving SCIM groups: %v", err)
	}
	if len(scimGroups) == 0 {
		t.Skipf("No SCIM groups found to test with for IDP ID: %s", testIdpId)
	}

	// Check if we have enough groups for the test, otherwise return an error
	if len(scimGroups) < 10 {
		t.Fatalf("Not enough SCIM groups available for testing. Required: 10, Found: %d", len(scimGroups))
	}

	// Use first group (deterministic for VCR compatibility)
	testScimName := scimGroups[0].Name

	// Test with both DESC and ASC sort orders
	for _, sortOrder := range []zscaler.SortOrder{zscaler.DESCSortOrder, zscaler.ASCSortOrder} {
		// Define sorting parameters
		sortField := zscaler.IDSortField

		// Call GetByName with sorting parameters and manually set them in the service struct
		service.SortBy = sortField
		service.SortOrder = sortOrder

		scimGroup, _, err := GetByName(context.Background(), service, testScimName, testIdpId)
		if err != nil {
			t.Errorf("Error getting SCIM group by name with sort order %s: %v", sortOrder, err)
			continue
		}

		if scimGroup == nil {
			t.Errorf("No SCIM group named '%s' found with sort order %s", testScimName, sortOrder)
		}
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "scimgroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	testIdpId := getTestIdpId(t, service)

	groups, _, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Errorf("Error getting scim group: %v", err)
		return
	}

	// Instead of failing the test, log a message and return successfully if no groups are found.
	if len(groups) == 0 {
		t.Logf("No SCIM Group found for tenant ID: %s. This is not necessarily an error, depending on tenant configuration.", testIdpId)
		return // Return successfully since the absence of SCIM Groups is not considered a failure condition.
	}

	// Validate each group
	for _, group := range groups {
		// Checking if essential fields are not empty
		if group.ID == 0 {
			t.Errorf("Scim group ID is empty")
		}
		if group.Name == "" {
			t.Errorf("Scim group Name is empty")
		}
	}
}

func TestNonExistentSCIMGroupName(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "scimgroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	testIdpId := getTestIdpId(t, service)
	_, _, err = GetByName(context.Background(), service, "NonExistentName", testIdpId)
	if err == nil {
		t.Errorf("Expected error when getting non-existent SCIM group by name, got none")
	}
}

func TestEmptyResponse(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "scimgroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	testIdpId := getTestIdpId(t, service)
	groups, _, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Errorf("Error getting SCIM Groups: %v", err)
		return
	}
	if groups == nil {
		t.Errorf("Received nil response for SCIM Groups")
		return
	}
}

func TestGetSCIMGroupByID(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "scimgroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	testIdpId := getTestIdpId(t, service)
	groups, _, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Errorf("Error getting all SCIM Groups: %v", err)
		return
	}

	// Instead of failing the test, log a message and return successfully if no groups are found.
	if len(groups) == 0 {
		t.Logf("No SCIM Group found for tenant ID: %s. This is not necessarily an error, depending on tenant configuration.", testIdpId)
		return // Return successfully since the absence of SCIM Groups is not considered a failure condition.
	}

	// Proceed with the test if there are groups.
	specificID := groups[0].ID
	group, _, err := Get(context.Background(), service, strconv.FormatInt(specificID, 10)) // Pass the service and specificID
	if err != nil {
		t.Errorf("Error getting SCIM Group by ID: %v", err)
		return
	}
	if group.ID != specificID {
		t.Errorf("Mismatch in group ID: expected '%d', got '%d'", specificID, group.ID)
		return
	}
}

func TestAllFieldsOfSCIMGroups(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "scimgroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	testIdpId := getTestIdpId(t, service)

	groups, _, err := GetAllByIdpId(context.Background(), service, testIdpId)
	if err != nil {
		t.Errorf("Error getting all SCIM Group: %v", err)
		return
	}

	// Instead of failing the test, log a message and return successfully if no groups are found.
	if len(groups) == 0 {
		t.Logf("No SCIM Group found for tenant ID: %s. This is not necessarily an error, depending on tenant configuration.", testIdpId)
		return // Return successfully since the absence of SCIM Groups is not considered a failure condition.
	}

	// Retrieve a specific SCIM Group by its ID
	specificID := groups[0].ID
	group, _, err := Get(context.Background(), service, strconv.FormatInt(specificID, 10)) // Pass both service and specificID
	if err != nil {
		t.Errorf("Error getting SCIM Group by ID: %v", err)
		return
	}

	// Now check each field
	if group.ID == 0 {
		t.Errorf("ID is empty")
	}
	if group.IdpID == 0 {
		t.Errorf("IdpID is empty")
	}
	if group.Name == "" {
		t.Errorf("Name is empty")
	}
	if group.InternalID == "" {
		t.Errorf("InternalID is empty")
	}
}
