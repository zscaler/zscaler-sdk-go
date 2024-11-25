package emergencyaccess

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestEmergencyAccessIntegration(t *testing.T) {
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Create new resource
	createdResource, _, err := Create(context.Background(), service, &EmergencyAccess{
		ActivatedOn:       "1",
		AllowedActivate:   true,
		AllowedDeactivate: true,
		EmailID:           randomName + "@securitygeek.io",
		FirstName:         "John",
		LastName:          "Smith",
		UserID:            "jsmith",
	})
	if err != nil {
		t.Fatalf("Failed to create emergency user: %v", err)
	}

	// *** New Step: Test GetByEmailID ***
	searchedResource, _, err := GetByEmailID(context.Background(), service, createdResource.EmailID)
	if err != nil {
		t.Errorf("Failed to get emergency user by EmailID: %v", err)
	} else {
		assert.Equal(t, createdResource.EmailID, searchedResource.EmailID, "EmailID does not match")
		t.Logf("Successfully found emergency access record by email ID: %s", searchedResource.EmailID)
	}

	// Test Get
	gotResource, _, err := Get(context.Background(), service, createdResource.UserID)
	if err != nil {
		t.Errorf("Failed to get emergency user by UserID: %v", err)
	}
	assert.Equal(t, createdResource.UserID, gotResource.UserID, "UserID does not match")

	time.Sleep(10 * time.Second)

	//Test Update
	updatedResource := *createdResource
	updatedResource.FirstName = randomName
	_, err = Update(context.Background(), service, createdResource.UserID, &updatedResource)
	if err != nil {
		t.Errorf("Failed to update emergency user: %v", err)
	}

	// Verify Update
	updated, _, err := Get(context.Background(), service, createdResource.UserID)
	if err != nil {
		t.Errorf("Failed to get updated emergency user: %v", err)
	}
	assert.Equal(t, randomName, updated.FirstName, "FirstName was not updated")

	// Test resources retrieval
	resources, _, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		t.Error("Expected retrieved resources to be non-empty, but got empty slice")
	}
	//check if the created resource is in the list
	found := false
	for _, resource := range resources {
		if resource.EmailID == createdResource.EmailID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved resources to contain created resource '%s', but it didn't", createdResource.EmailID)
	}

	// Test Emergency Access User Deactivation
	_, err = Deactivate(context.Background(), service, createdResource.UserID)
	if err != nil {
		t.Errorf("Failed to deactivate emergency user: %v", err)
	}

	// Wait 5 seconds before Activating the user again.
	time.Sleep(5 * time.Second)

	// Test Emergency Access User Activate
	_, err = Activate(context.Background(), service, createdResource.UserID)
	if err != nil {
		t.Errorf("Failed to activate emergency user: %v", err)
	}

	// Test Emergency Access User Deactivation
	_, err = Deactivate(context.Background(), service, createdResource.UserID)
	if err != nil {
		t.Errorf("Failed to deactivate emergency user: %v", err)
	}

	// Simulate delay after deactivation in Cloud Service 1 before proceeding to Okta deletion
	time.Sleep(10 * time.Second) // Adjust the delay as necessary

	// Begin Okta deletion process
	deleteUserInOkta(t, []string{createdResource.UserID}) // Passing the UserID to be deleted in Okta
}

// deleteUserInOkta deletes a user (or users) in Okta based on provided user IDs
func deleteUserInOkta(t *testing.T, userIDs []string) {
	// Fetch Okta domain and API token from environment variables
	oktaDomain := os.Getenv("OKTA_CLIENT_ORGURL")
	apiToken := os.Getenv("OKTA_CLIENT_TOKEN")

	// Initialize Okta client with environment variables
	ctx, client, err := okta.NewClient(
		context.TODO(),
		okta.WithOrgUrl(fmt.Sprintf("https://%s", oktaDomain)),
		okta.WithToken(apiToken),
	)
	if err != nil {
		t.Errorf("Error initializing Okta client: %v", err)
		return
	}

	for _, userID := range userIDs {
		_, err := client.User.DeactivateOrDeleteUser(ctx, userID, nil)
		if err != nil {
			t.Errorf("Failed to delete user %s in Okta: %v", userID, err)
		} else {
			fmt.Printf("User %s deleted successfully in Okta\n", userID)
		}
	}
}
