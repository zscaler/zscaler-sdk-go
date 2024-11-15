package services

import (
	"github.com/zscaler/zscaler-sdk-go/v2/zpa"
)

type Service struct {
	Client        *zpa.Client
	microTenantID *string
}

type ScimService struct {
	ScimClient *zpa.ScimClient
}

// NewScimService initializes a SCIM-based ZPA Service with *zpa.ScimConfig
func NewScimService(scimClient *zpa.ScimClient) *ScimService {
	return &ScimService{ScimClient: scimClient}
}

func New(c *zpa.Client) *Service {
	return &Service{Client: c}
}

func (service *Service) WithMicroTenant(microTenantID string) *Service {
	var mid *string
	if microTenantID != "" {
		mid_ := microTenantID
		mid = &mid_
	}
	return &Service{
		Client:        service.Client,
		microTenantID: mid,
	}
}

func (service *Service) MicroTenantID() *string {
	return service.microTenantID
}
