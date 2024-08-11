package vpncredentials

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	vpnCredentialsEndpoint = "/vpnCredentials"
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

func Get(service *services.Service, vpnCredentialID int) (*VPNCredentials, error) {
	var vpnCredentials VPNCredentials
	err := service.Client.Read(fmt.Sprintf("%s/%d", vpnCredentialsEndpoint, vpnCredentialID), &vpnCredentials)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning VPN Credentials from Get: %d", vpnCredentials.ID)
	return &vpnCredentials, nil
}

func GetVPNByType(service *services.Service, vpnType string) (*VPNCredentials, error) {
	var vpnTypes []VPNCredentials
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?type=%s", vpnCredentialsEndpoint, url.QueryEscape(vpnType)), &vpnTypes)
	if err != nil {
		return nil, err
	}
	for _, vpn := range vpnTypes {
		if strings.EqualFold(vpn.Type, vpnType) {
			return &vpn, nil
		}
	}
	return nil, fmt.Errorf("no VPN found with type: %s", vpnType)
}

func GetByFQDN(service *services.Service, vpnCredentialName string) (*VPNCredentials, error) {
	var vpnCredentials []VPNCredentials

	err := common.ReadAllPages(service.Client, vpnCredentialsEndpoint, &vpnCredentials)
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

func GetByIP(service *services.Service, vpnCredentialIP string) (*VPNCredentials, error) {
	var vpnCredentials []VPNCredentials

	err := common.ReadAllPages(service.Client, vpnCredentialsEndpoint, &vpnCredentials)
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

func Create(service *services.Service, vpnCredentials *VPNCredentials) (*VPNCredentials, *http.Response, error) {
	resp, err := service.Client.Create(vpnCredentialsEndpoint, *vpnCredentials)
	if err != nil {
		return nil, nil, err
	}

	createdVpnCredentials, ok := resp.(*VPNCredentials)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a vpn credential pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning vpn credential from create: %d", createdVpnCredentials.ID)
	return createdVpnCredentials, nil, nil
}

func Update(service *services.Service, vpnCredentialID int, vpnCredentials *VPNCredentials) (*VPNCredentials, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", vpnCredentialsEndpoint, vpnCredentialID), *vpnCredentials)
	if err != nil {
		return nil, nil, err
	}
	updatedVpnCredentials, _ := resp.(*VPNCredentials)

	service.Client.Logger.Printf("[DEBUG]returning vpn credential from Update: %d", updatedVpnCredentials.ID)
	return updatedVpnCredentials, nil, nil
}

func Delete(service *services.Service, vpnCredentialID int) error {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", vpnCredentialsEndpoint, vpnCredentialID))
	if err != nil {
		return err
	}

	return nil
}

// BulkDeleteVPNCredentials sends a bulk delete request for VPN credentials.
func BulkDelete(service *services.Service, ids []int) (*http.Response, error) {
	if len(ids) > maxBulkDeleteIDs {
		// Truncate the list to the first 100 IDs
		ids = ids[:maxBulkDeleteIDs]
		service.Client.Logger.Printf("[INFO] Truncating IDs list to the first %d items", maxBulkDeleteIDs)
	}

	// Define the payload
	payload := map[string][]int{
		"ids": ids,
	}

	// Call the generalized BulkDelete function from the client
	return service.Client.BulkDelete(vpnCredentialsEndpoint+"/bulkDelete", payload)
}

func GetAll(service *services.Service) ([]VPNCredentials, error) {
	var vpnTypes []VPNCredentials
	err := common.ReadAllPages(service.Client, vpnCredentialsEndpoint, &vpnTypes)
	return vpnTypes, err
}
