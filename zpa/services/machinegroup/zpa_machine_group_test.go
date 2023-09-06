package machinegroup

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestMachineGroup(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	groups, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting machine groups: %v", err)
		return
	}
	if len(groups) == 0 {
		t.Errorf("No machine groups found")
		return
	}
	name := groups[0].Name
	t.Log("Getting machine group by name:" + name)
	group, _, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting machine group by name: %v", err)
		return
	}
	if group.Name != name {
		t.Errorf("machine group name does not match: expected %s, got %s", name, group.Name)
		return
	}
}
