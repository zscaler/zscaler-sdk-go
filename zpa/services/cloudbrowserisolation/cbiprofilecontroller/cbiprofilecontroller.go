package cbiprofilecontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	cbiConfig          = "/cbiconfig/cbi/api/customers/"
	cbiProfileEndpoint = "/profiles"
)

type IsolationProfile struct {
	ID               string            `json:"id,omitempty"`
	Name             string            `json:"name,omitempty"`
	Enabled          bool              `json:"enabled"`
	IsDefault        bool              `json:"isDefault"`
	HREF             string            `json:"href,omitempty"`
	SecurityControls *SecurityControls `json:"securityControls,omitempty"`
	Regions          []Regions         `json:"regions,omitempty"`
}

type Regions struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SecurityControls struct {
	DocumentViewer     bool   `json:"documentViewer,omitempty"`
	UploadDownload     string `json:"uploadDownload,omitempty"`
	CopyPaste          string `json:"copyPaste,omitempty"`
	LocalRender        bool   `json:"localRender,omitempty"`
	AllowPrinting      bool   `json:"allowPrinting,omitempty"`
	RestrictKeystrokes bool   `json:"restrictKeystrokes,omitempty"`
}

func (service *Service) Get(profileID string) (*IsolationProfile, *http.Response, error) {
	v := new(IsolationProfile)
	relativeURL := fmt.Sprintf("%s/%s", cbiConfig+service.Client.Config.CustomerID+cbiProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(profileName string) (*IsolationProfile, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.Config.CustomerID + cbiProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[IsolationProfile](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, profile := range list {
		if strings.EqualFold(profile.Name, profileName) {
			return &profile, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no isolation profile named '%s' was found", profileName)
}

func (service *Service) Create(cbiProfile *IsolationProfile) (*IsolationProfile, *http.Response, error) {
	v := new(IsolationProfile)
	resp, err := service.Client.NewRequestDo("POST", cbiConfig+service.Client.Config.CustomerID+cbiProfileEndpoint, nil, cbiProfile, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(segmentGroupId string, segmentGroupRequest *IsolationProfile) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.Config.CustomerID+cbiProfileEndpoint, segmentGroupId)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, segmentGroupRequest, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(profileID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.Config.CustomerID+cbiProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]IsolationProfile, *http.Response, error) {
	relativeURL := cbiProfileEndpoint + service.Client.Config.CustomerID + cbiProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[IsolationProfile](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
