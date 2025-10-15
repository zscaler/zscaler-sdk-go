package roles

/*
func TestAdminRoles_data(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	roles, err := GetAllAdminRoles(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting admin roles: %v", err)
		return
	}
	if len(roles) == 0 {
		t.Errorf("No admin roles found")
		return
	}
	name := roles[0].Name
	t.Log("Getting admin roles by name: " + name)

	// Test GetByName function
	role, err := GetByName(context.Background(), service, name)
	if err != nil {
		t.Errorf("Error getting admin role by name: %v", err)
		return
	}
	if role.Name != name {
		t.Errorf("Admin role name does not match: expected %s, got %s", name, role.Name)
		return
	}

	// Test GetAPIRole function with includeApiRole=true
	apiRole, err := GetAPIRole(context.Background(), service, name, "true")
	if err != nil {
		t.Errorf("Error getting API role: %v", err)
		return
	}
	if apiRole.Name != name {
		t.Errorf("API role name does not match: expected %s, got %s", name, apiRole.Name)
		return
	}

	// Test GetAuditorRole function with includeAuditorRole=true
	auditorRole, err := GetAuditorRole(context.Background(), service, name, "true")
	if err != nil {
		t.Errorf("Error getting Auditor role: %v", err)
		return
	}
	if auditorRole.Name != name {
		t.Errorf("Auditor role name does not match: expected %s, got %s", name, auditorRole.Name)
		return
	}

	// Test GetPartnerRole function with includePartnerRole=true
	partnerRole, err := GetPartnerRole(context.Background(), service, name, "true")
	if err != nil {
		t.Errorf("Error getting Partner role: %v", err)
		return
	}
	if partnerRole.Name != name {
		t.Errorf("Partner role name does not match: expected %s, got %s", name, partnerRole.Name)
		return
	}

	// Negative Test: Try to retrieve an admin role with a non-existent name
	nonExistentName := "ThisAdminRoleDoesNotExist"
	_, err = GetByName(context.Background(), service, nonExistentName)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestResponseFormatValidation(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	roles, err := GetAllAdminRoles(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting admin role: %v", err)
		return
	}
	if len(roles) == 0 {
		t.Errorf("No admin role found")
		return
	}

	// Validate admin role
	for _, role := range roles {
		// Checking if essential fields are not empty
		if role.ID == 0 {
			t.Errorf("admin role ID is empty")
		}
		if role.Name == "" {
			t.Errorf("admin role Name is empty")
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Assuming a role with the name "Engineering" exists
	knownName := "Super Admin"

	// Case variations to test
	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, variation := range variations {
		t.Logf("Attempting to retrieve role with name variation: %s", variation)
		role, err := GetByName(context.Background(), service, variation)
		if err != nil {
			t.Errorf("Error getting role with name variation '%s': %v", variation, err)
			continue
		}

		// Check if the group's actual name matches the known name
		if role.Name != knownName {
			t.Errorf("Expected role name to be '%s' for variation '%s', but got '%s'", knownName, variation, role.Name)
		}
	}
}
*/
