package integration

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/postureprofile"
)

func TestPostureProfiles(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := postureprofile.New(client)

	profiles, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting posture profiles: %v", err)
		return
	}
	if len(profiles) == 0 {
		t.Errorf("No posture profiles found")
		return
	}
	name := profiles[0].Name
	t.Log("Getting posture profile by name:" + name)
	net, _, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting posture profile by name: %v", err)
		return
	}
	if net.Name != name {
		t.Errorf("Posture profile name does not match: expected %s, got %s", name, net.Name)
		return
	}
}
