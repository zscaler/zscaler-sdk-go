package cbicertificatecontroller

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestCBICertificates(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

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
		_, _, err := Create(context.Background(), service, &invalidCert)
		if err == nil {
			t.Errorf("Expected error while uploading invalid certificate, got nil")
		}
	})

	// Upload the certificate
	createdCert, _, err := Create(context.Background(), service, &cbiCertificate)
	if err != nil {
		t.Fatalf("Error uploading certificate: %v", err)
	}

	// Test 2: Verify the certificate is present in the GetAll list
	t.Run("TestGetAllCertificates", func(t *testing.T) {
		allCerts, _, err := GetAll(context.Background(), service)
		if err != nil {
			t.Fatalf("Error retrieving all certificates: %v", err)
		}
		found := false
		for _, cert := range allCerts {
			if cert.ID == createdCert.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Certificate not found in GetAll response")
		}
	})

	// Test 3: Update the certificate name
	t.Run("TestCertificateUpdate", func(t *testing.T) {
		// Generate a new random name for updating the certificate
		updatedName, err := generateRandomString(10)
		if err != nil {
			t.Fatalf("Failed to generate random string for updated certificate name: %v", err)
		}
		updatedCertName := fmt.Sprintf("updated-rootCA %s", updatedName)

		// Update the certificate with the new name
		cbiCertificate.Name = updatedCertName
		_, err = Update(context.Background(), service, createdCert.ID, &cbiCertificate)
		if err != nil {
			t.Fatalf("Error updating certificate: %v", err)
		}

		// Verify the update by retrieving the certificate again
		updatedCert, _, err := Get(context.Background(), service, createdCert.ID)
		if err != nil {
			t.Fatalf("Error retrieving updated certificate: %v", err)
		}
		if updatedCert.Name != updatedCertName {
			t.Errorf("Updated certificate name mismatch. Expected: %s, Got: %s", updatedCertName, updatedCert.Name)
		}
	})

	// Test 4: Retrieve the certificate by name
	t.Run("TestGetByName", func(t *testing.T) {
		retrievedCertByName, _, err := GetByNameOrID(context.Background(), service, cbiCertificate.Name)
		if err != nil {
			t.Fatalf("Error retrieving uploaded certificate by name: %v", err)
		}
		if retrievedCertByName.Name != cbiCertificate.Name {
			t.Errorf("Retrieved by name certificate name mismatch. Expected: %s, Got: %s", cbiCertificate.Name, retrievedCertByName.Name)
		}
	})

	//Test 5: Delete the certificate
	// t.Run("TestDeleteCertificate", func(t *testing.T) {
	// 	resp, err := Delete(context.Background(), service, createdCert.ID)
	// 	if err != nil {
	// 		t.Fatalf("Error deleting certificate: %v", err)
	// 	}
	// 	if resp.StatusCode != http.StatusOK {
	// 		t.Errorf("Expected status 204 No Content, got %d", resp.StatusCode)
	// 	}
	// })

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
