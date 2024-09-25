package appservicegroups

/*
import (
	"log"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestAppServiceGroups_data(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	appServices, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting app service groups: %v", err)
		return
	}
	if len(appServices) == 0 {
		t.Errorf("No app service groups found")
		return
	}
	appServiceName := appServices[0].Name
	t.Log("Getting app service group by name:" + appServiceName)
	appService, err := GetByName(context.Background(), service, appServiceName)
	if err != nil {
		t.Errorf("Error getting app service groups by name: %v", err)
		return
	}
	if appService.Name != appServiceName {
		t.Errorf("app service group name does not match: expected %s, got %s", appServiceName, appService.Name)
		return
	}
	// Negative Test: Try to retrieve a app service group with a non-existent name
	nonExistentName := "ThisAppServiceGroupDoesNotExist"
	_, err = GetByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}

	appServices, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting app service group : %v", err)
		return
	}
	if len(appServices) == 0 {
		t.Errorf("No app service group  found")
		return
	}

	// Validate app service group
	for _, appService := range appServices {
		// Checking if essential fields are not empty
		if appService.ID == 0 {
			t.Errorf("app service group  ID is empty")
		}
		if appService.Name == "" {
			t.Errorf("app service group Name is empty")
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}

	// Assuming a service with the name "ZOOM" exists
	knownName := "ZOOM"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve group with name variation: %s", variation)
		group, err := GetByName(context.Background(), service, variation)
		if err != nil {
			t.Errorf("Error getting service with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if group.Name != knownName {
			t.Errorf("Expected role name to be '%s' for variation '%s', but got '%s'", knownName, variation, group.Name)
		}
	}
}
*/
