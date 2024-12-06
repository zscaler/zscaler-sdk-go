package workloadgroups

/*
import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
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
		t.Log("No workload group found. Moving on with other tests.")
	} else {
		name := groups[0].Name
		t.Log("Getting workload group by name:" + name)
		group, err := service.GetByName(name)
		if err != nil {
			t.Errorf("Error getting workload group by name: %v", err)
			return
		}
		if group.Name != name {
			t.Errorf("workload group name does not match: expected %s, got %s", name, group.Name)
		}
	}

	nonExistentName := "ThisWorkloadGroupDoesNotExist"
	_, err = service.GetByName(nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
	} else {
		t.Log("Correctly received error when attempting to get non-existent workload group")
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
		t.Errorf("Error getting workload groups: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Log("No workload group found. Skipping validation.")
		return
	}

	for _, group := range groups {
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

	knownName := "BD_WORKLOAD_GROUP01"
	_, err = service.GetByName(knownName)
	if err != nil {
		t.Logf("Known workload group '%s' does not exist. Skipping case sensitivity tests.", knownName)
		return
	}

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

		if group.Name != knownName {
			t.Errorf("Expected group name to be '%s' for variation '%s', but got '%s'", knownName, variation, group.Name)
		}
	}
}
*/
