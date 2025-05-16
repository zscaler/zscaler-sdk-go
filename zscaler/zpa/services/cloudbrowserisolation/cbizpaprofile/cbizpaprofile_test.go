package cbizpaprofile

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestCBIZPAProfile(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	// Test to retrieve all profiles
	profiles, _, err := GetAll(context.Background(), service)
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
	t.Log("Getting isolation profile by name: " + name)
	profile, _, err := GetByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error getting isolation profile by name: %v", err)
		return
	}
	if profile.Name != name {
		t.Errorf("Isolation profile name does not match: expected %s, got %s", name, profile.Name)
		return
	}

	// Sleep for 1 second
	time.Sleep(1 * time.Second)

	// New test to retrieve a profile by its ID
	t.Run("TestGetProfileByID", func(t *testing.T) {
		id := profiles[0].ID
		t.Log("Getting isolation profile by ID: " + id)
		profileByID, _, err := Get(context.Background(), service, id)
		if err != nil {
			t.Errorf("Error getting isolation profile by ID: %v", err)
			return
		}
		if profileByID.ID != id {
			t.Errorf("Isolation profile ID does not match: expected %s, got %s", id, profileByID.ID)
		}
	})

	// Sleep for 1 second
	time.Sleep(1 * time.Second)

	// Negative Test: Try to retrieve a profile with a non-existent name
	nonExistentName := "ThisProfileNameDoesNotExist"
	_, _, err = GetByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	profiles, _, err := GetAll(context.Background(), service)
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
		if profile.CBIURL == "" {
			t.Errorf("IsolationProfile CBI URL is empty")
		}
	}

	// Sleep for 1 second
	time.Sleep(1 * time.Second)
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

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
				t.Logf("Attempting to retrieve customer version profile with name variation: %s", variation)
				version, _, err := GetByName(context.Background(), service, variation)
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

		// Sleep for 1 second
		time.Sleep(1 * time.Second)
	}
}

func TestProfileNamesWithSpaces(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }
	// Assuming that there are profiles with the following name variations
	variations := []string{
		"BD SA Profile",     // Single space
		"BD  SA Profile",    // Double space
		"BD   SA   Profile", // Multiple spaces
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve profile with name: %s", variation)
		profile, _, err := GetByName(context.Background(), service, variation)
		if err != nil {
			t.Errorf("Error getting isolation profile with name '%s': %v", variation, err)
			continue
		}

		// Verify if the profile's actual name matches the expected variation
		if profile.Name != variation {
			t.Errorf("Expected profile name to be '%s' but got '%s'", variation, profile.Name)
		}

		// Sleep for 1 second
		time.Sleep(1 * time.Second)
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, _, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}

	// Sleep for 1 second
	time.Sleep(1 * time.Second)
}
