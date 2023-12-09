package sandbox_submission

import (
	"net/http"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestSandboxSubmission(t *testing.T) {
	runSandboxTest(t, true)
}

func TestSandboxDiscan(t *testing.T) {
	runSandboxTest(t, false)
}

func runSandboxTest(t *testing.T, isSubmit bool) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := &Service{Client: client}

	baseURL := "https://github.com/SecurityGeekIO/malware-samples/raw/main/"
	fileNames := []string{
		"2a961d4e5a2100570c942ed20a29735b.bin",
		"327bd8a60fb54aaaba8718c890dda09d.bin",
		"7665f6ee9017276dd817d15212e99ca7.bin",
		"cefb4323ba4deb9dea94dcbe3faa139f.bin",
		"8356bd54e47b000c5fdcf8dc5f6a69fa.apk",
		"841abdc66ea1e208f63d717ebd11a5e9.apk",
		"test-pe-file.exe",
	}

	for _, fileName := range fileNames {
		fileURL := baseURL + fileName

		// Download the file
		resp, err := http.Get(fileURL)
		if err != nil {
			t.Fatalf("Failed to download the test file: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to download the test file: HTTP status code %d", resp.StatusCode)
		}

		var scanResult *ScanResult
		if isSubmit {
			scanResult, err = service.SubmitFile(fileName, resp.Body, "")
		} else {
			scanResult, err = service.Discan(fileName, resp.Body)
		}

		if err != nil {
			t.Fatalf("Error submitting file for scanning: %v", err)
		}

		t.Logf("File submitted successfully: %+v", scanResult)
	}
}
