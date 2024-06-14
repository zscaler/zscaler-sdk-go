package unit

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/isolationprofile"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/postureprofile"
)

// You can write similar tests for other functions like GetByName, Update, Delete, and GetAll.
func TestIsolationProfile_GetByName(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/isolation/profiles", func(w http.ResponseWriter, r *http.Request) {
		// Get the query parameter "name" from the request
		query := r.URL.Query()
		name := query.Get("search")

		// Check if the name matches the expected value
		if name == "CBIProfile1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"list":[
					{"id": "123", "name": "CBIProfile1"}
				],
				"totalPages":1
				}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "CBI Profile1 not found"}`))
		}
	})

	// Make the GetByName request
	profile, _, err := isolationprofile.GetByName(service, "CBIProfile1")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetByName request: %v", err)
	}

	// Check if the posture ID and name match the expected values
	if profile.ID != "123" {
		t.Errorf("Expected cbi profile ID '123', but got '%s'", profile.ID)
	}
	if profile.Name != "CBIProfile1" {
		t.Errorf("Expected cbi profile name 'CBIProfile1', but got '%s'", profile.Name)
	}
}

func TestIsolationProfile_GetAll(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/isolation/profiles", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
				{"id": "123", "name": "CBIProfile1"},
				{"id": "456", "name": "CBIProfile2"}
			],
			"totalPages":1
			}`))
	})

	// Make the GetAll request
	profiles, _, err := isolationprofile.GetAll(service)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check the returned profiles
	expectedProfiles := []*postureprofile.PostureProfile{
		{ID: "123", Name: "CBIProfile1"},
		{ID: "456", Name: "CBIProfile2"},
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
