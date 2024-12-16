package organization_details

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestOrganizationInformation(t *testing.T) {
	// Create the client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	// Step 1: Test GetSubscriptions
	t.Run("GetSubscriptions", func(t *testing.T) {
		subscriptions, err := GetSubscriptions(ctx, service)
		if err != nil {
			t.Fatalf("Error retrieving subscriptions: %v", err)
		}

		// Assert that subscriptions were returned
		assert.NotEmpty(t, subscriptions, "Expected non-empty subscriptions")
		t.Logf("Subscriptions retrieved: %+v", subscriptions)
	})

	// Step 2: Test GetOrgInformation
	t.Run("GetOrgInformation", func(t *testing.T) {
		orgInfo, err := GetOrgInformation(ctx, service)
		if err != nil {
			t.Fatalf("Error retrieving organization information: %v", err)
		}

		// Assert that organization information was returned
		assert.NotNil(t, orgInfo, "Expected non-nil organization information")
		assert.NotEmpty(t, orgInfo.Name, "Expected organization name to be set")
		t.Logf("Organization information retrieved: %+v", orgInfo)
	})

	// Step 3: Test GetOrgInformationLite
	t.Run("GetOrgInformationLite", func(t *testing.T) {
		orgInfoLite, err := GetOrgInformationLite(ctx, service)
		if err != nil {
			t.Fatalf("Error retrieving lite organization information: %v", err)
		}

		// Assert that lite organization information was returned
		assert.NotNil(t, orgInfoLite, "Expected non-nil lite organization information")
		assert.NotEmpty(t, orgInfoLite.Name, "Expected lite organization name to be set")
		t.Logf("Lite organization information retrieved: %+v", orgInfoLite)
	})
}
