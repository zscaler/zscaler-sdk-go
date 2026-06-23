package adaptive_access

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	adaptiveAccessEndpoint = "/zia/api/v1/adaptiveAccessProfiles"
	profileRulesHandler    = "/profiles/rules"
)

type AdaptiveAccess struct {
	// The Adaptive Access profile ID
	ID int `json:"id,omitempty"`

	// The Adaptive Access profile name
	Name string `json:"name,omitempty"`

	// The Adaptive Access profile type
	Type string `json:"type,omitempty"`

	// The Adaptive Access profile index
	AapIndex int `json:"aapIndex,omitempty"`

	// The Adaptive Access profile ID that is used by the API for policy configuration.
	// This field allows you to specify which Adaptive Access profiles are applied in the access policy criteria.
	IamAapID string `json:"iamAapId,omitempty"`

	// A Boolean value that indicates whether the Adaptive Access profile is deleted
	Deleted bool `json:"deleted,omitempty"`
}

type GetFilterOptions struct {
	// Filters based on the Adaptive Access profile IDs.
	IAMAapIDs []string
	// Filters based on the organization ID.
	OrgID *int
}

func GetByName(ctx context.Context, service *zscaler.Service, profileName string) (*AdaptiveAccess, error) {
	adaptiveAccessProfiles, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	for _, adaptiveAccess := range adaptiveAccessProfiles {
		if strings.EqualFold(adaptiveAccess.Name, profileName) {
			return &adaptiveAccess, nil
		}
	}
	return nil, fmt.Errorf("no adaptive access profile found with name: %s", profileName)
}

// GetProfileRules retrieves the Adaptive Access profile information based on
// the optional filtering parameters. This endpoint is not paginated; it returns
// the results directly, optionally filtered by Adaptive Access profile IDs
// (iamAapIds) and/or organization ID (orgId).
func GetProfileRules(ctx context.Context, service *zscaler.Service, opts *GetFilterOptions) ([]AdaptiveAccess, error) {
	var adaptiveAccessProfiles []AdaptiveAccess
	endpoint := adaptiveAccessEndpoint + profileRulesHandler

	queryParams := url.Values{}
	if opts != nil {
		for _, id := range opts.IAMAapIDs {
			queryParams.Add("iamAapIds", id)
		}
		if opts.OrgID != nil {
			queryParams.Set("orgId", strconv.Itoa(*opts.OrgID))
		}
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	err := service.Client.Read(ctx, endpoint, &adaptiveAccessProfiles)
	return adaptiveAccessProfiles, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]AdaptiveAccess, error) {
	var adaptiveAccessProfiles []AdaptiveAccess
	err := common.ReadAllPages(ctx, service.Client, adaptiveAccessEndpoint, &adaptiveAccessProfiles)
	return adaptiveAccessProfiles, err
}
