package unit

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/idpcontroller"
)

func TestIdpController_Get(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/idp/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Idp 1"}`))
	})
	service := &idpcontroller.Service{
		Client: client,
	}

	// Make the GET request
	idp, _, err := service.Get("123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	// Check if the idp ID and name match the expected values
	if idp.ID != "123" {
		t.Errorf("Expected idp ID '123', but got '%s'", idp.ID)
	}
	if idp.Name != "Idp 1" {
		t.Errorf("Expected idp name 'Idp 1', but got '%s'", idp.Name)
	}
}

// You can write similar tests for other functions like GetByName, Update, Delete, and GetAll.
/*
func TestIdpController_GetByName(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/idp", func(w http.ResponseWriter, r *http.Request) {
		// Get the query parameter "name" from the request
		query := r.URL.Query()
		name := query.Get("search")

		// Check if the name matches the expected value
		if name == "Idp1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"list":[
					{"id": "123", "name": "Idp1"}
				],
				"totalPages":1
				}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "IDP not found"}`))
		}
	})
	service := &idpcontroller.Service{
		Client: client,
	}

	// Make the GetByName request
	idp, _, err := service.GetByName("Idp1")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetByName request: %v", err)
	}

	// Check if the idp ID and name match the expected values
	if idp.ID != "123" {
		t.Errorf("Expected idp ID '123', but got '%s'", idp.ID)
	}
	if idp.Name != "Idp1" {
		t.Errorf("Expected idp name 'Idp1', but got '%s'", idp.Name)
	}
}
*/
func TestIdpController_GetAll(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/idp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
				{"id": "123", "name": "Idp 1"},
				{"id": "456", "name": "Idp 2"}
			],
			"totalPages":1
			}`))
	})
	service := &idpcontroller.Service{
		Client: client,
	}

	// Make the GetAll request
	idps, _, err := service.GetAll()
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check the returned idps
	expectedIdps := []*idpcontroller.IdpController{
		{ID: "123", Name: "Idp 1"},
		{ID: "456", Name: "Idp 2"},
	}
	if len(idps) != len(expectedIdps) {
		t.Errorf("Expected %d idps, but got %d", len(expectedIdps), len(idps))
	}
	for i, expectedIdp := range expectedIdps {
		idp := idps[i]
		if idp.ID != expectedIdp.ID {
			t.Errorf("Expected idp ID '%s', but got '%s'", expectedIdp.ID, idp.ID)
		}
		if idp.Name != expectedIdp.Name {
			t.Errorf("Expected idp name '%s', but got '%s'", expectedIdp.Name, idp.Name)
		}
	}
}
