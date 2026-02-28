package groups

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestGroups(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	// Create new resource
	createdResource, _, err := Create(context.Background(), service, &Groups{
		Name:                      name,
		Description:               name,
		Source:                    "SCIM",
		AdminEntitlementEnabled:   true,
		ServiceEntitlementEnabled: true,
		DynamicGroup:              true,
	})
	if err != nil {
		t.Fatalf("Error creating resource: %v", err)
	}

	t.Run("TestResourceCreation", func(t *testing.T) {
		if createdResource.ID == "" {
			t.Error("Expected created resource ID to be non-empty, but got ''")
		}
		if createdResource.Name != name {
			t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
		}
	})

	t.Run("TestResourceRetrieval", func(t *testing.T) {
		retrievedResource, err := Get(context.Background(), service, createdResource.ID)
		if err != nil {
			t.Fatalf("Error retrieving resource: %v", err)
		}
		if retrievedResource.ID != createdResource.ID {
			t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
		}
		if retrievedResource.Name != name {
			t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, createdResource.Name)
		}
	})

	t.Run("TestResourceUpdate", func(t *testing.T) {
		updatedResource := *createdResource
		updatedResource.Name = updateName
		idInt, err := strconv.Atoi(createdResource.ID)
		if err != nil {
			// If conversion fails, skip this test since Update expects int
			t.Skipf("Skipping test - ID '%s' cannot be converted to int: %v", createdResource.ID, err)
		}
		_, _, err = Update(context.Background(), service, idInt, &updatedResource)
		if err != nil {
			t.Fatalf("Error updating resource: %v", err)
		}
	})

	t.Run("TestResourceRetrievalByName", func(t *testing.T) {
		retrievedResources, err := GetByName(context.Background(), service, name)
		if err != nil {
			t.Fatalf("Error retrieving resource by name: %v", err)
		}
		if len(retrievedResources) == 0 {
			t.Fatalf("Expected at least one resource to be retrieved by name, but got none")
		}
		retrievedResource := retrievedResources[0]
		if retrievedResource.ID != createdResource.ID {
			t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
		}
		if retrievedResource.Name != name {
			t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, retrievedResource.Name)
		}
	})

	t.Run("TestAllResourcesRetrieval", func(t *testing.T) {
		// Pass nil or &common.PaginationQueryParams{} as the third argument
		resources, err := GetAll(context.Background(), service, nil)
		if err != nil {
			t.Fatalf("Error retrieving groups: %v", err)
		}
		if len(resources) == 0 {
			t.Error("Expected retrieved resources to be non-empty, but got empty slice")
		}
		found := false
		for _, resource := range resources {
			if resource.ID == createdResource.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected retrieved groups to contain created resource '%s', but it didn't", createdResource.ID)
		}
	})

	t.Run("TestResourceRemoval", func(t *testing.T) {
		_, err := Delete(context.Background(), service, createdResource.ID)
		if err != nil {
			t.Fatalf("Error deleting resource: %v", err)
		}
	})
}
