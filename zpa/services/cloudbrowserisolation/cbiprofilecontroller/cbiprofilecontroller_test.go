package cbiprofilecontroller

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbibannercontroller"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/cloudbrowserisolation/cbiregions"
)

func TestCBIProfileController(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	service := services.New(client)

	cbiRegionsList, _, err := cbiregions.GetAll(service)
	if err != nil {
		t.Errorf("Error getting cbi regions: %v", err)
		return
	}
	if len(cbiRegionsList) == 0 {
		t.Error("Expected retrieved cbi regions to be non-empty, but got empty slice")
	}

	// create application connector group for testing
	cbiBanner := cbibannercontroller.CBIBannerController{
		Name:              name,
		PrimaryColor:      "#0076BE",
		TextColor:         "#FFFFFF",
		NotificationTitle: "Heads up, you’ve been redirected to Browser Isolation!",
		NotificationText:  "The website you were trying to access is now rendered in a fully isolated environment to protect you from malicious content.",
		Banner:            true,
		Persist:           true,
		Logo:              "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAYQAAABQCAMAAAAuu/JsAAADAFBMVEUAAAAAgL8Ad8MAdr8Ad78Ad78Ad78Adr4Ad74Adr4Ad74Adr8Adr4Ad74Ad74Adr4Adr4Ad74Ad78AeL8Ad8AAd8AAeb4AgL8AktsAesAAd74Adr8Adr4Adr4Adr4Ad78AeMAAgMYAeL8Ad74Adr8Ad74AgP8Ae8EAd8AAdr8Ad8EA//8Ad78Ad74Aeb4Ae8YAdr4Ad74Adr8Ad74Ad78Adr8Ad78Adr4Ad78Adr8Ad78Ad74Ad8EAd78Adr8Ad78AeL4AecAAeckAeL8Ad78Ad74Ad74Ad78AecIAdr8Ad78Ad74Aer4AicQAd74Ad74AecAAi9EAfMEAd8AAeL8Ad78Ad78Ad78Ad74Ad78Ad8QAdsIAgMUAeL8Ad78Ad78AfMEAd8AAd78Ad78Ad78Adr8Ad74Adr8Ad74AeL8AeL4Ad74Adr4Aer8AecAAdr4AeMEAqv8Adr4Adr8Adr8AeL8Ad74Ad74Adr8Adr8Ad78AdsAAd78AeMEAdr8AgL8Aeb8AdsEAd78Adr8Ad78AgMgAd8AAeL4Ad78Adr4Ad78AecIAe78AeL8AgMQAd74Adr8AeMAAdr8AecEAgL8AeL4AmcwAdr8Adr8Ad74Ad78AesgAdr8AeL8AfMQAeL8AjsYAd78Ad78AecMAd78AdsEAeMAAd78Ad78AesIAdr8Ad74Aeb8AeMEAeMMAd8QAd8wAeL8Ad74Ad78AdsAAfMUAd78Adr4AeMMAgNUAgMwAd78Ad8AAd74AdsAAeL4Adr8Ad78Aeb4AesIAeL8Adr8AeL8Ad78AgL8Adr8Ad8EAdr8Ad78Ad74Adr8Ad74Adr4AeL4Ad78Ad78Adr8Ad74Ad78Ad78Ad74AgL8AeMAAdr8Ad8EAd74AeMAAd74Ad78Ad8AAdsAAd74Adr8AesIAdr8Adr8AeL8Ae8UAeb8AeMEAeL8AeMAAd78Ad8EAdsAAdr8AecEAer8Adr4Adr4Ad74Adr8Ad8AAd78Adr8AeL8Ad74Ad8EAeMAAdr5kTMBRAAABAHRSTlMABC9UeJy/0d3p9f78+fHl2My7oIVpPxAHQX689P/go1kSSPC3SwIdjeteAfvJNxum/X966tPIsa7C16Ut1s+QZj0TaPfknVgVqPOhQw29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",
	}

	cbiBannerController, _, err := cbibannercontroller.Create(service, &cbiBanner)
	if err != nil || cbiBannerController == nil || cbiBannerController.ID == "" {
		t.Fatalf("Error creating cbi banner controller or ID is empty")
		return
	}

	defer func() {
		if cbiBannerController != nil && cbiBannerController.ID != "" {
			existingCbiBanner, _, errCheck := cbibannercontroller.Get(service, cbiBannerController.ID)
			if errCheck == nil && existingCbiBanner != nil {
				_, errDelete := cbibannercontroller.Delete(service, cbiBannerController.ID)
				if errDelete != nil {
					t.Errorf("Error deleting cbi banner controller: %v", errDelete)
				}
			}
		}
	}()

	cbiProfile := IsolationProfile{
		Name:        name,
		Description: name,
		BannerID:    cbiBannerController.ID,
		RegionIDs:   []string{cbiRegionsList[0].ID, cbiRegionsList[1].ID},
		UserExperience: &UserExperience{
			SessionPersistence: true,
			BrowserInBrowser:   true,
		},
		SecurityControls: &SecurityControls{
			CopyPaste:          "all",
			UploadDownload:     "all",
			DocumentViewer:     true,
			LocalRender:        true,
			AllowPrinting:      true,
			RestrictKeystrokes: false,
		},
	}

	createdResource, _, err := Create(service, &cbiProfile)
	if err != nil || createdResource == nil {
		t.Fatalf("Error making POST request: %v or createdResource is nil", err)
	}

	// Fetch the resource again to get full details
	createdResource, _, err = Get(service, createdResource.ID)
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
	retrievedResource, _, err := Get(service, createdResource.ID)
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
	_, err = Update(service, createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := Get(service, createdResource.ID)
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
	retrievedResource, _, err = GetByName(service, updateName)
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
	resources, _, err := GetAll(service)
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
	_, err = Delete(service, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = Get(service, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}

/*
func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.Get("non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Delete("non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, err = service.Update("non_existent_id", &IsolationProfile{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	_, _, err = service.GetByName("non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
*/
