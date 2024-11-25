package microtenants

import (
	"context"
	"strings"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestMicrotenants(t *testing.T) {
	// Define the list of all possible domain names
	domains := []string{
		"216196257331281920.zpa-customer.com",
		"securitygeek.io",
		// "144124980601290752.zpa-customer.com",
		// "public-api-sdk-testing.com",
		// "144124981675032576.zpa-customer.com",
		// "public-api-sdk-testing1.com",
		// "bd-hashicorp.com",
		// "bd-redhat.com",
		// "216199618143191040.zpa-customer.com",
		// "securitygeek.io",
		// "72058304855015424.zpa-customer.com",
		// "securitygeekio.ca",
		// "72057604775346176.zpa-customer.com",
		// "72059901509107712.zpa-customer.com",
		// "72059899361624064.zpa-customer.com",
	}

	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	var createdResource *MicroTenant

	// Loop through each domain until successful creation or all domains exhausted
	for _, domain := range domains {
		createdResource, _, err = Create(context.Background(), service, MicroTenant{
			Name:                    name,
			Description:             name,
			Enabled:                 true,
			CriteriaAttribute:       "AuthDomain",
			CriteriaAttributeValues: []string{domain},
		})
		if err != nil {
			// Check for specific error message and continue if found
			if strings.Contains(err.Error(), "domains.does.not.belong.to.customer") {
				continue
			} else {
				t.Fatalf("Error creating resource with domain '%s': %v", domain, err)
			}
		} else {
			// If successfully created, break out of loop
			break
		}
	}

	// If we've exhausted all domains without a successful creation, fail the test
	if err != nil && strings.Contains(err.Error(), "domains.does.not.belong.to.customer") {
		t.Fatalf("Error creating resource: all domains exhausted, none belong to customer.")
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

	// t.Run("TestRetrieveLoginInfo", func(t *testing.T) {
	// 	if createdResource.UserResource == nil {
	// 		t.Fatalf("Expected user details in created resource, but got nil")
	// 	}

	// 	username := createdResource.UserResource.Username
	// 	if username == "" {
	// 		t.Error("Expected username to be non-empty, but got ''")
	// 	} else {
	// 		if !strings.Contains(username, "@") {
	// 			t.Errorf("Expected valid username format containing '@', but got '%s'", username)
	// 		}
	// 	}

	// 	password := createdResource.UserResource.Password
	// 	if password == "" {
	// 		t.Error("Expected password to be non-empty, but got ''")
	// 	}
	// })

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
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, _, err = Get(context.Background(), service, "non-existent-id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = Delete(context.Background(), service, "non-existent-id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = Update(context.Background(), service, "non-existent-id", &MicroTenant{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, _, err = GetByName(context.Background(), service, "non-existent-name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
