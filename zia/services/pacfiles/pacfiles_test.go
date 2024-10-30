package pacfiles

/*
import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
)

const (
	maxRetries    = 3
	retryInterval = 2 * time.Second
)

// Constants for conflict retries
const (
	maxConflictRetries    = 5
	conflictRetryInterval = 1 * time.Second
)

// retryOnConflict retries the operation when a conflict (EDIT_LOCK_NOT_AVAILABLE) occurs.
func retryOnConflict(operation func() error) error {
	var lastErr error
	for i := 0; i < maxConflictRetries; i++ {
		lastErr = operation()
		if lastErr == nil {
			return nil
		}

		if strings.Contains(lastErr.Error(), `"code":"EDIT_LOCK_NOT_AVAILABLE"`) {
			log.Printf("Conflict error detected, retrying in %v... (Attempt %d/%d)", conflictRetryInterval, i+1, maxConflictRetries)
			time.Sleep(conflictRetryInterval)
			continue
		}

		return lastErr
	}
	return lastErr
}

// Integration test for PAC file operations
// Integration test for PAC file operations
func TestPacFiles(t *testing.T) {
	// Randomly generated names for testing
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "updated-tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	// Initialize ZIA client
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Define the PAC file object for creation
	rule := PACFileConfig{
		Name:                  name + ".pac",
		Description:           "Test PAC file description",
		Domain:                "bd-hashicorp.com",
		PACContent:            generateTestPACContent(),
		PACCommitMessage:      "Initial version via integration test", // Mandatory
		PACVersionStatus:      "DEPLOYED",                             // Mandatory
		PACVerificationStatus: "VERIFY_NOERR",                         // Mandatory
	}

	// Step 1: Create the PAC file
	var createdResource *PACFileConfig
	err = retryOnConflict(func() error {
		createdResource, err = CreatePacFile(service, &rule)
		return err
	})
	if err != nil {
		t.Fatalf("Error creating PAC file: %v", err)
	}

	// Assertions after creation
	if createdResource.ID == 0 {
		t.Fatal("Expected created resource ID to be non-empty, but got 0")
	}
	if createdResource.Name != rule.Name {
		t.Errorf("Expected created resource name '%s', but got '%s'", rule.Name, createdResource.Name)
	}

	// Step 2: Clone the PAC file twice
	firstClonedPacFile, err := CreateClonedPacFileVersion(service, createdResource.ID, createdResource.PACVersion, nil, &PACFileConfig{
		Name:                  updateName + "-cloned-1",
		Description:           "Cloned PAC file 1",
		PACCommitMessage:      "First clone via integration test",
		PACVerificationStatus: "VERIFY_NOERR",
		PACVersionStatus:      "DEPLOYED",
		PACContent:            generateTestPACContent(),
	})
	if err != nil {
		t.Fatalf("Error cloning PAC file: %v", err)
	}
	t.Logf("First cloned PAC file version: %d, Name: %s", firstClonedPacFile.PACVersion, firstClonedPacFile.Name)

	secondClonedPacFile, err := CreateClonedPacFileVersion(service, createdResource.ID, createdResource.PACVersion, nil, &PACFileConfig{
		Name:                  updateName + "-cloned-2",
		Description:           "Cloned PAC file 2",
		PACCommitMessage:      "Second clone via integration test",
		PACVerificationStatus: "VERIFY_NOERR",
		PACVersionStatus:      "DEPLOYED",
		PACContent:            generateTestPACContent(),
	})
	if err != nil {
		t.Fatalf("Error cloning second PAC file: %v", err)
	}
	t.Logf("Second cloned PAC file version: %d, Name: %s", secondClonedPacFile.PACVersion, secondClonedPacFile.Name)

	// Step 3: Validate both cloned PAC files
	_, err = ValidatePacFile(service, firstClonedPacFile.PACContent)
	if err != nil {
		t.Fatalf("Error validating first cloned PAC file content: %v", err)
	}

	_, err = ValidatePacFile(service, secondClonedPacFile.PACContent)
	if err != nil {
		t.Fatalf("Error validating second cloned PAC file content: %v", err)
	}

	// Step 4: Update the second cloned PAC file
	// Retrieve the second cloned PAC file to get its version
	retrievedSecondClonedPacFile, err := GetPacFileByName(service, secondClonedPacFile.Name)
	if err != nil {
		t.Fatalf("Error retrieving second cloned PAC file: %v", err)
	}
	secondClonedPacFile.PACVersion = retrievedSecondClonedPacFile.PACVersion // Correctly set the PAC version

	secondClonedPacFile.Name = updateName + "-updated"
	t.Logf("Updating PAC file with PACVersion: %d, Name: %s", secondClonedPacFile.PACVersion, secondClonedPacFile.Name)
	err = retryOnConflict(func() error {
		_, err = UpdatePacFile(service, createdResource.ID, secondClonedPacFile.PACVersion, "DEPLOY", secondClonedPacFile, nil)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating second cloned PAC file: %v", err)
	}

	// Step 5: Verify the update by retrieving the updated PAC file
	updatedResources, err := GetPacFileVersion(service, createdResource.ID, "")
	if err != nil {
		t.Fatalf("Error retrieving updated PAC file: %v", err)
	}
	if len(updatedResources) == 0 {
		t.Fatal("Expected at least one PAC file version after update")
	}

	// Find the second version and check if the name was updated
	var updatedResource *PACFileConfig
	for _, res := range updatedResources {
		if res.PACVersion == secondClonedPacFile.PACVersion {
			updatedResource = &res
			break
		}
	}
	if updatedResource == nil {
		t.Fatal("Expected to find updated PAC file version, but it was not found")
	}

	if updatedResource.Name != updateName+"-updated" {
		t.Errorf("Expected updated PAC file name '%s', but got '%s'", updateName+"-updated", updatedResource.Name)
	}

	// Step 6: Test GetPacFiles, GetPacFileByName, GetPacFileVersion, and GetPacVersionID
	// Retrieve by name
	retrievedByNameResource, err := GetPacFileByName(service, updatedResource.Name)
	if err != nil {
		t.Fatalf("Error retrieving PAC file by name: %v", err)
	}
	if retrievedByNameResource.ID != updatedResource.ID {
		t.Errorf("Expected retrieved PAC file ID '%d', but got '%d'", updatedResource.ID, retrievedByNameResource.ID)
	}

	// Retrieve all PAC files
	allPacFiles, err := GetPacFiles(service, "")
	if err != nil {
		t.Fatalf("Error retrieving PAC files: %v", err)
	}
	found := false
	for _, pacFile := range allPacFiles {
		if pacFile.ID == updatedResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected PAC file ID '%d' in the list, but it wasn't found", updatedResource.ID)
	}

	// Retrieve a specific PAC file version
	retrievedPacVersion, err := GetPacVersionID(service, createdResource.ID, updatedResource.PACVersion, "")
	if err != nil {
		t.Fatalf("Error retrieving specific PAC file version: %v", err)
	}
	if retrievedPacVersion.ID != updatedResource.ID {
		t.Errorf("Expected PAC file version ID '%d', but got '%d'", updatedResource.ID, retrievedPacVersion.ID)
	}

	// Step 7: Test deletion of the PAC file
	err = retryOnConflict(func() error {
		_, err = DeletePacFile(service, createdResource.ID)
		return err
	})
	if err != nil {
		t.Fatalf("Error deleting PAC file: %v", err)
	}

	// Step 8: Verify deletion by attempting to retrieve the deleted PAC file
	_, err = GetPacFileVersion(service, createdResource.ID, "")
	if err == nil {
		t.Fatalf("Expected error when retrieving deleted PAC file, but got nil")
	}
}
*/

// generateTestPACContent generates sample PAC file content for testing.
/*
// func generateTestPACContent() string {
// 	return `
// 		function FindProxyForURL(url, host) {
// 			var privateIP = /^(0|10|127|192\.168|172\.1[6789]|172\.2[0-9]|172\.3[01]|169\.254|192\.88\.99)\.[0-9.]+$/;
// 			var resolved_ip = dnsResolve(host);

// 			/* Don't send non-FQDN or private IP auths to us */
// 			if (isPlainHostName(host) || isInNet(resolved_ip, "192.0.2.0", "255.255.255.0") || privateIP.test(resolved_ip))
// 				return "DIRECT";

// 			/* FTP goes directly */
// 			if (url.substring(0, 4) == "ftp:")
// 				return "DIRECT";

// 			/* Test with ZPA */
// 			if (isInNet(resolved_ip, "100.64.0.0", "255.255.0.0"))
// 				return "DIRECT";

// 			/* Updates are directly accessible */
// 			if (((localHostOrDomainIs(host, "trust.zscaler.com")) ||
// 				(localHostOrDomainIs(host, "trust.zscaler.net")) ||
// 				(localHostOrDomainIs(host, "trust.zscalerone.net")) ||
// 				(localHostOrDomainIs(host, "trust.zscalertwo.net")) ||
// 				(localHostOrDomainIs(host, "trust.zscalerthree.net")) ||
// 				(localHostOrDomainIs(host, "trust.zscalergov.net")) ||
// 				(localHostOrDomainIs(host, "trust.zsdemo.net")) ||
// 				(localHostOrDomainIs(host, "trust.zscloud.net")) ||
// 				(localHostOrDomainIs(host, "trust.zsfalcon.net")) ||
// 				(localHostOrDomainIs(host, "trust.zdxcloud.net")) ||
// 				(localHostOrDomainIs(host, "trust.zdxpreview.net")) ||
// 				(localHostOrDomainIs(host, "trust.zdxbeta.net")) ||
// 				(localHostOrDomainIs(host, "trust.zsdevel.net")) ||
// 				(localHostOrDomainIs(host, "trust.zsbetagov.net")) ||
// 				(localHostOrDomainIs(host, "trust.zspreview.net")) ||
// 				(localHostOrDomainIs(host, "trust.zscalerten.net")) ||
// 				(localHostOrDomainIs(host, "trust.zdxten.net"))) &&
// 				(url.substring(0, 5) == "http:" || url.substring(0, 6) == "https:"))
// 				return "DIRECT";

// 			/* Default Traffic Forwarding. Forwarding to Zen on port 80, but you can use port 9400 also */
// 			return "PROXY ${GATEWAY_FX}:80; PROXY ${SECONDARY_GATEWAY_FX}:80; DIRECT";
// 		}`
// }
