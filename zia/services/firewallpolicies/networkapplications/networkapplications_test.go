package networkapplications

import (
	"fmt"
	"math/rand"
	"net/url"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
)

func TestNetworkApplications(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	// Fetching the first page of network applications
	nwApplications, err := GetFirstPage(service, "")
	if err != nil {
		t.Errorf("Error getting network applications: %v", err)
		return
	}
	if len(nwApplications) == 0 {
		t.Errorf("No network application found")
		return
	}

	// Selecting one application to test GetByName
	nwApplicationID := nwApplications[0].ID
	locale := "en-US" // Replace with the desired locale
	t.Log("Getting network application by ID: " + nwApplicationID)

	// Testing GetNetworkApplication with the selected application ID and locale
	nwApplication, err := GetNetworkApplication(service, nwApplicationID, locale)
	if err != nil {
		t.Errorf("Error getting network application by ID: %v", err)
		return
	}
	if nwApplication.ID != nwApplicationID {
		t.Errorf("Network application ID does not match: expected %s, got %s", nwApplicationID, nwApplication.ID)
		return
	}

	// Negative Test: Try to retrieve a network application with a non-existent ID
	nonExistentID := "ThisApplicationDoesNotExist"
	_, err = GetNetworkApplication(service, nonExistentID, locale)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent ID, got nil")
		return
	}
}

func TestFilteringByParentCategory(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	// Fetching only the first page of network applications
	nwApplications, err := GetFirstPage(service, "")
	if err != nil {
		t.Fatalf("Error getting the first page of network applications: %v", err)
	}
	if len(nwApplications) == 0 {
		t.Fatalf("No network applications found on the first page")
	}

	// Selecting a parentCategory from the first application
	filterCategory := nwApplications[0].ParentCategory
	locale := "en-US" // or any other locale you wish to use

	// Fetching applications filtered by parentCategory
	filteredApplication, err := GetByName(service, filterCategory, locale)
	if err != nil {
		t.Fatalf("Error fetching application by category: %v", err)
	}

	// Validating the filtered application
	if filteredApplication.ParentCategory != filterCategory {
		t.Errorf("Filtered application does not match the category: expected %s, got %s", filterCategory, filteredApplication.ParentCategory)
	}
}

func TestLocaleSpecificResponse(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)
	locales := []string{"en-US", "de-DE", "es-ES", "fr-FR", "ja-JP", "zh-CN"}

	for _, locale := range locales {
		t.Run("Locale: "+locale, func(t *testing.T) {
			// Using GetFirstPage instead of GetAll
			applications, err := GetFirstPage(service, locale)
			if err != nil {
				t.Errorf("Error fetching applications for locale %s: %v", locale, err)
				return
			}

			if len(applications) == 0 {
				t.Errorf("No applications found for locale %s", locale)
				return
			}

			// Here, you can add additional validations specific to the locale, if necessary
		})
	}
}

func TestDeprecatedApplications(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	nwApplications, err := GetFirstPage(service, "")
	if err != nil {
		t.Fatalf("Error getting network applications: %v", err)
	}

	var foundDeprecated bool
	for _, app := range nwApplications {
		if app.Deprecated {
			foundDeprecated = true
			t.Logf("Found deprecated application: %s", app.ID)
			break
		}
	}

	if !foundDeprecated {
		t.Logf("No deprecated applications found in the first page")
	}
}

func TestDescriptionField(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	nwApplications, err := GetFirstPage(service, "")
	if err != nil {
		t.Fatalf("Error getting network applications: %v", err)
	}

	for _, app := range nwApplications {
		if app.Description == "" {
			t.Errorf("Description is missing for application ID: %s", app.ID)
		}
	}
}

func TestInvalidLocaleResponses(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	invalidLocales := []string{"abc", "xyz", "123"}
	for _, locale := range invalidLocales {
		t.Run("Invalid Locale: "+locale, func(t *testing.T) {
			_, err := GetFirstPage(service, locale)
			if err == nil {
				t.Errorf("Expected error for invalid locale %s, but got none", locale)
			}
		})
	}
}

func TestRandomizedLocaleSpecificResponse(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	locales := []string{"en-US", "de-DE", "es-ES", "fr-FR", "ja-JP", "zh-CN"}

	for _, locale := range locales {
		t.Run("Locale: "+locale, func(t *testing.T) {
			applications, err := GetFirstPage(service, locale)
			if err != nil {
				t.Errorf("Error fetching applications for locale %s: %v", locale, err)
				return
			}

			if len(applications) == 0 {
				t.Errorf("No applications found for locale %s", locale)
				return
			}

			// Randomly select an application from the first page
			randomApp := applications[rand.Intn(len(applications))]
			t.Log("Testing application: " + randomApp.ID + " in locale " + locale)

			// Validate the randomly selected application
			// Additional test logic for the selected application can be added here
		})
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	// Fetching only the first page of network applications
	nwApplications, err := GetFirstPage(service, "")
	if err != nil {
		t.Errorf("Error getting the first page of network applications: %v", err)
		return
	}
	if len(nwApplications) == 0 {
		t.Errorf("No network applications found on the first page")
		return
	}

	// Validate the response format
	for _, nwApplication := range nwApplications {
		// Checking if essential fields are not empty
		if nwApplication.ID == "" {
			t.Errorf("network application ID is empty")
		}
	}
}

// GetFirstPage fetches the first page of network applications for a specific locale
func GetFirstPage(service *services.Service, locale string) ([]NetworkApplications, error) {
	var networkApplications []NetworkApplications
	endpoint := networkApplicationsEndpoint
	if locale != "" {
		endpoint = fmt.Sprintf("%s?locale=%s", networkApplicationsEndpoint, url.QueryEscape(locale))
	}
	err := service.Client.Read(endpoint, &networkApplications)
	return networkApplications, err
}
