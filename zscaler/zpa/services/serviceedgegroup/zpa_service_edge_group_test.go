package serviceedgegroup

import (
	"context"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestServiceEdgeGroup_Create(t *testing.T) {
	name := tests.GetTestName("tests-svcedge")
	updateName := tests.GetTestName("tests-svcedge")
	client, err := tests.NewVCRTestClient(t, "serviceedgegroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	// Create new resource
	createdResource, _, err := Create(context.Background(), service, ServiceEdgeGroup{
		Name:                   name,
		Description:            name,
		Enabled:                true,
		Latitude:               "37.33874",
		Longitude:              "-121.8852525",
		Location:               "San Jose, CA, USA",
		IsPublic:               "TRUE",
		UpgradeDay:             "SUNDAY",
		UpgradeTimeInSecs:      "66600",
		OverrideVersionProfile: true,
		VersionProfileName:     "Default",
		VersionProfileID:       "0",
		GraceDistanceEnabled:   true,
		GraceDistanceValue:     "10",
		GraceDistanceValueUnit: "MILES",
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
		retrievedResource, _, err := Get(context.Background(), service, createdResource.ID)
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
		_, err = Update(context.Background(), service, createdResource.ID, &updatedResource)
		if err != nil {
			t.Fatalf("Error updating resource: %v", err)
		}
	})

	t.Run("TestResourceRetrievalByName", func(t *testing.T) {
		retrievedResource, _, err := GetByName(context.Background(), service, updateName)
		if err != nil {
			t.Fatalf("Error retrieving resource by name: %v", err)
		}
		if retrievedResource.ID != createdResource.ID {
			t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
		}
		if retrievedResource.Name != updateName {
			t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, retrievedResource.Name)
		}
	})

	t.Run("TestAllResourcesRetrieval", func(t *testing.T) {
		resources, _, err := GetAll(context.Background(), service)
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

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "serviceedgegroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, _, err = Get(context.Background(), service, "non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "serviceedgegroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, err = Delete(context.Background(), service, "non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "serviceedgegroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, err = Update(context.Background(), service, "non_existent_id", &ServiceEdgeGroup{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "serviceedgegroup", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, _, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
