package locationlite

import (
	"context"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestLocationLite_data(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	servers, err := GetAll(context.Background(), service, nil)
	if err != nil {
		t.Errorf("Error getting locations: %v", err)
		return
	}
	if len(servers) == 0 {
		t.Errorf("No location lite found")
		return
	}
	name := servers[0].Name
	t.Log("Getting location lite by name:" + name)
	server, err := GetLocationLiteByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error getting location lite by name: %v", err)
		return
	}
	if server.Name != name {
		t.Errorf("location lite name does not match: expected %s, got %s", name, server.Name)
		return
	}
	// Negative Test: Try to retrieve an location lite with a non-existent name
	nonExistentName := "ThisLocationDoesNotExist"
	_, err = GetLocationLiteByName(context.Background(), service, nonExistentName)
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
	servers, err := GetAll(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error getting all location lites: %v", err)
	}
	if len(servers) == 0 {
		t.Fatalf("No location lites found for testing")
	}

	// Choose the first server's ID for testing
	testID := servers[0].ID

	// Retrieve the server by ID
	server, err := GetLocationLiteID(context.Background(), service, testID)
	if err != nil {
		t.Errorf("Error retrieving location lite with ID %d: %v", testID, err)
		return
	}

	// Verify the retrieved server
	if server == nil {
		t.Errorf("No server returned for ID %d", testID)
		return
	}

	if server.ID != testID {
		t.Errorf("Retrieved server ID mismatch: expected %d, got %d", testID, server.ID)
	}
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	locations, err := GetAll(context.Background(), service, nil)
	if err != nil {
		t.Errorf("Error getting location lite: %v", err)
		return
	}
	if len(locations) == 0 {
		t.Errorf("No location lite found")
		return
	}

	// Validate location lite
	for _, location := range locations {
		// Checking if essential fields are not empty
		if location.ID == 0 {
			t.Errorf("location lite ID is empty")
		}
		if location.Name == "" {
			t.Errorf("location lite Name is empty")
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Assuming a group with the name "Road Warrior" exists
	knownName := "Road Warrior"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		group, err := GetLocationLiteByName(context.Background(), service, variation)
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
