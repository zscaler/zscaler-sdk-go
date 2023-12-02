package sandbox_submission

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestSandboxSubmission(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := &Service{Client: client}

	// Define the URL of the file to be tested
	fileURL := "https://github.com/SecurityGeekIO/malware-samples/raw/main/test-pe-file.exe"

	// Download the file
	resp, err := http.Get(fileURL)
	if err != nil {
		t.Fatalf("Failed to download the test file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to download the test file: HTTP status code %d", resp.StatusCode)
	}

	// Submit the file for scanning
	scanResult, err := service.SubmitFile("test-pe-file.exe", resp.Body, "")
	if err != nil {
		t.Fatalf("Error submitting file for scanning: %v", err)
	}

	expectedFileType := "exe"
	if scanResult.FileType != expectedFileType {
		t.Errorf("File type does not match. Expected %s, got %s", expectedFileType, scanResult.FileType)
	}

	t.Logf("File submitted successfully: %+v", scanResult)
}

func TestSandboxDiscan(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := &Service{Client: client}

	// Define the URL of the file to be tested
	fileURL := "https://github.com/SecurityGeekIO/malware-samples/raw/main/test-pe-file.exe"

	// Download the file
	resp, err := http.Get(fileURL)
	if err != nil {
		t.Fatalf("Failed to download the test file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to download the test file: HTTP status code %d", resp.StatusCode)
	}

	// Submit the file for scanning
	scanResult, err := service.Discan("test-pe-file.exe", resp.Body)
	if err != nil {
		t.Fatalf("Error submitting file for scanning: %v", err)
	}

	expectedFileType := "exe"
	if scanResult.FileType != expectedFileType {
		t.Errorf("File type does not match. Expected %s, got %s", expectedFileType, scanResult.FileType)
	}

	t.Logf("File submitted successfully: %+v", scanResult)
}
