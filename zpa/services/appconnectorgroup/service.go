package appconnectorgroup

import (
	"github.com/zscaler/zscaler-sdk-go/zpa"
)

type Service struct {
	Client        *zpa.Client
	microTenantID *string
}

func New(c *zpa.Client) *Service {
	return &Service{Client: c}
}

func (service *Service) WithMircoTenant(microTenantID string) *Service {
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
