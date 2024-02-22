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
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestEmergencyAccessIntegration(t *testing.T) {
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	service := New(client)

	// Create new resource
	createdResource, _, err := service.Create(&EmergencyAccess{
		ActivatedOn:       "1",
		AllowedActivate:   true,
		AllowedDeactivate: true,
		EmailID:           randomName + "@bd-hashicorp.com",
		FirstName:         "John",
		LastName:          "Smith",
		UserID:            "jsmith",
	})
	if err != nil {
		t.Fatalf("Failed to create emergency user: %v", err)
	}

	// Test Get
	gotResource, _, err := service.Get(createdResource.UserID)
	if err != nil {
		t.Errorf("Failed to get emergency user by UserID: %v", err)
	}
	assert.Equal(t, createdResource.UserID, gotResource.UserID, "UserID does not match")

	//Test Update
	updatedResource := *createdResource
	updatedResource.FirstName = randomName
	_, err = service.Update(createdResource.UserID, &updatedResource)
	if err != nil {
		t.Errorf("Failed to update emergency user: %v", err)
	}

	// Verify Update
	updated, _, err := service.Get(createdResource.UserID)
	if err != nil {
		t.Errorf("Failed to get updated emergency user: %v", err)
	}
	assert.Equal(t, randomName, updated.FirstName, "FirstName was not updated")

	// Test Emergency Access User Deactivation
	_, err = service.Deactivate(createdResource.UserID)
	if err != nil {
		t.Errorf("Failed to deactivate emergency user: %v", err)
	}

	// Test Emergency Access User Activate
	_, err = service.Activate(createdResource.UserID)
	if err != nil {
		t.Errorf("Failed to activate emergency user: %v", err)
	}

	// Test Emergency Access User Deactivation
	_, err = service.Deactivate(createdResource.UserID)
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
	oktaDomain := os.Getenv("OKTA_DOMAIN")
	apiToken := os.Getenv("OKTA_API_TOKEN")

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
