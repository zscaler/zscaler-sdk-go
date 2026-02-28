package email_profiles

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	emailProfilesEndpoint = "/zia/api/v1/emailRecipientProfile"
)

type EmailProfiles struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Emails      []string `json:"emails,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, profileID int) (*EmailProfiles, error) {
	var emailProfile EmailProfiles
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", emailProfilesEndpoint, profileID), &emailProfile)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning email profilefrom Get: %d", emailProfile.ID)
	return &emailProfile, nil
}

// GetAllFilterOptions represents optional filter parameters for GetAll.
// Page and pageSize are handled internally by common.ReadAllPages.
type GetAllFilterOptions struct {
	// Filters profiles based on the recipient profile name.
	Name *string
	// Filters based on the recipient email included in each profile.
	EmailRecipient *string
	// Filters profiles based on the description field.
	Description *string
}

// GetAllLiteFilterOptions represents optional filter parameters for GetAllLite.
// Page and pageSize are handled internally by common.ReadAllPages.
// Consumer values: EMAIL_DLP, EXTERNAL_DLP, INTERNAL_DLP, ALL_DLP_CONSUMER_TYPES.
type GetAllLiteFilterOptions struct {
	// Filters profiles based on the recipient profile name.
	Name *string
	// Filters profiles based on the description field.
	Description *string
	// Filters profiles based on the type of consumer.
	Consumer *string
	// Filters profiles based on tenant ID(s).
	TenantIDs []int
}

func GetEmailProfileByName(ctx context.Context, service *zscaler.Service, profileName string) (*EmailProfiles, error) {
	// Use GetAll with name filter to leverage API filtering
	opts := &GetAllFilterOptions{Name: &profileName}
	emailProfiles, err := GetAll(ctx, service, opts)
	if err != nil {
		return nil, err
	}
	// Search for exact match (case-insensitive)
	for _, emailProfile := range emailProfiles {
		if strings.EqualFold(emailProfile.Name, profileName) {
			return &emailProfile, nil
		}
	}
	return nil, fmt.Errorf("no email profile found with name: %s", profileName)
}

func Create(ctx context.Context, service *zscaler.Service, profiles *EmailProfiles) (*EmailProfiles, *http.Response, error) {
	resp, err := service.Client.Create(ctx, emailProfilesEndpoint, *profiles)
	if err != nil {
		return nil, nil, err
	}

	createdEmailProfile, ok := resp.(*EmailProfiles)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a email profile pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new email profile from create: %d", createdEmailProfile.ID)
	return createdEmailProfile, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, profileID int, profiles *EmailProfiles) (*EmailProfiles, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", emailProfilesEndpoint, profileID), *profiles)
	if err != nil {
		return nil, nil, err
	}
	updatedProfile, _ := resp.(*EmailProfiles)

	service.Client.GetLogger().Printf("[DEBUG]returning updates email profile from update: %d", updatedProfile.ID)
	return updatedProfile, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, profileID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", emailProfilesEndpoint, profileID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAllLite retrieves all email recipient profiles in lite format with optional filters.
// Page and pageSize are handled by common.ReadAllPages.
func GetAllLite(ctx context.Context, service *zscaler.Service, opts *GetAllLiteFilterOptions) ([]EmailProfiles, error) {
	var emailProfiles []EmailProfiles
	endpoint := emailProfilesEndpoint + "/lite"

	queryParams := url.Values{}
	if opts != nil {
		if opts.Name != nil && *opts.Name != "" {
			queryParams.Set("name", *opts.Name)
		}
		if opts.Description != nil && *opts.Description != "" {
			queryParams.Set("description", *opts.Description)
		}
		if opts.Consumer != nil && *opts.Consumer != "" {
			queryParams.Set("consumer", *opts.Consumer)
		}
		for _, id := range opts.TenantIDs {
			queryParams.Add("tenantIds", strconv.Itoa(id))
		}
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &emailProfiles)
	return emailProfiles, err
}

// GetAll retrieves all email recipient profiles with optional filters.
// Page and pageSize are handled by common.ReadAllPages.
func GetAll(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) ([]EmailProfiles, error) {
	var emailProfiles []EmailProfiles
	endpoint := emailProfilesEndpoint

	queryParams := url.Values{}
	if opts != nil {
		if opts.Name != nil && *opts.Name != "" {
			queryParams.Set("name", *opts.Name)
		}
		if opts.EmailRecipient != nil && *opts.EmailRecipient != "" {
			queryParams.Set("emailRecipient", *opts.EmailRecipient)
		}
		if opts.Description != nil && *opts.Description != "" {
			queryParams.Set("description", *opts.Description)
		}
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &emailProfiles)
	return emailProfiles, err
}

// GetCount retrieves the count of recipient email profiles with optional filters.
func GetCount(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) (int, error) {
	var count int
	endpoint := emailProfilesEndpoint + "/count"

	queryParams := url.Values{}
	if opts != nil {
		if opts.Name != nil && *opts.Name != "" {
			queryParams.Set("name", *opts.Name)
		}
		if opts.EmailRecipient != nil && *opts.EmailRecipient != "" {
			queryParams.Set("emailRecipient", *opts.EmailRecipient)
		}
		if opts.Description != nil && *opts.Description != "" {
			queryParams.Set("description", *opts.Description)
		}
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	err := service.Client.Read(ctx, endpoint, &count)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve email recipient profile count: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning email recipient profile count: %d", count)
	return count, nil
}
