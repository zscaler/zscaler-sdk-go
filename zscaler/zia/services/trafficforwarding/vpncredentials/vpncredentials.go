package vpncredentials

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	vpnCredentialsEndpoint = "/zia/api/v1/vpnCredentials"
	maxBulkDeleteIDs       = 100
)

type VPNCredentials struct {
	// VPN credential id
	ID int `json:"id"`

	// VPN authentication type (i.e., how the VPN credential is sent to the server). It is not modifiable after VpnCredential is created.
	// Note: Zscaler no longer supports adding a new XAUTH VPN credential, but existing entries can be edited or deleted using the respective endpoints.
	Type string `json:"type,omitempty"`

	// Fully Qualified Domain Name. Applicable only to UFQDN or XAUTH (or HOSTED_MOBILE_USERS) auth type.
	FQDN string `json:"fqdn,omitempty"`

	// Static IP address for VPN that is self-provisioned or provisioned by Zscaler. This is a required field for IP auth type and is not applicable to other auth types.
	// Note: If you want Zscaler to provision static IP addresses for your organization, contact Zscaler Support.
	IPAddress string `json:"ipAddress,omitempty"`

	// Pre-shared key. This is a required field for UFQDN and IP auth type.
	PreSharedKey string `json:"preSharedKey,omitempty"`

	// Additional information about this VPN credential.
	Comments string `json:"comments,omitempty"`

	// Location that is associated to this VPN credential. Non-existence means not associated to any location.
	Location *Location `json:"location,omitempty"`

	// SD-WAN Partner that manages the location. If a partner does not manage the location, this is set to Self.
	ManagedBy *ManagedBy `json:"managedBy,omitempty"`
}

type Location struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type ManagedBy struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, vpnCredentialID int) (*VPNCredentials, error) {
	var vpnCredentials VPNCredentials
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", vpnCredentialsEndpoint, vpnCredentialID), &vpnCredentials)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning VPN Credentials from Get: %d", vpnCredentials.ID)
	return &vpnCredentials, nil
}

func GetVPNByType(ctx context.Context, service *zscaler.Service, vpnType string, includeOnlyWithoutLocation *bool, locationId *int, managedBy *int) ([]VPNCredentials, error) {
	queryParams := url.Values{}
	queryParams.Set("type", vpnType)

	if includeOnlyWithoutLocation != nil {
		queryParams.Set("includeOnlyWithoutLocation", strconv.FormatBool(*includeOnlyWithoutLocation))
	}
	if locationId != nil {
		queryParams.Set("locationId", strconv.Itoa(*locationId))
	}
	if managedBy != nil {
		queryParams.Set("managedBy", strconv.Itoa(*managedBy))
	}

	var vpnTypes []VPNCredentials
	err := common.ReadAllPages(ctx, service.Client, fmt.Sprintf("%s?%s", vpnCredentialsEndpoint, queryParams.Encode()), &vpnTypes)
	if err != nil {
		return nil, err
	}
	return vpnTypes, nil
}

func GetByFQDN(ctx context.Context, service *zscaler.Service, vpnCredentialName string) (*VPNCredentials, error) {
	var vpnCredentials []VPNCredentials

	err := common.ReadAllPages(ctx, service.Client, vpnCredentialsEndpoint, &vpnCredentials)
	if err != nil {
		return nil, err
	}
	for _, vpnCredential := range vpnCredentials {
		if strings.EqualFold(vpnCredential.FQDN, vpnCredentialName) {
			return &vpnCredential, nil
		}
	}
	return nil, fmt.Errorf("no vpn credentials found with fqdn: %s", vpnCredentialName)
}

func GetByIP(ctx context.Context, service *zscaler.Service, vpnCredentialIP string) (*VPNCredentials, error) {
	var vpnCredentials []VPNCredentials

	err := common.ReadAllPages(ctx, service.Client, vpnCredentialsEndpoint, &vpnCredentials)
	if err != nil {
		return nil, err
	}
	for _, vpnCredential := range vpnCredentials {
		if strings.EqualFold(vpnCredential.IPAddress, vpnCredentialIP) {
			return &vpnCredential, nil
		}
	}
	return nil, fmt.Errorf("no vpn credentials found with ip: %s", vpnCredentialIP)
}

func Create(ctx context.Context, service *zscaler.Service, vpnCredentials *VPNCredentials) (*VPNCredentials, *http.Response, error) {
	resp, err := service.Client.Create(ctx, vpnCredentialsEndpoint, *vpnCredentials)
	if err != nil {
		return nil, nil, err
	}

	createdVpnCredentials, ok := resp.(*VPNCredentials)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a vpn credential pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning vpn credential from create: %d", createdVpnCredentials.ID)
	return createdVpnCredentials, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, vpnCredentialID int, vpnCredentials *VPNCredentials) (*VPNCredentials, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", vpnCredentialsEndpoint, vpnCredentialID), *vpnCredentials)
	if err != nil {
		return nil, nil, err
	}
	updatedVpnCredentials, _ := resp.(*VPNCredentials)

	service.Client.GetLogger().Printf("[DEBUG]returning vpn credential from Update: %d", updatedVpnCredentials.ID)
	return updatedVpnCredentials, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, vpnCredentialID int) error {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", vpnCredentialsEndpoint, vpnCredentialID))
	if err != nil {
		return err
	}

	return nil
}

// BulkDeleteVPNCredentials sends a bulk delete request for VPN credentials.
func BulkDelete(ctx context.Context, service *zscaler.Service, ids []int) (*http.Response, error) {
	if len(ids) > maxBulkDeleteIDs {
		// Truncate the list to the first 100 IDs
		ids = ids[:maxBulkDeleteIDs]
		service.Client.GetLogger().Printf("[INFO] Truncating IDs list to the first %d items", maxBulkDeleteIDs)
	}

	// Define the payload
	payload := map[string][]int{
		"ids": ids,
	}

	// Call the generalized BulkDelete function from the client
	return service.Client.BulkDelete(ctx, vpnCredentialsEndpoint+"/bulkDelete", payload)
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]VPNCredentials, error) {
	var vpnTypes []VPNCredentials
	err := common.ReadAllPages(ctx, service.Client, vpnCredentialsEndpoint, &vpnTypes)
	return vpnTypes, err
}
