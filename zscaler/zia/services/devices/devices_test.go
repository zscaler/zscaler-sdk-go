package devices

import (
	"context"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

// TestDevicesGetAll validates the unfiltered list endpoint and, when data is
// available, exercises Get and GetByName using a record returned by GetAll.
func TestDevicesGetAll(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	devices, err := GetAll(ctx, service, nil)
	if err != nil {
		t.Fatalf("Error retrieving devices: %v", err)
	}

	if len(devices) == 0 {
		t.Log("No devices returned; skipping Get/GetByName assertions")
		return
	}

	first := devices[0]
	if first.ID == 0 {
		t.Errorf("Expected device to have a non-zero ID")
	}

	// Get by ID.
	byID, err := Get(ctx, service, first.ID)
	if err != nil {
		t.Fatalf("Error retrieving device by ID %d: %v", first.ID, err)
	}
	if byID == nil || byID.ID != first.ID {
		t.Fatalf("Expected device with ID %d, got %+v", first.ID, byID)
	}

	// Get by name (only when the device exposes a name).
	if first.Name == "" {
		t.Log("First device has no name; skipping GetByName assertion")
		return
	}
	byName, err := GetByName(ctx, service, first.Name)
	if err != nil {
		t.Fatalf("Error retrieving device by name '%s': %v", first.Name, err)
	}
	if byName == nil || !strings.EqualFold(byName.Name, first.Name) {
		t.Errorf("Expected device name '%s', got %+v", first.Name, byName)
	}
}

// TestDevicesGetAllWithFilters exercises the optional filtering parameters
// (id, search, valid, includeCbiDevices) against a live tenant. The filters are
// derived from the data returned by the unfiltered call.
func TestDevicesGetAllWithFilters(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	all, err := GetAll(ctx, service, nil)
	if err != nil {
		t.Fatalf("Error retrieving devices: %v", err)
	}

	if len(all) == 0 {
		t.Log("No devices returned; skipping filter assertions")
		return
	}

	// Filter by ID using the first returned device.
	idFiltered, err := GetAll(ctx, service, &GetAllFilterOptions{ID: []int{all[0].ID}})
	if err != nil {
		t.Fatalf("Error retrieving devices with id filter: %v", err)
	}
	for _, d := range idFiltered {
		if d.ID != all[0].ID {
			t.Errorf("Expected only device ID %d, got %d", all[0].ID, d.ID)
		}
	}

	// Filter by valid devices.
	valid := true
	if _, err := GetAll(ctx, service, &GetAllFilterOptions{Valid: &valid}); err != nil {
		t.Fatalf("Error retrieving devices with valid filter: %v", err)
	}

	// Filter including CBI devices.
	includeCbi := true
	if _, err := GetAll(ctx, service, &GetAllFilterOptions{IncludeCbiDevices: &includeCbi}); err != nil {
		t.Fatalf("Error retrieving devices with includeCbiDevices filter: %v", err)
	}

	// Filter by search using the first device name, if present.
	if all[0].Name != "" {
		search := all[0].Name
		if _, err := GetAll(ctx, service, &GetAllFilterOptions{Search: &search}); err != nil {
			t.Fatalf("Error retrieving devices with search filter: %v", err)
		}
	}
}

func TestDevicesGetByNameNonExistent(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = GetByName(context.Background(), service, "non_existent_device_name")
	if err == nil {
		t.Error("Expected error retrieving device by non-existent name, but got nil")
	}
}

func TestDevicesGetNonExistentID(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = Get(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error retrieving device with non-existent ID, but got nil")
	}
}
