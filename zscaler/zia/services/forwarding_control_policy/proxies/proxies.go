package proxies

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	proxiesEndpoint   = "/zia/api/v1/proxies"
	ipGatewayEndpoint = "/zia/api/v1/dedicatedIPGateways/lite"
)

type Proxies struct {

	// Proxy ID
	ID int `json:"id,omitempty"`

	// Proxy name
	Name string `json:"name,omitempty"`

	// Gateway type - Supported Values: PROXYCHAIN, ZIA, ECSELF
	Type string `json:"type,omitempty"`

	// The IP address or the FQDN of the third-party proxy service
	Address string `json:"address,omitempty"`

	// The port number on which the third-party proxy service listens to the requests forwarded from Zscaler
	Port int `json:"port,omitempty"`

	// The root certificate used by the third-party proxy to perform SSL inspection.
	// This root certificate is used by Zscaler to validate the SSL leaf certificates signed by the upstream proxy.
	// The required root certificate appears in this drop-down menu only if it is uploaded from the Administration > Root Certificates page.
	Cert *common.IDNameExternalID `json:"cert,omitempty"`

	// Additional notes or information
	Description string `json:"description,omitempty"`

	// Flag indicating whether X-Authenticated-User header is added by the proxy. Enable to automatically insert authenticated user ID to the HTTP header, X-Authenticated-User.
	InsertXauHeader bool `json:"insertXauHeader,omitempty"`

	// Flag indicating whether the added X-Authenticated-User header is Base64 encoded. When enabled, the user ID is encoded using the Base64 encoding method.
	Base64EncodeXauHeader bool `json:"base64EncodeXauHeader,omitempty"`

	// Last user that modified the proxy
	LastModifiedBy *common.IDNameExternalID `json:"lastModifiedBy,omitempty"`

	// Timestamp of when the proxy was last modified
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`
}

type DedicatedIPGateways struct {
	Id                  int                      `json:"id,omitempty"`
	Name                string                   `json:"name,omitempty"`
	Description         string                   `json:"description,omitempty"`
	PrimaryDataCenter   *common.IDNameExtensions `json:"primaryDataCenter,omitempty"`
	SecondaryDataCenter *common.IDNameExtensions `json:"secondaryDataCenter,omitempty"`
	CreateTime          int                      `json:"createTime,omitempty"`
	LastModifiedTime    int                      `json:"lastModifiedTime,omitempty"`
	LastModifiedBy      *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`
	Default             bool                     `json:"default,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, proxyID int) (*Proxies, error) {
	var proxies Proxies
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", proxiesEndpoint, proxyID), &proxies)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning proxies from Get: %d", proxies.ID)
	return &proxies, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, proxyName string) (*Proxies, error) {
	var proxies []Proxies
	err := common.ReadAllPages(ctx, service.Client, proxiesEndpoint, &proxies)
	if err != nil {
		return nil, err
	}
	for _, proxy := range proxies {
		if strings.EqualFold(proxy.Name, proxyName) {
			return &proxy, nil
		}
	}
	return nil, fmt.Errorf("no proxy found with name: %s", proxyName)
}

func Create(ctx context.Context, service *zscaler.Service, proxyID *Proxies) (*Proxies, *http.Response, error) {
	resp, err := service.Client.Create(ctx, proxiesEndpoint, *proxyID)
	if err != nil {
		return nil, nil, err
	}

	createdProxies, ok := resp.(*Proxies)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a proxies pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new proxies from create: %d", createdProxies.ID)
	return createdProxies, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, proxyID int, proxies *Proxies) (*Proxies, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", proxiesEndpoint, proxyID), *proxies)
	if err != nil {
		return nil, nil, err
	}
	updatedProxies, _ := resp.(*Proxies)

	service.Client.GetLogger().Printf("[DEBUG]returning updates proxies from update: %d", updatedProxies.ID)
	return updatedProxies, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, proxyID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", proxiesEndpoint, proxyID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]Proxies, error) {
	var proxies []Proxies
	err := common.ReadAllPages(ctx, service.Client, proxiesEndpoint+"/lite", &proxies)
	return proxies, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]Proxies, error) {
	var proxies []Proxies
	err := common.ReadAllPages(ctx, service.Client, proxiesEndpoint, &proxies)
	return proxies, err
}

func GetDedicatedIPGWLite(ctx context.Context, service *zscaler.Service) ([]DedicatedIPGateways, error) {
	var gws []DedicatedIPGateways
	err := common.ReadAllPages(ctx, service.Client, ipGatewayEndpoint, &gws)
	return gws, err
}
