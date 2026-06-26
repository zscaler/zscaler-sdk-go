package http_header_action_profile

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	headerProfileActionEndpoint = "/zia/api/v1/httpHeaderActionProfile"
)

type HttpHeaderActionProfile struct {
	ID                          int                           `json:"id,omitempty"`
	Name                        string                        `json:"name,omitempty"`
	SlotId                      int                           `json:"slotId,omitempty"`
	Deleted                     bool                          `json:"deleted,omitempty"`
	ProfileReadyForUse          bool                          `json:"profileReadyForUse,omitempty"`
	Description                 string                        `json:"description,omitempty"`
	HttpHeaderActionProfileKeys []HttpHeaderActionProfileKeys `json:"httpHeaderActionProfileKeys,omitempty"`
}

type HttpHeaderActionProfileKeys struct {
	ID    int    `json:"id,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func Create(ctx context.Context, service *zscaler.Service, alertDefinitions *HttpHeaderActionProfile) (*HttpHeaderActionProfile, *http.Response, error) {
	resp, err := service.Client.Create(ctx, headerProfileActionEndpoint, *alertDefinitions)
	if err != nil {
		return nil, nil, err
	}

	createdHeaderProfile, ok := resp.(*HttpHeaderActionProfile)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a header action profile pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new header action profile from create: %d", createdHeaderProfile.ID)
	return createdHeaderProfile, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, profileID int, profile *HttpHeaderActionProfile) (*HttpHeaderActionProfile, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", headerProfileActionEndpoint, profileID), *profile)
	if err != nil {
		return nil, nil, err
	}
	updatedHeaderProfile, _ := resp.(*HttpHeaderActionProfile)

	service.Client.GetLogger().Printf("[DEBUG]returning updates header action profile from update: %d", updatedHeaderProfile.ID)
	return updatedHeaderProfile, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, profileID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", headerProfileActionEndpoint, profileID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, profileName string) (*HttpHeaderActionProfile, error) {
	headerProfiles, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	// Search for exact match (case-insensitive)
	for _, headerProfile := range headerProfiles {
		if strings.EqualFold(headerProfile.Name, profileName) {
			return &headerProfile, nil
		}
	}
	return nil, fmt.Errorf("no header profile found with name: %s", profileName)
}

// GetAll retrieves all HTTP header action profiles.
//
// This endpoint returns a flat list and does not support pagination, so it is
// read directly without the page/pageSize pagination helper.
func GetAll(ctx context.Context, service *zscaler.Service) ([]HttpHeaderActionProfile, error) {
	var headerProfiles []HttpHeaderActionProfile
	err := service.Client.Read(ctx, headerProfileActionEndpoint, &headerProfiles)
	return headerProfiles, err
}
