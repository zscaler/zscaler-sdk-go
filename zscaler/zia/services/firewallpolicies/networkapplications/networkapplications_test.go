package networkapplications

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

func TestNetworkApplications(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "networkapplications", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Test basic GetFirstPage and GetNetworkApplication
	t.Run("GetFirstPage and GetNetworkApplication", func(t *testing.T) {
		nwApplications, err := GetFirstPage(context.Background(), service, "")
		if err != nil {
			t.Errorf("Error getting network applications: %v", err)
			return
		}
		if len(nwApplications) == 0 {
			t.Errorf("No network application found")
			return
		}

		nwApplicationID := nwApplications[0].ID
		locale := "en-US"
		t.Log("Getting network application by ID: " + nwApplicationID)

		nwApplication, err := GetNetworkApplication(context.Background(), service, nwApplicationID, locale)
		if err != nil {
			t.Errorf("Error getting network application by ID: %v", err)
			return
		}
		if nwApplication.ID != nwApplicationID {
			t.Errorf("Network application ID does not match: expected %s, got %s", nwApplicationID, nwApplication.ID)
			return
		}

		// Negative Test
		nonExistentID := "ThisApplicationDoesNotExist"
		_, err = GetNetworkApplication(context.Background(), service, nonExistentID, locale)
		if err == nil {
			t.Errorf("Expected error when getting by non-existent ID, got nil")
			return
		}
	})

	// Test filtering by parent category
	t.Run("FilteringByParentCategory", func(t *testing.T) {
		nwApplications, err := GetFirstPage(context.Background(), service, "")
		if err != nil {
			t.Fatalf("Error getting the first page of network applications: %v", err)
		}
		if len(nwApplications) == 0 {
			t.Fatalf("No network applications found on the first page")
		}

		filterCategory := nwApplications[0].ParentCategory
		locale := "en-US"

		filteredApplication, err := GetByName(context.Background(), service, filterCategory, locale)
		if err != nil {
			t.Fatalf("Error fetching application by category: %v", err)
		}
		if filteredApplication.ParentCategory != filterCategory {
			t.Errorf("Filtered application does not match the category: expected %s, got %s", filterCategory, filteredApplication.ParentCategory)
		}
	})

	// Test locale specific responses
	t.Run("LocaleSpecificResponse", func(t *testing.T) {
		locales := []string{"en-US", "de-DE", "es-ES", "fr-FR", "ja-JP", "zh-CN"}
		for _, locale := range locales {
			applications, err := GetFirstPage(context.Background(), service, locale)
			if err != nil {
				t.Errorf("Error fetching applications for locale %s: %v", locale, err)
				continue
			}
			if len(applications) == 0 {
				t.Errorf("No applications found for locale %s", locale)
			}
		}
	})

	// Test deprecated applications
	t.Run("DeprecatedApplications", func(t *testing.T) {
		nwApplications, err := GetFirstPage(context.Background(), service, "")
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
	})

	// Test description field
	t.Run("DescriptionField", func(t *testing.T) {
		nwApplications, err := GetFirstPage(context.Background(), service, "")
		if err != nil {
			t.Fatalf("Error getting network applications: %v", err)
		}
		for _, app := range nwApplications {
			if app.Description == "" {
				t.Errorf("Description is missing for application ID: %s", app.ID)
			}
		}
	})

	// Test invalid locale responses
	t.Run("InvalidLocaleResponses", func(t *testing.T) {
		invalidLocales := []string{"abc", "xyz", "123"}
		for _, locale := range invalidLocales {
			_, err := GetFirstPage(context.Background(), service, locale)
			if err == nil {
				t.Errorf("Expected error for invalid locale %s, but got none", locale)
			}
		}
	})

	// Test response format validation
	t.Run("ResponseFormatValidation", func(t *testing.T) {
		nwApplications, err := GetFirstPage(context.Background(), service, "")
		if err != nil {
			t.Errorf("Error getting the first page of network applications: %v", err)
			return
		}
		if len(nwApplications) == 0 {
			t.Errorf("No network applications found on the first page")
			return
		}
		for _, nwApplication := range nwApplications {
			if nwApplication.ID == "" {
				t.Errorf("network application ID is empty")
			}
		}
	})
}

// GetFirstPage fetches the first page of network applications for a specific locale
func GetFirstPage(ctx context.Context, service *zscaler.Service, locale string) ([]NetworkApplications, error) {
	var networkApplications []NetworkApplications
	endpoint := networkApplicationsEndpoint
	if locale != "" {
		endpoint = fmt.Sprintf("%s?locale=%s", networkApplicationsEndpoint, url.QueryEscape(locale))
	}
	err := service.Client.Read(ctx, endpoint, &networkApplications)
	return networkApplications, err
}
