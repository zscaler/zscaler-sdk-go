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
	//The universally unique identifier (UUID) for the browser isolation profile
	ID string `json:"id,omitempty"`

	// Name of the browser isolation profile
	Name string `json:"name,omitempty"`

	// The browser isolation profile URL
	URL string `json:"url,omitempty"`

	// (Optional) Indicates whether this is a default browser isolation profile. Zscaler sets this field.
	DefaultProfile bool `json:"defaultProfile,omitempty"`
}

func (service *Service) Get(profileID int) (*IsolationProfile, error) {
	var cbiProfile IsolationProfile
	err := service.Client.Read(fmt.Sprintf("%s/%d", cbiProfileEndpoint, profileID), &cbiProfile)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning cloud browser isolation from Get: %d", cbiProfile.ID)
	return &cbiProfile, nil
}

func (service *Service) GetByName(profileName string) (*IsolationProfile, error) {
	var cbiProfiles []IsolationProfile
	err := common.ReadAllPages(service.Client, cbiProfileEndpoint, &cbiProfiles)
	if err != nil {
		return nil, err
	}
	for _, cbi := range cbiProfiles {
		if strings.EqualFold(cbi.Name, profileName) {
			return &cbi, nil
		}
	}
	return nil, fmt.Errorf("no cloud browser isolation profile found with name: %s", profileName)
}

func (service *Service) GetAll() ([]IsolationProfile, error) {
	var cbiProfiles []IsolationProfile
	err := common.ReadAllPages(service.Client, cbiProfileEndpoint, &cbiProfiles)
	return cbiProfiles, err
}
