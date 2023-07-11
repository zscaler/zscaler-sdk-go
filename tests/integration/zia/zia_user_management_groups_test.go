package integration

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zia/services/usermanagement"
)

func TestAccGroupManagement(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := usermanagement.New(client)

	groups, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting groups: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No groups found")
		return
	}
	name := groups[0].Name
	t.Log("Getting groups by name:" + name)
	group, err := service.GetGroupByName(name)
	if err != nil {
		t.Errorf("Error getting groups by name: %v", err)
		return
	}
	if group.Name != name {
		t.Errorf("group name does not match: expected %s, got %s", name, group.Name)
		return
	}
}
