package main

/*
import (
	"os"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
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

	service := New(client)
	activationService := activation.New(service.Client)

	_, err = activationService.CreateActivation(activation.Activation{
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
*/
