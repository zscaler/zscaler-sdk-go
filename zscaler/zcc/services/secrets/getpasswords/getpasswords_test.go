package getpasswords

/*
import (
	"fmt"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zcc/services"
)

func TestGetPasswords(t *testing.T) {
	client, err := tests.NewZccClient()
	if err != nil {
		t.Fatalf("Failed to create ZCC client: %v", err)
	}
	service := services.New(client)

	testCases := []struct {
		username string
	}{
		{""},         // No filters
		{"testuser"}, // Username only
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("username=%s", tc.username), func(t *testing.T) {
			passwords, err := GetPasswords(service, tc.username, "")
			if err != nil {
				t.Fatalf("Error retrieving passwords: %v", err)
			}

			// Log the raw response for debugging
			t.Logf("Raw passwords response: %+v", passwords)

			// Check if the response is not nil
			if passwords == nil {
				t.Errorf("Expected non-nil response")
				return
			}

			// Check specific fields in the returned structure if necessary
			if passwords.ExitPass == "" {
				t.Errorf("Expected non-empty ExitPass")
			}
			if passwords.LogoutPass == "" {
				t.Errorf("Expected non-empty LogoutPass")
			}
			if passwords.UninstallPass == "" {
				t.Errorf("Expected non-empty UninstallPass")
			}
			if passwords.ZdSettingsAccessPass == "" {
				t.Errorf("Expected non-empty ZdSettingsAccessPass")
			}
			if passwords.ZdxDisablePass == "" {
				t.Errorf("Expected non-empty ZdxDisablePass")
			}
			if passwords.ZiaDisablePass == "" {
				t.Errorf("Expected non-empty ZiaDisablePass")
			}
			if passwords.ZpaDisablePass == "" {
				t.Errorf("Expected non-empty ZpaDisablePass")
			}
		})
	}
}
*/
