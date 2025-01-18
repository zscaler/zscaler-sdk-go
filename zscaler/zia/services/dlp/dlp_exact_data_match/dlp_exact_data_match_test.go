package dlp_exact_data_match

import (
	"context"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestDLPEDM_data(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	templates, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting idm profiles: %v", err)
		return
	}
	if len(templates) == 0 {
		t.Errorf("No idm profile found")
		return
	}
	name := templates[0].ProjectName
	t.Log("Getting edm template by name:" + name)
	template, err := GetDLPEDMByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error getting edm template by name: %v", err)
		return
	}
	if template.ProjectName != name {
		t.Errorf("edm template name does not match: expected %s, got %s", name, template.ProjectName)
		return
	}
	// Negative Test: Try to retrieve an edm template with a non-existent name
	nonExistentName := "ThisEDMTemplateDoesNotExist"
	_, err = GetDLPEDMByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestGetById(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Get all servers to find a valid ID
	templates, err := GetAll(context.Background(), service)
	if err != nil {
		t.Fatalf("Error getting all edm templates: %v", err)
	}
	if len(templates) == 0 {
		t.Fatalf("No edm templates found for testing")
	}

	// Choose the first server's ID for testing
	testID := templates[0].SchemaID

	// Retrieve the server by ID
	template, err := GetDLPEDMSchemaID(context.Background(), service, testID)
	if err != nil {
		t.Errorf("Error retrieving edm template with ID %d: %v", testID, err)
		return
	}

	// Verify the retrieved server
	if template == nil {
		t.Errorf("No server returned for ID %d", testID)
		return
	}

	if template.SchemaID != testID {
		t.Errorf("Retrieved server ID mismatch: expected %d, got %d", testID, template.SchemaID)
	}
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	templates, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting edm template: %v", err)
		return
	}
	if len(templates) == 0 {
		t.Errorf("No edm template found")
		return
	}

	// Validate edm template
	for _, template := range templates {
		// Checking if essential fields are not empty
		if template.SchemaID == 0 {
			t.Errorf("edm template ID is empty")
		}
		if template.ProjectName == "" {
			t.Errorf("edm template Name is empty")
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Assuming a edm template with the name "BD_EDM_TEMPLATE01" exists
	knownName := "BD_EDM_TEMPLATE01"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		template, err := GetDLPEDMByName(context.Background(), service, variation)
		if err != nil {
			t.Errorf("Error getting edm template with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if template.ProjectName != knownName {
			t.Errorf("Expected edm template name to be '%s' for variation '%s', but got '%s'", knownName, variation, template.ProjectName)
		}
	}
}

func TestEDMFields(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Retrieve all EDM profiles
	edmProfiles, err := GetAll(context.Background(), service) // Assuming appropriate method name and parameters
	if err != nil {
		t.Fatalf("Error getting all EDM profiles: %v", err)
	}
	if len(edmProfiles) == 0 {
		t.Fatalf("No EDM profiles found for testing")
	}

	// Iterate through each EDM profile and check various fields
	for _, profile := range edmProfiles {
		if profile.SchemaID == 0 {
			t.Errorf("SchemaID field is empty")
		}
		if profile.ProjectName == "" {
			t.Errorf("ProjectName field is empty")
		}
		if !profile.SchemaActive {
			t.Errorf("SchemaActive field is not active")
		}

		// Asserting elements in the TokenList
		primaryKeyFound := false
		for _, token := range profile.TokenList {
			if token.Name == "" || token.Type == "" {
				t.Errorf("Token fields Name or Type are not properly populated")
			}
			if token.PrimaryKey {
				primaryKeyFound = true
				break // Break after finding the first primary key
			}
		}
		if !primaryKeyFound {
			t.Errorf("No primary key token found in profile: %s", profile.ProjectName)
		}
	}
}
