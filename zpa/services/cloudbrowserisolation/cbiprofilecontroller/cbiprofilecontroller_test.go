package cbiprofilecontroller

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/cloudbrowserisolation/cbibannercontroller"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/cloudbrowserisolation/cbicertificatecontroller"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/cloudbrowserisolation/cbiregions"
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
		if !strings.HasPrefix(r.Name, "tests-") {
			continue
		}
		log.Printf("Deleting resource with ID: %s, Name: %s", r.ID, r.Name)
		_, _ = service.Delete(r.ID)
	}
}

func readFileContent(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	// Prepend the data URI prefix
	return "data:image/png;base64," + string(data), nil
}
func TestCBIProfileController(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	cbiRegionsService := cbiregions.New(client)
	cbiRegionsList, _, err := cbiRegionsService.GetAll()
	if err != nil {
		t.Errorf("Error getting cbi regions: %v", err)
		return
	}
	if len(cbiRegionsList) == 0 {
		t.Error("Expected retrieved cbi regions to be non-empty, but got empty slice")
	}

	cbiCertificateService := cbicertificatecontroller.New(client)
	cbiCertificate, _, err := cbiCertificateService.GetByName("Zscaler Root Certificate")
	if err != nil {
		t.Errorf("Error getting cbi certificate: %v", err)
		return
	}
	if cbiCertificate == nil {
		t.Error("Expected to retrieve a cbi certificate, but got nil")
	}

	cbiLogo, err := readFileContent("../cbibannercontroller/cbiLogo")
	if err != nil {
		t.Fatalf("Error reading CBI Banner content: %v", err)
	}
	// create application connector group for testing
	cbiBannerService := cbibannercontroller.New(client)
	cbiBanner := cbibannercontroller.CBIBannerController{
		Name:              name,
		PrimaryColor:      "#0076BE",
		TextColor:         "#FFFFFF",
		NotificationTitle: "Heads up, you’ve been redirected to Browser Isolation!",
		NotificationText:  "The website you were trying to access is now rendered in a fully isolated environment to protect you from malicious content.",
		Banner:            true,
		Persist:           true,
		Logo:              cbiLogo,
	}

	cbiBannerController, _, err := cbiBannerService.Create(&cbiBanner)
	if err != nil || cbiBannerController == nil || cbiBannerController.ID == "" {
		t.Fatalf("Error creating cbi banner controller or ID is empty")
		return
	}

	defer func() {
		if cbiBannerController != nil && cbiBannerController.ID != "" {
			existingCbiBanner, _, errCheck := cbiBannerService.Get(cbiBannerController.ID)
			if errCheck == nil && existingCbiBanner != nil {
				_, errDelete := cbiBannerService.Delete(cbiBannerController.ID)
				if errDelete != nil {
					t.Errorf("Error deleting cbi banner controller: %v", errDelete)
				}
			}
		}
	}()

	service := New(client)
	cbiProfile := IsolationProfile{
		Name:           name,
		Description:    name,
		BannerID:       cbiBannerController.ID,
		RegionIDs:      []string{cbiRegionsList[0].ID, cbiRegionsList[1].ID}, // Ensure at least 3 regions are assigned
		CertificateIDs: []string{cbiCertificate.ID},
		UserExperience: &UserExperience{ // <--- This seems to be a struct literal; you might be missing a field name here
			SessionPersistence: true,
			BrowserInBrowser:   true,
		},
		SecurityControls: &SecurityControls{ // <--- Similarly, this might be missing a field name
			CopyPaste:          "all",
			UploadDownload:     "all",
			DocumentViewer:     true,
			LocalRender:        true,
			AllowPrinting:      true,
			RestrictKeystrokes: false,
		},
	}

	createdResource, _, err := service.Create(&cbiProfile)
	if err != nil || createdResource == nil {
		t.Fatalf("Error making POST request: %v or createdResource is nil", err)
	}

	// Fetch the resource again to get full details
	createdResource, _, err = service.Get(createdResource.ID)
	if err != nil {
		t.Fatalf("Error fetching the created resource: %v", err)
	}
	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}

	// Test resource retrieval
	retrievedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource update
	retrievedResource.Name = updateName
	_, err = service.Update(createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}
	// Test resource retrieval by name
	retrievedResource, _, err = service.GetByName(updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.Name)
	}
	// Test resources retrieval
	resources, _, err := service.GetAll()
	if err != nil {
		t.Errorf("Error retrieving groups: %v", err)
	}
	if len(resources) == 0 {
		t.Error("Expected retrieved resources to be non-empty, but got empty slice")
	}
	// check if the created resource is in the list
	found := false
	for _, resource := range resources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved groups to contain created resource '%s', but it didn't", createdResource.ID)
	}
	// Test resource removal
	_, err = service.Delete(createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = service.Get(createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}
