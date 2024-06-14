package ecgroup

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zcon/services"
)

func TestECGroup(t *testing.T) {
	client, err := tests.NewZConClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	ecgroups, err := GetAll(service)
	if err != nil {
		t.Errorf("Error getting Cloud & Branch Connector Group: %v", err)
		return
	}
	if len(ecgroups) == 0 {
		t.Log("No Cloud & Branch Connector Group found. Moving on with other tests.")
	} else {
		// Proceed with tests that require at least one EC Group
		name := ecgroups[0].Name
		t.Log("Getting Cloud & Branch Connector Group by name: " + name)
		ecgroup, err := GetByName(service, name)
		if err != nil {
			t.Errorf("Error getting Cloud & Branch Connector Group by name: %v", err)
			return
		}
		if ecgroup.Name != name {
			t.Errorf("Cloud & Branch Connector Group name does not match: expected %s, got %s", name, ecgroup.Name)
			return
		}

		ecgroupLite, err := GetEcGroupLiteByName(service, name)
		if err != nil {
			t.Errorf("Error getting Cloud & Branch Connector Group by name: %v", err)
			return
		}
		if ecgroupLite.Name != name {
			t.Errorf("Cloud & Branch Connector Group name does not match: expected %s, got %s", name, ecgroupLite.Name)
			return
		}
	}

	// Negative Test: Try to retrieve a EcGroup with a non-existent name
	nonExistentName := "ThisEcGroupDoesNotExist"
	_, err = GetByName(service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	} else {
		t.Log("Correctly received error when attempting to get non-existent Cloud & Branch Connector Group")
	}
}
