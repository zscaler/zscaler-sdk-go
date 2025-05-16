package machinegroup

import (
	"context"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestMachineGroup(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	// Test to retrieve all machine groups
	groups, _, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting machine groups: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No machine group found")
		return
	}

	// Test to retrieve a group by its name
	name := groups[0].Name
	t.Log("Getting machine group by name:" + name)
	group, _, err := GetByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error getting machine group by name: %v", err)
		return
	}
	if group.Name != name {
		t.Errorf("Machine group name does not match: expected %s, got %s", name, group.Name)
		return
	}

	// Additional step: Use the ID of the first machine group to test the Get function
	t.Log("Getting machine group by ID:" + groups[0].ID)
	groupByID, _, err := Get(context.Background(), service, groups[0].ID)
	if err != nil {
		t.Errorf("Error getting machine group by ID: %v", err)
		return
	}
	if groupByID.ID != groups[0].ID {
		t.Errorf("Machine group ID does not match: expected %s, got %s", groups[0].ID, groupByID.ID)
		return
	}

	// Negative Test: Try to retrieve a group with a non-existent name
	nonExistentName := "ThisMachineGroupNameDoesNotExist"
	_, _, err = GetByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}

	// Negative Test: Try to retrieve a group with a non-existent ID
	nonExistentID := "non_existent_id"
	_, _, err = Get(context.Background(), service, nonExistentID)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent ID, got nil")
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

	groups, _, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting machine groups: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No machine group found")
		return
	}

	// Validate each group
	for _, group := range groups {
		// Checking if essential fields are not empty
		if group.ID == "" {
			t.Errorf("Machine Group ID is empty")
		}
		if group.Name == "" {
			t.Errorf("Machine Group Name is empty")
		}
	}
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
	// Assuming a group with the name "BD-MGR01" exists
	knownName := "BD-MGR01"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		group, _, err := GetByName(context.Background(), service, variation)
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

func TestMachineGroupNamesWithSpaces(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	// Assuming that there are groups with the following name variations
	variations := []string{
		"BD-MGR01", "BD-MGR02", "BD MGR 03", "BD  MGR  04", "BD   MGR   05",
		"BD    MGR06", "BD  MGR 07", "BD  M GR   08", "BD   M   GR 09",
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name: %s", variation)
		group, _, err := GetByName(context.Background(), service, variation)
		if err != nil {
			t.Errorf("Error getting machine group with name '%s': %v", variation, err)
			continue
		}

		// Verify if the group's actual name matches the expected variation
		if group.Name != variation {
			t.Errorf("Expected group name to be '%s' but got '%s'", variation, group.Name)
		}
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
		t.Error("Expected error retrieving resource by non-existent_name name, but got nil")
	}
}
