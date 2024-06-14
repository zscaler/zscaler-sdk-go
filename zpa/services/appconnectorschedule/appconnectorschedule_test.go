package appconnectorschedule

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

func TestAppConnectorSchedule(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := services.New(client)

	// Retrieve CustomerID from environment variable
	customerID := os.Getenv("ZPA_CUSTOMER_ID")
	if customerID == "" {
		t.Fatal("ZPA_CUSTOMER_ID environment variable is not set")
	}

	// Test 1: CreateSchedule (without ID)
	newSchedule := AssistantSchedule{
		CustomerID:        customerID,
		DeleteDisabled:    true,
		Enabled:           true,
		Frequency:         "days",
		FrequencyInterval: "5",
	}
	_, createResp, err := CreateSchedule(service, newSchedule)
	if err != nil {
		if strings.Contains(err.Error(), "resource.already.exist") {
			t.Log("Assistance Scheduler already enabled")
		} else {
			t.Fatalf("Error creating schedule: %v", err)
		}
	} else if createResp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got: %v", createResp.StatusCode)
	}

	// Test 2: GetSchedule (Initial fetch)
	schedule, resp, err := GetSchedule(service)
	if err != nil {
		t.Fatalf("Error getting schedule: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got: %v", resp.StatusCode)
	}
	if schedule == nil || schedule.ID == "" {
		t.Fatal("Expected non-nil schedule with valid ID")
	}
	t.Logf("Got schedule: %+v", schedule)

	// Ensure the schedule is enabled before updating
	if !schedule.Enabled {
		schedule.Enabled = true
		schedule.FrequencyInterval = "5" // Set a valid interval when enabling
		_, err = UpdateSchedule(service, schedule.ID, schedule)
		if err != nil {
			t.Fatalf("Error enabling schedule: %v", err)
		}
		t.Log("Schedule was disabled and has been enabled for the test.")
	}

	// Test 3: UpdateSchedule with various frequency intervals
	intervals := []string{"7", "14", "30", "60", "90"}
	for _, interval := range intervals {
		schedule.FrequencyInterval = interval
		updateResp, err := UpdateSchedule(service, schedule.ID, schedule)
		if err != nil {
			t.Fatalf("Error updating schedule with interval %s: %v", interval, err)
		}
		if updateResp.StatusCode != 204 {
			t.Errorf("Expected status code 204 for interval %s, got: %v", interval, updateResp.StatusCode)
		}
		t.Logf("Updated schedule with interval: %s", interval)
	}

	// Test 4: GetSchedule (Post-update fetch)
	updatedSchedule, resp, err := GetSchedule(service)
	if err != nil {
		t.Fatalf("Error getting updated schedule: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got: %v", resp.StatusCode)
	}
	if updatedSchedule == nil {
		t.Fatal("Expected non-nil updated schedule")
	}
	t.Logf("Got updated schedule: %+v", updatedSchedule)
}

func TestUpdateScheduleWhenDisabled(t *testing.T) {
	client, err := tests.NewZpaClient()
	require.NoError(t, err, "Error creating client")

	service := services.New(client)
	schedule, _, err := GetSchedule(service)
	require.NoError(t, err, "Error getting schedule")
	require.NotNil(t, schedule, "Schedule should not be nil")

	// Temporarily disable the schedule for testing
	schedule.Enabled = false
	schedule.FrequencyInterval = "7"

	// Check if update fails when the schedule is disabled
	_, err = UpdateSchedule(service, schedule.ID, schedule)
	require.Error(t, err, "Update should fail when Enabled is false")
	require.Contains(t, err.Error(), "cannot update a disabled schedule", "Expected error message when updating a disabled schedule")
}

func TestFrequencyIntervalBoundaries(t *testing.T) {
	client, err := tests.NewZpaClient()
	require.NoError(t, err, "Error creating client")

	service := services.New(client)
	schedule, _, err := GetSchedule(service)
	require.NoError(t, err, "Error getting schedule")
	require.NotNil(t, schedule, "Schedule should not be nil")

	validIntervals := []string{"5", "7", "14", "30", "60", "90"}
	invalidIntervals := []string{"1", "6", "15", "29", "100"}

	// Test invalid intervals with delay to avoid rate limiting
	for _, interval := range invalidIntervals {
		schedule.FrequencyInterval = interval
		_, err := UpdateSchedule(service, schedule.ID, schedule)
		require.Error(t, err, "Invalid interval %s should be rejected", interval)
		time.Sleep(1 * time.Second) // Delay to avoid rate limiting
	}

	// Test valid intervals with delay to avoid rate limiting
	for _, interval := range validIntervals {
		schedule.FrequencyInterval = interval
		_, err := UpdateSchedule(service, schedule.ID, schedule)
		require.NoError(t, err, "Valid interval %s should be accepted", interval)
		time.Sleep(1 * time.Second) // Delay to avoid rate limiting
	}
}

func TestCustomerIDValidation(t *testing.T) {
	client, err := tests.NewZpaClient()
	require.NoError(t, err, "Error creating client")

	service := services.New(client)
	schedule := AssistantSchedule{
		CustomerID:        "", // Intentionally left blank
		DeleteDisabled:    true,
		Enabled:           true,
		Frequency:         "days",
		FrequencyInterval: "5",
	}

	_, _, err = CreateSchedule(service, schedule)
	require.Error(t, err, "Schedule creation should fail with empty CustomerID")
}
