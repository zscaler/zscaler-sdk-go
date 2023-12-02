package appservicegroups

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestAppServiceGroups_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	appServices, err := service.GetAll()
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
	appService, err := service.GetByName(appServiceName)
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
