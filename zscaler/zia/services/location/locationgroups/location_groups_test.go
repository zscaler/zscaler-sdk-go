package locationgroups

import (
	"context"
	"fmt"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestLocationGroups(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "locationgroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Test resources retrieval
	resources, err := GetAll(context.Background(), service, nil)
	if err != nil {
		t.Errorf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		return
	}
	name := resources[0].Name
	resourceByName, err := GetLocationGroupByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}

	_, err = GetLocationGroup(context.Background(), service, resourceByName.ID)
	if err != nil {
		t.Errorf("expected resource to exist: %v", err)
	}

	// Test GetGroupType with both DYNAMIC_GROUP and STATIC_GROUP
	groupTypes := []string{"DYNAMIC_GROUP", "STATIC_GROUP"}
	for _, gType := range groupTypes {
		t.Run(fmt.Sprintf("GroupType-%s", gType), func(t *testing.T) {
			resourceByType, err := GetGroupType(context.Background(), service, gType)
			if err != nil {
				t.Errorf("Error retrieving resource by type '%s': %v", gType, err)
			} else {
				_, err = GetLocationGroup(context.Background(), service, resourceByType.ID)
				if err != nil {
					t.Errorf("expected resource to exist for group type '%s': %v", gType, err)
				}
			}
		})
	}
}

func TestLocationGroupCount(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "locationgroups", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	ctx := context.Background()

	// Test 1: Count without any filters (nil options)
	t.Run("CountWithoutFilters", func(t *testing.T) {
		count, err := GetLocationGroupCount(ctx, service, nil)
		if err != nil {
			t.Errorf("Error retrieving location group count without filters: %v", err)
			return
		}
		if count < 0 {
			t.Errorf("Expected count to be >= 0, got: %d", count)
		}
		t.Logf("Total location group count: %d", count)
	})

	// Test 2: Count with Name filter
	t.Run("CountWithNameFilter", func(t *testing.T) {
		// First get a resource to use its name
		resources, err := GetAll(ctx, service, nil)
		if err != nil || len(resources) == 0 {
			t.Skip("Skipping test: no resources available to test name filter")
			return
		}
		testName := resources[0].Name
		nameFilter := &GetAllFilterOptions{
			Name: &testName,
		}
		count, err := GetLocationGroupCount(ctx, service, nameFilter)
		if err != nil {
			t.Errorf("Error retrieving location group count with name filter: %v", err)
			return
		}
		if count < 0 {
			t.Errorf("Expected count to be >= 0, got: %d", count)
		}
		t.Logf("Location group count with name filter '%s': %d", testName, count)
	})

	// Test 3: Count with GroupType filter
	t.Run("CountWithGroupTypeFilter", func(t *testing.T) {
		groupType := "STATIC_GROUP"
		groupTypeFilter := &GetAllFilterOptions{
			GroupType: &groupType,
		}
		count, err := GetLocationGroupCount(ctx, service, groupTypeFilter)
		if err != nil {
			t.Errorf("Error retrieving location group count with groupType filter: %v", err)
			return
		}
		if count < 0 {
			t.Errorf("Expected count to be >= 0, got: %d", count)
		}
		t.Logf("Location group count with groupType filter '%s': %d", groupType, count)
	})

	// Test 4: Count with FetchLocations filter
	t.Run("CountWithFetchLocationsFilter", func(t *testing.T) {
		fetchLocations := false
		fetchLocationsFilter := &GetAllFilterOptions{
			FetchLocations: &fetchLocations,
		}
		count, err := GetLocationGroupCount(ctx, service, fetchLocationsFilter)
		if err != nil {
			t.Errorf("Error retrieving location group count with fetchLocations filter: %v", err)
			return
		}
		if count < 0 {
			t.Errorf("Expected count to be >= 0, got: %d", count)
		}
		t.Logf("Location group count with fetchLocations=false: %d", count)
	})

	// Test 5: Count with Comments filter
	t.Run("CountWithCommentsFilter", func(t *testing.T) {
		comments := "test"
		commentsFilter := &GetAllFilterOptions{
			Comments: &comments,
		}
		count, err := GetLocationGroupCount(ctx, service, commentsFilter)
		if err != nil {
			t.Errorf("Error retrieving location group count with comments filter: %v", err)
			return
		}
		if count < 0 {
			t.Errorf("Expected count to be >= 0, got: %d", count)
		}
		t.Logf("Location group count with comments filter '%s': %d", comments, count)
	})

	// Test 6: Count with LastModUser filter
	t.Run("CountWithLastModUserFilter", func(t *testing.T) {
		lastModUser := "admin"
		lastModUserFilter := &GetAllFilterOptions{
			LastModUser: &lastModUser,
		}
		count, err := GetLocationGroupCount(ctx, service, lastModUserFilter)
		if err != nil {
			t.Errorf("Error retrieving location group count with lastModUser filter: %v", err)
			return
		}
		if count < 0 {
			t.Errorf("Expected count to be >= 0, got: %d", count)
		}
		t.Logf("Location group count with lastModUser filter '%s': %d", lastModUser, count)
	})

	// Test 7: Count with LocationID filter
	t.Run("CountWithLocationIDFilter", func(t *testing.T) {
		// First get a resource to use its location ID if available
		resources, err := GetAll(ctx, service, nil)
		if err != nil || len(resources) == 0 {
			t.Skip("Skipping test: no resources available to test locationID filter")
			return
		}
		// Use a test location ID (0 if not available)
		locationID := 0
		if len(resources[0].Locations) > 0 {
			locationID = resources[0].Locations[0].ID
		}
		locationIDFilter := &GetAllFilterOptions{
			LocationID: &locationID,
		}
		count, err := GetLocationGroupCount(ctx, service, locationIDFilter)
		if err != nil {
			t.Errorf("Error retrieving location group count with locationID filter: %v", err)
			return
		}
		if count < 0 {
			t.Errorf("Expected count to be >= 0, got: %d", count)
		}
		t.Logf("Location group count with locationID filter %d: %d", locationID, count)
	})

	// Test 8: Count with Version filter
	t.Run("CountWithVersionFilter", func(t *testing.T) {
		version := 1
		versionFilter := &GetAllFilterOptions{
			Version: &version,
		}
		count, err := GetLocationGroupCount(ctx, service, versionFilter)
		if err != nil {
			t.Errorf("Error retrieving location group count with version filter: %v", err)
			return
		}
		if count < 0 {
			t.Errorf("Expected count to be >= 0, got: %d", count)
		}
		t.Logf("Location group count with version filter %d: %d", version, count)
	})

	// Test 9: Count with multiple filters combined
	t.Run("CountWithMultipleFilters", func(t *testing.T) {
		fetchLocations := false
		groupType := "STATIC_GROUP"
		combinedFilter := &GetAllFilterOptions{
			GroupType:      &groupType,
			FetchLocations: &fetchLocations,
		}
		count, err := GetLocationGroupCount(ctx, service, combinedFilter)
		if err != nil {
			t.Errorf("Error retrieving location group count with multiple filters: %v", err)
			return
		}
		if count < 0 {
			t.Errorf("Expected count to be >= 0, got: %d", count)
		}
		t.Logf("Location group count with multiple filters (groupType=%s, fetchLocations=%v): %d", groupType, fetchLocations, count)
	})
}
