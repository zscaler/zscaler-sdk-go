package cbizpaprofile

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
)

func TestCBIZPAProfiles(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	profiles, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting cbi zpa profile: %v", err)
		return
	}
	if len(profiles) == 0 {
		t.Errorf("No cbi zpa profiles found")
		return
	}
	name := profiles[0].Name
	t.Log("Getting cbi zpa profiles by name:" + name)
	rg, _, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting cbi zpa profile by name: %v", err)
		return
	}
	if rg.Name != name {
		t.Errorf("cbi region name does not match: expected %s, got %s", name, rg.Name)
		return
	}
}
