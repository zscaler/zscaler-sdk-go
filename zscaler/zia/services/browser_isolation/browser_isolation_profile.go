package browser_isolation

import (
	"context"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	cbiProfileEndpoint = "/zia/api/v1/browserIsolation/profiles"
)

type CBIProfile struct {
	ID string `json:"id,omitempty"`

	// Name of the browser isolation profile
	Name string `json:"name,omitempty"`

	// The browser isolation profile URL
	URL string `json:"url,omitempty"`

	// (Optional) Indicates whether this is a default browser isolation profile. Zscaler sets this field.
	DefaultProfile bool `json:"defaultProfile,omitempty"`
}

// GetAll retrieves all cloud browser isolation profiles.
// Note: This endpoint does not support pagination.
func GetAll(ctx context.Context, service *zscaler.Service) ([]CBIProfile, error) {
	var cbiProfiles []CBIProfile
	err := service.Client.Read(ctx, cbiProfileEndpoint, &cbiProfiles)
	return cbiProfiles, checkNotSubscribedError(err)
}

// GetByName retrieves a cloud browser isolation profile by name.
func GetByName(ctx context.Context, service *zscaler.Service, profileName string) (*CBIProfile, error) {
	cbiProfiles, err := GetAll(ctx, service)
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
