package cbizpaprofile

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
)

func TestCBIZPAProfile(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	profiles, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting isolation profiles: %v", err)
		return
	}
	if len(profiles) == 0 {
		t.Errorf("No isolation profile found")
		return
	}
	name := profiles[0].Name
	t.Log("Getting isolation profile by name:" + name)
	profile, _, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting isolation profile by name: %v", err)
		return
	}
	if profile.Name != name {
		t.Errorf("isolation profile name does not match: expected %s, got %s", name, profile.Name)
		return
	}
}
