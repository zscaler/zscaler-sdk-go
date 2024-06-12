package bacertificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

func TestBACertificates(t *testing.T) {
	// Initialize the ZPA client
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Certificate generation
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	randomString, err := generateRandomString(10)
	if err != nil {
		t.Fatalf("Failed to generate random string for common name: %v", err)
	}
	commonName := fmt.Sprintf("tests-%s.bd-hashicorp.com", randomString)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{"US"},
			Province:           []string{"California"},
			Locality:           []string{"San Jose"},
			Organization:       []string{"BD-HashiCorp"},
			OrganizationalUnit: []string{"ITDepartment"},
			CommonName:         commonName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	// Encode to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	// Combine key and cert
	fullCert := string(certPEM) + string(keyPEM)

	// Create the certificate object
	service := services.New(client)
	baCertificate := BaCertificate{
		CertBlob:    fullCert,
		Name:        template.Subject.CommonName,
		Description: "Integration Test Certificate",
	}

	// Test 1: Upload Certificate with Invalid Data
	t.Run("TestInvalidCertificateUpload", func(t *testing.T) {
		invalidCert := BaCertificate{
			CertBlob:    "", // Invalid as it's empty
			Name:        "invalid-cert",
			Description: "Invalid Test Certificate",
		}
		_, _, err := Create(service, invalidCert)
		if err == nil {
			t.Errorf("Expected error while uploading invalid certificate, got nil")
		}
	})

	// Upload the certificate
	createdCert, _, err := Create(service, baCertificate)
	if err != nil {
		t.Fatalf("Error uploading certificate: %v", err)
	}

	// Test 2: Retrieve Non-Existent Certificate
	t.Run("TestRetrieveNonExistentCert", func(t *testing.T) {
		_, _, err := Get(service, "non_existent_id")
		if err == nil {
			t.Errorf("Expected error while retrieving non-existent certificate, got nil")
		}
	})

	// Verify the upload by retrieving the certificate by ID
	retrievedCert, _, err := Get(service, createdCert.ID)
	if err != nil {
		t.Fatalf("Error retrieving uploaded certificate: %v", err)
	}
	if retrievedCert.Name != baCertificate.Name {
		t.Errorf("Retrieved certificate name mismatch. Expected: %s, Got: %s", baCertificate.Name, retrievedCert.Name)
	}

	// Retrieve the certificate by name
	retrievedCertByName, _, err := GetIssuedByName(service, createdCert.Name)
	if err != nil {
		t.Fatalf("Error retrieving uploaded certificate by name: %v", err)
	}
	if retrievedCertByName.Name != baCertificate.Name {
		t.Errorf("Retrieved by name certificate name mismatch. Expected: %s, Got: %s", baCertificate.Name, retrievedCertByName.Name)
	}

	// Verify GetAll function
	t.Run("TestGetAllCertificates", func(t *testing.T) {
		certificates, _, err := GetAll(service)
		if err != nil {
			t.Fatalf("Error retrieving all certificates: %v", err)
		}
		found := false
		for _, cert := range certificates {
			if cert.ID == createdCert.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Uploaded certificate not found in GetAll response")
		}
	})

	// Delete the certificate
	_, err = Delete(service, createdCert.ID)
	if err != nil {
		t.Fatalf("Error deleting certificate: %v", err)
	}

	// Test 3: Attempt Retrieval After Deletion
	t.Run("TestRetrieveAfterDeletion", func(t *testing.T) {
		_, _, err := Get(service, createdCert.ID)
		if err == nil {
			t.Errorf("Expected error while retrieving deleted certificate, got nil")
		}
	})

	// Verify deletion
	_, _, err = Get(service, createdCert.ID)
	if err == nil || !strings.Contains(err.Error(), "400") {
		t.Errorf("Certificate still exists after deletion or unexpected error: %v", err)
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
