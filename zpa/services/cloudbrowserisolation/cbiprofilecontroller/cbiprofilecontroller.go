package cbiprofilecontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

const (
	cbiConfig          = "/cbiconfig/cbi/api/customers/"
	cbiProfileEndpoint = "/profiles"
)

type IsolationProfile struct {
	ID               string            `json:"id,omitempty"`
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	Enabled          bool              `json:"enabled,omitempty"`
	CreationTime     string            `json:"creationTime,omitempty"`
	ModifiedBy       string            `json:"modifiedBy,omitempty"`
	ModifiedTime     string            `json:"modifiedTime,omitempty"`
	CBITenantID      string            `json:"cbiTenantId,omitempty"`
	CBIProfileID     string            `json:"cbiProfileId,omitempty"`
	CBIURL           string            `json:"cbiUrl,omitempty"`
	BannerID         string            `json:"bannerId,omitempty"`
	SecurityControls *SecurityControls `json:"securityControls,omitempty"`
	IsDefault        bool              `json:"isDefault,omitempty"`
	Regions          []Regions         `json:"regions,omitempty"`
	RegionIDs        []string          `json:"regionIds,omitempty"`
	Href             string            `json:"href,omitempty"`
	UserExperience   *UserExperience   `json:"userExperience,omitempty"`
	Certificates     []Certificates    `json:"certificates,omitempty"`
	CertificateIDs   []string          `json:"certificateIds,omitempty"`
	Banner           *Banner           `json:"banner,omitempty"`
	DebugMode        *DebugMode        `json:"debugMode,omitempty"`
}

type Certificates struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty"`
}

type Banner struct {
	ID string `json:"id,omitempty"`
}

type UserExperience struct {
	SessionPersistence  bool          `json:"sessionPersistence,omitempty"`
	BrowserInBrowser    bool          `json:"browserInBrowser,omitempty"`
	PersistIsolationBar bool          `json:"persistIsolationBar,omitempty"`
	ForwardToZia        *ForwardToZia `json:"forwardToZia,omitempty"`
}

type ForwardToZia struct {
	Enabled        bool   `json:"enabled,omitempty"`
	OrganizationID bool   `json:"organizationId,omitempty"`
	CloudName      string `json:"cloudName,omitempty"`
	PacFileUrl     string `json:"pacFileUrl,omitempty"`
}

type Watermark struct {
	Enabled       bool   `json:"enabled,omitempty"`
	ShowUserID    bool   `json:"showUserId,omitempty"`
	ShowTimestamp bool   `json:"showTimestamp,omitempty"`
	ShowMessage   bool   `json:"showMessage,omitempty"`
	Message       string `json:"message,omitempty"`
}
type SecurityControls struct {
	DocumentViewer     bool      `json:"documentViewer,omitempty"`
	AllowPrinting      bool      `json:"allowPrinting,omitempty"`
	Watermark          Watermark `json:"watermark,omitempty"`
	FlattenedPdf       bool      `json:"flattenedPdf,omitempty"`
	UploadDownload     string    `json:"uploadDownload,omitempty"`
	RestrictKeystrokes bool      `json:"restrictKeystrokes,omitempty"`
	CopyPaste          string    `json:"copyPaste,omitempty"`
	LocalRender        bool      `json:"localRender,omitempty"`
}

type Regions struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type DebugMode struct {
	Allowed      bool   `json:"allowed,omitempty"`
	FilePassword string `json:"filePassword,omitempty"`
}

func Get(service *services.Service, profileID string) (*IsolationProfile, *http.Response, error) {
	v := new(IsolationProfile)
	relativeURL := fmt.Sprintf("%s/%s", cbiConfig+service.Client.Config.CustomerID+cbiProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(service *services.Service, profileName string) (*IsolationProfile, *http.Response, error) {
	list, resp, err := GetAll(service)
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

func Create(service *services.Service, cbiProfile *IsolationProfile) (*IsolationProfile, *http.Response, error) {
	v := new(IsolationProfile)
	resp, err := service.Client.NewRequestDo("POST", cbiConfig+service.Client.Config.CustomerID+cbiProfileEndpoint, nil, cbiProfile, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(service *services.Service, profileID string, segmentGroupRequest *IsolationProfile) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.Config.CustomerID+cbiProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, segmentGroupRequest, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(service *services.Service, profileID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.Config.CustomerID+cbiProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(service *services.Service) ([]IsolationProfile, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.Config.CustomerID + cbiProfileEndpoint
	var list []IsolationProfile
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &list)
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
