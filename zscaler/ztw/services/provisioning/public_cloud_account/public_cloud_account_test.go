package public_cloud_account

/*
// TestGetAccountID verifies the retrieval of a specific account by ID
func TestGetAccountID(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "public_cloud_account", "ztw")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Replace this with an actual account ID if known
	testAccountID := 12345

	account, err := GetAccountID(context.Background(), service, testAccountID)
	if err != nil {
		t.Logf("No account found for ID %d: %v", testAccountID, err)
	} else {
		if account.ID != testAccountID {
			t.Errorf("Expected account ID %d but got %d", testAccountID, account.ID)
		} else {
			log.Printf("Successfully retrieved account with ID %d", account.ID)
		}
	}
}

// TestGetLite verifies that all public cloud accounts are retrieved correctly
func TestGetLite(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "public_cloud_account", "ztw")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	accounts, err := GetLite(context.Background(), service)
	if err != nil {
		t.Errorf("Error retrieving lite accounts: %v", err)
	} else if len(accounts) == 0 {
		t.Logf("No accounts found in lite data")
	} else {
		t.Logf("Successfully retrieved %d lite accounts", len(accounts))
	}
}

// TestGetAccountStatus verifies the retrieval of the account status
func TestGetAccountStatus(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "public_cloud_account", "ztw")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service
	status, err := GetAccountStatus(context.Background(), service)
	if err != nil {
		t.Errorf("Error retrieving account status: %v", err)
	} else {
		t.Logf("Retrieved account status: AccountIdEnabled: %v, SubIDEnabled: %v", status.AccountIdEnabled, status.SubIDEnabled)
	}
}
*/
