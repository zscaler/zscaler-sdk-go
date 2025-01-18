package admin_roles

import (
	"context"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	adminRolesEndpoint = "/zcc/papi/public/v1/getAdminRoles"
)

type AdminRole struct {
	AdminManagement              string `json:"adminManagement"`
	AdministratorGroup           string `json:"administratorGroup"`
	AndroidProfile               string `json:"androidProfile"`
	AppBypass                    string `json:"appBypass"`
	AppProfileGroup              string `json:"appProfileGroup"`
	AuditLogs                    string `json:"auditLogs"`
	AuthSetting                  string `json:"authSetting"`
	ClientConnectorAppStore      string `json:"clientConnectorAppStore"`
	ClientConnectorIDP           string `json:"clientConnectorIdp"`
	ClientConnectorNotifications string `json:"clientConnectorNotifications"`
	ClientConnectorSupport       string `json:"clientConnectorSupport"`
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

func GetAdminRoles(ctx context.Context, service *zscaler.Service, pageSize ...int) ([]AdminRole, error) {
	// Determine the pageSize to use (default if not provided)
	effectivePageSize := 0 // Default to let `NewPagination` handle it
	if len(pageSize) > 0 && pageSize[0] > 0 {
		effectivePageSize = pageSize[0]
	}

	// Construct empty query parameters (no userType required)
	queryParams := struct{}{}

	// Leverage ReadAllPages to handle pagination
	return common.ReadAllPages[AdminRole](ctx, service.Client, adminRolesEndpoint, queryParams, effectivePageSize)
}
