package http_header_profile

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	headerProfileEndpoint = "/zia/api/v1/httpHeaderProfile"
)

type HttpHeaderProfile struct {
	ID                        int                         `json:"id,omitempty"`
	Name                      string                      `json:"name,omitempty"`
	Description               string                      `json:"description,omitempty"`
	SlotId                    int                         `json:"slotId,omitempty"`
	Deleted                   bool                        `json:"deleted,omitempty"`
	ProfileReadyForUse        bool                        `json:"profileReadyForUse,omitempty"`
	HttpHeaderProfileCriteria []HttpHeaderProfileCriteria `json:"httpHeaderProfileCriteria,omitempty"`
}

type HttpHeaderProfileCriteria struct {
	Id               int      `json:"id,omitempty"`
	Header           string   `json:"header,omitempty"`
	Operator         string   `json:"operator,omitempty"`
	CategoryBitmap   []string `json:"categoryBitmap,omitempty"`
	CloudAppBitmap   []string `json:"cloudAppBitmap,omitempty"`
	UserAgent        string   `json:"userAgent,omitempty"`
	UserAgentBitmap  string   `json:"userAgentBitmap,omitempty"`
	UserAgentVersion string   `json:"userAgentVersion,omitempty"`
}

type CategoryBitmap struct {
	URLSupercategory   string `json:"urlSupercategory,omitempty"`
	Deprecated         bool   `json:"deprecated,omitempty"`
	BackendName        string `json:"backendName,omitempty"`
	Name               string `json:"name,omitempty"`
	UserConfiguredName string `json:"userConfiguredName,omitempty"`
	Comments           string `json:"comments,omitempty"`
}

type CloudAppBitmap struct {
	Val                 int    `json:"val"`
	WebApplicationClass string `json:"webApplicationClass"`
	BackendName         string `json:"backendName"`
	OriginalName        string `json:"originalName"`
	Name                string `json:"name"`
	Deprecated          bool   `json:"deprecated"`
	Misc                bool   `json:"misc"`
	AppNotReady         bool   `json:"appNotReady"`
	UnderMigration      bool   `json:"underMigration"`
	AppCatModified      bool   `json:"appCatModified"`
}

func Create(ctx context.Context, service *zscaler.Service, alertDefinitions *HttpHeaderProfile) (*HttpHeaderProfile, *http.Response, error) {
	resp, err := service.Client.Create(ctx, headerProfileEndpoint, *alertDefinitions)
	if err != nil {
		return nil, nil, err
	}

	createdHeaderProfile, ok := resp.(*HttpHeaderProfile)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a header profile pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new header profile from create: %d", createdHeaderProfile.ID)
	return createdHeaderProfile, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, profileID int, profile *HttpHeaderProfile) (*HttpHeaderProfile, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", headerProfileEndpoint, profileID), *profile)
	if err != nil {
		return nil, nil, err
	}
	updatedHeaderProfile, _ := resp.(*HttpHeaderProfile)

	service.Client.GetLogger().Printf("[DEBUG]returning updates header profile from update: %d", updatedHeaderProfile.ID)
	return updatedHeaderProfile, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, profileID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", headerProfileEndpoint, profileID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, profileName string) (*HttpHeaderProfile, error) {
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

// GetAll retrieves all HTTP header profiles.
//
// This endpoint returns a flat list and does not support pagination, so it is
// read directly without the page/pageSize pagination helper.
func GetAll(ctx context.Context, service *zscaler.Service) ([]HttpHeaderProfile, error) {
	var headerProfiles []HttpHeaderProfile
	err := service.Client.Read(ctx, headerProfileEndpoint, &headerProfiles)
	return headerProfiles, err
}
