package alerts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

// DELETE METHOD IS MISSING BUT AVAILABLE THROUGH THE UI
func TestAlertSubscriptions(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()
	description := "Zscaler Subscription Alert"
	updatedDescription := description + " - Updated"

	alert := &AlertSubscriptions{
		Description:      description,
		Email:            "alert@securitygeek.io",
		Pt0Severities:    []string{"CRITICAL", "MAJOR", "INFO", "MINOR", "DEBUG"},
		SecureSeverities: []string{"CRITICAL", "MAJOR", "INFO", "MINOR", "DEBUG"},
		ManageSeverities: []string{"CRITICAL", "MAJOR", "INFO", "MINOR", "DEBUG"},
		ComplySeverities: []string{"CRITICAL", "MAJOR", "INFO", "MINOR", "DEBUG"},
		SystemSeverities: []string{"CRITICAL", "MAJOR", "INFO", "MINOR", "DEBUG"},
		Deleted:          false,
	}

	// Step 1: Create alert subscription
	createdAlert, _, err := Create(ctx, service, alert)
	if err != nil {
		t.Fatalf("Error creating alert subscription: %v", err)
	}
	if createdAlert.ID == 0 {
		t.Fatal("Expected non-zero ID after creation")
	}
	assert.Equal(t, alert.Email, createdAlert.Email, "Email should match")
	assert.Equal(t, alert.Description, createdAlert.Description, "Description should match")

	// Step 2: Update alert subscription
	createdAlert.Description = updatedDescription
	updatedAlert, _, err := Update(ctx, service, createdAlert.ID, createdAlert)
	if err != nil {
		t.Fatalf("Error updating alert subscription: %v", err)
	}
	assert.Equal(t, updatedDescription, updatedAlert.Description, "Updated description should match")

	// Step 3: Retrieve alert subscription by ID
	retrieved, err := Get(ctx, service, updatedAlert.ID)
	if err != nil {
		t.Fatalf("Error retrieving alert subscription: %v", err)
	}
	assert.Equal(t, updatedAlert.ID, retrieved.ID, "Retrieved ID should match")
	assert.Equal(t, updatedAlert.Description, retrieved.Description, "Retrieved description should match")

	// Step 4: Retrieve all alert subscriptions and check if present
	allAlerts, err := GetAll(ctx, service)
	if err != nil {
		t.Fatalf("Error retrieving all alert subscriptions: %v", err)
	}
	found := false
	for _, a := range allAlerts {
		if a.ID == updatedAlert.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected alert subscription with ID %d to be in the list, but it wasn't", updatedAlert.ID)
	}
}
