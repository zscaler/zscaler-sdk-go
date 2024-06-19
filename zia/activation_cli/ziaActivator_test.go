package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/activation"
)

func TestActivationCLI(t *testing.T) {
	// Check that necessary environment variables are set
	checkEnvVarForTest(t, "ZIA_USERNAME")
	checkEnvVarForTest(t, "ZIA_PASSWORD")
	checkEnvVarForTest(t, "ZIA_API_KEY")
	checkEnvVarForTest(t, "ZIA_CLOUD")

	// Construct the client
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	_, err = activation.CreateActivation(service, activation.Activation{
		Status: "active",
	})
	if err != nil {
		t.Fatalf("[ERROR] Activation Failed: %v", err)
	}

	// Destroy the session
	if err := client.Logout(); err != nil {
		t.Fatalf("[ERROR] Failed destroying session: %v", err)
	}
}

func checkEnvVarForTest(t *testing.T, k string) {
	if v := os.Getenv(k); v == "" {
		t.Fatalf("[ERROR] Couldn't find environment variable %s", k)
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

	// Unset required environment variables
	os.Unsetenv("ZIA_USERNAME")
	os.Unsetenv("ZIA_PASSWORD")
	os.Unsetenv("ZIA_API_KEY")
	os.Unsetenv("ZIA_CLOUD")

	requireEnvVars := []string{"ZIA_USERNAME", "ZIA_PASSWORD", "ZIA_API_KEY", "ZIA_CLOUD"}
	for _, envVar := range requireEnvVars {
		t.Run(fmt.Sprintf("Missing %s", envVar), func(t *testing.T) {
			originalValue := os.Getenv(envVar)
			os.Unsetenv(envVar)

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
	// Check that necessary environment variables are set
	checkEnvVarForTest(t, "ZIA_USERNAME")
	checkEnvVarForTest(t, "ZIA_PASSWORD")
	checkEnvVarForTest(t, "ZIA_API_KEY")
	checkEnvVarForTest(t, "ZIA_CLOUD")

	// Construct the client
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	statuses := []string{"active"}
	for _, status := range statuses {
		t.Run(fmt.Sprintf("Activation status %s", status), func(t *testing.T) {
			_, err := activation.CreateActivation(service, activation.Activation{
				Status: status,
			})
			if err != nil {
				t.Fatalf("[ERROR] Activation Failed for status %s: %v", status, err)
			}
		})
	}

	// Destroy the session
	if err := client.Logout(); err != nil {
		t.Fatalf("[ERROR] Failed destroying session: %v", err)
	}
}

func TestSuccessfulActivationAndLogout(t *testing.T) {
	// Check that necessary environment variables are set
	checkEnvVarForTest(t, "ZIA_USERNAME")
	checkEnvVarForTest(t, "ZIA_PASSWORD")
	checkEnvVarForTest(t, "ZIA_API_KEY")
	checkEnvVarForTest(t, "ZIA_CLOUD")

	// Construct the client
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Perform activation
	_, err = activation.CreateActivation(service, activation.Activation{
		Status: "active",
	})
	if err != nil {
		t.Fatalf("[ERROR] Activation Failed: %v", err)
	}

	// Destroy the session
	if err := client.Logout(); err != nil {
		t.Fatalf("[ERROR] Failed destroying session: %v", err)
	}
}
