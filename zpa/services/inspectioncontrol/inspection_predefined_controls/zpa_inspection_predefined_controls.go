package inspection_predefined_controls

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig           = "/mgmtconfig/v1/admin/customers/"
	predControlsEndpoint = "/inspectionControls/predefined"
)

type PredefinedControls struct {
	ID                               string                          `json:"id,omitempty"`
	Name                             string                          `json:"name,omitempty"`
	Action                           string                          `json:"action,omitempty"`
	ActionValue                      string                          `json:"actionValue,omitempty"`
	AssociatedInspectionProfileNames []common.AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
	Attachment                       string                          `json:"attachment,omitempty"`
	ControlGroup                     string                          `json:"controlGroup,omitempty"`
	ControlType                      string                          `json:"controlType,omitempty"`
	ControlNumber                    string                          `json:"controlNumber,omitempty"`
	CreationTime                     string                          `json:"creationTime,omitempty"`
	DefaultAction                    string                          `json:"defaultAction,omitempty"`
	DefaultActionValue               string                          `json:"defaultActionValue,omitempty"`
	Description                      string                          `json:"description,omitempty"`
	ModifiedBy                       string                          `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                          `json:"modifiedTime,omitempty"`
	ParanoiaLevel                    string                          `json:"paranoiaLevel,omitempty"`
	ProtocolType                     string                          `json:"protocolType,omitempty"`
	Severity                         string                          `json:"severity,omitempty"`
	Version                          string                          `json:"version,omitempty"`
}

type ControlGroupItem struct {
	ControlGroup                 string               `json:"controlGroup,omitempty"`
	PredefinedInspectionControls []PredefinedControls `json:"predefinedInspectionControls,omitempty"`
	DefaultGroup                 bool                 `json:"defaultGroup,omitempty"`
}

type ControlsRequestFilters struct {
	Version  string `url:"version,omitempty"`
	Search   string `url:"search,omitempty"`
	PageSize int    `url:"pagesize,omitempty"`
	Page     int    `url:"page,omitempty"`
}

// Get Predefined Controls by ID
// https://help.zscaler.com/zpa/api-reference#/inspection-control-controller/getPredefinedControlById
func Get(service *services.Service, controlID string) (*PredefinedControls, *http.Response, error) {
	v := new(PredefinedControls)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+predControlsEndpoint, controlID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetAll(service *services.Service, version string) ([]PredefinedControls, error) {
	v := []ControlGroupItem{}
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + predControlsEndpoint)
	_, err := service.Client.NewRequestDo("GET", relativeURL, struct {
		Version string `url:"version"`
	}{Version: version}, nil, &v)
	if err != nil {
		return nil, err
	}
	predefinedControls := []PredefinedControls{}
	for _, group := range v {
		predefinedControls = append(predefinedControls, group.PredefinedInspectionControls...)
	}
	return predefinedControls, nil
}

func GetByName(service *services.Service, name, version string) (*PredefinedControls, *http.Response, error) {
	queryParams := url.Values{}
	queryParams.Set("version", version)

	if name != "" {
		search := fmt.Sprintf("name+EQ+%s", url.QueryEscape(name))
		queryParams.Set("search", search)
	}

	relativeURL := fmt.Sprintf("%s%s%s?%s", mgmtConfig, service.Client.Config.CustomerID, predControlsEndpoint, queryParams.Encode())

	var v []ControlGroupItem
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	for _, group := range v {
		for _, control := range group.PredefinedInspectionControls {
			if strings.EqualFold(control.Name, name) {
				service.Client.Config.Logger.Printf("[INFO] got predefined controls:%#v", v)
				return &control, resp, nil
			}
		}
	}
	service.Client.Config.Logger.Printf("[ERROR] no predefined control named '%s' found", name)
	return nil, resp, fmt.Errorf("no predefined control named '%s' found", name)
}

func GetAllByGroup(service *services.Service, version, groupName string) ([]PredefinedControls, error) {
	queryParams := url.Values{}
	queryParams.Set("version", version)

	if groupName != "" {
		// Properly encode the groupName component
		encodedGroupName := url.QueryEscape(groupName)
		search := fmt.Sprintf("controlGroup+EQ+%s", encodedGroupName)
		queryParams.Set("search", search)
	}

	relativeURL := fmt.Sprintf("%s%s%s?%s", mgmtConfig, service.Client.Config.CustomerID, predControlsEndpoint, queryParams.Encode())

	var v []ControlGroupItem
	_, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, err
	}

	for _, group := range v {
		if strings.EqualFold(group.ControlGroup, groupName) {
			return group.PredefinedInspectionControls, nil
		}
	}
	return []PredefinedControls{}, nil
}
