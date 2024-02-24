package postureprofile

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestPostureProfiles(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Test to retrieve all profiles
	profiles, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting posture profiles: %v", err)
		return
	}
	if len(profiles) == 0 {
		t.Errorf("No posture profile found")
		return
	}

	// Test to retrieve a profile by its name
	name := profiles[0].Name
	adaptedName := common.RemoveCloudSuffix(name)
	t.Log("Getting posture profile by name:" + adaptedName)
	profile, _, err := service.GetByName(adaptedName)
	if err != nil {
		t.Errorf("Error getting posture profile by name: %v", err)
		return
	}
	if common.RemoveCloudSuffix(profile.Name) != adaptedName {
		t.Errorf("posture profile name does not match: expected %s, got %s", adaptedName, profile.Name)
		return
	}

	// Negative Test: Try to retrieve a profile with a non-existent name
	nonExistentName := "ThisPostureProfileNameDoesNotExist"
	_, _, err = service.GetByName(nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	profiles, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting posture profiles: %v", err)
		return
	}
	if len(profiles) == 0 {
		t.Errorf("No posture profile found")
		return
	}

	// Validate each profile
	for _, profile := range profiles {
		// Checking if essential fields are not empty
		if profile.ID == "" {
			t.Errorf("Posture Profile ID is empty")
		}
		if profile.Name == "" {
			t.Errorf("Posture Profile Name is empty")
		}
		if profile.PostureudID == "" {
			t.Errorf("Posture Profile UDID is empty")
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Assuming a profile with the name "CrowdStrike_ZPA_ZTA_40" exists
	knownName := "CrowdStrike_ZPA_ZTA_40"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve profile with name variation: %s", variation)
		profile, _, err := service.GetByName(variation)
		if err != nil {
			t.Errorf("Error getting posture profile with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the profile's actual name matches the known name
		if common.RemoveCloudSuffix(profile.Name) != knownName {
			t.Errorf("Expected posture profile name to be '%s' for variation '%s', but got '%s'", knownName, variation, profile.Name)
		}
	}
}

/*
	func TestPostureProfileNamesWithSpaces(t *testing.T) {
		client, err := tests.NewZpaClient()
		if err != nil {
			t.Errorf("Error creating client: %v", err)
			return
		}

		service := New(client)

		// Assuming that there are profiles with the following name variations
		variations := []string{
			"CrowdStrike ZPA ZTA 40", "CrowdStrike  ZPAZTA  40", "CrowdStrike   ZPAZTA   40",
			"CrowdStrike    ZPAZTA40", "CrowdStrike  ZPAZTA 40", "CrowdStrike  ZPA ZTA   40",
			"CrowdStrike   ZPA   ZTA 40",
		}

		for _, variation := range variations {
			t.Logf("Attempting to retrieve profile with name: %s", variation)
			profile, _, err := service.GetByName(variation)
			if err != nil {
				t.Errorf("Error getting posture profile with name '%s': %v", variation, err)
				continue
			}

			// Verify if the profile's actual name matches the expected variation
			if common.RemoveCloudSuffix(profile.Name) != variation {
				t.Errorf("Expected posture profile name to be '%s' but got '%s'", variation, profile.Name)
			}
		}
	}
*/
func TestPostureProfileByPostureUDID(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Use GetByName to fetch a known Posture Profile
	knownName := "CrowdStrike_ZPA_ZTA_40"
	posture, _, err := service.GetByName(knownName)
	if err != nil || posture == nil {
		t.Errorf("Error getting posture profile with name '%s': %v", knownName, err)
		return
	}

	// Use the PostureudID from the above posture profile to test GetByPostureUDID
	t.Logf("Attempting to retrieve posture with PostureudID: %s", posture.PostureudID)
	postureByUDID, _, err := service.GetByPostureUDID(posture.PostureudID)
	if err != nil {
		t.Errorf("Error getting posture profile with PostureudID '%s': %v", posture.PostureudID, err)
		return
	}

	// Check if the posture profile's actual PostureudID matches the known PostureudID
	if postureByUDID.PostureudID != posture.PostureudID {
		t.Errorf("Expected posture profile UDID to be '%s', but got '%s'", posture.PostureudID, postureByUDID.PostureudID)
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.GetByName("non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
