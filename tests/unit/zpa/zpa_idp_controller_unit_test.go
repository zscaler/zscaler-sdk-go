package unit

/*
import (
	"context"
	"net/http"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
)

func TestIdpController_Get(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/idp/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Idp 1"}`))
	})

	// Make the GET request
	idp, _, err := idpcontroller.Get(context.Background(), service, "123")
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

func TestIdpController_GetByName(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/idp", func(w http.ResponseWriter, r *http.Request) {
		// Get the query parameter "page" from the request
		query := r.URL.Query()
		page := query.Get("page")

		// Check if the name matches the expected value
		if page == "1" {
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

	// Make the GetByName request
	idp, _, err := idpcontroller.GetByName(context.Background(), service, "Idp1")
	// Check if the request was successful
	if err != nil {
		t.Fatalf("Error making GetByName request: %v", err)
	}

	// Check if the Idp1 and name match the expected values
	if idp == nil || idp.ID != "123" {
		t.Errorf("Expected Idp1 '123', but got '%v'", idp)
	}
	if idp == nil || idp.Name != "Idp1" {
		t.Errorf("Expected idp name 'Idp1', but got '%v'", idp)
	}
}

func TestIdpController_GetAll(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

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

	// Make the GetAll request
	idps, _, err := idpcontroller.GetAll(context.Background(), service)
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
*/
