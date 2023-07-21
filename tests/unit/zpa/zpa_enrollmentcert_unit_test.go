package unit

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/enrollmentcert"
)

func TestEnrollmentCert_Get(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/enrollmentCert/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Connector"}`))
	})
	service := &enrollmentcert.Service{
		Client: client,
	}

	// Make the GET request
	enrollmentCert, _, err := service.Get("123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
	}

	// Check if the Connector certificate ID and name match the expected values
	if enrollmentCert.ID != "123" {
		t.Errorf("Expected ba Connector certificate ID '123', but got '%s'", enrollmentCert.ID)
	}
	if enrollmentCert.Name != "Connector" {
		t.Errorf("Expected ba certificate name 'Connector', but got '%s'", enrollmentCert.Name)
	}
}

/*
	func TestEnrollmentCert_GetByName(t *testing.T) {
		client, mux, server := tests.NewZpaClientMock()
		defer server.Close()
		mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/enrollmentCert", func(w http.ResponseWriter, r *http.Request) {
			// Get the query parameter "name" from the request
			query := r.URL.Query()
			name := query.Get("search")

			// Check if the name matches the expected value
			if name == "Connector" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"list":[
						{"id": "123", "name": "Connector"}
					],
					"totalPages":1
					}`))
			} else {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error": "Connector Certificate not found"}`))
			}
		})
		service := &enrollmentcert.Service{
			Client: client,
		}

		// Make the GetByName request
		enrollmentCert, _, err := service.GetByName("Connector")
		// Check if the request was successful
		if err != nil {
			t.Errorf("Error making GetByName request: %v", err)
		}

		// Check if the certificate ID and name match the expected values
		if enrollmentCert.ID != "123" {
			t.Errorf("Expected connector certificate ID '123', but got '%s'", enrollmentCert.ID)
		}
		if enrollmentCert.Name != "Connector" {
			t.Errorf("Expected connector certificate name 'Connector', but got '%s'", enrollmentCert.Name)
		}
	}
*/
func TestEnrollmentCert_GetAll(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v2/admin/customers/customerid/enrollmentCert", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
				{"id": "123", "name": "Connector"},
				{"id": "456", "name": "Service Edge"},
				{"id": "789", "name": "Client"},
				{"id": "1011", "name": "Isolation Client"},
				{"id": "1213", "name": "Root"}
			],
			"totalPages":1
			}`))
	})
	service := &enrollmentcert.Service{
		Client: client,
	}

	// Make the GetAll request
	enrollmentCerts, _, err := service.GetAll()
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check the returned certificates
	expectedEnrollmentCerts := []*enrollmentcert.EnrollmentCert{
		{ID: "123", Name: "Connector"},
		{ID: "456", Name: "Service Edge"},
		{ID: "789", Name: "Client"},
		{ID: "1011", Name: "Isolation Client"},
		{ID: "1213", Name: "Root"},
	}
	if len(enrollmentCerts) != len(expectedEnrollmentCerts) {
		t.Errorf("Expected %d enrollment certificates, but got %d", len(expectedEnrollmentCerts), len(enrollmentCerts))
	}
	for i, expectedEnrollmentCert := range expectedEnrollmentCerts {
		enrollmentCert := enrollmentCerts[i]
		if enrollmentCert.ID != expectedEnrollmentCert.ID {
			t.Errorf("Expected enrollment certificate ID '%s', but got '%s'", expectedEnrollmentCert.ID, enrollmentCert.ID)
		}
		if enrollmentCert.Name != expectedEnrollmentCert.Name {
			t.Errorf("Expected enrollment certificate name '%s', but got '%s'", expectedEnrollmentCert.Name, enrollmentCert.Name)
		}
	}
}
