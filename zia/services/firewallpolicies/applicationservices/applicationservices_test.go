package applicationservices

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestApplicationServices_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	appServices, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting application services: %v", err)
		return
	}
	if len(appServices) == 0 {
		t.Errorf("No application service found")
		return
	}
	appServiceName := appServices[0].Name
	t.Log("Getting application service by name:" + appServiceName)
	appService, err := service.GetByName(appServiceName)
	if err != nil {
		t.Errorf("Error getting application service by name: %v", err)
		return
	}
	if appService.Name != appServiceName {
		t.Errorf("application service name does not match: expected %s, got %s", appServiceName, appService.Name)
		return
	}
	// Negative Test: Try to retrieve a application service with a non-existent name
	nonExistentName := "ThisApplicationServiceDoesNotExist"
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

	appServices, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting application service : %v", err)
		return
	}
	if len(appServices) == 0 {
		t.Errorf("No application service found")
		return
	}

	// Validate time window
	for _, appService := range appServices {
		// Checking if essential fields are not empty
		if appService.ID == 0 {
			t.Errorf("application service ID is empty")
		}
		if appService.Name == "" {
			t.Errorf("application service Name is empty")
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

	// Assuming a service with the name "SKYPEFORBUSINESS" exists
	knownName := "SKYPEFORBUSINESS"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve service with name variation: %s", variation)
		service, err := service.GetByName(variation)
		if err != nil {
			t.Errorf("Error getting service with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if service.Name != knownName {
			t.Errorf("Expected role name to be '%s' for variation '%s', but got '%s'", knownName, variation, service.Name)
		}
	}
}
