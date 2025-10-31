package oauth2_user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig = "/zpa/mgmtconfig/v1/admin/customers/"
)

// Specifies the Provisioning Key type for App Connectors or ZPA Private Service Edges or ZPA Private Cloud Controller. The supported values are CONNECTOR_GRP, NP_ASSISTANT_GRP, SITE_CONTROLLER_GRP and SERVICE_EDGE_GRP.
type OauthUser struct {
	ComponentGroupID     string   `json:"componentGroupId,omitempty"`
	ConfigCloudName      string   `json:"configCloudName,omitempty"`
	EnrollmentServer     string   `json:"enrollmentServer,omitempty"`
	NonceAssociationType string   `json:"nonceAssociationType,omitempty"`
	TenantID             string   `json:"tenantId,omitempty"`
	UserCodes            []string `json:"userCodes,omitempty"`
	ZcomponentID         string   `json:"zcomponentId,omitempty"`
}

type UserCodeStatusRequest struct {
	UserCodes []string `json:"userCodes"`
}

// Verifies the provided list of user codes for a given component provisioning.
func VerifyUserCodes(ctx context.Context, service *zscaler.Service, associationType string, oauthUser *OauthUser) (*OauthUser, *http.Response, error) {
	v := new(OauthUser)
	relativeURL := fmt.Sprintf("%s%s/%s/usercodes", mgmtConfig, service.Client.GetCustomerID(), associationType)
	resp, err := service.Client.NewRequestDo(ctx, "POST", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, oauthUser, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Adds a new Provisioning Key for the specified customer.
func VerifyUserCodeStatus(ctx context.Context, service *zscaler.Service, associationType string, userCodes []string) (*OauthUser, *http.Response, error) {
	v := new(OauthUser)
	relativeURL := fmt.Sprintf("%s%s/%s/usercodes/status", mgmtConfig, service.Client.GetCustomerID(), associationType)
	requestBody := UserCodeStatusRequest{UserCodes: userCodes}
	resp, err := service.Client.NewRequestDo(ctx, "POST", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, requestBody, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
