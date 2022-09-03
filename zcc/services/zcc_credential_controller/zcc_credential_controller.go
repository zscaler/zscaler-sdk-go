package credential_controller

import (
	"fmt"
	"net/http"
)

const (
	credentialEndpoint = "/cred/v1"
)

type CredentialController struct {
	CompanyCredentialsResponses []CompanyCredentialsResponses `json:"companyCredentialsResponses,omitempty"`
}

type CompanyCredentialsResponses struct {
	APIKey          string `json:"apiKey,omitempty"`
	ErrorCodeEnum   string `json:"errorCodeEnum,omitempty"`
	ID              int    `json:"id,omitempty"`
	JwtExpirySecs   int    `json:"jwtExpirySecs,omitempty"`
	Name            string `json:"name,omitempty"`
	Role            int    `json:"role,omitempty"`
	SecretKey       string `json:"secretKey,omitempty"`
	Status          int    `json:"status,omitempty"`
	UpdateTimestamp int    `json:"updateTimestamp,omitempty"`
}

func (service *Service) Get(credentialId int) (*CredentialController, error) {
	v := new(CredentialController)
	relativeURL := fmt.Sprintf("%s/%d", credentialEndpoint, credentialId)
	err := service.Client.Read(relativeURL, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (service *Service) GetAllErrors() ([]CredentialController, error) {
	var basicErrors []CredentialController
	err := service.Client.Read(credentialEndpoint, &basicErrors)
	return basicErrors, err
}

func (service *Service) Create(basicError CredentialController) (*CredentialController, error) {
	resp, err := service.Client.Create(credentialEndpoint, basicError)
	if err != nil {
		return nil, err
	}
	res, ok := resp.(*BasicErrorController)
	if !ok {
		return nil, fmt.Errorf("couldn't marshal response to a valid object: %#v", resp)
	}
	return res, nil
}

func (service *Service) Patch(credentialId int, basicError CredentialController) (*CredentialController, error) {
	path := fmt.Sprintf("%s/%d", credentialEndpoint, credentialId)
	resp, err := service.Client.UpdateWithPatch(path, basicError)
	if err != nil {
		return nil, err
	}
	res, _ := resp.(BasicErrorController)
	return &res, err
}

func (service *Service) Delete(credentialId int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", credentialEndpoint, credentialId))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
