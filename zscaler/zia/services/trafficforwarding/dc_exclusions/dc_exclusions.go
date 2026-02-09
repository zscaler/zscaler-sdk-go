package dc_exclusions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	dcExclusionsEndpoint = "/zia/api/v1/dcExclusions"
	datacentersEndpoint  = "/zia/api/v1/datacenters"
)

type DCExclusions struct {
	DcID        int                      `json:"dcid,omitempty"`
	Expired     bool                     `json:"expired,omitempty"`
	StartTime   int                      `json:"startTime,omitempty"`
	EndTime     int                      `json:"endTime,omitempty"`
	Description string                   `json:"description,omitempty"`
	DcName      *common.IDNameExtensions `json:"dcName,omitempty"`
}

type Datacenter struct {
	ID                int     `json:"id,omitempty"`
	Name              string  `json:"name,omitempty"`
	Provider          string  `json:"provider,omitempty"`
	City              string  `json:"city,omitempty"`
	Timezone          string  `json:"timezone,omitempty"`
	Lat               int     `json:"lat,omitempty"`
	Longi             int     `json:"longi,omitempty"`
	Latitude          float64 `json:"latitude,omitempty"`
	Longitude         float64 `json:"longitude,omitempty"`
	GovOnly           bool    `json:"govOnly,omitempty"`
	ThirdPartyCloud   bool    `json:"thirdPartyCloud,omitempty"`
	UploadBandwidth   int     `json:"uploadBandwidth,omitempty"`
	DownloadBandwidth int     `json:"downloadBandwidth,omitempty"`
	OwnedByCustomer   bool    `json:"ownedByCustomer,omitempty"`
	ManagedBcp        bool    `json:"managedBcp,omitempty"`
	DontPublish       bool    `json:"dontPublish,omitempty"`
	DontProvision     bool    `json:"dontProvision,omitempty"`
	NotReadyForUse    bool    `json:"notReadyForUse,omitempty"`
	ForFutureUse      bool    `json:"forFutureUse,omitempty"`
	RegionalSurcharge bool    `json:"regionalSurcharge,omitempty"`
	CreateTime        int     `json:"createTime,omitempty"`
	LastModifiedTime  int     `json:"lastModifiedTime,omitempty"`
	Virtual           bool    `json:"virtual,omitempty"`
	// Legacy field for backward compatibility
	Datacenter string `json:"datacenter,omitempty"`
}

// GetAll returns all DC exclusions. The API returns a flat list and does not support pagination.
func GetAll(ctx context.Context, service *zscaler.Service) ([]DCExclusions, error) {
	var list []DCExclusions
	err := service.Client.Read(ctx, dcExclusionsEndpoint, &list)
	return list, err
}

func GetByName(ctx context.Context, service *zscaler.Service, dcName string) (*DCExclusions, error) {
	dcExclusions, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	for i := range dcExclusions {
		if strings.EqualFold(dcExclusions[i].DcName.Name, dcName) {
			return &dcExclusions[i], nil
		}
	}
	return nil, fmt.Errorf("no dc exclusion found with name: %s", dcName)
}

// Create sends the DC exclusion as an array payload. The API expects a JSON array
// with one object containing only: description, dcid, startTime, endTime (no dcName).
func Create(ctx context.Context, service *zscaler.Service, dc *DCExclusions) (*DCExclusions, *http.Response, error) {
	// Build request body: array with one item; do not include DcName (API does not accept it on create)
	payload := []DCExclusions{{
		DcID:        dc.DcID,
		StartTime:   dc.StartTime,
		EndTime:     dc.EndTime,
		Description: dc.Description,
		// DcName and Expired are omitted
	}}
	respBody, err := service.Client.CreateWithSlicePayload(ctx, dcExclusionsEndpoint, payload)
	if err != nil {
		return nil, nil, err
	}
	if len(respBody) == 0 {
		return dc, nil, nil
	}
	// Response may be array or single object
	var created DCExclusions
	if err := json.Unmarshal(respBody, &created); err == nil && created.DcID != 0 {
		service.Client.GetLogger().Printf("[DEBUG]returning new dc exclusion from create: %d", created.DcID)
		return &created, nil, nil
	}
	var list []DCExclusions
	if err := json.Unmarshal(respBody, &list); err == nil && len(list) > 0 {
		service.Client.GetLogger().Printf("[DEBUG]returning new dc exclusion from create: %d", list[0].DcID)
		return &list[0], nil, nil
	}
	return nil, nil, errors.New("api response could not be parsed as dc exclusion")
}

// Update sends a PUT to the base dcExclusions endpoint with an array payload. The API expects
// a JSON array with one object containing: description, dcid, startTime, endTime.
func Update(ctx context.Context, service *zscaler.Service, dcExclusions *DCExclusions) (*DCExclusions, *http.Response, error) {
	payload := []DCExclusions{{
		DcID:        dcExclusions.DcID,
		StartTime:   dcExclusions.StartTime,
		EndTime:     dcExclusions.EndTime,
		Description: dcExclusions.Description,
	}}
	respBody, err := service.Client.UpdateWithSlicePayload(ctx, dcExclusionsEndpoint, payload)
	if err != nil {
		return nil, nil, err
	}
	if len(respBody) == 0 {
		return dcExclusions, nil, nil
	}
	// Response may be array or single object
	var updated DCExclusions
	if err := json.Unmarshal(respBody, &updated); err == nil && updated.DcID != 0 {
		service.Client.GetLogger().Printf("[DEBUG]returning updated dc exclusion from update: %d", updated.DcID)
		return &updated, nil, nil
	}
	var list []DCExclusions
	if err := json.Unmarshal(respBody, &list); err == nil && len(list) > 0 {
		service.Client.GetLogger().Printf("[DEBUG]returning updated dc exclusion from update: %d", list[0].DcID)
		return &list[0], nil, nil
	}
	return nil, nil, errors.New("api response could not be parsed as dc exclusion")
}

func Delete(ctx context.Context, service *zscaler.Service, dcID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", dcExclusionsEndpoint, dcID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetDatacenters returns all datacenters. The API returns a flat list; uses Read (no pagination).
func GetDatacenters(ctx context.Context, service *zscaler.Service) ([]Datacenter, error) {
	var list []Datacenter
	err := service.Client.Read(ctx, datacentersEndpoint, &list)
	return list, err
}
