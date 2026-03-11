package tag_key_controller

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_namespace"
)

func TestTagKeyCRUD(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	namespaceName := "tests-ns-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	tagVal1 := "test-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	tagVal2 := "test-" + acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)

	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// 1. Create namespace first; tag key is created within a namespace
	createdNamespace, _, err := tag_namespace.Create(context.Background(), service, tag_namespace.Namespace{
		Name:        namespaceName,
		Description: namespaceName,
		Enabled:     true,
		Origin:      "CUSTOM",
		Type:        "STATIC",
	})
	if err != nil {
		t.Fatalf("Error creating namespace for tag key test: %v", err)
	}
	defer func() {
		time.Sleep(time.Second * 2)
		_, _, getErr := tag_namespace.Get(context.Background(), service, createdNamespace.ID)
		if getErr != nil {
			t.Logf("Namespace might have already been deleted: %v", getErr)
			return
		}
		_, err := tag_namespace.Delete(context.Background(), service, createdNamespace.ID)
		if err != nil {
			t.Logf("Error deleting namespace: %v", err)
		}
	}()

	namespaceID := createdNamespace.ID

	// 2. Create tag key within the namespace
	createdResource, _, err := Create(context.Background(), service, namespaceID, TagKey{
		Name:        name,
		Description: name,
		Enabled:     true,
		Origin:      "CUSTOM",
		Type:        "STATIC",
		TagValues: []TagValue{
			{Name: tagVal1},
			{Name: tagVal2},
		},
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
		retrievedResource, _, err := Get(context.Background(), service, namespaceID, createdResource.ID)
		if err != nil {
			t.Fatalf("Error retrieving resource: %v", err)
		}
		if retrievedResource.ID != createdResource.ID {
			t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
		}
		if retrievedResource.Name != name {
			t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, retrievedResource.Name)
		}
	})

	t.Run("TestResourceUpdate", func(t *testing.T) {
		updatedResource := *createdResource
		updatedResource.Name = updateName
		_, err = Update(context.Background(), service, namespaceID, createdResource.ID, &updatedResource)
		if err != nil {
			t.Fatalf("Error updating resource: %v", err)
		}
	})

	t.Run("TestResourceRetrievalByName", func(t *testing.T) {
		retrievedResource, _, err := GetByName(context.Background(), service, namespaceID, updateName)
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
		resources, _, err := GetAll(context.Background(), service, namespaceID)
		if err != nil {
			t.Fatalf("Error retrieving resources: %v", err)
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
			t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.ID)
		}
	})

	t.Run("TestBulkUpdateStatus", func(t *testing.T) {
		_, err := BulkUpdateStatus(context.Background(), service, namespaceID, BulkUpdateStatusRequest{
			Enabled:   false,
			TagKeyIDs: []string{createdResource.ID},
		})
		if err != nil {
			t.Fatalf("Error bulk updating status: %v", err)
		}

		retrievedResource, _, err := Get(context.Background(), service, namespaceID, createdResource.ID)
		if err != nil {
			t.Fatalf("Error retrieving resource after bulk update: %v", err)
		}
		if retrievedResource.Enabled {
			t.Errorf("Expected tag key to be disabled after bulk update, but it is still enabled")
		}
	})

	t.Run("TestResourceRemoval", func(t *testing.T) {
		_, err := Delete(context.Background(), service, namespaceID, createdResource.ID)
		if err != nil {
			t.Fatalf("Error deleting resource: %v", err)
		}
	})
}

func TestRetrieveNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	namespaceName := "tests-ns-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	createdNamespace, _, err := tag_namespace.Create(context.Background(), service, tag_namespace.Namespace{
		Name:        namespaceName,
		Description: namespaceName,
		Enabled:     true,
		Origin:      "CUSTOM",
		Type:        "STATIC",
	})
	if err != nil {
		t.Fatalf("Error creating namespace: %v", err)
	}
	defer tag_namespace.Delete(context.Background(), service, createdNamespace.ID)

	_, _, err = Get(context.Background(), service, createdNamespace.ID, "non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	namespaceName := "tests-ns-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	createdNamespace, _, err := tag_namespace.Create(context.Background(), service, tag_namespace.Namespace{
		Name:        namespaceName,
		Description: namespaceName,
		Enabled:     true,
		Origin:      "CUSTOM",
		Type:        "STATIC",
	})
	if err != nil {
		t.Fatalf("Error creating namespace: %v", err)
	}
	defer tag_namespace.Delete(context.Background(), service, createdNamespace.ID)

	_, err = Delete(context.Background(), service, createdNamespace.ID, "non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	namespaceName := "tests-ns-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	createdNamespace, _, err := tag_namespace.Create(context.Background(), service, tag_namespace.Namespace{
		Name:        namespaceName,
		Description: namespaceName,
		Enabled:     true,
		Origin:      "CUSTOM",
		Type:        "STATIC",
	})
	if err != nil {
		t.Fatalf("Error creating namespace: %v", err)
	}
	defer tag_namespace.Delete(context.Background(), service, createdNamespace.ID)

	_, err = Update(context.Background(), service, createdNamespace.ID, "non_existent_id", &TagKey{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	namespaceName := "tests-ns-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	createdNamespace, _, err := tag_namespace.Create(context.Background(), service, tag_namespace.Namespace{
		Name:        namespaceName,
		Description: namespaceName,
		Enabled:     true,
		Origin:      "CUSTOM",
		Type:        "STATIC",
	})
	if err != nil {
		t.Fatalf("Error creating namespace: %v", err)
	}
	defer tag_namespace.Delete(context.Background(), service, createdNamespace.ID)

	_, _, err = GetByName(context.Background(), service, createdNamespace.ID, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
