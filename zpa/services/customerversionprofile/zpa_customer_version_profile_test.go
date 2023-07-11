package customerversionprofile

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
)

func TestCustomerVersionProfile(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	profiles, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting customer version profiles: %v", err)
		return
	}
	if len(profiles) == 0 {
		t.Errorf("No customer version profile found")
		return
	}
	name := profiles[0].Name
	t.Log("Getting customer version profile by name:" + name)
	profile, _, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting customer version profile by name: %v", err)
		return
	}
	if profile.Name != name {
		t.Errorf("customer version profile name does not match: expected %s, got %s", name, profile.Name)
		return
	}
}
