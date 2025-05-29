package administrator_controller

/*
func TestAdministratorController(t *testing.T) {
	username := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateDisplayName := "Updated ZPA Admin"

	service, err := tests.NewZPAClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Create new resource
	createdResource, _, err := Create(context.Background(), service, &AdministratorController{
		Username:           username + "@bd-hashicorp.com",
		DisplayName:        "ZPA Admin " + username,
		Email:              username + "@bd-hashicorp.com",
		OperationType:      "UPSERT",
		RoleId:             "12",
		Eula:               "0",
		IsEnabled:          true,
		ForcePwdChange:     false,
		PinSession:         true,
		LocalLoginDisabled: false,
		IsLocked:           false,
		Password:           tests.TestPassword(10),
		Role: Role{
			ID: "12",
		},
	})
	if err != nil {
		t.Fatalf("Error creating resource: %v", err)
	}

	expectedUsername := username + "@bd-hashicorp.com"

	t.Run("TestResourceCreation", func(t *testing.T) {
		if createdResource.ID == "" {
			t.Error("Expected created resource ID to be non-empty, but got ''")
		}
		if createdResource.Username != expectedUsername {
			t.Errorf("Expected created resource username '%s', but got '%s'", expectedUsername, createdResource.Username)
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
		if retrievedResource.Username != expectedUsername {
			t.Errorf("Expected retrieved resource username '%s', but got '%s'", expectedUsername, retrievedResource.Username)
		}
	})

	t.Run("TestResourceUpdate", func(t *testing.T) {
		updatedResource := *createdResource
		updatedResource.DisplayName = updateDisplayName
		_, err = Update(context.Background(), service, createdResource.ID, &updatedResource)
		if err != nil {
			t.Fatalf("Error updating resource: %v", err)
		}

		// Optional: re-fetch and verify
		refetched, _, err := Get(context.Background(), service, createdResource.ID)
		if err != nil {
			t.Fatalf("Error refetching updated resource: %v", err)
		}
		if refetched.DisplayName != updateDisplayName {
			t.Errorf("Expected updated display name '%s', got '%s'", updateDisplayName, refetched.DisplayName)
		}
	})

	t.Run("TestResourceRetrievalByName", func(t *testing.T) {
		retrievedResource, _, err := GetByName(context.Background(), service, expectedUsername)
		if err != nil {
			t.Fatalf("Error retrieving resource by name: %v", err)
		}
		if retrievedResource.ID != createdResource.ID {
			t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
		}
		if retrievedResource.Username != expectedUsername {
			t.Errorf("Expected retrieved resource username '%s', but got '%s'", expectedUsername, retrievedResource.Username)
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
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

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
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, err = Update(context.Background(), service, "non_existent_id", &AdministratorController{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, _, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
*/
