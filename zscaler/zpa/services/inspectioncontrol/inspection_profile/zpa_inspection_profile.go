package inspection_profile

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
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
	ID                                string                    `json:"id,omitempty"`
	CommonGlobalOverrideActionsConfig map[string]interface{}    `json:"commonGlobalOverrideActionsConfig,omitempty"`
	CreationTime                      string                    `json:"creationTime,omitempty"`
	ZSDefinedControlChoice            string                    `json:"zsDefinedControlChoice,omitempty"`
	Description                       string                    `json:"description,omitempty"`
	GlobalControlActions              []string                  `json:"globalControlActions,omitempty"`
	IncarnationNumber                 string                    `json:"incarnationNumber,omitempty"`
	ModifiedBy                        string                    `json:"modifiedBy,omitempty"`
	ModifiedTime                      string                    `json:"modifiedTime,omitempty"`
	Name                              string                    `json:"name,omitempty"`
	ParanoiaLevel                     string                    `json:"paranoiaLevel,omitempty"`
	PredefinedControlsVersion         string                    `json:"predefinedControlsVersion,omitempty"`
	CheckControlDeploymentStatus      bool                      `json:"checkControlDeploymentStatus,omitempty"`
	ControlInfoResource               []ControlInfoResource     `json:"controlsInfo,omitempty"`
	CustomControls                    []InspectionCustomControl `json:"customControls"`
	PredefinedControls                []CustomCommonControls    `json:"predefinedControls"`
	WebSocketControls                 []CustomCommonControls    `json:"websocketControls"`
	ThreatLabzControls                []ThreatLabzControls      `json:"threatlabzControls"`
}

type ControlInfoResource struct {
	ControlType string `json:"controlType,omitempty"`
	Count       string `json:"count,omitempty"`
}

type InspectionCustomControl struct {
	Action                           string                   `json:"action,omitempty"`
	ActionValue                      string                   `json:"actionValue,omitempty"`
	AssociatedInspectionProfileNames []AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
	Rules                            []common.Rules           `json:"rules,omitempty"`
	ControlNumber                    string                   `json:"controlNumber,omitempty"`
	ControlRuleJson                  string                   `json:"controlRuleJson,omitempty"`
	ControlType                      string                   `json:"controlType,omitempty"`
	CreationTime                     string                   `json:"creationTime,omitempty"`
	DefaultAction                    string                   `json:"defaultAction,omitempty"`
	DefaultActionValue               string                   `json:"defaultActionValue,omitempty"`
	Description                      string                   `json:"description,omitempty"`
	ID                               string                   `json:"id,omitempty"`
	ModifiedBy                       string                   `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                   `json:"modifiedTime,omitempty"`
	Name                             string                   `json:"name,omitempty"`
	ProtocolType                     string                   `json:"protocolType,omitempty"`
	ParanoiaLevel                    string                   `json:"paranoiaLevel,omitempty"`
	Severity                         string                   `json:"severity,omitempty"`
	Type                             string                   `json:"type,omitempty"`
	Version                          string                   `json:"version,omitempty"`
}

type CustomCommonControls struct {
	ID                               string                   `json:"id,omitempty"`
	Name                             string                   `json:"name,omitempty"`
	Action                           string                   `json:"action,omitempty"`
	ActionValue                      string                   `json:"actionValue,omitempty"`
	AssociatedInspectionProfileNames []AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
	Attachment                       string                   `json:"attachment,omitempty"`
	ControlGroup                     string                   `json:"controlGroup,omitempty"`
	ControlNumber                    string                   `json:"controlNumber,omitempty"`
	ControlType                      string                   `json:"controlType,omitempty"`
	CreationTime                     string                   `json:"creationTime,omitempty"`
	DefaultAction                    string                   `json:"defaultAction,omitempty"`
	DefaultActionValue               string                   `json:"defaultActionValue,omitempty"`
	Description                      string                   `json:"description,omitempty"`
	ModifiedBy                       string                   `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                   `json:"modifiedTime,omitempty"`
	ParanoiaLevel                    string                   `json:"paranoiaLevel,omitempty"`
	ProtocolType                     string                   `json:"protocolType,omitempty"`
	Severity                         string                   `json:"severity,omitempty"`
	Version                          string                   `json:"version,omitempty"`
}

type AssociatedProfileNames struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type AssociatedCustomers struct {
	CustomerID           string `json:"customerId,omitempty"`
	ExcludeConstellation bool   `json:"excludeConstellation,omitempty"`
	IsPartner            bool   `json:"isPartner,omitempty"`
	Name                 string `json:"name,omitempty"`
}

type ThreatLabzControls struct {
	ID                               string                   `json:"id,omitempty"`
	Name                             string                   `json:"name,omitempty"`
	Enabled                          bool                     `json:"enabled,omitempty"`
	Action                           string                   `json:"action,omitempty"`
	ActionValue                      string                   `json:"actionValue,omitempty"`
	AssociatedCustomers              []AssociatedCustomers    `json:"associatedCustomers,omitempty"`
	AssociatedInspectionProfileNames []AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
	Attachment                       string                   `json:"attachment,omitempty"`
	ControlGroup                     string                   `json:"controlGroup,omitempty"`
	ControlNumber                    string                   `json:"controlNumber,omitempty"`
	ControlType                      string                   `json:"controlType,omitempty"`
	CreationTime                     string                   `json:"creationTime,omitempty"`
	DefaultAction                    string                   `json:"defaultAction,omitempty"`
	DefaultActionValue               string                   `json:"defaultActionValue,omitempty"`
	Description                      string                   `json:"description,omitempty"`
	ModifiedBy                       string                   `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                   `json:"modifiedTime,omitempty"`
	ParanoiaLevel                    string                   `json:"paranoiaLevel,omitempty"`
	ProtocolType                     string                   `json:"protocolType,omitempty"`
	Severity                         string                   `json:"severity,omitempty"`
	Version                          string                   `json:"version,omitempty"`
	EngineVersion                    string                   `json:"engineVersion,omitempty"`
	LastDeploymentTime               string                   `json:"lastDeploymentTime,omitempty"`
	RuleDeploymentState              string                   `json:"ruleDeploymentState,omitempty"`
	RuleMetadata                     string                   `json:"ruleMetadata,omitempty"`
	RuleProcessor                    string                   `json:"ruleProcessor,omitempty"`
	RulesetName                      string                   `json:"rulesetName,omitempty"`
	RulesetVersion                   string                   `json:"rulesetVersion,omitempty"`
	ZscalerInfoUrl                   string                   `json:"zscalerInfoUrl,omitempty"`
}

func Get(service *zscaler.Service, profileID string) (*InspectionProfile, *http.Response, error) {
	v := new(InspectionProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
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

func GetByName(service *zscaler.Service, profileName string) (*InspectionProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + inspectionProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[InspectionProfile](service.Client, relativeURL, "")
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

func Create(service *zscaler.Service, inspectionProfile InspectionProfile) (*InspectionProfile, *http.Response, error) {
	setVersion(&inspectionProfile)
	v := new(InspectionProfile)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, nil, inspectionProfile, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(service *zscaler.Service, profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	setVersion(inspectionProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, nil, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func PutAssociate(service *zscaler.Service, profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	setVersion(inspectionProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID+"/associateAllPredefinedControls")
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, PatchQuery{Version: inspectionProfile.PredefinedControlsVersion}, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func PutDeassociate(service *zscaler.Service, profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	setVersion(inspectionProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID+"/deAssociateAllPredefinedControls")
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, PatchQuery{Version: inspectionProfile.PredefinedControlsVersion}, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Patch(service *zscaler.Service, profileID string, inspectionProfile *InspectionProfile) (*http.Response, error) {
	setVersion(inspectionProfile)
	relativeURL := fmt.Sprintf("%s/%s/patch", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("PATCH", relativeURL, PatchQuery{Version: inspectionProfile.PredefinedControlsVersion}, inspectionProfile, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(service *zscaler.Service, profileID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+inspectionProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(service *zscaler.Service) ([]InspectionProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + inspectionProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[InspectionProfile](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
