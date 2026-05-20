package zia_posture

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	ziaPostureEndpointV2 = "/zcc/papi/public/v2/zia-posture-profiles"
)

// GetAllFilterOptions models the documented optional query parameters for
// GET /zcc/papi/public/v2/zia-posture-profiles. Pagination (skip/perPage)
// is handled by the pagination helper; callers only supply filters here.
type GetAllFilterOptions struct {
	// Keyword filters records by name (substring match).
	Keyword string
	// PlatformType filters by device platform (0 = all platforms).
	// Use the common.DeviceType* constants (1=iOS, 2=Android, 3=Windows,
	// 4=macOS, 5=Linux).
	PlatformType int
}

type ZIAPosture struct {
	ID                  int                 `json:"id,omitempty"`
	Name                string              `json:"name,omitempty"`
	Platform            int                 `json:"platform,omitempty"`
	HighTrustCriteria   HighTrustCriteria   `json:"highTrustCriteria,omitempty"`
	MediumTrustCriteria MediumTrustCriteria `json:"mediumTrustCriteria,omitempty"`
	LowTrustCriteria    LowTrustCriteria    `json:"lowTrustCriteria,omitempty"`
}

type HighTrustCriteria struct {
	Cs []TrustCriteriaSet `json:"cs,omitempty"`
}

type MediumTrustCriteria struct {
	Cs []TrustCriteriaSet `json:"cs,omitempty"`
}

type LowTrustCriteria struct {
	Cs []TrustCriteriaSet `json:"cs,omitempty"`
}

type TrustCriteriaSet struct {
	Cn []TrustCriterion `json:"cn,omitempty"`
}

type TrustCriterion struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	UDID string `json:"udid,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, postureID int) (*ZIAPosture, error) {
	var ziaPosture ZIAPosture
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", ziaPostureEndpointV2, postureID), &ziaPosture)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning zia posture from Get: %d", ziaPosture.ID)
	return &ziaPosture, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, postureName string) (*ZIAPosture, error) {
	// Narrow the server-side search via keyword; final exact match is
	// still done client-side because keyword is substring, not equality.
	postures, err := GetAll(ctx, service, &GetAllFilterOptions{Keyword: postureName})
	if err != nil {
		return nil, err
	}
	for _, posture := range postures {
		if strings.EqualFold(posture.Name, postureName) {
			return &posture, nil
		}
	}
	return nil, fmt.Errorf("no zia posture found with name: %s", postureName)
}

func Create(ctx context.Context, service *zscaler.Service, posture *ZIAPosture) (*ZIAPosture, *http.Response, error) {
	if posture == nil {
		return nil, nil, errors.New("zia posture is required")
	}

	var created ZIAPosture
	resp, err := service.Client.NewZccRequestDo(ctx, "POST", ziaPostureEndpointV2, nil, posture, &created)
	if err != nil {
		return nil, resp, err
	}

	service.Client.GetLogger().Printf("[DEBUG] returning new zia posture from create: %d", created.ID)
	return &created, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, postureID int, posture *ZIAPosture) (*ZIAPosture, *http.Response, error) {
	if posture == nil {
		return nil, nil, errors.New("zia posture is required")
	}

	endpoint := fmt.Sprintf("%s/%d", ziaPostureEndpointV2, postureID)
	var updated ZIAPosture
	resp, err := service.Client.NewZccRequestDo(ctx, "PUT", endpoint, nil, posture, &updated)
	if err != nil {
		return nil, resp, err
	}

	if updated.ID == 0 {
		updated.ID = postureID
	}

	service.Client.GetLogger().Printf("[DEBUG] returning updated zia posture from update: %d", updated.ID)
	return &updated, resp, nil
}

func PartialUpdate(ctx context.Context, service *zscaler.Service, postureID int, posture *ZIAPosture) (*ZIAPosture, *http.Response, error) {
	if posture == nil {
		return nil, nil, errors.New("zia posture is required")
	}

	endpoint := fmt.Sprintf("%s/%d", ziaPostureEndpointV2, postureID)
	var updated ZIAPosture
	resp, err := service.Client.NewZccRequestDo(ctx, "PATCH", endpoint, nil, posture, &updated)
	if err != nil {
		return nil, resp, err
	}
	if updated.ID == 0 {
		updated.ID = postureID
	}

	service.Client.GetLogger().Printf("[DEBUG] returning zia posture from partial update: %d", updated.ID)
	return &updated, resp, nil
}

func Delete(ctx context.Context, service *zscaler.Service, postureID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", ziaPostureEndpointV2, postureID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) ([]ZIAPosture, error) {
	params := common.QueryParamsV2{}
	if opts != nil {
		params.Keyword = opts.Keyword
		params.PlatformType = opts.PlatformType
	}
	return common.ReadAllPagesV2[ZIAPosture](ctx, service.Client, ziaPostureEndpointV2, params, common.DefaultPageSize)
}
