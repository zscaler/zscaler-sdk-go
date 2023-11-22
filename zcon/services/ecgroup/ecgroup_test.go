package ecgroup

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestECGroup(t *testing.T) {
	client, err := tests.NewZConClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	ecgroups, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting Cloud & Branch Connector Group: %v", err)
		return
	}
	if len(ecgroups) == 0 {
		t.Errorf("No Cloud & Branch Connector Group found")
		return
	}
	name := ecgroups[0].Name
	t.Log("Getting Cloud & Branch Connector Group by name:" + name)
	ecgroup, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting Cloud & Branch Connector Group by name: %v", err)
		return
	}
	if ecgroup.Name != name {
		t.Errorf("Cloud & Branch Connector Group name does not match: expected %s, got %s", name, ecgroup.Name)
		return
	}

	ecgroupLite, err := service.GetEcGroupLiteByName(name)
	if err != nil {
		t.Errorf("Error getting Cloud & Branch Connector Group by name: %v", err)
		return
	}
	if ecgroupLite.Name != name {
		t.Errorf("Cloud & Branch Connector Group name does not match: expected %s, got %s", name, ecgroupLite.Name)
		return
	}
}
