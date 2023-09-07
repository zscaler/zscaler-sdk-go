package unit

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/postureprofile"
)

func TestPostureProfile_Get(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/posture/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Posture 1"}`))
	})
	service := &postureprofile.Service{
		Client: client,
	}

	// Make the GET request
	posture, _, err := service.Get("123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	// Check if the posture ID and name match the expected values
	if posture.ID != "123" {
		t.Errorf("Expected posture ID '123', but got '%s'", posture.ID)
	}
	if posture.Name != "Posture 1" {
		t.Errorf("Expected posture name 'Posture 1', but got '%s'", posture.Name)
	}
}

// You can write similar tests for other functions like GetByName, Update, Delete, and GetAll.
func TestPostureProfile_GetByName(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/posture", func(w http.ResponseWriter, r *http.Request) {
		// Get the query parameter "name" from the request
		query := r.URL.Query()
		name := query.Get("search")

		// Check if the name matches the expected value
		if name == "Posture1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"list":[
					{"id": "123", "name": "Posture1"}
				],
				"totalPages":1
				}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Posture Profile not found"}`))
		}
	})
	service := &postureprofile.Service{
		Client: client,
	}

	// Make the GetByName request
	profile, _, err := service.GetByName("Posture1")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetByName request: %v", err)
	}

	// Check if the posture ID and name match the expected values
	if profile.ID != "123" {
		t.Errorf("Expected posture ID '123', but got '%s'", profile.ID)
	}
	if profile.Name != "Posture1" {
		t.Errorf("Expected posture name 'Posture1', but got '%s'", profile.Name)
	}
}

func TestPostureProfile_GetAll(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/posture", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
				{"id": "123", "name": "Posture 1"},
				{"id": "456", "name": "Posture 2"}
			],
			"totalPages":1
			}`))
	})
	service := &postureprofile.Service{
		Client: client,
	}

	// Make the GetAll request
	profiles, _, err := service.GetAll()
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check the returned profiles
	expectedProfiles := []*postureprofile.PostureProfile{
		{ID: "123", Name: "Posture 1"},
		{ID: "456", Name: "Posture 2"},
	}
	if len(profiles) != len(expectedProfiles) {
		t.Errorf("Expected %d profiles, but got %d", len(expectedProfiles), len(profiles))
	}
	for i, expectedProfile := range expectedProfiles {
		profile := profiles[i]
		if profile.ID != expectedProfile.ID {
			t.Errorf("Expected profile ID '%s', but got '%s'", expectedProfile.ID, profile.ID)
		}
		if profile.Name != expectedProfile.Name {
			t.Errorf("Expected profile name '%s', but got '%s'", expectedProfile.Name, profile.Name)
		}
	}
}
