package emergencyaccess

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestEmergencyAccess(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Create new resource
	createdResource := EmergencyAccess{
		ActivatedOn:       "1",
		AllowedActivate:   true,
		AllowedDeactivate: false,
		EmailId:           "jdoe@bd-hashicorp.com",
		FirstName:         "John",
		LastName:          "Doe",
		UserId:            "jdoe",
	}

	if err != nil {
		t.Fatalf("Error creating resource: %v", err)
	}

	t.Run("TestResourceCreation", func(t *testing.T) {
		if createdResource.EmailId == "" {
			t.Error("Expected created resource ID to be non-empty, but got ''")
		}

	})

	t.Run("TestResourceRetrieval", func(t *testing.T) {
		retrievedResource, _, err := service.Get(createdResource.EmailId)
		if err != nil {
			t.Fatalf("Error retrieving resource: %v", err)
		}
		if retrievedResource.EmailId != createdResource.EmailId {
			t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.EmailId, retrievedResource.EmailId)
		}

	})

	t.Run("TestAllResourcesRetrieval", func(t *testing.T) {
		resources, _, err := service.GetAll()
		if err != nil {
			t.Fatalf("Error retrieving groups: %v", err)
		}
		if len(resources) == 0 {
			t.Error("Expected retrieved resources to be non-empty, but got empty slice")
		}
		found := false
		for _, resource := range resources {
			if resource.EmailId == createdResource.EmailId {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected retrieved groups to contain created resource '%s', but it didn't", createdResource.EmailId)
		}
	})
}

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.Get("non-existent-id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}
