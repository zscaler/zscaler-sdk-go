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

// UserCodeRequest is the request structure for verifying user codes
type UserCodeRequest struct {
	ComponentGroupID string   `json:"componentGroupId,omitempty"`
	UserCodes        []string `json:"user_codes,omitempty"`
}

// UserCodeInfo represents a single user code with its status in the response
type UserCodeInfo struct {
	Code                 string `json:"code,omitempty"`
	UserCode             string `json:"userCode,omitempty"`
	Status               string `json:"status,omitempty"`
	ConfigCloudName      string `json:"configCloudName,omitempty"`
	EnrollmentServer     string `json:"enrollmentServer,omitempty"`
	NonceAssociationType string `json:"nonceAssociationType,omitempty"`
	TenantID             string `json:"tenantId,omitempty"`
	ZcomponentID         string `json:"zcomponentId,omitempty"`
}

// UserCodeResponse is the response structure from the user codes API
type UserCodeResponse struct {
	ComponentGroupID     string         `json:"componentGroupId,omitempty"`
	ConfigCloudName      string         `json:"configCloudName,omitempty"`
	EnrollmentServer     string         `json:"enrollmentServer,omitempty"`
	NonceAssociationType string         `json:"nonceAssociationType,omitempty"`
	TenantID             string         `json:"tenantId,omitempty"`
	UserCodes            []UserCodeInfo `json:"userCodes,omitempty"`
	ZcomponentID         string         `json:"zcomponentId,omitempty"`
}

// OauthUser is kept for backward compatibility
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

// VerifyUserCodes verifies the provided list of user codes for a given component provisioning.
func VerifyUserCodes(ctx context.Context, service *zscaler.Service, associationType string, request *UserCodeRequest) (*UserCodeResponse, *http.Response, error) {
	v := new(UserCodeResponse)
	relativeURL := fmt.Sprintf("%s%s/%s/usercodes", mgmtConfig, service.Client.GetCustomerID(), associationType)
	resp, err := service.Client.NewRequestDo(ctx, "POST", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, request, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// VerifyUserCodeStatus checks the status of user codes for the specified customer.
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
