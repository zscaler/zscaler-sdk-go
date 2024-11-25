package dlp_exact_data_match_lite

import (
	"context"
	"strings"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestDLPEDMLite(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	for _, activeOnly := range []bool{true, false} {
		profiles, err := GetAllEDMSchema(context.Background(), service, activeOnly, true) // Assuming fetchTokens is relevant here
		if err != nil {
			t.Errorf("Error getting idm profiles with activeOnly %t: %v", activeOnly, err)
			return
		}
		if len(profiles) == 0 {
			t.Errorf("No idm profile found with activeOnly %t", activeOnly)
			return
		}
		schemaName := profiles[0].Schema.Name
		t.Log("Getting idm profile by schema name:", schemaName, "with activeOnly:", activeOnly)
		profileByName, err := GetBySchemaName(context.Background(), service, schemaName, activeOnly, true)
		if err != nil {
			t.Errorf("Error getting idm profile by schema name with activeOnly %t: %v", activeOnly, err)
			return
		}
		if len(profileByName) == 0 || profileByName[0].Schema.Name != schemaName {
			t.Errorf("idm profile name does not match with activeOnly %t: expected %s", activeOnly, schemaName)
			return
		}
	}
	// Negative Test: Try to retrieve an edm template with a non-existent name
	nonExistentName := "ThisEdmDoesNotExist"
	profiles, err := GetBySchemaName(context.Background(), service, nonExistentName, false, true)
	if err != nil {
		t.Errorf("Error when getting by non-existent name: %v", err)
		return
	}

	// Check if the result set is empty
	if len(profiles) != 0 {
		t.Errorf("Expected no profiles for non-existent name, but found some")
	}
}

func TestGetDLPProfileLiteById(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	for _, activeOnly := range []bool{true, false} {
		templates, err := GetAllEDMSchema(context.Background(), service, activeOnly, true)
		if err != nil {
			t.Fatalf("Error getting all EDM templates with activeOnly %t: %v", activeOnly, err)
		}
		if len(templates) == 0 {
			t.Fatalf("No EDM template found for testing with activeOnly %t", activeOnly)
		}

		t.Logf("Total EDM templates found with activeOnly %t: %d", activeOnly, len(templates))

		testID := templates[0].Schema.ID
		if testID == 0 {
			t.Errorf("The Schema ID of the first profile is empty with activeOnly %t", activeOnly)
		} else {
			t.Logf("Valid Schema ID found with activeOnly %t: %d", activeOnly, testID)
		}
	}
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	for _, activeOnly := range []bool{true, false} {
		templates, err := GetAllEDMSchema(context.Background(), service, activeOnly, true)
		if err != nil {
			t.Errorf("Error getting edm template with activeOnly %t: %v", activeOnly, err)
			return
		}
		if len(templates) == 0 {
			t.Errorf("No edm template found with activeOnly %t", activeOnly)
			return
		}

		for _, template := range templates {
			if template.Schema.ID == 0 {
				t.Errorf("edm template Schema ID is empty with activeOnly %t", activeOnly)
			}
			if template.Schema.Name == "" {
				t.Errorf("edm template Schema Name is empty with activeOnly %t", activeOnly)
			}
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

	for _, activeOnly := range []bool{true, false} {
		for _, variation := range variations {
			t.Logf("Attempting to retrieve group with schema name variation: %s with activeOnly %t", variation, activeOnly)
			profiles, err := GetBySchemaName(context.Background(), service, variation, activeOnly, true)
			if err != nil {
				t.Errorf("Error getting idm profile with schema name variation '%s' and activeOnly %t: %v", variation, activeOnly, err)
				continue
			}
			if len(profiles) == 0 || profiles[0].Schema.Name != knownName {
				t.Errorf("Expected group name to be '%s' for variation '%s' with activeOnly %t, but got none or incorrect name", knownName, variation, activeOnly)
			}
		}
	}
}

func TestEDMFields(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Retrieve all EDM profiles
	edmProfiles, err := GetAllEDMSchema(context.Background(), service, true, false) // Assuming appropriate method name and parameters
	if err != nil {
		t.Fatalf("Error getting all EDM profiles: %v", err)
	}
	if len(edmProfiles) == 0 {
		t.Fatalf("No EDM profiles found for testing")
	}

	// Iterate through each EDM profile and check various fields
	for _, profile := range edmProfiles {
		if profile.Schema.ID == 0 {
			t.Errorf("Schema ID field is empty")
		}
		if profile.Schema.Name == "" {
			t.Errorf("Schema Name field is empty")
		}

		// Asserting elements in the TokenList
		for _, token := range profile.TokenList {
			if token.Name == "" || token.Type == "" {
				t.Errorf("Token fields Name or Type are not properly populated")
			}
			if !token.PrimaryKey {
				t.Errorf("Expected a primary key token, but found none")
			}
		}
	}
}
