package cbiprofilecontroller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	cbiConfig          = "/zpa/cbiconfig/cbi/api/customers/"
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
	SessionPersistence  bool          `json:"sessionPersistence"`
	BrowserInBrowser    bool          `json:"browserInBrowser"`
	PersistIsolationBar bool          `json:"persistIsolationBar"`
	Translate           bool          `json:"translate"`
	ZGPU                bool          `json:"zgpu,omitempty"`
	ForwardToZia        *ForwardToZia `json:"forwardToZia,omitempty"`
}

type ForwardToZia struct {
	Enabled        bool   `json:"enabled"`
	OrganizationID string `json:"organizationId"`
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
	DocumentViewer     bool       `json:"documentViewer,omitempty"`
	AllowPrinting      bool       `json:"allowPrinting,omitempty"`
	Watermark          *Watermark `json:"watermark,omitempty"`
	FlattenedPdf       bool       `json:"flattenedPdf,omitempty"`
	UploadDownload     string     `json:"uploadDownload,omitempty"`
	RestrictKeystrokes bool       `json:"restrictKeystrokes,omitempty"`
	CopyPaste          string     `json:"copyPaste,omitempty"`
	LocalRender        bool       `json:"localRender,omitempty"`
	DeepLink           *DeepLink  `json:"deepLink,omitempty"`
}

type DeepLink struct {
	Enabled      bool     `json:"enabled,omitempty"`
	Applications []string `json:"applications,omitempty"`
}

type Regions struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type DebugMode struct {
	Allowed      bool   `json:"allowed,omitempty"`
	FilePassword string `json:"filePassword,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, profileID string) (*IsolationProfile, *http.Response, error) {
	v := new(IsolationProfile)
	relativeURL := fmt.Sprintf("%s/%s", cbiConfig+service.Client.GetCustomerID()+cbiProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByNameOrID(ctx context.Context, service *zscaler.Service, identifier string) (*IsolationProfile, *http.Response, error) {
	// Retrieve all profiles
	list, resp, err := GetAll(ctx, service)
	if err != nil {
		return nil, nil, err
	}

	// Try to find by ID
	for _, profile := range list {
		if profile.ID == identifier {
			return Get(ctx, service, profile.ID)
		}
	}

	// Try to find by name
	for _, profile := range list {
		if strings.EqualFold(profile.Name, identifier) {
			return Get(ctx, service, profile.ID)
		}
	}

	return nil, resp, fmt.Errorf("no isolation profile named or with ID '%s' was found", identifier)
}

func Create(ctx context.Context, service *zscaler.Service, cbiProfile *IsolationProfile) (*IsolationProfile, *http.Response, error) {
	v := new(IsolationProfile)
	resp, err := service.Client.NewRequestDo(ctx, "POST", cbiConfig+service.Client.GetCustomerID()+cbiProfileEndpoint, nil, cbiProfile, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, profileID string, segmentGroupRequest *IsolationProfile) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.GetCustomerID()+cbiProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, nil, segmentGroupRequest, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, profileID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.GetCustomerID()+cbiProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]IsolationProfile, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.GetCustomerID() + cbiProfileEndpoint
	var list []IsolationProfile
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &list)
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
