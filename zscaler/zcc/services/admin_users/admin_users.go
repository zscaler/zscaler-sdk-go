package admin_users

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	editAdminUserEndpoint        = "/zcc/papi/public/v1/editAdminUser"
	getAdminUserEndpoint         = "/zcc/papi/public/v1/getAdminUsers"
	getAdminUsersSyncEndpoint    = "/zcc/papi/public/v1/getAdminUsersSyncInfo"
	adminUsersZpaSyncEndpoint    = "/zcc/papi/public/v1/syncZpaAdminUsers"
	adminUsersZiaZdxSyncEndpoint = "/zcc/papi/public/v1/syncZiaZdxAdminUsers"
)

type AdminUser struct {
	AccountEnabled string `json:"accountEnabled"`
	CompanyID      string `json:"companyId"`
	CompanyRole    Role   `json:"companyRole"`
	EditEnabled    string `json:"editEnabled"`
	ID             int    `json:"id"`
	IsDefaultAdmin string `json:"isDefaultAdmin"`
	ServiceType    string `json:"serviceType"`
	UserName       string `json:"userName"`
}

type Role struct {
	AdminManagement              string `json:"adminManagement"`
	AdministratorGroup           string `json:"administratorGroup"`
	AndroidProfile               string `json:"androidProfile"`
	AppBypass                    string `json:"appBypass"`
	AppProfileGroup              string `json:"appProfileGroup"`
	AuditLogs                    string `json:"auditLogs"`
	AuthSetting                  string `json:"authSetting"`
	ClientConnectorAppStore      string `json:"clientConnectorAppStore"`
	ClientConnectorIDP           string `json:"clientConnectorIdp"`
	ClientConnectorSupport       string `json:"clientConnectorSupport"`
	ClientConnectorNotifications string `json:"clientConnectorNotifications"`
	CompanyID                    string `json:"companyId"`
	CreatedBy                    string `json:"createdBy"`
	Dashboard                    string `json:"dashboard"`
	DDILConfiguration            string `json:"ddilConfiguration"`
	DedicatedProxyPorts          string `json:"dedicatedProxyPorts"`
	DeviceGroups                 string `json:"deviceGroups"`
	DeviceOverview               string `json:"deviceOverview"`
	DevicePosture                string `json:"devicePosture"`
	EnrolledDevicesGroup         string `json:"enrolledDevicesGroup"`
	ForwardingProfile            string `json:"forwardingProfile"`
	ID                           string `json:"id"`
	IOSProfile                   string `json:"iosProfile"`
	IsEditable                   bool   `json:"isEditable"`
	LinuxProfile                 string `json:"linuxProfile"`
	MACProfile                   string `json:"macProfile"`
	MachineTunnel                string `json:"machineTunnel"`
	ObfuscateData                string `json:"obfuscateData"`
	PartnerDeviceOverview        string `json:"partnerDeviceOverview"`
	PublicAPI                    string `json:"publicApi"`
	RoleName                     string `json:"roleName"`
	TrustedNetwork               string `json:"trustedNetwork"`
	UpdatedBy                    string `json:"updatedBy"`
	UserAgent                    string `json:"userAgent"`
	WindowsProfile               string `json:"windowsProfile"`
	ZPAPartnerLogin              string `json:"zpaPartnerLogin"`
	ZscalerDeception             string `json:"zscalerDeception"`
	ZscalerEntitlement           string `json:"zscalerEntitlement"`
}

type SyncZiaZdxZpaAdminUsers struct {
	CompanyIDs         []int                  `json:"companyIds"`
	ErrorCode          string                 `json:"errorCode"`
	ErrorInfoArguments []string               `json:"errorInfoArguments"`
	ErrorMessage       string                 `json:"errorMessage"`
	ResponseData       map[string]interface{} `json:"responseData"`
	Success            string                 `json:"success"`
}

func GetAdminUsers(ctx context.Context, service *zscaler.Service, userType string, pageSize ...int) ([]AdminUser, error) {
	// Determine the pageSize to use (default if not provided)
	effectivePageSize := 0 // Default to let `NewPagination` handle it
	if len(pageSize) > 0 && pageSize[0] > 0 {
		effectivePageSize = pageSize[0]
	}

	// Construct query parameters with optional userType
	queryParams := struct {
		UserType string `url:"userType,omitempty"`
	}{
		UserType: userType,
	}

	// Leverage ReadAllPages to handle pagination
	return common.ReadAllPages[AdminUser](ctx, service.Client, getAdminUserEndpoint, queryParams, effectivePageSize)
}

func UpdateAdminUser(ctx context.Context, service *zscaler.Service, adminUser *AdminUser) (*AdminUser, error) {
	if adminUser == nil {
		return nil, errors.New("adminUser is required")
	}

	// Marshal the AdminUser struct into JSON
	body, err := json.Marshal(adminUser)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal admin user request: %w", err)
	}

	// Make the PUT request
	resp, err := service.Client.NewZccRequestDo(ctx, "PUT", editAdminUserEndpoint, nil, bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update admin user: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP response codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update admin user: received status code %d", resp.StatusCode)
	}

	// Decode the response body into an AdminUser struct
	var response AdminUser
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetAdminUserSyncInfo retrieves synchronization information for admin users.
func GetAdminUserSyncInfo(ctx context.Context, service *zscaler.Service) error {
	// Make the GET request
	resp, err := service.Client.NewZccRequestDo(ctx, "GET", getAdminUsersSyncEndpoint, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve admin user sync info: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP response codes
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to retrieve admin user sync info: received status code %d", resp.StatusCode)
	}

	// Since the API returns an empty JSON, simply return nil to indicate success
	return nil
}

func GetSyncZiaZdxAdminUsers(ctx context.Context, service *zscaler.Service) (*SyncZiaZdxZpaAdminUsers, error) {
	// Make the POST request
	resp, err := service.Client.NewZccRequestDo(ctx, "POST", adminUsersZiaZdxSyncEndpoint, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sync information: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP response codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve sync information: received status code %d", resp.StatusCode)
	}

	// Parse the response body into the SyncZiaZdxZpaAdminUsers struct
	var syncInfo SyncZiaZdxZpaAdminUsers
	err = json.NewDecoder(resp.Body).Decode(&syncInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to decode sync information response: %w", err)
	}

	return &syncInfo, nil
}

func GetSyncZpaAdminUsers(ctx context.Context, service *zscaler.Service) (*SyncZiaZdxZpaAdminUsers, error) {
	// Make the POST request
	resp, err := service.Client.NewZccRequestDo(ctx, "POST", adminUsersZpaSyncEndpoint, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sync information: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP response codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve sync information: received status code %d", resp.StatusCode)
	}

	// Parse the response body into the SyncZiaZdxZpaAdminUsers struct
	var syncInfo SyncZiaZdxZpaAdminUsers
	err = json.NewDecoder(resp.Body).Decode(&syncInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to decode sync information response: %w", err)
	}

	return &syncInfo, nil
}
