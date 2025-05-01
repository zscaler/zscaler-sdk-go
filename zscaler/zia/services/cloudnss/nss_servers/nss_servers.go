package nss_servers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	nssServersEndpoint = "/zia/api/v1/nssServers"
)

var (
	supportedServerTypes = map[string]bool{
		"NONE": true, "SOFTWARE_AA_FLAG": true, "NSS_FOR_WEB": true, "NSS_FOR_FIREWALL": true,
		"VZEN": true, "VZEN_SME": true, "VZEN_SMLB": true, "PINNED_NSS": true,
		"MD5_CAPABLE": true, "ADP": true, "ZIRSVR": true, "NSS_FOR_ZPA": true,
	}
)

type NSSServers struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
	State     string `json:"state,omitempty"`
	Type      string `json:"type,omitempty"`
	IcapSvrId int    `json:"icapSvrId,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, nssID int) (*NSSServers, error) {
	var nss NSSServers
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", nssServersEndpoint, nssID), &nss)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning nss server from Get: %d", nss.ID)
	return &nss, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, serverName string) (*NSSServers, error) {
	var nssServers []NSSServers
	err := common.ReadAllPages(ctx, service.Client, nssServersEndpoint, &nssServers)
	if err != nil {
		return nil, err
	}
	for _, nss := range nssServers {
		if strings.EqualFold(nss.Name, serverName) {
			return &nss, nil
		}
	}
	return nil, fmt.Errorf("no nss server found with name: %s", serverName)
}

func Create(ctx context.Context, service *zscaler.Service, nssServer *NSSServers) (*NSSServers, error) {

	resp, err := service.Client.Create(ctx, nssServersEndpoint, *nssServer)
	if err != nil {
		return nil, err
	}

	createdServers, ok := resp.(*NSSServers)
	if !ok {
		return nil, errors.New("object returned from api was not a nss server Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning rule from create: %d", createdServers.ID)
	return createdServers, nil
}

func Update(ctx context.Context, service *zscaler.Service, nssID int, nssServers *NSSServers) (*NSSServers, error) {

	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", nssServersEndpoint, nssID), *nssServers)
	if err != nil {
		return nil, err
	}

	updatedServers, ok := resp.(*NSSServers)
	if !ok {
		return nil, errors.New("object returned from api was not a nss server Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG] returning nss server from update: %d", updatedServers.ID)
	return updatedServers, nil
}

func Delete(ctx context.Context, service *zscaler.Service, nssID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", nssServersEndpoint, nssID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service, serverType *string) ([]NSSServers, error) {
	endpoint := nssServersEndpoint

	if serverType != nil {
		t := strings.ToUpper(strings.TrimSpace(*serverType))
		if !supportedServerTypes[t] {
			return nil, fmt.Errorf("invalid server type: %s", t)
		}
		endpoint += "?type=" + url.QueryEscape(t)
	}

	var nssServers []NSSServers
	err := common.ReadAllPages(ctx, service.Client, endpoint, &nssServers)
	return nssServers, err
}
