package cloudbrowserisolation

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestIsolationProfile(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Test to retrieve all profiles
	profiles, err := service.GetAll()
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
	profile, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting isolation profile by name: %v", err)
		return
	}
	if profile.Name != name {
		t.Errorf("isolation profile name does not match: expected %s, got %s", name, profile.Name)
		return
	}

	// Negative Test: Try to retrieve a profile with a non-existent name
	nonExistentName := "ThisProfileNameDoesNotExist"
	_, err = service.GetByName(nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	profiles, err := service.GetAll()
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
		if profile.URL == "" {
			t.Errorf("IsolationProfile Isolation URL is empty")
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	service := New(client)

	requiredNames := []string{"BD_SA_Profile1_ZIA", "BD_SA_Profile2_ZIA", "BD SA Profile ZIA", "BD  SA Profile ZIA", "BD   SA   Profile  ZIA"}

	for _, knownName := range requiredNames {
		// Case variations to test for each knownName
		variations := []string{
			strings.ToUpper(knownName),
			strings.ToLower(knownName),
			cases.Title(language.English).String(knownName),
		}

		for _, variation := range variations {
			t.Run(fmt.Sprintf("GetByName case sensitivity test for %s", variation), func(t *testing.T) {
				t.Logf("Attempting to retrieve customer version profile with name variation: %s", variation)
				version, err := service.GetByName(variation)
				if err != nil {
					t.Errorf("Error getting customer version profile with name variation '%s': %v", variation, err)
					return
				}

				// Check if the customer version profile's actual name matches the known name
				if version.Name != knownName {
					t.Errorf("Expected customer version profile name to be '%s' for variation '%s', but got '%s'", knownName, variation, version.Name)
				}
			})
		}
	}
}

func TestProfileNamesWithSpaces(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Assuming that there are profiles with the following name variations
	variations := []string{
		"BD SA Profile ZIA",      // Single space
		"BD  SA Profile ZIA",     // Double space
		"BD   SA   Profile  ZIA", // Multiple spaces
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve profile with name: %s", variation)
		profile, err := service.GetByName(variation)
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
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.GetByName("non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
