package basic_error_controller

import (
	"fmt"
	"net/http"
)

const (
	basicErrorEndpoint = "/error"
)

type BasicErrorController struct {
	Empty     bool                   `json:"empty,omitempty"`
	Model     map[string]interface{} `json:"model,omitempty"`
	ModelMap  map[string]interface{} `json:"modelMap,omitempty"`
	Reference bool                   `json:"reference,omitempty"`
	Status    []string               `json:"status,omitempty"`
	View      []string               `json:"view,omitempty"`
	ViewName  string                 `json:"viewName,omitempty"`
}

func (service *Service) Get(errorId int) (*BasicErrorController, error) {
	v := new(BasicErrorController)
	relativeURL := fmt.Sprintf("%s/%d", basicErrorEndpoint, errorId)
	err := service.Client.Read(relativeURL, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (service *Service) GetAllErrors() ([]BasicErrorController, error) {
	var basicErrors []BasicErrorController
	err := service.Client.Read(basicErrorEndpoint, &basicErrors)
	return basicErrors, err
}

func (service *Service) Create(basicError BasicErrorController) (*BasicErrorController, error) {
	resp, err := service.Client.Create(basicErrorEndpoint, basicError)
	if err != nil {
		return nil, err
	}
	res, ok := resp.(*BasicErrorController)
	if !ok {
		return nil, fmt.Errorf("couldn't marshal response to a valid object: %#v", resp)
	}
	return res, nil
}

/*
	func (service *Service) Patch(errorId int, basicError BasicErrorController) (*BasicErrorController, error) {
		path := fmt.Sprintf("%s/%d", basicErrorEndpoint, errorId)
		resp, err := service.Client.UpdateWithPatch(path, basicError)
		if err != nil {
			return nil, err
		}
		res, _ := resp.(BasicErrorController)
		return &res, err
	}
*/
func (service *Service) Delete(errorId int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", basicErrorEndpoint, errorId))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
