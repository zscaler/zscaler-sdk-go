package dlp_idm_profile_lite

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestDLPIDMProfileLite_data(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	for _, activeOnly := range []bool{true, false} {
		profiles, err := GetAll(service, activeOnly)
		if err != nil {
			t.Errorf("Error getting idm profiles with activeOnly %t: %v", activeOnly, err)
			return
		}
		if len(profiles) == 0 {
			t.Errorf("No idm profile found with activeOnly %t", activeOnly)
			return
		}
		name := profiles[0].TemplateName
		t.Log("Getting idm profile by name:", name, "with activeOnly:", activeOnly)
		profile, err := GetDLPProfileLiteByName(service, name, activeOnly)
		if err != nil {
			t.Errorf("Error getting idm profile by name with activeOnly %t: %v", activeOnly, err)
			return
		}
		if profile.TemplateName != name {
			t.Errorf("idm profile name does not match with activeOnly %t: expected %s, got %s", activeOnly, name, profile.TemplateName)
			return
		}
	}

	// Negative Test remains the same
}

func TestGetDLPProfileLiteById(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := services.New(client)

	for _, activeOnly := range []bool{true, false} {
		profiles, err := GetAll(service, activeOnly)
		if err != nil {
			t.Fatalf("Error getting all IDM profiles with activeOnly %t: %v", activeOnly, err)
		}
		if len(profiles) == 0 {
			t.Fatalf("No IDM profiles found for testing with activeOnly %t", activeOnly)
		}

		t.Logf("Total IDM profiles found with activeOnly %t: %d", activeOnly, len(profiles))

		testID := profiles[0].ProfileID
		if testID == 0 {
			t.Errorf("The ProfileID of the first profile is empty with activeOnly %t", activeOnly)
		} else {
			t.Logf("Valid ProfileID found with activeOnly %t: %d", activeOnly, testID)
		}
	}
}

func TestResponseFormatValidation(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)

	for _, activeOnly := range []bool{true, false} {
		profiles, err := GetAll(service, activeOnly)
		if err != nil {
			t.Errorf("Error getting idm profile with activeOnly %t: %v", activeOnly, err)
			return
		}
		if len(profiles) == 0 {
			t.Errorf("No idm profile found with activeOnly %t", activeOnly)
			return
		}

		for _, profile := range profiles {
			if profile.ProfileID == 0 {
				t.Errorf("idm profile ID is empty with activeOnly %t", activeOnly)
			}
			if profile.TemplateName == "" {
				t.Errorf("idm profile Name is empty with activeOnly %t", activeOnly)
			}
		}
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	service := services.New(client)
	knownName := "BD_IDM_TEMPLATE01"

	variations := []string{
		strings.ToUpper(knownName),
		strings.ToLower(knownName),
		cases.Title(language.English).String(knownName),
	}

	for _, activeOnly := range []bool{true, false} {
		for _, variation := range variations {
			t.Logf("Attempting to retrieve group with name variation: %s with activeOnly %t", variation, activeOnly)
			profile, err := GetDLPProfileLiteByName(service, variation, activeOnly)
			if err != nil {
				t.Errorf("Error getting idm profile with name variation '%s' and activeOnly %t: %v", variation, activeOnly, err)
				continue
			}

			if profile.TemplateName != knownName {
				t.Errorf("Expected group name to be '%s' for variation '%s' with activeOnly %t, but got '%s'", knownName, variation, activeOnly, profile.TemplateName)
			}
		}
	}
}
