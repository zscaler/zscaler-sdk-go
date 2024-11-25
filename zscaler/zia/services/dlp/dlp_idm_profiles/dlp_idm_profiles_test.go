package dlp_idm_profiles

import (
	"context"
	"strings"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestDLPIDMProfile_data(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	profiles, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting idm profiles: %v", err)
		return
	}
	if len(profiles) == 0 {
		t.Errorf("No idm profile found")
		return
	}
	name := profiles[0].ProfileName
	t.Log("Getting idm profile by name:" + name)
	profile, err := GetByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error getting idm profile by name: %v", err)
		return
	}
	if profile.ProfileName != name {
		t.Errorf("idm profile name does not match: expected %s, got %s", name, profile.ProfileName)
		return
	}
	// Negative Test: Try to retrieve an idm profile with a non-existent name
	nonExistentName := "ThisIDMProfileDoesNotExist"
	_, err = GetByName(context.Background(), service, nonExistentName)
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
	profiles, err := GetAll(context.Background(), service)
	if err != nil {
		t.Fatalf("Error getting all idm profiles: %v", err)
	}
	if len(profiles) == 0 {
		t.Fatalf("No idm profiles found for testing")
	}

	// Choose the first server's ID for testing
	testID := profiles[0].ProfileID

	// Retrieve the server by ID
	profile, err := Get(context.Background(), service, testID)
	if err != nil {
		t.Errorf("Error retrieving idm profile with ID %d: %v", testID, err)
		return
	}

	// Verify the retrieved server
	if profile == nil {
		t.Errorf("No server returned for ID %d", testID)
		return
	}

	if profile.ProfileID != testID {
		t.Errorf("Retrieved server ID mismatch: expected %d, got %d", testID, profile.ProfileID)
	}
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	profiles, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting idm profile: %v", err)
		return
	}
	if len(profiles) == 0 {
		t.Errorf("No idm profile found")
		return
	}

	// Validate idm profile
	for _, profile := range profiles {
		// Checking if essential fields are not empty
		if profile.ProfileID == 0 {
			t.Errorf("idm profile ID is empty")
		}
		if profile.ProfileName == "" {
			t.Errorf("idm profile Name is empty")
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Assuming a idm with the name "BD_IDM_TEMPLATE01" exists
	knownName := "BD_IDM_TEMPLATE01"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		profile, err := GetByName(context.Background(), service, variation)
		if err != nil {
			t.Errorf("Error getting idm profile with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if profile.ProfileName != knownName {
			t.Errorf("Expected group name to be '%s' for variation '%s', but got '%s'", knownName, variation, profile.ProfileName)
		}
	}
}
