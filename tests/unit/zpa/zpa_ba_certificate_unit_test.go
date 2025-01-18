package unit

/*
import (
	"context"
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
)

func TestBaCertificate_Get(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/certificate/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Certificate 1"}`))
	})

	// Make the GET request
	certificate, _, err := bacertificate.Get(context.Background(), service, "123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	// Check if the certificate ID and name match the expected values
	if certificate.ID != "123" {
		t.Errorf("Expected ba certificate ID '123', but got '%s'", certificate.ID)
	}
	if certificate.Name != "Certificate 1" {
		t.Errorf("Expected ba certificate name 'Certificate 1', but got '%s'", certificate.Name)
	}
}

// You can write similar tests for other functions like GetByName, Update, Delete, and GetAll.
/*
func TestBaCertificate_GetByName(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/clientlessCertificate/issued", func(w http.ResponseWriter, r *http.Request) {
		// Get the query parameter "name" from the request
		query := r.URL.Query()
		name := query.Get("search")

		// Check if the name matches the expected value
		if name == "Certificate1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"list":[
					{"id": "123", "name": "Certificate1"}
				],
				"totalPages":1
				}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Certificate not found"}`))
		}
	})
	service := &bacertificate.Service{
		Client: client,
	}

	// Make the GetByName request
	certificate, _, err := service.GetIssuedByName("Certificate1")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetByName request: %v", err)
	}

	// Check if the certificate ID and name match the expected values
	if certificate.ID != "123" {
		t.Errorf("Expected certificate ID '123', but got '%s'", certificate.ID)
	}
	if certificate.Name != "Certificate1" {
		t.Errorf("Expected certificate name 'Certificate1', but got '%s'", certificate.Name)
	}
}

func TestBaCertificate_GetAll(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/clientlessCertificate/issued", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
				{"id": "123", "name": "Certificate 1"},
				{"id": "456", "name": "Certificate 2"}
			],
			"totalPages":1
			}`))
	})

	// Make the GetAll request
	certificates, _, err := bacertificate.GetAll(context.Background(), service)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check the returned certificates
	expectedCertificates := []*bacertificate.BaCertificate{
		{ID: "123", Name: "Certificate 1"},
		{ID: "456", Name: "Certificate 2"},
	}
	if len(certificates) != len(expectedCertificates) {
		t.Errorf("Expected %d certificates, but got %d", len(expectedCertificates), len(certificates))
	}
	for i, expectedCertificate := range expectedCertificates {
		certificate := certificates[i]
		if certificate.ID != expectedCertificate.ID {
			t.Errorf("Expected certificate ID '%s', but got '%s'", expectedCertificate.ID, certificate.ID)
		}
		if certificate.Name != expectedCertificate.Name {
			t.Errorf("Expected certificate name '%s', but got '%s'", expectedCertificate.Name, certificate.Name)
		}
	}
}
*/
