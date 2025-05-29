package authdomain

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	mgmtConfig             = "/zpa/mgmtconfig/v1/admin/customers/"
	authDomainEndpoint     = "/authDomains"
	ancestorPolicyEndpoint = "/ancestorPolicy"
)

type AuthDomain struct {
	AuthDomains []string `json:"authDomains"`
}

type AncestorPolicy struct {
	AccessType     string           `json:"accessType,omitempty"`
	AccessMappings []AccessMappings `json:"accessMappings,omitempty"`
}

type AccessMappings struct {
	ID                 string `json:"id,omitempty"`
	CreationTime       string `json:"creationTime,omitempty"`
	ModifiedBy         string `json:"modifiedBy,omitempty"`
	ModifiedTime       string `json:"modifiedTime,omitempty"`
	AncestorCustomerID string `json:"ancestorCustomerId,omitempty"`
	RoleID             string `json:"roleId,omitempty"`
	CustomerID         string `json:"customerId,omitempty"`
}

func GetAllAuthDomains(ctx context.Context, service *zscaler.Service) (*AuthDomain, *http.Response, error) {
	v := new(AuthDomain)
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + authDomainEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetAncestorPolicy(ctx context.Context, service *zscaler.Service) (*AncestorPolicy, *http.Response, error) {
	v := new(AncestorPolicy)
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + ancestorPolicyEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Create(ctx context.Context, service *zscaler.Service, policy *AncestorPolicy) (*AncestorPolicy, *http.Response, error) {
	v := new(AncestorPolicy)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+ancestorPolicyEndpoint, policy, v, nil)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
