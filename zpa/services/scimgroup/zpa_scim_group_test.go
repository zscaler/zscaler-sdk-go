package scimgroup

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

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

func TestSCIMGroupGetByNameWithSort(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	testIdpId := getTestIdpId(t)

	scimGroupService := New(client)

	// Retrieve a list of SCIM groups
	scimGroups, _, err := scimGroupService.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Fatalf("Error retrieving SCIM groups: %v", err)
	}
	if len(scimGroups) == 0 {
		t.Fatalf("No SCIM groups found to test with")
	}

	// Check if we have enough groups for the test, otherwise return an error
	if len(scimGroups) < 100 {
		t.Fatalf("Not enough SCIM groups available for testing. Required: 50, Found: %d", len(scimGroups))
	}

	// Randomly pick a group name from the first 50 groups
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(500)
	testScimName := scimGroups[randomIndex].Name

	// Test with both DESC and ASC sort orders
	for _, sortOrder := range []SortOrder{DESCSortOrder, ASCSortOrder} {
		// Define sorting parameters
		sortField := IDSortField

		// Call GetByName with sorting parameters
		scimGroup, _, err := scimGroupService.WithSort(sortField, sortOrder).GetByName(testScimName, testIdpId)
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
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	testIdpId := getTestIdpId(t)
	service := New(client)

	groups, _, err := service.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Errorf("Error getting scim group: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No scim group found")
		return
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
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	testIdpId := getTestIdpId(t)
	service := New(client)
	_, _, err = service.GetByName("NonExistentName", testIdpId)
	if err == nil {
		t.Errorf("Expected error when getting non-existent SCIM group by name, got none")
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
		t.Errorf("Error getting SCIM Groups: %v", err)
		return
	}
	if groups == nil {
		t.Errorf("Received nil response for SCIM Groups")
		return
	}
}

func TestGetSCIMGroupByID(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	testIdpId := getTestIdpId(t)
	groups, _, err := service.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Errorf("Error getting all SCIM Groups: %v", err)
		return
	}

	if len(groups) == 0 {
		t.Errorf("No SCIM Group found")
		return
	}

	specificID := groups[0].ID
	group, _, err := service.Get(strconv.FormatInt(specificID, 10))
	if err != nil {
		t.Errorf("Error getting SCIM Group by ID: %v", err)
		return
	}
	if group.ID != specificID {
		t.Errorf("Mismatch in group ID: expected '%d', got %d", specificID, group.ID)
		return
	}
}

func TestAllFieldsOfSCIMGroups(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)
	testIdpId := getTestIdpId(t)
	groups, _, err := service.GetAllByIdpId(testIdpId)
	if err != nil {
		t.Errorf("Error getting all SCIM Group: %v", err)
		return
	}

	if len(groups) == 0 {
		t.Errorf("No SCIM Group found")
		return
	}

	specificID := groups[0].ID
	group, _, err := service.Get(strconv.FormatInt(specificID, 10))
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
		t.Errorf("Error getting SCIM Groups: %v", err)
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
