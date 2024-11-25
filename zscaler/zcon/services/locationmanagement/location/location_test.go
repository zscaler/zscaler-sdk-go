package location

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcon/services"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestZConLocation(t *testing.T) {
	client, err := tests.NewZConClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	locations, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting locations: %v", err)
		return
	}
	if len(locations) == 0 {
		t.Errorf("No locations found")
		return
	}
	name := locations[0].Name
	t.Log("Getting locations by name:" + name)
	location, err := GetLocationByName(service, name)
	if err != nil {
		t.Errorf("Error getting location by name: %v", err)
		return
	}
	if location.Name != name {
		t.Errorf("admin locations does not match: expected %s, got %s", name, location.Name)
		return
	}

	locationName, err := GetLocationByName(service, name)
	if err != nil {
		t.Errorf("Error getting admin roles by name: %v", err)
		return
	}
	if locationName.Name != name {
		t.Errorf("admin roles name does not match: expected %s, got %s", name, locationName.Name)
		return
	}
	// Negative Test: Try to retrieve a Location with a non-existent name
	nonExistentName := "ThisLocationNotExist"
	_, err = GetLocationByName(service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZConClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	locations, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting location: %v", err)
		return
	}
	if len(locations) == 0 {
		t.Errorf("No machine location found")
		return
	}

	// Validate each group
	for _, location := range locations {
		// Checking if essential fields are not empty
		if location.ID == 0 {
			t.Errorf("Location ID is empty")
		}
		if location.Name == "" {
			t.Errorf("LocationName is empty")
		}
		if location.Country == "" {
			t.Errorf("Location Description is empty")
		}
		if location.TZ == "" {
			t.Errorf("Location Description is empty")
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZConClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	// Assuming a group with the name "BD-MGR01" exists
	knownName := "BD_CC01_US"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		group, err := GetLocationByName(service, variation)
		if err != nil {
			t.Errorf("Error getting machine group with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if group.Name != knownName {
			t.Errorf("Expected group name to be '%s' for variation '%s', but got '%s'", knownName, variation, group.Name)
		}
	}
}
