package cbicertificatecontroller

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

// clean all resources
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources() // clean up at the beginning
}

func teardown() {
	cleanResources() // clean up at the end
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	return !present || (present && (val == "" || val == "true")) // simplified for clarity
}

func cleanResources() {
	if !shouldClean() {
		return
	}

	client, err := tests.NewZpaClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	service := New(client)
	resources, _, _ := service.GetAll()
	for _, r := range resources {
		if strings.HasPrefix(r.Name, "tests-") || strings.HasPrefix(r.Name, "updated-") {
			log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
			_, _ = service.Delete(r.ID)
		}
	}
}

func TestCBICertificates(t *testing.T) {
	// Initialize the ZPA client
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Generate private key for root certificate
	rootKey, err := rsa.GenerateKey(rand.Reader, 4096) // Use 4096 bits for the key
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Create root certificate template
	rootCertTmpl := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{"US"},
			Province:           []string{"California"},
			Locality:           []string{"San Jose"},
			Organization:       []string{"BD-HashiCorp"},
			OrganizationalUnit: []string{"ITDepartment"},
			CommonName:         "bd-hashicorp.com",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 10 years validity
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Create the root certificate
	rootCertBytes, err := x509.CreateCertificate(rand.Reader, &rootCertTmpl, &rootCertTmpl, &rootKey.PublicKey, rootKey)
	if err != nil {
		t.Fatalf("Failed to create root certificate: %v", err)
	}

	// Encode the root certificate to PEM
	rootCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootCertBytes})

	// Generate a random name for the certificate name
	randomName, err := generateRandomString(10)
	if err != nil {
		t.Fatalf("Failed to generate random string for certificate name: %v", err)
	}
	certName := fmt.Sprintf("test-rootCA %s", randomName)

	// Create the certificate object
	service := New(client)
	cbiCertificate := CBICertificate{
		PEM:  string(rootCertPEM),
		Name: certName,
	}

	// Test 1: Upload Certificate with Invalid Data
	t.Run("TestInvalidCertificateUpload", func(t *testing.T) {
		invalidCert := CBICertificate{
			PEM:  "", // Invalid as it's empty
			Name: "invalid-cert",
		}
		_, _, err := service.Create(&invalidCert)
		if err == nil {
			t.Errorf("Expected error while uploading invalid certificate, got nil")
		}
	})

	// Upload the certificate
	createdCert, _, err := service.Create(&cbiCertificate)
	if err != nil {
		t.Fatalf("Error uploading certificate: %v", err)
	}

	// Test for updating the certificate name
	t.Run("TestCertificateUpdate", func(t *testing.T) {
		// Generate a new random name for updating the certificate
		updatedName, err := generateRandomString(10)
		if err != nil {
			t.Fatalf("Failed to generate random string for updated certificate name: %v", err)
		}
		updatedCertName := fmt.Sprintf("updated-rootCA %s", updatedName)

		// Update the certificate with the new name
		cbiCertificate.Name = updatedCertName
		_, err = service.Update(createdCert.ID, &cbiCertificate)
		if err != nil {
			t.Fatalf("Error updating certificate: %v", err)
		}

		// Verify the update by retrieving the certificate again
		updatedCert, _, err := service.Get(createdCert.ID)
		if err != nil {
			t.Fatalf("Error retrieving updated certificate: %v", err)
		}
		if updatedCert.Name != updatedCertName {
			t.Errorf("Updated certificate name mismatch. Expected: %s, Got: %s", updatedCertName, updatedCert.Name)
		}
	})
	// Verify the upload by retrieving the certificate by ID
	retrievedCert, _, err := service.Get(createdCert.ID)
	if err != nil {
		t.Fatalf("Error retrieving uploaded certificate: %v", err)
	}
	if retrievedCert.Name != cbiCertificate.Name {

		// Verify the upload by retrieving the certificate by ID
		retrievedCert, _, err := service.Get(createdCert.ID)
		if err != nil {
			t.Fatalf("Error retrieving uploaded certificate: %v", err)
		}
		if retrievedCert.Name != cbiCertificate.Name {
			t.Errorf("Retrieved certificate name mismatch. Expected: %s, Got: %s", cbiCertificate.Name, retrievedCert.Name)
		}

		// Retrieve the certificate by name
		retrievedCertByName, _, err := service.GetByName(createdCert.Name)
		if err != nil {
			t.Fatalf("Error retrieving uploaded certificate by name: %v", err)
		}
		if retrievedCertByName.Name != cbiCertificate.Name {
			t.Errorf("Retrieved by name certificate name mismatch. Expected: %s, Got: %s", cbiCertificate.Name, retrievedCertByName.Name)
		}

		// Delete the certificate
		_, err = service.Delete(createdCert.ID)
		if err != nil {
			t.Fatalf("Error deleting certificate: %v", err)
		}

		// Test 3: Attempt Retrieval After Deletion
		t.Run("TestRetrieveAfterDeletion", func(t *testing.T) {
			_, _, err := service.Get(createdCert.ID)
			if err == nil {
				t.Errorf("Expected error while retrieving deleted certificate, got nil")
			}
		})

		// Verify deletion
		_, _, err = service.Get(createdCert.ID)
		if err == nil || !strings.Contains(err.Error(), "404") {
			t.Errorf("Certificate still exists after deletion or unexpected error: %v", err)
		}

	}
}

// generateRandomString generates a random string of the given length
func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b) // This reads len(b) random bytes into b
	if err != nil {
		return "", err
	}

	for i, byteVal := range b {
		b[i] = charset[byteVal%byte(len(charset))]
	}

	return string(b), nil
}
