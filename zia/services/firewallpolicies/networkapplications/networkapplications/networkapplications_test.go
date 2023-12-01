package networkapplications

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestNetworkApplications(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Fetching the first page of network applications
	nwApplications, err := service.GetFirstPage("")
	if err != nil {
		t.Errorf("Error getting network applications: %v", err)
		return
	}
	if len(nwApplications) == 0 {
		t.Errorf("No network application found")
		return
	}

	// Selecting one application to test GetByName
	nwApplicationName := nwApplications[0].ID
	locale := "en-US" // Replace with the desired locale
	t.Log("Getting network application by name: " + nwApplicationName)

	// Testing GetByName with the selected application and locale
	nwApplication, err := service.GetByName(nwApplicationName, locale)
	if err != nil {
		t.Errorf("Error getting network application by name: %v", err)
		return
	}
	if nwApplication.ID != nwApplicationName {
		t.Errorf("Network application ID does not match: expected %s, got %s", nwApplicationName, nwApplication.ID)
		return
	}

	// Negative Test: Try to retrieve a network application with a non-existent name
	nonExistentName := "ThisApplicationDoesNotExist"
	_, err = service.GetByName(nonExistentName, locale)
	if err == nil {
		t.Errorf("Expected error when getting by non-existent name, got nil")
		return
	}
}

func TestFilteringByParentCategory(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := New(client)

	// Fetching only the first page of network applications
	nwApplications, err := service.GetFirstPage("")
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
	filteredApplication, err := service.GetByName(filterCategory, locale)
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

	service := New(client)
	locales := []string{"en-US", "de-DE", "es-ES", "fr-FR", "ja-JP", "zh-CN"}

	for _, locale := range locales {
		t.Run("Locale: "+locale, func(t *testing.T) {
			// Using GetFirstPage instead of GetAll
			applications, err := service.GetFirstPage(locale)
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

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := New(client)

	// Fetching only the first page of network applications
	nwApplications, err := service.GetFirstPage("")
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
func (service *Service) GetFirstPage(locale string) ([]NetworkApplications, error) {
	var networkApplications []NetworkApplications
	endpoint := networkApplicationsEndpoint
	if locale != "" {
		endpoint = fmt.Sprintf("%s?locale=%s", networkApplicationsEndpoint, url.QueryEscape(locale))
	}
	err := service.Client.Read(endpoint, &networkApplications)
	return networkApplications, err
}
