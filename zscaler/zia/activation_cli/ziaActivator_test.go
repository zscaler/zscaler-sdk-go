package main

/*
import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/activation"
)

// Check if environment variables exist for the test
func checkEnvVarForTest(t *testing.T, k string) {
	if v := os.Getenv(k); v == "" {
		t.Fatalf("[ERROR] Couldn't find environment variable %s", k)
	}
}

func TestActivationCLI(t *testing.T) {
	// Ensure the required environment variables are set for the test
	checkEnvVarForTest(t, "ZSCALER_CLIENT_ID")
	checkEnvVarForTest(t, "ZSCALER_CLIENT_SECRET")
	checkEnvVarForTest(t, "ZSCALER_VANITY_DOMAIN")

	// Create the ZIA client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Perform activation
	_, err = activation.CreateActivation(context.Background(), service, activation.Activation{
		Status: "active",
	})
	if err != nil {
		t.Fatalf("[ERROR] Activation Failed: %v", err)
	}
}

func TestMissingEnvVars(t *testing.T) {
	originalEnv := os.Environ()
	defer func() {
		// Restore original environment after test
		for _, env := range originalEnv {
			parts := strings.SplitN(env, "=", 2)
			os.Setenv(parts[0], parts[1])
		}
	}()

	requireEnvVars := []string{"ZSCALER_CLIENT_ID", "ZSCALER_CLIENT_SECRET", "ZSCALER_VANITY_DOMAIN"}
	for _, envVar := range requireEnvVars {
		t.Run(fmt.Sprintf("Missing %s", envVar), func(t *testing.T) {
			originalValue := os.Getenv(envVar)
			os.Unsetenv(envVar)

			// Attempt to run the check for missing environment variable
			if v := os.Getenv(envVar); v != "" {
				t.Fatalf("[ERROR] Environment variable %s should not be set", envVar)
			}

			// Restore the original value for the next iteration
			if originalValue != "" {
				os.Setenv(envVar, originalValue)
			}
		})
	}
}

func TestActivationStatuses(t *testing.T) {
	// Ensure the required environment variables are set for the test
	checkEnvVarForTest(t, "ZSCALER_CLIENT_ID")
	checkEnvVarForTest(t, "ZSCALER_CLIENT_SECRET")
	checkEnvVarForTest(t, "ZSCALER_VANITY_DOMAIN")

	// Create the ZIA client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Test different activation statuses
	statuses := []string{"active", "inactive", "suspended"}
	for _, status := range statuses {
		t.Run(fmt.Sprintf("Activation status %s", status), func(t *testing.T) {
			_, err := activation.CreateActivation(context.Background(), service, activation.Activation{
				Status: status,
			})
			if err != nil {
				t.Fatalf("[ERROR] Activation Failed for status %s: %v", status, err)
			}
		})
	}
}

func TestSuccessfulActivationAndLogout(t *testing.T) {
	// Ensure the required environment variables are set for the test
	checkEnvVarForTest(t, "ZSCALER_CLIENT_ID")
	checkEnvVarForTest(t, "ZSCALER_CLIENT_SECRET")
	checkEnvVarForTest(t, "ZSCALER_VANITY_DOMAIN")

	// Create the ZIA client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Perform activation
	_, err = activation.CreateActivation(context.Background(), service, activation.Activation{
		Status: "active",
	})
	if err != nil {
		t.Fatalf("[ERROR] Activation Failed: %v", err)
	}

}
*/
