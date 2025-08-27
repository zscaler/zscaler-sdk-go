package c2cincidentreceiver

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	c2cIRReceiverEndpoint = "/zia/api/v1/cloudToCloudIR"
)

type C2CIncidentReceiver struct {
	ID                       int                      `json:"id,omitempty"`
	Name                     string                   `json:"name,omitempty"`
	Status                   []string                 `json:"status,omitempty"`
	ModifiedTime             int                      `json:"modifiedTime,omitempty"`
	LastTenantValidationTime int                      `json:"lastTenantValidationTime,omitempty"`
	LastValidationMsg        LastValidationMsg        `json:"lastValidationMsg,omitempty"`
	LastModifiedBy           *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`
	OnboardableEntity        *OnboardableEntity       `json:"onboardableEntity,omitempty"`
}

type LastValidationMsg struct {
	ErrorMsg  string `json:"errorMsg,omitempty"`
	ErrorCode string `json:"errorCode,omitempty"`
}

type OnboardableEntity struct {
	ID                      int                      `json:"id,omitempty"`
	Name                    string                   `json:"name,omitempty"`
	Type                    string                   `json:"type,omitempty"`
	EnterpriseTenantID      string                   `json:"enterpriseTenantId,omitempty"`
	Application             string                   `json:"application,omitempty"`
	LastValidationMsg       LastValidationMsg        `json:"lastValidationMsg,omitempty"`
	TenantAuthorizationInfo TenantAuthorizationInfo  `json:"tenantAuthorizationInfo,omitempty"`
	ZscalerAppTenantID      *common.IDNameExtensions `json:"zscalerAppTenantId,omitempty"`
}

type TenantAuthorizationInfo struct {
	AccessToken          string                    `json:"accessToken,omitempty"`
	BotToken             string                    `json:"botToken,omitempty"`
	RedirectUrl          string                    `json:"redirectUrl,omitempty"`
	Type                 string                    `json:"type,omitempty"`
	Env                  string                    `json:"env,omitempty"`
	TempAuthCode         string                    `json:"tempAuthCode,omitempty"`
	Subdomain            string                    `json:"subdomain,omitempty"`
	Apicp                string                    `json:"apicp,omitempty"`
	ClientID             string                    `json:"clientId,omitempty"`
	ClientSecret         string                    `json:"clientSecret,omitempty"`
	SecretToken          string                    `json:"secretToken,omitempty"`
	UserName             string                    `json:"userName,omitempty"`
	UserPwd              string                    `json:"userPwd,omitempty"`
	InstanceUrl          string                    `json:"instanceUrl,omitempty"`
	RoleArn              string                    `json:"roleArn,omitempty"`
	QuarantineBucketName string                    `json:"quarantineBucketName,omitempty"`
	CloudTrailBucketName string                    `json:"cloudTrailBucketName,omitempty"`
	BotID                string                    `json:"botId,omitempty"`
	OrgApiKey            string                    `json:"orgApiKey,omitempty"`
	ExternalID           string                    `json:"externalId,omitempty"`
	EnterpriseID         string                    `json:"enterpriseId,omitempty"`
	CredJson             string                    `json:"credJson,omitempty"`
	Role                 string                    `json:"role,omitempty"`
	OrganizationID       string                    `json:"organizationId,omitempty"`
	WorkspaceName        string                    `json:"workspaceName,omitempty"`
	WorkspaceID          string                    `json:"workspaceId,omitempty"`
	QtnChannelUrl        string                    `json:"qtnChannelUrl,omitempty"`
	FeaturesSupported    []string                  `json:"featuresSupported,omitempty"`
	MalQtnLibName        string                    `json:"malQtnLibName,omitempty"`
	DlpQtnLibName        string                    `json:"dlpQtnLibName,omitempty"`
	Credentials          string                    `json:"credentials,omitempty"`
	TokenEndpoint        string                    `json:"tokenEndpoint,omitempty"`
	RestApiEndpoint      string                    `json:"restApiEndpoint,omitempty"`
	SmirBucketConfig     []common.IDNameExtensions `json:"smirBucketConfig,omitempty"`
	QtnInfo              []interface{}             `json:"qtnInfo,omitempty"`
	QtnInfoCleared       bool                      `json:"qtnInfoCleared,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, receiverID int) (*C2CIncidentReceiver, error) {
	var c2cIRReceiver C2CIncidentReceiver
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", c2cIRReceiverEndpoint, receiverID), &c2cIRReceiver)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning incident receiver from Get: %d", c2cIRReceiver.ID)
	return &c2cIRReceiver, nil
}

func GetC2CIRName(ctx context.Context, service *zscaler.Service, recieverName string) (*C2CIncidentReceiver, error) {
	var c2cIRReceivers []C2CIncidentReceiver
	err := common.ReadAllPages(ctx, service.Client, c2cIRReceiverEndpoint, &c2cIRReceivers)
	if err != nil {
		return nil, err
	}
	for _, c2cIRReceiver := range c2cIRReceivers {
		if strings.EqualFold(c2cIRReceiver.Name, recieverName) {
			return &c2cIRReceiver, nil
		}
	}
	return nil, fmt.Errorf("no incident receiver found with name: %s", recieverName)
}

func ValidateDelete(ctx context.Context, service *zscaler.Service, receiverID int) (*C2CIncidentReceiver, error) {
	var c2cIRReceiver C2CIncidentReceiver
	err := service.Client.Read(ctx, fmt.Sprintf("%s/config/%d/validateDelete", c2cIRReceiverEndpoint, receiverID), &c2cIRReceiver)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning incident receiver from ValidateDelete: %d", c2cIRReceiver.ID)
	return &c2cIRReceiver, nil
}

func C2CIRCount(ctx context.Context, service *zscaler.Service, search string) (int, error) {
	var count int
	endpoint := c2cIRReceiverEndpoint + "/count"
	if search != "" {
		endpoint += "?search=" + url.QueryEscape(search)
	}
	err := service.Client.Read(ctx, endpoint, &count)
	return count, err
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]C2CIncidentReceiver, error) {
	var c2cIRReceivers []C2CIncidentReceiver
	err := common.ReadAllPages(ctx, service.Client, c2cIRReceiverEndpoint+"/lite", &c2cIRReceivers)
	return c2cIRReceivers, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]C2CIncidentReceiver, error) {
	var c2cIRReceivers []C2CIncidentReceiver
	err := common.ReadAllPages(ctx, service.Client, c2cIRReceiverEndpoint, &c2cIRReceivers)
	return c2cIRReceivers, err
}
