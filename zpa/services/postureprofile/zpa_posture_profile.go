package postureprofile

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v1/zpa/services/common"
)

const (
	mgmtConfig             = "/mgmtconfig/v2/admin/customers/"
	postureProfileEndpoint = "/posture"
)

// PostureProfile ...
type PostureProfile struct {
	CreationTime      string `json:"creationTime,omitempty"`
	Domain            string `json:"domain,omitempty"`
	ID                string `json:"id,omitempty"`
	MasterCustomerID  string `json:"masterCustomerId,omitempty"`
	ModifiedBy        string `json:"modifiedBy,omitempty"`
	ModifiedTime      string `json:"modifiedTime,omitempty"`
	Name              string `json:"name,omitempty"`
	PostureudID       string `json:"postureUdid,omitempty"`
	ZscalerCloud      string `json:"zscalerCloud,omitempty"`
	ZscalerCustomerID string `json:"zscalerCustomerId,omitempty"`
}

func (service *Service) Get(id string) (*PostureProfile, *http.Response, error) {
	v := new(PostureProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+postureProfileEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByPostureUDID(postureUDID string) (*PostureProfile, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + postureProfileEndpoint)
	list, resp, err := common.GetAllPagesGeneric[PostureProfile](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, postureProfile := range list {
		if postureProfile.PostureudID == postureUDID {
			return &postureProfile, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no posture profile with postureUDID '%s' was found", postureUDID)
}

func (service *Service) GetByName(postureName string) (*PostureProfile, *http.Response, error) {
	adaptedPostureName := common.RemoveCloudSuffix(postureName)
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + postureProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[PostureProfile](service.Client, relativeURL, adaptedPostureName)
	if err != nil {
		return nil, nil, err
	}
	for _, postureProfile := range list {
		if strings.EqualFold(common.RemoveCloudSuffix(postureProfile.Name), adaptedPostureName) {
			return &postureProfile, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no posture profile named '%s' was found", postureName)
}

func (service *Service) GetAll() ([]PostureProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + postureProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[PostureProfile](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}

	return list, resp, nil
}
