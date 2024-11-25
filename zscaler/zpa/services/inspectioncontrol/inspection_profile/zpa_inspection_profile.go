package inspection_profile

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                = "/zpa/mgmtconfig/v1/admin/customers/"
	inspectionProfileEndpoint = "/inspectionProfile"
	defaultVersion            = "OWASP_CRS/3.3.0"
)

type PatchQuery struct {
	Version string `json:"version,omitempty" url:"version,omitempty"`
}

type InspectionProfile struct {
	ID                                string                        `json:"id,omitempty"`
	Name                              string                        `json:"name,omitempty"`
	Description                       string                        `json:"description,omitempty"`
	APIProfile                        bool                          `json:"apiProfile,omitempty"`
	OverrideAction                    string                        `json:"overrideAction,omitempty"`
	CommonGlobalOverrideActionsConfig map[string]interface{}        `json:"commonGlobalOverrideActionsConfig,omitempty"`
	CreationTime                      string                        `json:"creationTime,omitempty"`
	ZSDefinedControlChoice            string                        `json:"zsDefinedControlChoice,omitempty"`
	GlobalControlActions              []string                      `json:"globalControlActions,omitempty"`
	IncarnationNumber                 string                        `json:"incarnationNumber,omitempty"`
	ModifiedBy                        string                        `json:"modifiedBy,omitempty"`
	ModifiedTime                      string                        `json:"modifiedTime,omitempty"`
	ParanoiaLevel                     string                        `json:"paranoiaLevel,omitempty"`
	PredefinedControlsVersion         string                        `json:"predefinedControlsVersion,omitempty"`
	CheckControlDeploymentStatus      bool                          `json:"checkControlDeploymentStatus,omitempty"`
	ControlInfoResource               []ControlInfoResource         `json:"controlsInfo,omitempty"`
	CustomControls                    []InspectionCustomControl     `json:"customControls,omitempty"`
	PredefinedAPIControls             []common.CustomCommonControls `json:"predefinedApiControls,omitempty"`
	PredefinedControls                []common.CustomCommonControls `json:"predefinedControls,omitempty"`
	WebSocketControls                 []WebSocketControls           `json:"websocketControls,omitempty"`
	ThreatLabzControls                []ThreatLabzControls          `json:"threatlabzControls,omitempty"`
}

type ControlInfoResource struct {
	ControlType string `json:"controlType,omitempty"`
	Count       string `json:"count,omitempty"`
}

type InspectionCustomControl struct {
	Action                           string                          `json:"action,omitempty"`
	ActionValue                      string                          `json:"actionValue,omitempty"`
	ControlNumber                    string                          `json:"controlNumber,omitempty"`
	ControlRuleJson                  string                          `json:"controlRuleJson,omitempty"`
	ControlType                      string                          `json:"controlType,omitempty"`
	CreationTime                     string                          `json:"creationTime,omitempty"`
	DefaultAction                    string                          `json:"defaultAction,omitempty"`
	DefaultActionValue               string                          `json:"defaultActionValue,omitempty"`
	Description                      string                          `json:"description,omitempty"`
	ID                               string                          `json:"id,omitempty"`
	ModifiedBy                       string                          `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                          `json:"modifiedTime,omitempty"`
	Name                             string                          `json:"name,omitempty"`
	ProtocolType                     string                          `json:"protocolType,omitempty"`
	ParanoiaLevel                    string                          `json:"paranoiaLevel,omitempty"`
	Severity                         string                          `json:"severity,omitempty"`
	Type                             string                          `json:"type,omitempty"`
	Version                          string                          `json:"version,omitempty"`
	AssociatedInspectionProfileNames []common.AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
	Rules                            []common.Rules                  `json:"rules,omitempty"`
}

type AssociatedCustomers struct {
	CustomerID           string `json:"customerId,omitempty"`
	ExcludeConstellation bool   `json:"excludeConstellation,omitempty"`
	IsPartner            bool   `json:"isPartner,omitempty"`
	Name                 string `json:"name,omitempty"`
}

type ThreatLabzControls struct {
	ID                               string                          `json:"id,omitempty"`
	Name                             string                          `json:"name,omitempty"`
	Description                      string                          `json:"description,omitempty"`
	Enabled                          bool                            `json:"enabled,omitempty"`
	Action                           string                          `json:"action,omitempty"`
	ActionValue                      string                          `json:"actionValue,omitempty"`
	Attachment                       string                          `json:"attachment,omitempty"`
	ControlGroup                     string                          `json:"controlGroup,omitempty"`
	ControlNumber                    string                          `json:"controlNumber,omitempty"`
	ControlType                      string                          `json:"controlType,omitempty"`
	CreationTime                     string                          `json:"creationTime,omitempty"`
	DefaultAction                    string                          `json:"defaultAction,omitempty"`
	DefaultActionValue               string                          `json:"defaultActionValue,omitempty"`
	ModifiedBy                       string                          `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                          `json:"modifiedTime,omitempty"`
	ParanoiaLevel                    string                          `json:"paranoiaLevel,omitempty"`
	Severity                         string                          `json:"severity,omitempty"`
	Version                          string                          `json:"version,omitempty"`
	EngineVersion                    string                          `json:"engineVersion,omitempty"`
	LastDeploymentTime               string                          `json:"lastDeploymentTime,omitempty"`
	RuleDeploymentState              string                          `json:"ruleDeploymentState,omitempty"`
	RuleMetadata                     string                          `json:"ruleMetadata,omitempty"`
	RuleProcessor                    string                          `json:"ruleProcessor,omitempty"`
	RulesetName                      string                          `json:"rulesetName,omitempty"`
	RulesetVersion                   string                          `json:"rulesetVersion,omitempty"`
	ZscalerInfoUrl                   string                          `json:"zscalerInfoUrl,omitempty"`
	AssociatedCustomers              []AssociatedCustomers           `json:"associatedCustomers,omitempty"`
	AssociatedInspectionProfileNames []common.AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
}

type WebSocketControls struct {
	ID                               string                          `json:"id,omitempty"`
	Name                             string                          `json:"name,omitempty"`
	Description                      string                          `json:"description,omitempty"`
	Action                           string                          `json:"action,omitempty"`
	ActionValue                      string                          `json:"actionValue,omitempty"`
	ControlNumber                    string                          `json:"controlNumber,omitempty"`
	ControlType                      string                          `json:"controlType,omitempty"`
	CreationTime                     string                          `json:"creationTime,omitempty"`
	DefaultAction                    string                          `json:"defaultAction,omitempty"`
	DefaultActionValue               string                          `json:"defaultActionValue,omitempty"`
	ModifiedBy                       string                          `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                          `json:"modifiedTime,omitempty"`
	ParanoiaLevel                    string                          `json:"paranoiaLevel,omitempty"`
	Severity                         string                          `json:"severity,omitempty"`
	Version                          string                          `json:"version,omitempty"`
	ZSDefinedControlChoice           string                          `json:"zsDefinedControlChoice,omitempty"`
	AssociatedInspectionProfileNames []common.AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, profileID string) (*InspectionProfile, *http.Response, error) {
	v := new(InspectionProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func setVersion(inspectionProfile *InspectionProfile) {
	// make sure to set version
	if inspectionProfile.PredefinedControlsVersion == "" {
		found := false
		for _, control := range inspectionProfile.PredefinedControls {
			if control.Version != "" {
				found = true
				inspectionProfile.PredefinedControlsVersion = control.Version
				break
			}
		}
		if !found {
			inspectionProfile.PredefinedControlsVersion = defaultVersion
		}
	}
}

func GetByName(ctx context.Context, service *zscaler.Service, profileName string) (*InspectionProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + inspectionProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[InspectionProfile](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, inspection := range list {
		if strings.EqualFold(inspection.Name, profileName) {
			return &inspection, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no inspection profile named '%s' was found", profileName)
}

func Create(ctx context.Context, service *zscaler.Service, inspectionProfile InspectionProfile) (*InspectionProfile, *http.Response, error) {
	setVersion(&inspectionProfile)
	v := new(InspectionProfile)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, nil, inspectionProfile, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	setVersion(inspectionProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, nil, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func PutAssociate(ctx context.Context, service *zscaler.Service, profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	setVersion(inspectionProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID+"/associateAllPredefinedControls")
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, PatchQuery{Version: inspectionProfile.PredefinedControlsVersion}, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func PutDeassociate(ctx context.Context, service *zscaler.Service, profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	setVersion(inspectionProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID+"/deAssociateAllPredefinedControls")
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, PatchQuery{Version: inspectionProfile.PredefinedControlsVersion}, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Patch(ctx context.Context, service *zscaler.Service, profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	setVersion(inspectionProfile)
	relativeURL := fmt.Sprintf("%s/%s/patch", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo(ctx, "PATCH", relativeURL, PatchQuery{Version: inspectionProfile.PredefinedControlsVersion}, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, profileID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", relativeURL, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]InspectionProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + inspectionProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[InspectionProfile](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
