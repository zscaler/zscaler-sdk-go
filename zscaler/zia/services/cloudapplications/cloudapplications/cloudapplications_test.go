package cloudapplications

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestCloudApplications(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "cloudapplications", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Define test parameters for Cloud Application Policy
	appPolicyParams := map[string]interface{}{
		"appClass": []interface{}{"WEB_MAIL"},
	}

	// Fetch Cloud Application Policy for WEB_MAIL
	cloudAppPolicies, err := GetCloudApplicationPolicy(context.Background(), service, appPolicyParams)
	if err != nil {
		t.Fatalf("Error fetching Cloud Application Policy: %v", err)
	}
	assert.NotEmpty(t, cloudAppPolicies, "Expected non-empty response for Cloud Application Policy")

	// Log result for debugging
	// for _, app := range cloudAppPolicies {
	// 	t.Logf("Policy: %s, %s", app.AppName, app.ParentName)
	// }

	// Define test parameters for Cloud Application SSL Policy
	sslPolicyParams := map[string]interface{}{
		"appClass": []interface{}{"SOCIAL_NETWORKING"},
	}

	// Fetch Cloud Application SSL Policy for SOCIAL_NETWORKING
	cloudAppSSLPolicies, err := GetCloudApplicationSSLPolicy(context.Background(), service, sslPolicyParams)
	if err != nil {
		t.Fatalf("Error fetching Cloud Application SSL Policy: %v", err)
	}
	assert.NotEmpty(t, cloudAppSSLPolicies, "Expected non-empty response for Cloud Application SSL Policy")

	// Log result for debugging
	// for _, app := range cloudAppSSLPolicies {
	// 	t.Logf("SSL Policy: %s, %s", app.AppName, app.ParentName)
	// }
}
