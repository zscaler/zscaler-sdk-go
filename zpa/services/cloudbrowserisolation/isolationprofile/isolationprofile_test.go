package isolationprofile

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

func TestIsolationProfile(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	// Test to retrieve all profiles
	profiles, _, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting isolation profiles: %v", err)
		return
	}
	if len(profiles) == 0 {
		t.Errorf("No isolation profile found")
		return
	}

	// Test to retrieve a profile by its name
	name := profiles[0].Name
	t.Log("Getting isolation profile by name:" + name)
	profile, _, err := GetByName(service, name)
	if err != nil {
		t.Errorf("Error getting isolation profile by name: %v", err)
		return
	}
	if profile.Name != name {
		t.Errorf("Isolation profile name does not match: expected %s, got %s", name, profile.Name)
		return
	}

	// Negative Test: Try to retrieve a profile with a non-existent name
	nonExistentName := "ThisProfileNameDoesNotExist"
	_, _, err = GetByName(service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}

	// Cover the GetAll function error
	service.Client.Config.CustomerID = "invalid_id"
	_, _, err = GetAll(service)
	if err == nil {
		t.Errorf("Expected error when getting all profiles with invalid CustomerID, got nil")
		return
	}
	// Restore valid CustomerID for further tests
	service.Client.Config.CustomerID = client.Config.CustomerID
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	profiles, _, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting isolation profiles: %v", err)
		return
	}
	if len(profiles) == 0 {
		t.Errorf("No isolation profile found")
		return
	}

	// Validate each profile
	for _, profile := range profiles {
		// Checking if essential fields are not empty
		if profile.ID == "" {
			t.Errorf("IsolationProfile ID is empty")
		}
		if profile.Name == "" {
			t.Errorf("IsolationProfile Name is empty")
		}
		if profile.IsolationURL == "" {
			t.Errorf("IsolationProfile IsolationURL is empty")
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	requiredNames := []string{"BD_SA_Profile1", "BD SA Profile", "BD  SA Profile", "BD   SA   Profile"}

	for _, knownName := range requiredNames {
		// Case variations to test for each knownName
		variations := []string{
			strings.ToUpper(knownName),
			strings.ToLower(knownName),
			cases.Title(language.English).String(knownName),
		}

		for _, variation := range variations {
			t.Run(fmt.Sprintf("GetByName case sensitivity test for %s", variation), func(t *testing.T) {
				t.Logf("Attempting to retrieve isolation profile with name variation: %s", variation)
				version, _, err := GetByName(service, variation)
				if err != nil {
					t.Errorf("Error getting isolation profile with name variation '%s': %v", variation, err)
					return
				}

				// Check if the isolation profile's actual name matches the known name
				if version.Name != knownName {
					t.Errorf("Expected isolation profile name to be '%s' for variation '%s', but got '%s'", knownName, variation, version.Name)
				}
			})
		}
	}
}

func TestProfileNamesWithSpaces(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	// Assuming that there are profiles with the following name variations
	variations := []string{
		"BD SA Profile",     // Single space
		"BD  SA Profile",    // Double space
		"BD   SA   Profile", // Multiple spaces
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve profile with name: %s", variation)
		profile, _, err := GetByName(service, variation)
		if err != nil {
			t.Errorf("Error getting isolation profile with name '%s': %v", variation, err)
			continue
		}

		// Verify if the profile's actual name matches the expected variation
		if profile.Name != variation {
			t.Errorf("Expected profile name to be '%s' but got '%s'", variation, profile.Name)
		}
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	_, _, err = GetByName(service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non_existent_name name, but got nil")
	}
}
