package cloudbrowserisolation

import (
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	cbiProfileEndpoint = "/browserIsolation/profiles"
)

type IsolationProfile struct {
	ID string `json:"id,omitempty"`

	// Name of the browser isolation profile
	Name string `json:"name,omitempty"`

	// The browser isolation profile URL
	URL string `json:"url,omitempty"`

	// (Optional) Indicates whether this is a default browser isolation profile. Zscaler sets this field.
	DefaultProfile bool `json:"defaultProfile,omitempty"`
}

// Updated GetByName function
func (service *Service) GetByName(profileName string) (*IsolationProfile, error) {
	var cbiProfiles []IsolationProfile
	err := common.ReadAllPages(service.Client, cbiProfileEndpoint, &cbiProfiles)
	if err != nil {
		return nil, checkNotSubscribedError(err)
	}
	for _, cbi := range cbiProfiles {
		if strings.EqualFold(cbi.Name, profileName) {
			return &cbi, nil
		}
	}
	return nil, fmt.Errorf("no cloud browser isolation profile found with name: %s", profileName)
}

// Updated GetAll function
func (service *Service) GetAll() ([]IsolationProfile, error) {
	var cbiProfiles []IsolationProfile
	err := common.ReadAllPages(service.Client, cbiProfileEndpoint, &cbiProfiles)
	return cbiProfiles, checkNotSubscribedError(err)
}

type NotSubscribedError struct {
	message string
}

func (e *NotSubscribedError) Error() string {
	return e.message
}

// Helper function to check and wrap the "Not Subscribed" error
func checkNotSubscribedError(err error) error {
	if err != nil && strings.Contains(err.Error(), "Cloud Browser Isolation subscription is required") {
		return &NotSubscribedError{message: "NOT_SUBSCRIBED: Cloud Browser Isolation subscription is required"}
	}
	return err
}
