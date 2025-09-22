package c2c_ip_ranges

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig       = "/zpa/mgmtconfig/v1/admin/customers/"
	ipRangesEndpoint = "/v2/ipRanges"
)

type IPRanges struct {
	AvailableIps  string `json:"availableIps,omitempty"`
	CountryCode   string `json:"countryCode,omitempty"`
	CreationTime  string `json:"creationTime,omitempty"`
	CustomerId    string `json:"customerId,omitempty"`
	Description   string `json:"description,omitempty"`
	Enabled       bool   `json:"enabled,omitempty"`
	ID            string `json:"id,omitempty"`
	IpRangeBegin  string `json:"ipRangeBegin,omitempty"`
	IpRangeEnd    string `json:"ipRangeEnd,omitempty"`
	IsDeleted     string `json:"isDeleted,omitempty"`
	LatitudeInDb  string `json:"latitudeInDb,omitempty"`
	Location      string `json:"location,omitempty"`
	LocationHint  string `json:"locationHint,omitempty"`
	LongitudeInDb string `json:"longitudeInDb,omitempty"`
	ModifiedBy    string `json:"modifiedBy,omitempty"`
	ModifiedTime  string `json:"modifiedTime,omitempty"`
	Name          string `json:"name,omitempty"`
	SccmFlag      bool   `json:"sccmFlag,omitempty"`
	SubnetCidr    string `json:"subnetCidr,omitempty"`
	TotalIps      string `json:"totalIps,omitempty"`
	UsedIps       string `json:"usedIps,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, ipRangeID string) (*IPRanges, *http.Response, error) {
	v := new(IPRanges)
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+ipRangesEndpoint, ipRangeID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, portalName string) (*IPRanges, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + ipRangesEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[IPRanges](ctx, service.Client, relativeURL, common.Filter{Search: portalName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, portalName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no user portal named '%s' was found", portalName)
}

func Create(ctx context.Context, service *zscaler.Service, ipRange *IPRanges) (*IPRanges, *http.Response, error) {
	v := new(IPRanges)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+ipRangesEndpoint, nil, ipRange, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, ipRangeID string, ipRange *IPRanges) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+ipRangesEndpoint, ipRangeID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, nil, ipRange, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, ipRangeID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+ipRangesEndpoint, ipRangeID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]*IPRanges, *http.Response, error) {
	var v []*IPRanges
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + ipRangesEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
