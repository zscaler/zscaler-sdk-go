package cbiregions

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestGetAllRegions(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Failed to create ZPA client: %v", err)
	}
	service := New(client)

	// 1. First GetAll regions and ensure a response is returned.
	regions, resp, err := service.GetAll()
	if err != nil || resp.StatusCode >= 400 || len(regions) == 0 {
		t.Fatalf("Failed to fetch regions: %v", err)
	}

	// To ensure that subsequent tests don't fail due to an empty regions list
	if len(regions) == 0 {
		t.Fatal("No regions returned. Can't proceed with further tests.")
		return
	}

	// 3. Test the GetByName method by querying the Name of any of the returned regions from GetAll.
	firstRegionName := regions[0].Name
	singleRegionByName, resp, err := service.GetByName(firstRegionName)
	if err != nil || resp.StatusCode >= 400 || singleRegionByName == nil {
		t.Errorf("Failed to fetch region by Name %s: %v", firstRegionName, err)
	} else if singleRegionByName.Name != firstRegionName {
		t.Errorf("Mismatch in region Name. Expected %s, got %s", firstRegionName, singleRegionByName.Name)
	}
}
