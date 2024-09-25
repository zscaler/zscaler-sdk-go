package inspection_custom_controls

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig             = "/zpa/mgmtconfig/v1/admin/customers/"
	customControlsEndpoint = "/inspectionControls/custom"
)

type InspectionCustomControl struct {
	ID                               string                   `json:"id,omitempty"`
	Action                           string                   `json:"action,omitempty"`
	ActionValue                      string                   `json:"actionValue,omitempty"`
	AssociatedInspectionProfileNames []AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
	Rules                            []Rules                  `json:"rules,omitempty"`
	ControlNumber                    string                   `json:"controlNumber,omitempty"`
	ControlType                      string                   `json:"controlType,omitempty"`
	ControlRuleJson                  string                   `json:"controlRuleJson,omitempty"`
	CreationTime                     string                   `json:"creationTime,omitempty"`
	DefaultAction                    string                   `json:"defaultAction,omitempty"`
	DefaultActionValue               string                   `json:"defaultActionValue,omitempty"`
	Description                      string                   `json:"description,omitempty"`
	ModifiedBy                       string                   `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                   `json:"modifiedTime,omitempty"`
	Name                             string                   `json:"name,omitempty"`
	ParanoiaLevel                    string                   `json:"paranoiaLevel,omitempty"`
	ProtocolType                     string                   `json:"protocolType,omitempty"`
	Severity                         string                   `json:"severity,omitempty"`
	Type                             string                   `json:"type,omitempty"`
	Version                          string                   `json:"version,omitempty"`
}

type Rules struct {
	Conditions []Conditions `json:"conditions,omitempty"`
	Names      []string     `json:"names,omitempty"`
	Type       string       `json:"type,omitempty"`
}

type Conditions struct {
	LHS string `json:"lhs,omitempty"`
	OP  string `json:"op,omitempty"`
	RHS string `json:"rhs,omitempty"`
}

type AssociatedProfileNames struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func unmarshalRulesJson(rulesJsonStr string) ([]Rules, error) {
	var rules []Rules
	err := json.Unmarshal([]byte(rulesJsonStr), &rules)
	return rules, err
}

func Get(ctx context.Context, service *zscaler.Service, customID string) (*InspectionCustomControl, *http.Response, error) {
	v := new(InspectionCustomControl)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+customControlsEndpoint, customID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	rules, err := unmarshalRulesJson(v.ControlRuleJson)
	v.Rules = rules
	return v, resp, err
}

func GetByName(ctx context.Context, service *zscaler.Service, controlName string) (*InspectionCustomControl, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + customControlsEndpoint
	list, resp, err := common.GetAllPagesGeneric[InspectionCustomControl](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, control := range list {
		if strings.EqualFold(control.Name, controlName) {
			rules, err := unmarshalRulesJson(control.ControlRuleJson)
			control.Rules = rules
			return &control, resp, err
		}
	}
	return nil, resp, fmt.Errorf("no custom inspection control named '%s' was found", controlName)
}

func Create(ctx context.Context, service *zscaler.Service, customControls InspectionCustomControl) (*InspectionCustomControl, *http.Response, error) {
	v := new(InspectionCustomControl)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+customControlsEndpoint, nil, customControls, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, customID string, customControls *InspectionCustomControl) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+customControlsEndpoint, customID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, nil, customControls, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, customID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+customControlsEndpoint, customID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", relativeURL, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]InspectionCustomControl, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + customControlsEndpoint
	list, resp, err := common.GetAllPagesGeneric[InspectionCustomControl](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
