package sandbox_submission

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

func TestSandboxSubmission(t *testing.T) {
	runSandboxTest(t, true)
}

func TestSandboxDiscan(t *testing.T) {
	runSandboxTest(t, false)
}

func runSandboxTest(t *testing.T, isSubmit bool) {
	// Retrieve environment variables for Zscaler Cloud and Sandbox Token
	sandboxCloud := os.Getenv("ZSCALER_SANDBOX_CLOUD")
	sandboxToken := os.Getenv("ZSCALER_SANDBOX_TOKEN")

	// Ensure environment variables are set, otherwise fail the test
	if sandboxCloud == "" {
		t.Fatalf("Environment variable ZSCALER_CLOUD is not set")
	}
	if sandboxToken == "" {
		t.Fatalf("Environment variable ZSCALER_SANDBOX_TOKEN is not set")
	}

	// Step 1: Configure Zscaler Sandbox client with cloud and token from environment variables
	config, err := zscaler.NewConfiguration(
		zscaler.WithSandboxToken(sandboxToken),
		zscaler.WithZscalerCloud(sandboxCloud), // Pass Zscaler Cloud configuration
	)
	if err != nil {
		t.Fatalf("Error creating configuration: %v", err)
	}

	// Step 2: Create the ZIA OneAPI client using the configuration
	service, err := zscaler.NewOneAPIClient(config) // Specify the service (in this case "zia")
	if err != nil {
		t.Fatalf("Error creating OneAPI client: %v", err)
	}

	baseURL := "https://github.com/zscaler/malware-samples/raw/main/"
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
			// Use the service object for the SubmitFile test
			scanResult, err = SubmitFile(context.Background(), service, fileName, resp.Body, "1") // Force set to "0"
		} else {
			// Use the service object for the Discan test
			scanResult, err = Discan(context.Background(), service, fileName, resp.Body)
		}

		if err != nil {
			t.Fatalf("Error submitting file for scanning: %v", err)
		}

		t.Logf("File submitted successfully: %+v", scanResult)
	}
}
