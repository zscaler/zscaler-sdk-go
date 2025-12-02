package portal_link

import (
	"context"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
)

func TestUserPortalLink(t *testing.T) {
	tests.ResetTestNameCounter()

	name := tests.GetTestName("tests-plink")
	updateName := tests.GetTestName("tests-plink")
	// Use unique ExtLabel for VCR determinism - portal_link test uses "portallink"
	extLabel := "portallink02"

	client, err := tests.NewVCRTestClient(t, "portal_link", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	userPortal, _, err := portal_controller.Create(context.Background(), service, portal_controller.UserPortalController{
		Name:                    name,
		Description:             name,
		Enabled:                 true,
		UserNotification:        "Created via GO SDK",
		UserNotificationEnabled: true,
		ManagedByZS:             true,
		ExtLabel:                extLabel,
		ExtDomainName:           "-securitygeek-io.b.zscalerportal.net",
		ExtDomain:               "securitygeek.io",
	})
	// Check if the request was successful - use Fatalf to stop test if creation fails
	if err != nil {
		t.Fatalf("Error creating user portal for testing portal link: %v", err)
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := portal_controller.Get(context.Background(), service, userPortal.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := portal_controller.Delete(context.Background(), service, userPortal.ID)
			if err != nil {
				t.Errorf("Error deleting user portal: %v", err)
			}
		}
	}()

	// Create new resource
	createdResource, _, err := Create(context.Background(), service, UserPortalLink{
		Name:        name,
		Description: name,
		Enabled:     true,
		Link:        "server1.example.com",
		Protocol:    "https://",
		UserPortals: []portal_controller.UserPortalController{
			{
				ID: userPortal.ID,
			},
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
