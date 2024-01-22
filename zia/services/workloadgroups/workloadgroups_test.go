package workloadgroups

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestWorkloadGroups(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	groups, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting workload groups: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No workload group found")
		return
	}
	name := groups[0].Name
	t.Log("Getting workload group by name:" + name)
	group, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting workload group by name: %v", err)
		return
	}
	if group.Name != name {
		t.Errorf("workload group name does not match: expected %s, got %s", name, group.Name)
		return
	}
	// Negative Test: Try to retrieve an workload group with a non-existent name
	nonExistentName := "ThisWorkloadGroupNotExist"
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

	groups, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting workload group: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No workload group found")
		return
	}

	// Validate workload group
	for _, group := range groups {
		// Checking if essential fields are not empty
		if group.ID == 0 {
			t.Errorf("workload group ID is empty")
		}
		if group.Name == "" {
			t.Errorf("workload group Name is empty")
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Assuming a workload group with the name "BD_WORKLOAD_GROUP01" exists
	knownName := "BD_WORKLOAD_GROUP01"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		group, err := service.GetByName(variation)
		if err != nil {
			t.Errorf("Error getting workload group with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if group.Name != knownName {
			t.Errorf("Expected group name to be '%s' for variation '%s', but got '%s'", knownName, variation, group.Name)
		}
	}
}
