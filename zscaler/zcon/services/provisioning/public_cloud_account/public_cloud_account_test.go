package public_cloud_account

/*
import (
	"log"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zcon/services"
)

// TestGetAccountID verifies the retrieval of a specific account by ID
func TestGetAccountID(t *testing.T) {
	client, err := tests.NewZConClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	// Replace this with an actual account ID if known
	testAccountID := 12345

	account, err := GetAccountID(service, testAccountID)
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
	client, err := tests.NewZConClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	accounts, err := GetLite(service)
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
	client, err := tests.NewZConClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	status, err := GetAccountStatus(service)
	if err != nil {
		t.Errorf("Error retrieving account status: %v", err)
	} else {
		t.Logf("Retrieved account status: AccountIdEnabled: %v, SubIDEnabled: %v", status.AccountIdEnabled, status.SubIDEnabled)
	}
}
*/
