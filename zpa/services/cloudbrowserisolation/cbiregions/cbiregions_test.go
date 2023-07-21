package cbiregions

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
)

func TestCBIRegions(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	regions, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error getting cbi regions: %v", err)
		return
	}
	if len(regions) == 0 {
		t.Errorf("No cbi regions found")
		return
	}
	name := regions[0].Name
	t.Log("Getting cbi region by name:" + name)
	rg, _, err := service.GetByName(name)
	if err != nil {
		t.Errorf("Error getting cbi region by name: %v", err)
		return
	}
	if rg.Name != name {
		t.Errorf("cbi region name does not match: expected %s, got %s", name, rg.Name)
		return
	}
}
