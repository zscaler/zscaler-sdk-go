package departments

/*
// {"code":"INVALID_INPUT_ARGUMENT","message":"Authentication type is not hosted"}
// This test only works for Non Zidentity Tenants
const (
	maxRetries    = 3
	retryInterval = 2 * time.Second
)

// Constants for conflict retries
const (
	maxConflictRetries    = 5
	conflictRetryInterval = 1 * time.Second
)

func retryOnConflict(operation func() error) error {
	var lastErr error
	for i := 0; i < maxConflictRetries; i++ {
		lastErr = operation()
		if lastErr == nil {
			return nil
		}

		if strings.Contains(lastErr.Error(), `"code":"EDIT_LOCK_NOT_AVAILABLE"`) {
			log.Printf("Conflict error detected, retrying in %v... (Attempt %d/%d)", conflictRetryInterval, i+1, maxConflictRetries)
			time.Sleep(conflictRetryInterval)
			continue
		}

		return lastErr
	}
	return lastErr
}

func TestAccDepartmentManagement(t *testing.T) {
	name := acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	comments := acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)
	updateComments := acctest.RandStringFromCharSet(30, acctest.CharSetAlpha)

	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	department := Department{
		Name:     name,
		Comments: comments,
	}

	var createdResource *Department

	// Test resource creation
	err = retryOnConflict(func() error {
		createdResource, _, err = Create(context.Background(), service, &department)
		return err
	})
	if err != nil {
		t.Fatalf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-zero, but got 0")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created department '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource retrieval
	retrievedResource, err := tryRetrieveResource(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved department '%s', but got '%s'", name, retrievedResource.Name)
	}

	// Test resource update
	retrievedResource.Comments = updateComments
	err = retryOnConflict(func() error {
		_, _, err = Update(context.Background(), service, createdResource.ID, retrievedResource)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating resource: %v", err)
	}

	updatedResource, err := GetDepartments(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Fatalf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Comments != updateComments {
		t.Errorf("Expected retrieved updated resource description '%s', but got '%s'", updateComments, updatedResource.Comments)
	}

	// Test resource retrieval by name
	retrievedResource, err = GetDepartmentsByName(context.Background(), service, name)
	if err != nil {
		t.Fatalf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Comments != updateComments {
		t.Errorf("Expected retrieved resource comment '%s', but got '%s'", updateComments, createdResource.Comments)
	}
	// Test resources retrieval
	resources, err := GetAll(context.Background(), service, nil)
	if err != nil {
		t.Fatalf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		t.Fatal("Expected retrieved resources to be non-empty, but got empty slice")
	}
	// check if the created resource is in the list
	found := false
	for _, resource := range resources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved resources to contain created resource '%d', but it didn't", createdResource.ID)
	}
	// Test resource removal
	err = retryOnConflict(func() error {
		_, delErr := Delete(context.Background(), service, createdResource.ID)
		return delErr
	})
	_, err = GetDepartments(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Fatalf("Expected error retrieving deleted resource, but got nil")
	}
}

// tryRetrieveResource attempts to retrieve a resource with retry mechanism.
func tryRetrieveResource(ctx context.Context, service *zscaler.Service, id int) (*Department, error) {
	var resource *Department
	var err error

	for i := 0; i < maxRetries; i++ {
		// Use the passed context (ctx) instead of context.Background()
		resource, err = GetDepartments(ctx, service, id)
		if err == nil && resource != nil && resource.ID == id {
			return resource, nil
		}
		log.Printf("Attempt %d: Error retrieving resource, retrying in %v...", i+1, retryInterval)
		time.Sleep(retryInterval)
	}

	return nil, err
}

func TestRetrieveNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = GetDepartments(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	_, err = Delete(context.Background(), service, 0)
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, _, err = Update(context.Background(), service, 0, &Department{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = GetDepartmentsByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
*/
