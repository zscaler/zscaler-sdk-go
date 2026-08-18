/*
Copyright (c) 2023, Zscaler Inc.

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
*/

package federated_application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtconfigV1                         = "/zpa/mgmtconfig/v1/customers/"
	federatedApplicationHostEndpoint     = "/application/host"
	federatedApplicationFederateEndpoint = "/application/federate"
)

type FederatedApplication struct {
	DomainNames   []string            `json:"domainNames,omitempty"`
	Enabled       bool                `json:"enabled,omitempty"`
	HostInfo      *common.PartnerInfo `json:"hostInfo,omitempty"`
	ID            string              `json:"id"`
	Name          string              `json:"name,omitempty"`
	TcpPortRanges []string            `json:"tcpPortRanges,omitempty"`
	UdpPortRanges []string            `json:"udpPortRanges,omitempty"`
}

type FederateRequest struct {
	ApplicationGid string   `json:"applicationGid,omitempty"`
	GuestGids      []string `json:"guestGids,omitempty"`
}

func GetAllHost(ctx context.Context, service *zscaler.Service, hostID string) ([]FederatedApplication, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtconfigV1+service.Client.GetCustomerID()+federatedApplicationHostEndpoint, hostID)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[FederatedApplication](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func UpdateFederate(ctx context.Context, service *zscaler.Service, body *FederateRequest) (*http.Response, error) {
	path := mgmtconfigV1 + service.Client.GetCustomerID() + federatedApplicationFederateEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, body, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}
