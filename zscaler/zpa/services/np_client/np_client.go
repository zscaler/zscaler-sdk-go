package np_client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                = "/zpa/mgmtconfig/v1/admin/customers/"
	vpnConnectedUsersEndpoint = "/vpnConnectedUsers"
)

type NPClient struct {
	ClientIpAddress    string `json:"clientIpAddress,omitempty"`
	CommonName         string `json:"commonName,omitempty"`
	CreationTime       int    `json:"creationTime,omitempty"`
	DeviceState        int    `json:"deviceState,omitempty"`
	Id                 int    `json:"id,omitempty"`
	ModifiedBy         int    `json:"modifiedBy,omitempty"`
	ModifiedTime       int    `json:"modifiedTime,omitempty"`
	VpnServiceEdgeName string `json:"vpnServiceEdgeName,omitempty"`
	VpnServiceEdgeId   int    `json:"vpnServiceEdgeId,omitempty"`
	UserName           string `json:"UserName,omitempty"`
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]NPClient, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + vpnConnectedUsersEndpoint

	filter := common.Filter{}

	list, resp, err := common.GetAllPagesGenericWithCustomFilters[NPClient](ctx, service.Client, relativeURL, filter)
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, userName string) (*NPClient, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + vpnConnectedUsersEndpoint

	filter := common.Filter{}

	list, resp, err := common.GetAllPagesGenericWithCustomFilters[NPClient](ctx, service.Client, relativeURL, filter)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.UserName, userName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no username named '%s' was found", userName)
}
