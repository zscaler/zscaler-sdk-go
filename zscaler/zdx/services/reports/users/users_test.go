package users

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestGetAllUsers(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := GetUsersFilters{
		From: int(from),
		To:   int(to),
	}

	// Call GetAllUsers with the filters
	users, resp, err := GetAllUsers(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting all users: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(users) == 0 {
		t.Log("No users found.")
	} else {
		t.Logf("Retrieved %d users", len(users))
		for _, user := range users {
			t.Logf("User ID: %d, Name: %s, Email: %s", user.ID, user.Name, user.Email)
			for _, device := range user.Devices {
				t.Logf("Device ID: %d, Name: %s", device.ID, device.Name)
			}
		}
	}
}

func TestGetUser(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := GetUsersFilters{
		From: int(from),
		To:   int(to),
	}

	// Invoke GetAllUsers to retrieve the ID of the first user
	users, _, err := GetAllUsers(context.Background(), service, filters)
	if err != nil {
		t.Fatalf("Error getting all users: %v", err)
	}

	if len(users) == 0 {
		t.Log("No users found, skipping GetUser test.")
		return
	}

	firstUserID := users[0].ID

	// Call GetUser with the first user's ID
	user, resp, err := GetUser(context.Background(), service, strconv.Itoa(firstUserID))
	if err != nil {
		t.Fatalf("Error getting user: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	t.Logf("Retrieved user: ID: %d, Name: %s, Email: %s", user.ID, user.Name, user.Email)
	for _, device := range user.Devices {
		t.Logf("Device ID: %d, Name: %s", device.ID, device.Name)
		for _, loc := range device.UserLocation {
			t.Logf("User Location: ID: %s, City: %s, Country: %s", loc.ID, loc.City, loc.Country)
		}
		for _, loc := range device.ZSLocation {
			t.Logf("ZS Location: ID: %d, Name: %s", loc.ID, loc.Name)
		}
	}
}
